package ollama

import (
	"bufio"
	"encoding/json"
)

// Iterator provides an iterator for streamed responses.
type Iterator[T any] struct {
	current *T
	scanner *bufio.Scanner
	done    bool
	err     error
}

// Next advances the iterator and returns the next T.
// It returns false when the stream ends or an error occurs.
func (i *Iterator[T]) Next() bool {
	if i.done {
		return false
	}
	if !i.scanner.Scan() {
		if err := i.scanner.Err(); err != nil {
			i.err = err
		}
		i.done = true
		return false
	}

	var data T
	if err := json.Unmarshal([]byte(i.scanner.Text()), &data); err != nil {
		i.err = err
		i.done = true
		return false
	}

	i.current = &data
	return true
}

// Value returns the current T, or an error if one occurred.
func (it *Iterator[T]) Value() (*T, error) {
	if it.err != nil {
		return nil, it.err
	}
	return it.current, nil
}

// Current returns the current T.
func (it *Iterator[T]) Current() *T {
	return it.current
}

// Err returns the error encountered during iteration, if any.
func (it *Iterator[T]) Err() error {
	return it.err
}
