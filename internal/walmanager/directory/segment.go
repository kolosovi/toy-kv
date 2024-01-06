package directory

import (
	"fmt"
	"strconv"
)

type Segment struct {
	ID       SegmentID
	Path     string
	Status   SegmentStatus
	Size     int64
	Refcount int
}

type SegmentID int64

func segmentIDFromString(input string) (SegmentID, error) {
	wrapErr := func(err error) error {
		return fmt.Errorf("cannot parse SegmentID from [%s]: %w", input, err)
	}

	intID, err := strconv.ParseInt(input, 10, 64)
	if err != nil {
		return 0, wrapErr(fmt.Errorf("strconv.ParseInt: %w", err))
	}
	return SegmentID(intID), nil
}

type SegmentStatus uint8

const (
	SegmentStatusRetired SegmentStatus = iota
	SegmentStatusActive
)
