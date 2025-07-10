package prb

import (
	"errors"
	"sync"
)

var (
	ErrInvalidCapacity = errors.New("capacity must be positive")
	ErrInvalidWindow   = errors.New("bubbleWindow must be zero or positive and less than capacity")
	ErrBufferEmpty     = errors.New("buffer is empty")
	ErrBufferFull      = errors.New("buffer is full, refused to overwrite higher priority element")
)

type Element[T comparable] struct {
	Value          T
	Priority       int
	InsertionOrder int64
}

type PriorityRingBuffer[T comparable] struct {
	elements       []Element[T]
	capacity       int
	head, tail     int
	size           int
	bubbleWindow   int
	orderCounter   int64
	overwriteGuard bool
	mu             sync.RWMutex
}

type Config struct {
	Capacity       int
	BubbleWindow   int
	OverwriteGuard bool
	ThreadSafe     bool
}

func New[T comparable](config Config) (*PriorityRingBuffer[T], error) {
	if config.Capacity <= 0 {
		return nil, ErrInvalidCapacity
	}
	if config.BubbleWindow < 0 || config.BubbleWindow > config.Capacity-1 {
		return nil, ErrInvalidWindow
	}

	return &PriorityRingBuffer[T]{
		elements:       make([]Element[T], config.Capacity),
		capacity:       config.Capacity,
		bubbleWindow:   config.BubbleWindow,
		overwriteGuard: config.OverwriteGuard,
		mu:             sync.RWMutex{},
	}, nil
}

func (b *PriorityRingBuffer[T]) Insert(value T, priority int) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	element := Element[T]{
		Value:          value,
		Priority:       priority,
		InsertionOrder: b.orderCounter,
	}
	b.orderCounter++

	overwriting := b.size == b.capacity
	insertIndex := b.tail

	if overwriting && b.overwriteGuard {
		if priority <= b.elements[b.head].Priority {
			return ErrBufferFull
		}
	}

	b.elements[insertIndex] = element

	b.bubbleElement(insertIndex)

	b.tail = (b.tail + 1) % b.capacity
	if overwriting {
		b.head = b.tail
	} else {
		b.size++
	}

	return nil
}

func (b *PriorityRingBuffer[T]) bubbleElement(insertIndex int) {
	for i := 1; i <= b.bubbleWindow && b.size > 0; i++ {
		previousIndex := (insertIndex - 1 + b.capacity) % b.capacity

		if b.size == b.capacity && previousIndex == b.head {
			break
		}

		current := b.elements[insertIndex]
		previous := b.elements[previousIndex]

		if b.shouldSwap(current, previous) {
			b.elements[insertIndex], b.elements[previousIndex] = previous, current
			insertIndex = previousIndex
		} else {
			break
		}
	}
}

func (b *PriorityRingBuffer[T]) shouldSwap(current, previous Element[T]) bool {
	return current.Priority > previous.Priority ||
		(current.Priority == previous.Priority && current.InsertionOrder < previous.InsertionOrder)
}

func (b *PriorityRingBuffer[T]) Dequeue() (Element[T], error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.size == 0 {
		return Element[T]{}, ErrBufferEmpty
	}
	element := b.elements[b.head]
	b.head = (b.head + 1) % b.capacity
	b.size--
	return element, nil
}

func (b *PriorityRingBuffer[T]) Peek() (Element[T], error) {
	b.mu.RLock()
	defer b.mu.RUnlock()
	if b.size == 0 {
		return Element[T]{}, ErrBufferEmpty
	}
	return b.elements[b.head], nil
}
func (b *PriorityRingBuffer[T]) PeekMaxPriority() (Element[T], error) {
	b.mu.RLock()
	defer b.mu.RUnlock()
	if b.size == 0 {
		return Element[T]{}, errors.New("buffer is empty")
	}
	maxIndex := b.head
	for i := 1; i < b.size; i++ {
		index := (b.head + i) % b.capacity
		candidate := b.elements[index]
		currMax := b.elements[maxIndex]
		if candidate.Priority > currMax.Priority ||
			(candidate.Priority == currMax.Priority && candidate.InsertionOrder < currMax.InsertionOrder) {
			maxIndex = index
		}
	}
	return b.elements[maxIndex], nil
}

type SearchFilter[T comparable] func(Element[T]) bool

func (b *PriorityRingBuffer[T]) Search(filters ...SearchFilter[T]) []int {
	b.mu.RLock()
	defer b.mu.RUnlock()

	var result []int
	for i := 0; i < b.size; i++ {
		index := (b.head + i) % b.capacity
		element := b.elements[index]

		match := true
		for _, filter := range filters {
			if !filter(element) {
				match = false
				break
			}
		}

		if match {
			result = append(result, i)
		}
	}
	return result
}

func SearchByValue[T comparable](value T) SearchFilter[T] {
	return func(e Element[T]) bool {
		return e.Value == value
	}
}

func SearchByPriority[T comparable](priority int) SearchFilter[T] {
	return func(e Element[T]) bool {
		return e.Priority == priority
	}
}

func SearchByMinPriority[T comparable](minPriority int) SearchFilter[T] {
	return func(e Element[T]) bool {
		return e.Priority >= minPriority
	}
}

func (b *PriorityRingBuffer[T]) Len() int {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.size
}
