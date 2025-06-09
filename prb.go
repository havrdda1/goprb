package prb

import (
	"errors"
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
	orderCounter   int64
	overwriteGuard bool
	mu             *sync.Mutex
}

type Config struct {
	Capacity       int
	BubbleWindow   int
	OverwriteGuard bool
	ThreadSafe     bool
}

func New[T any](config Config) *PriorityRingBuffer[T] {
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
	b.lock()
	defer b.unlock()

	element := Element[T]{Value: value, Priority: priority, InsertionOrder: b.orderCounter}
	b.orderCounter++

	overwriting := b.size == b.capacity
	insertIndex := b.tail

	if overwriting && b.overwriteGuard {
		if priority <= b.elements[b.head].Priority {
			return errors.New("buffer is full, refused to overwrite higher priority element")
		}
	}

	b.elements[insertIndex] = element

	for i := 1; i <= b.bubbleWindow && b.size > 0; i++ {
		previousIndex := (insertIndex - 1 + b.capacity) % b.capacity
		if b.size == b.capacity && previousIndex == b.head {
			break
		}
		previousElement := b.elements[previousIndex]
		currentElement := b.elements[insertIndex]
		if currentElement.Priority > previousElement.Priority || currentElement.InsertionOrder < previousElement.InsertionOrder {
			b.elements[insertIndex], b.elements[previousIndex] = b.elements[previousIndex], b.elements[insertIndex]
			insertIndex = previousIndex
		} else {
			break
		}
	}

	b.tail = (b.tail + 1) % b.capacity
	if overwriting {
		b.head = b.tail
	} else {
		b.size++
	}
	return nil
}

func (b *PriorityRingBuffer[T]) Dequeue() (Element[T], error) {
	b.lock()
	defer b.unlock()
	if b.size == 0 {
		return Element[T]{}, errors.New("buffer is empty")
	}
	element := b.elements[b.head]
	b.head = (b.head + 1) % b.capacity
	b.size--
	return element, nil
}

func (b *PriorityRingBuffer[T]) PeekHead() (Element[T], error) {
	b.lock()
	defer b.unlock()
	if b.size == 0 {
		return Element[T]{}, errors.New("buffer is empty")
	}
	return b.elements[b.head], nil
}
func (b *PriorityRingBuffer[T]) PeekMaxPriority() (Element[T], error) {

}

func (b *PriorityRingBuffer[T]) Search(value *T, priority *int) []int {

}

func (b *PriorityRingBuffer[T]) Size() int {
	b.lock()
	defer b.unlock()
	return b.size
}
func (b *PriorityRingBuffer[T]) Capacity() int {
	return b.capacity
}
