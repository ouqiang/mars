package common

import (
	"container/list"
	"sync"
)

// Queue 队列, 保存固定大小的元素
type Queue struct {
	bucket  *list.List
	maxSize int
	mu      sync.Mutex
}

// New 创建queue
func NewQueue(size int) *Queue {
	q := &Queue{
		maxSize: size,
		bucket:  list.New(),
	}

	return q
}

// Add 添加元素到队列中
func (q *Queue) Add(value interface{}) (removeValue interface{}) {
	q.mu.Lock()
	defer q.mu.Unlock()

	if q.bucket.Len() == q.maxSize {
		firstEle := q.bucket.Front()
		if firstEle != nil {
			removeValue = q.bucket.Remove(firstEle)
		}
	}

	q.bucket.PushBack(value)

	return removeValue
}

// Len 获取队列大小
func (q *Queue) Len() int {
	q.mu.Lock()
	l := q.bucket.Len()
	q.mu.Unlock()

	return l
}
