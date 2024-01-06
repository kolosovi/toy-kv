package dao

import (
	"fmt"
	"math/rand"
	"os"
	"testing"
	"time"

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

func TestPersistence(t *testing.T) {
	var walFilename = newFilename()
	t.Run("record must be persisted to disk", func(t *testing.T) {
		dao := New(walmanager.New(walmanager.WithWALFilename(walFilename)))
		require.NoError(t, dao.Start())
		t.Cleanup(func() { require.NoError(t, dao.Stop()) })

		require.NoError(t, dao.Put(toykv.Record{K: "foo", V: "foo_value_old"}))
		require.NoError(t, dao.Put(toykv.Record{K: "foo", V: "foo_value"}))
		require.NoError(t, dao.Put(toykv.Record{K: "bar", V: "bar_value"}))
		require.NoError(t, dao.Delete("bar"))
	})

	t.Run("must be able to read persisted records and add new ones", func(t *testing.T) {
		dao := New(walmanager.New(walmanager.WithWALFilename(walFilename)))
		require.NoError(t, dao.Start())
		t.Cleanup(func() { require.NoError(t, dao.Stop()) })

		assertKeyExists(t, dao, "foo", "foo_value")
		assertKeyNotFound(t, dao, "bar")

		require.NoError(t, dao.Put(toykv.Record{K: "baz", V: "baz_value"}))
	})

	t.Run("must be able to read records added in all previous runs", func(t *testing.T) {
		dao := New(walmanager.New(walmanager.WithWALFilename(walFilename)))
		require.NoError(t, dao.Start())
		t.Cleanup(func() { require.NoError(t, dao.Stop()) })

		assertKeyExists(t, dao, "foo", "foo_value")
		assertKeyNotFound(t, dao, "bar")
		assertKeyExists(t, dao, "baz", "baz_value")
	})
}

func assertKeyExists(t *testing.T, dao *DAO, key, expectedValue string) {
	actualValue, err := dao.Get(toykv.K(key))
	require.NoError(t, err)
	require.Equal(t, toykv.V(expectedValue), actualValue)
}

func assertKeyNotFound(t *testing.T, dao *DAO, key string) {
	_, err := dao.Get("bar")
	require.ErrorIs(t, err, toykv.ErrNotFound, "not ErrNotFound: %v", err)
}

func TestStress(t *testing.T) {
	var walFilename = newFilename()
	const numGenerations = 10
	const numKeys = 100

	t.Run("populate storage", func(t *testing.T) {
		dao := New(walmanager.New(walmanager.WithWALFilename(walFilename)))
		require.NoError(t, dao.Start())
		t.Cleanup(func() { require.NoError(t, dao.Stop()) })
		for iter := 0; iter < numGenerations; iter++ {
			for i := 0; i < numKeys; i++ {
				key := toykv.K(fmt.Sprintf("key_%d", i))
				value := toykv.V(fmt.Sprintf("value_%d", i))
				require.NoError(t, dao.Put(toykv.Record{K: key, V: value}))
			}
		}
	})

	t.Run("startup", func(t *testing.T) {
		dao := New(walmanager.New(walmanager.WithWALFilename(walFilename)))
		startedAt := time.Now()
		require.NoError(t, dao.Start())
		t.Cleanup(func() { require.NoError(t, dao.Stop()) })
		startupDuration := time.Since(startedAt)
		require.Lessf(
			t,
			startupDuration, time.Second,
			"startup duration is %v > %v", startupDuration, time.Second,
		)
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
