package directory

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"sync"
)

type Directory struct {
	mu            *sync.Mutex
	nextSegmentID SegmentID

	Path     string
	Segments []Segment
}

func New(path string) *Directory {
	return &Directory{
		mu:   &sync.Mutex{},
		Path: path,
	}
}

// Init walks through segment files in the WAL directory
// on disk and initializes in-memory state (i.e the list
// of segments, their metadata).
// Must be called before any other operations with the
// directory.
func (d *Directory) Init() error {
	entries, err := os.ReadDir(d.Path)
	if err != nil {
		return fmt.Errorf("os.ReadDir: %w", err)
	}
	for _, entry := range entries {
		submatches := segmentNameRegexp.FindStringSubmatch(entry.Name())
		if submatches == nil {
			continue
		}
		segmentID, err := segmentIDFromString(submatches[1])
		if err != nil {
			return err
		}
		segmentPath := d.segmentPath(segmentID)
		stat, err := os.Stat(segmentPath)
		if err != nil {
			return fmt.Errorf("os.Stat: %w", err)
		}
		d.Segments = append(d.Segments, Segment{
			ID:   segmentID,
			Path: d.segmentPath(segmentID),
			Size: stat.Size(),
		})
	}
	if len(d.Segments) == 0 {
		segmentID := d.nextSegmentID
		d.Segments = append(d.Segments, Segment{
			ID:     segmentID,
			Path:   d.segmentPath(segmentID),
			Status: SegmentStatusActive,
		})
	} else {
		slices.SortFunc(d.Segments, func(x, y Segment) int {
			if x.ID < y.ID {
				return -1
			}
			return 1
		})
		activeSegment := &d.Segments[len(d.Segments)-1]
		activeSegment.Status = SegmentStatusActive
		d.nextSegmentID = activeSegment.ID + 1
	}
	return nil
}

func (d *Directory) WithActiveSegment(cb func(seg *Segment)) {
	activeSegment := func() *Segment {
		d.mu.Lock()
		defer d.mu.Unlock()

		var i int
		for i = 0; i < len(d.Segments); i++ {
			if d.Segments[i].Status == SegmentStatusActive {
				break
			}
		}
		d.Segments[i].Refcount++
		return &d.Segments[i]
	}()

	cleanup := func() {
		d.mu.Lock()
		defer d.mu.Unlock()

		activeSegment.Refcount--
	}
	defer cleanup()

	cb(activeSegment)
}

func (d *Directory) segmentPath(id SegmentID) string {
	return filepath.Join(d.Path, segmentFilename(id))
}

func segmentFilename(id SegmentID) string {
	return fmt.Sprintf(segmentNameTemplate, id)
}

// Which file operations are atomic and can be used for replacing
// multiple segment with a single merged one?
// https://rcrowley.org/2010/01/06/things-unix-can-do-atomically.html
// see man rename. The strategy is to replace the latest segment in
// merged range with the merged version, then delete older segments.
// This doesn't lose data and is retriable (all the segments are still valid
// and represent the index).
var segmentNameRegexp = regexp.MustCompile(`^wal_([0-9]+)$`)
var segmentNameTemplate = "wal_%d"

func init() {
	exampleSegmentFilename := segmentFilename(SegmentID(0))
	if !segmentNameRegexp.MatchString(exampleSegmentFilename) {
		panic(fmt.Sprintf(
			"exampleSegmentFilename doesn't match segmentNameRegexp (`%s` does not match `%s`)",
			exampleSegmentFilename,
			segmentNameRegexp,
		),
		)
	}
}
