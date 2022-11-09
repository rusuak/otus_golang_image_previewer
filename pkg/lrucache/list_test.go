package lrucache

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestList(t *testing.T) {
	t.Run("empty list", func(t *testing.T) {
		l := NewList()

		require.Equal(t, 0, l.Len())
		require.Nil(t, l.Front())
		require.Nil(t, l.Back())
	})

	t.Run("complex", func(t *testing.T) {
		l := NewList()

		l.PushFront(10) // [10]
		l.PushBack(20)  // [10, 20]
		l.PushBack(30)  // [10, 20, 30]
		require.Equal(t, 3, l.Len())

		middle := l.Front().Next // 20
		l.Remove(middle)         // [10, 30]
		require.Equal(t, 2, l.Len())

		for i, v := range [...]int{40, 50, 60, 70, 80} {
			if i%2 == 0 {
				l.PushFront(v)
			} else {
				l.PushBack(v)
			}
		} // [80, 60, 40, 10, 30, 50, 70]

		require.Equal(t, 7, l.Len())
		require.Equal(t, 80, l.Front().Value)
		require.Equal(t, 70, l.Back().Value)

		l.MoveToFront(l.Front()) // [80, 60, 40, 10, 30, 50, 70]
		l.MoveToFront(l.Back())  // [70, 80, 60, 40, 10, 30, 50]

		elems := make([]int, 0, l.Len())
		for i := l.Front(); i != nil; i = i.Next {
			elems = append(elems, i.Value.(int))
		}
		require.Equal(t, []int{70, 80, 60, 40, 10, 30, 50}, elems)
	})

	t.Run("single item", func(t *testing.T) {
		l := NewList()

		l.PushFront(10) // [10]
		require.Equal(t, 1, l.Len())
		require.Equal(t, 10, l.Front().Value)
		require.Equal(t, l.Front(), l.Back())
		require.Nil(t, l.Front().Next, l.Front().Prev)

		item := l.Front()
		l.Remove(item)
		require.Equal(t, 0, l.Len())
		require.Nil(t, l.Front(), l.Front())

		item = l.PushBack(20)    // [20]
		l.MoveToFront(l.Front()) // [20]
		require.Equal(t, 1, l.Len())
		require.Equal(t, 20, l.Front().Value)
		require.Equal(t, l.Front(), l.Back())
		require.Nil(t, l.Front().Next, l.Front().Prev)

		l.Remove(item)
		require.Equal(t, 0, l.Len())
	})

	t.Run("movements", func(t *testing.T) {
		l := NewList()

		item1 := l.PushFront(11) // [11]
		item2 := l.PushFront(22) // [22,11]
		item3 := l.PushFront(33) // [33,22,11]

		l.MoveToFront(item1)    // [11,33,22]
		l.MoveToFront(item2)    // [22,11,33]
		l.MoveToFront(item3)    // [33,22,11]
		l.MoveToFront(l.Back()) // [11,33,22]

		require.Equal(t, 11, l.Front().Value)
		require.Equal(t, 33, l.Front().Next.Value)
		require.Nil(t, l.Front().Prev)

		require.Equal(t, 22, l.Back().Value)
		require.Nil(t, l.Back().Next)
		require.Equal(t, 33, l.Back().Prev.Value)
	})
}
