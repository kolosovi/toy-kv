package walmanager

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"

	toykv "github.com/kolosovi/toy-kv"
	"github.com/kolosovi/toy-kv/internal/pb/wal"
	"google.golang.org/protobuf/proto"
)

type reader struct {
	f           *os.File
	buf         []byte
	nextLogSize int
}

func newReader(f *os.File) (*reader, error) {
	r := &reader{f: f, buf: []byte{}}
	if err := r.peek(); err != nil {
		return nil, fmt.Errorf("peek: %w", err)
	}
	return r, nil
}

func (r *reader) HasNext() bool {
	return r.nextLogSize > 0
}

/*
peek

Scan:
if logLen == 0: ErrNoLogs
otherwise:
1. read logLen bytes into buffer
2. parse log
3. reset buffer
4. peek (sets logLen)
*/
func (r *reader) Scan(log *Log) error {
	if !r.HasNext() {
		return ErrNoLogs
	}
	r.ensureBufLength(r.nextLogSize)
	dst := r.buf[:r.nextLogSize]
	for len(dst) > 0 {
		readCount, err := r.f.Read(dst)
		dst = dst[readCount:]
		if err == io.EOF {
			break
		}
	}
	if len(dst) != 0 {
		return fmt.Errorf(
			"incomplete log message (expected %d bytes, found %d)",
			r.nextLogSize,
			r.nextLogSize-len(dst),
		)
	}
	dto := &wal.Log{}
	if err := proto.Unmarshal(r.buf[:r.nextLogSize], dto); err != nil {
		return fmt.Errorf("proto.Unmarshal: %w", err)
	}
	switch dtoLog := dto.Log.(type) {
	case *wal.Log_Insert:
		log.Typ = LogTypeInsert
		log.Insert = Insert{
			Record: toykv.Record{
				K: toykv.K(dtoLog.Insert.Kv.K),
				V: toykv.V(dtoLog.Insert.Kv.V),
			},
		}
	case *wal.Log_Delete:
		log.Typ = LogTypeDelete
		log.Delete = Delete{
			K: toykv.K(dtoLog.Delete.K),
		}
	default:
		return fmt.Errorf("unknown log type %T", dto.Log)
	}
	if err := r.peek(); err != nil {
		return fmt.Errorf("peek: %w", err)
	}
	return nil
}

func (r *reader) peek() error {
	r.ensureBufLength(sizeofLogSize)
	dst := r.buf[:sizeofLogSize]
	for len(dst) > 0 {
		readCount, err := r.f.Read(dst)
		dst = dst[readCount:]
		if err == io.EOF {
			break
		}
	}
	if len(dst) == sizeofLogSize {
		r.nextLogSize = 0
		return nil
	}
	if len(dst) != 0 {
		return fmt.Errorf("incomplete log size (%d more bytes to read)", len(dst))
	}
	r.nextLogSize = int(binary.LittleEndian.Uint32(r.buf[:sizeofLogSize]))
	return nil
}

func (r *reader) ensureBufLength(n int) {
	if cap(r.buf) >= n {
		r.buf = r.buf[:n]
	} else {
		newBuf := make([]byte, n)
		copy(newBuf, r.buf)
		r.buf = newBuf
	}
}

func (r *reader) Close() error {
	if err := r.f.Close(); err != nil {
		return fmt.Errorf("cannot close WAL file opened for reading: %w", err)
	}
	return nil
}
