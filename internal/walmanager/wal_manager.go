package walmanager

import (
	"encoding/binary"
	"errors"
	"fmt"
	"os"

	"github.com/kolosovi/toy-kv/internal/pb/wal"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// Assumes that writes are atomic, even though they aren't.
// TODO error handling:
//  1. what should I do if I cannot write a particular log? It seems that I
//     must track last valid offset in the log & rewind to that offset on errors,
//     otherwise a bad write might affect all subsequent ones.
//  2. the following situation is possible: user tries to write & gets an error
//     because of an I/O problem, but the write actually reaches disk. Reads
//     of the same key return ErrNotFound because the in-memory index wasn't
//     updated. Later, when DAO is restarted, the key mysteriously appears.
//     What do people do with this problem? It seems that writes should be
//     waiting for their WAL commit in some queue. This implies that there
//     is a total order of all writes. Otherwise the following is possible:
//     user tries to make a write w1 and fails, then user successfuly makes
//     write w2, then w1 is finally acknowledged and becomes visible.
//     Not sure if this should be avoided.
type WALManager struct {
	config config

	f *os.File
}

func New(opts ...Option) *WALManager {
	return &WALManager{config: newConfig(opts...)}
}

func (m *WALManager) Start() (err error) {
	m.f, err = os.OpenFile(
		m.config.walFilename,
		os.O_RDWR|os.O_CREATE|os.O_APPEND,
		os.FileMode(0644),
	)
	if err != nil {
		return fmt.Errorf("cannot open WAL file: %w", err)
	}
	return nil
}

func (m *WALManager) Stop() error {
	if m.f == nil {
		return nil
	}
	if err := m.f.Sync(); err != nil {
		return fmt.Errorf("cannot sync WAL file to disk: %v", err)
	}
	if err := m.f.Close(); err != nil {
		return fmt.Errorf("cannot close WAL file: %v", err)
	}
	return nil
}

func (m *WALManager) WriteInsert(i Insert) error {
	dto := &wal.Log{
		Log: &wal.Log_Insert{
			Insert: &wal.Insert{
				Kv: &wal.KV{
					K: []byte(i.K),
					V: []byte(i.V),
				},
			},
		},
	}
	return m.writeLog(dto)
}

func (m *WALManager) WriteDelete(i Delete) error {
	dto := &wal.Log{
		Log: &wal.Log_Delete{
			Delete: &wal.Delete{K: []byte(i.K)},
		},
	}
	return m.writeLog(dto)
}

func (m *WALManager) writeLog(log protoreflect.ProtoMessage) error {
	logBytes, err := proto.Marshal(log)
	if err != nil {
		return fmt.Errorf("cannot marshal log: %w", err)
	}
	buf := make([]byte, sizeofLogSize+len(logBytes))
	binary.LittleEndian.PutUint32(buf, uint32(len(logBytes)))
	copy(buf[sizeofLogSize:], logBytes)
	for len(buf) != 0 {
		n, err := m.f.Write(buf)
		if err != nil {
			return fmt.Errorf("WAL file write: %w", err)
		}
		buf = buf[n:]
	}
	if err = m.f.Sync(); err != nil {
		return fmt.Errorf("WAL file sync: %w", err)
	}
	return nil
}

const sizeofLogSize = 4

func (m *WALManager) Reader() (WALReader, error) {
	f, err := os.OpenFile(
		m.config.walFilename,
		os.O_RDONLY,
		0,
	)
	if err != nil {
		return nil, fmt.Errorf("cannot open WAL file for reading: %w", err)
	}
	r, err := newReader(f)
	if err != nil {
		return nil, fmt.Errorf("newReader: %w", err)
	}
	return r, nil
}

type WALReader interface {
	HasNext() bool
	Scan(log *Log) error
	Close() error
}

var ErrNoLogs = errors.New("no more logs")
