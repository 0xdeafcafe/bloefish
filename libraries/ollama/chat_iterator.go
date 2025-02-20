package ollama

import "strings"

// StreamingChatIterator wraps an Iterator[StreamingChatEvent] to provide
// convenient access to accumulated messages and completion status.
type StreamingChatIterator struct {
	iterator *Iterator[StreamingChatEvent]
	builder  strings.Builder

	complete bool
	err      error
}

// NewStreamingChatIterator creates a new StreamingChatIterator from an Iterator.
func NewStreamingChatIterator(iterator *Iterator[StreamingChatEvent]) *StreamingChatIterator {
	return &StreamingChatIterator{
		iterator: iterator,
	}
}

// Next advances the stream to the next event.
// Returns false when the stream is complete or an error occurs.
func (s *StreamingChatIterator) Next() bool {
	if !s.iterator.Next() {
		s.err = s.iterator.Err()
		return false
	}

	event, err := s.iterator.Value()
	if err != nil {
		s.err = err
		return false
	}

	if event.Done {
		s.complete = true
		return false
	}

	if event.Message != nil {
		s.builder.WriteString(event.Message.Content)
	}

	return true
}

// Content returns the accumulated message content.
func (s *StreamingChatIterator) Content() string {
	return s.builder.String()
}

// IsComplete returns true if the stream has received a completion event.
func (s *StreamingChatIterator) IsComplete() bool {
	return s.complete
}

// Err returns any error that occurred during iteration.
func (s *StreamingChatIterator) Err() error {
	return s.err
}

func (s *StreamingChatIterator) Current() *StreamingChatEvent {
	return s.iterator.Current()
}
