package dao

import (
	"fmt"

	"github.com/kolosovi/toy-kv"
	"github.com/kolosovi/toy-kv/internal/walmanager"
)

type DAO struct {
	index map[toykv.K]toykv.V
	walManager *walmanager.WALManager
}

func New(walManager *walmanager.WALManager) *DAO {
	return &DAO{
		index: make(map[toykv.K]toykv.V, 0),
		walManager: walManager,
	}
}

func (d *DAO) Start() error {
	return d.walManager.Start()
}

func (d *DAO) Stop() error {
	return d.walManager.Stop()
}

func (d *DAO) Put(r toykv.Record) error {
	log := walmanager.Insert{Record: r}
	if err := d.walManager.WriteInsert(log); err != nil {
		return fmt.Errorf("walManager.WriteInsert: %w", err)
	}
	d.index[r.K] = r.V
	return nil
}

func (d *DAO) Get(k toykv.K) (toykv.V, error) {
	if v, ok := d.index[k]; ok {
		return v, nil
	}
	return "", toykv.ErrNotFound
}

func (d *DAO) Delete(k toykv.K) error {
	log := walmanager.Delete{K: k}
	if err := d.walManager.WriteDelete(log); err != nil {
		return fmt.Errorf("walManager.WriteDelete: %w", err)
	}
	if _, ok := d.index[k]; ok {
		delete(d.index, k)
		return nil
	}
	return toykv.ErrNotFound
}