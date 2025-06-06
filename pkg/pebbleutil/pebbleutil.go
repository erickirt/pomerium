package pebbleutil

import (
	"context"
	"iter"
	"slices"

	"github.com/cockroachdb/pebble/v2"
	"github.com/cockroachdb/pebble/v2/vfs"
)

// Iterate iterates over a pebble reader.
func Iterate[T any](src pebble.Reader, iterOptions *pebble.IterOptions, f func(it *pebble.Iterator) (T, error)) iter.Seq2[T, error] {
	var zero T
	return func(yield func(T, error) bool) {
		it, err := src.NewIter(iterOptions)
		if err != nil {
			yield(zero, err)
			return
		}

		for it.First(); it.Valid(); it.Next() {
			value, err := f(it)
			if err != nil {
				_ = it.Close()
				yield(zero, err)
				return
			}

			if !yield(value, nil) {
				_ = it.Close()
				return
			}
		}

		err = it.Error()
		if err != nil {
			_ = it.Close()
			yield(zero, err)
			return
		}

		err = it.Close()
		if err != nil {
			yield(zero, err)
			return
		}
	}
}

// IterateKeys yields the keys in a pebble reader.
func IterateKeys(src pebble.Reader, iterOptions *pebble.IterOptions) iter.Seq2[[]byte, error] {
	return Iterate(src, iterOptions, func(it *pebble.Iterator) ([]byte, error) {
		return slices.Clone(it.Key()), nil
	})
}

// IterateValues yields the values in a pebble reader.
func IterateValues(src pebble.Reader, iterOptions *pebble.IterOptions) iter.Seq2[[]byte, error] {
	return Iterate(src, iterOptions, func(it *pebble.Iterator) ([]byte, error) {
		value, err := it.ValueAndErr()
		if err != nil {
			return nil, err
		}
		return slices.Clone(value), nil
	})
}

// MustOpen opens a pebble database. It sets options useful for pomerium and panics if there is an error.
func MustOpen(dirname string, options *pebble.Options) *pebble.DB {
	db, err := Open(dirname, options)
	if err != nil {
		panic(err)
	}
	return db
}

// MustOpenMemory opens an in-memory pebble database. It panics if there is an error.
func MustOpenMemory(options *pebble.Options) *pebble.DB {
	if options == nil {
		options = new(pebble.Options)
	}
	options.FS = vfs.NewMem()
	return MustOpen("", options)
}

// Open opens a pebble database. It sets options useful for pomerium.
func Open(dirname string, options *pebble.Options) (*pebble.DB, error) {
	if options == nil {
		options = new(pebble.Options)
	}
	options.LoggerAndTracer = pebbleLogger{}
	return pebble.Open(dirname, options)
}

// PrefixToUpperBound returns an upper bound for the given prefix.
func PrefixToUpperBound(prefix []byte) []byte {
	upperBound := make([]byte, len(prefix))
	copy(upperBound, prefix)
	for i := len(upperBound) - 1; i >= 0; i-- {
		upperBound[i] = upperBound[i] + 1
		if upperBound[i] != 0 {
			return upperBound[:i+1]
		}
	}
	return nil // no upper-bound
}

type pebbleLogger struct{}

func (pebbleLogger) Infof(_ string, _ ...any)                     {}
func (pebbleLogger) Errorf(_ string, _ ...any)                    {}
func (pebbleLogger) Fatalf(_ string, _ ...any)                    {}
func (pebbleLogger) Eventf(_ context.Context, _ string, _ ...any) {}
func (pebbleLogger) IsTracingEnabled(_ context.Context) bool      { return false }
