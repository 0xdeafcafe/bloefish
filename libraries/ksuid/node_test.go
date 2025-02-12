package ksuid

import (
	"context"
	"testing"
)

func BenchmarkGenerate(b *testing.B) {
	for range b.N {
		Generate(context.Background(), "user")
	}
}
