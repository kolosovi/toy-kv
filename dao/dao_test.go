package dao

import (
	"fmt"
	"math/rand"
	"os"
	"testing"

	"github.com/kolosovi/toy-kv"
	"github.com/kolosovi/toy-kv/internal/walmanager"
	"github.com/stretchr/testify/require"
)

func TestOperations(t *testing.T) {
	dao := New(walmanager.New(walmanager.WithWALFilename(newFilename())))
	require.NoError(t, dao.Start())
	t.Cleanup(func() {
		require.NoError(t, dao.Stop())
	})

	// No keys, Get & Delete must return ErrNotFound
	_, err := dao.Get("foo")
	require.ErrorIs(t, err, toykv.ErrNotFound, "not ErrNotFound: %v", err)
	err = dao.Delete("foo")
	require.ErrorIs(t, err, toykv.ErrNotFound, "not ErrNotFound: %v", err)

	// After Put, Get must return inserted value
	require.NoError(t, dao.Put(toykv.Record{K: "foo", V: "foo_value"}))
	v, err := dao.Get("foo")
	require.NoError(t, err)
	require.Equal(t, toykv.V("foo_value"), v)

	// After second Put, Get must return updated value
	require.NoError(t, dao.Put(toykv.Record{K: "foo", V: "foo_value_new"}))
	v, err = dao.Get("foo")
	require.NoError(t, err)
	require.Equal(t, toykv.V("foo_value_new"), v)

	// Delete must return no error
	require.NoError(t, dao.Delete("foo"))

	// After Delete, Get and Delete of the same key must return ErrNotFound
	_, err = dao.Get("foo")
	require.ErrorIs(t, err, toykv.ErrNotFound, "not ErrNotFound: %v", err)
	err = dao.Delete("foo")
	require.ErrorIs(t, err, toykv.ErrNotFound, "not ErrNotFound: %v", err)
}

func TestReproduce(t *testing.T) {
	t.Skip()
	dao := New(walmanager.New(walmanager.WithWALFilename("test_4022139208788608928")))
	require.NoError(t, dao.Start())
	t.Cleanup(func() {
		require.NoError(t, dao.Stop())
	})

	v, err := dao.Get("foo")
	require.NoError(t, err)
	require.Equal(t, toykv.V("foo_value"), v)
}

func TestPersistence(t *testing.T) {
	var walFilename = newFilename()
	t.Run("record must be persisted to disk", func(t *testing.T) {
		dao := New(walmanager.New(walmanager.WithWALFilename(walFilename)))
		require.NoError(t, dao.Start())
		t.Cleanup(func() {
			require.NoError(t, dao.Stop())
		})

		require.NoError(t, dao.Put(toykv.Record{K: "foo", V: "foo_value_old"}))
		require.NoError(t, dao.Put(toykv.Record{K: "foo", V: "foo_value"}))
		require.NoError(t, dao.Put(toykv.Record{K: "bar", V: "bar_value"}))
		require.NoError(t, dao.Delete("bar"))
	})

	t.Run("must be able to read persisted record", func(t *testing.T) {
		dao := New(walmanager.New(walmanager.WithWALFilename(walFilename)))
		require.NoError(t, dao.Start())
		t.Cleanup(func() {
			require.NoError(t, dao.Stop())
		})

		v, err := dao.Get("foo")
		require.NoError(t, err)
		require.Equal(t, toykv.V("foo_value"), v)

		_, err = dao.Get("bar")
		require.ErrorIs(t, err, toykv.ErrNotFound, "not ErrNotFound: %v", err)
	})
}

var filenames []string

const filenamePattern = "test_%d"

func newFilename() string {
	filename := fmt.Sprintf(filenamePattern, rand.Int63())
	filenames = append(filenames, filename)
	return filename
}

func TestMain(m *testing.M) {
	mainSetup()
	code := m.Run()
	mainTeardown()
	os.Exit(code)
}

func mainSetup() {
	filenames = filenames[:0]
}

func mainTeardown() {
	for _, filename := range filenames {
		if err := os.Remove(filename); err != nil {
			fmt.Printf("cannot remove filename %v: %v", filename, err)
		}
	}
	filenames = filenames[:0]
}
