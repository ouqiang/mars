package common

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestQueue_Add(t *testing.T) {
	q := NewQueue(10)
	for i := 0; i < 100; i++ {
		value := q.Add(i)
		if i >= 10 {
			require.Equal(t, i-10, value.(int))
		}
	}
	require.Equal(t, 10, q.Len())
}
