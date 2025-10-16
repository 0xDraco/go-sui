package grpc

import (
	"context"
	"errors"
	"io"
)

type pageFetcher[T any] func(ctx context.Context, token []byte) ([]T, []byte, error)

type pageIterator[T any] struct {
	token []byte
	fetch pageFetcher[T]
	done  bool
}

func newPageIterator[T any](initial []byte, fetch pageFetcher[T]) (*pageIterator[T], error) {
	if fetch == nil {
		return nil, errors.New("nil fetch function")
	}
	return &pageIterator[T]{
		token: cloneBytes(initial),
		fetch: fetch,
		done:  false,
	}, nil
}

func (it *pageIterator[T]) Next(ctx context.Context) ([]T, error) {
	if it == nil {
		return nil, errors.New("nil iterator")
	}

	for {
		if it.done {
			return nil, io.EOF
		}

		items, next, err := it.fetch(ctx, cloneBytes(it.token))
		if err != nil {
			return nil, err
		}

		it.token = cloneBytes(next)
		if len(next) == 0 {
			it.done = true
		}

		if len(items) == 0 {
			if it.done {
				return nil, io.EOF
			}
			continue
		}

		return items, nil
	}
}

func (it *pageIterator[T]) ForEach(ctx context.Context, fn func(T) error) error {
	if it == nil {
		return errors.New("nil iterator")
	}
	if fn == nil {
		return errors.New("nil function")
	}

	for {
		items, err := it.Next(ctx)
		if err != nil {
			if errors.Is(err, io.EOF) {
				return nil
			}
			return err
		}
		for _, item := range items {
			if err := fn(item); err != nil {
				return err
			}
		}
	}
}

func (it *pageIterator[T]) Collect(ctx context.Context) ([]T, error) {
	if it == nil {
		return nil, errors.New("nil iterator")
	}
	var out []T
	err := it.ForEach(ctx, func(item T) error {
		out = append(out, item)
		return nil
	})
	return out, err
}

func cloneBytes(b []byte) []byte {
	if len(b) == 0 {
		return nil
	}
	cp := make([]byte, len(b))
	copy(cp, b)
	return cp
}
