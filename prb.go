package prb

import (
	"sync"
)

type Element[T any] struct {
	Value          T
	Priority       int
	InsertionOrder int64
}

type PriorityRingBuffer[T any] struct {
	elements       []Element[T]
	capacity       int
	head, tail     int
	size           int
	bubbleWindow   int
	orderCounter   int
	overwriteGuard bool
	mu             *sync.Mutex
}

type Config struct {
	Capacity       int
	BubbleWindow   int
	OverwriteGuard bool
	ThreadSafe     bool
}

func NewPRB[T any](config Config) *PriorityRingBuffer[T] {
	var mu *sync.Mutex
	if config.ThreadSafe {
		mu = &sync.Mutex{}
	}
	return &PriorityRingBuffer[T]{
		elements:       make([]Element[T], config.Capacity),
		capacity:       config.Capacity,
		bubbleWindow:   config.BubbleWindow,
		overwriteGuard: config.OverwriteGuard,
		mu:             mu,
	}
}

func (b *PriorityRingBuffer[T]) lock() {
	if b.mu != nil {
		b.mu.Lock()
	}
}

func (b *PriorityRingBuffer[T]) unlock() {
	if b.mu != nil {
		b.mu.Unlock()
	}
}
func (b *PriorityRingBuffer[T]) Insert(value T, priority int) error {

}

func (b *PriorityRingBuffer[T]) Dequeue() (Element[T], error) {

}

func (b *PriorityRingBuffer[T]) PeekHead() (Element[T], error) {

}
func (b *PriorityRingBuffer[T]) PeekMaxPriority() (Element[T], error) {

}

func (b *PriorityRingBuffer[T]) Search(value *T, priority *int) []int {

}
func (b *PriorityRingBuffer[T]) Size() int {

}
func (b *PriorityRingBuffer[T]) Capacity() int {

}
