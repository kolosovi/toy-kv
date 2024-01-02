package walmanager

import (
	"encoding/binary"
	"fmt"
	"os"

	"github.com/kolosovi/toy-kv/internal/pb/wal"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// Assumes that writes are atomic, even though they aren't.
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
	dto := &wal.Insert{
		Kv: &wal.KV{
			K: []byte(i.K),
			V: []byte(i.V),
		},
	}
	return m.writeLog(dto)
}

func (m *WALManager) WriteDelete(i Delete) error {
	dto := &wal.Delete{K: []byte(i.K)}
	return m.writeLog(dto)
}

func (m *WALManager) writeLog(log protoreflect.ProtoMessage) error {
	logBytes, err := proto.Marshal(log)
	if err != nil {
		return fmt.Errorf("cannot marshal log: %w", err)
	}
	buf := make([]byte, sizeofLogSize + len(logBytes))
	binary.LittleEndian.PutUint32(buf, uint32(len(logBytes)))
	copy(buf[sizeofLogSize:], logBytes)
	for len(buf) != 0 {
		n, err := m.f.Write(buf)
		if err != nil {
			return fmt.Errorf("WAL file write: %w", err)
		}
		buf = buf[n:]
	}
	return nil
}

const sizeofLogSize = 4

func (m *WALManager) Reader() (WALReader, error) {
	return &reader{}, nil
}

type WALReader interface {
	HasNext() bool
	Scan(log *Log) error
	Err() error
}