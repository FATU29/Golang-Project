package main

import "fmt"

type MaxHeap struct {
	Data []int
}

func NewMaxHeap() *MaxHeap {
	return &MaxHeap{
		Data: []int{},
	}
}

func (h *MaxHeap) Insert(value int) {
	h.Data = append(h.Data, value)
	h.ShiftUp(len(h.Data) - 1)
}

func (h *MaxHeap) ShiftUp(index int) {
	for index > 0 {
		parent := (index - 1) / 2
		if h.Data[index] > h.Data[parent] {
			h.Data[index], h.Data[parent] = h.Data[parent], h.Data[index]
			index = parent
		} else {
			break
		}
	}
}

func (h *MaxHeap) Pop() int {
	size := len(h.Data)
	if size == 0 {
		return -1
	}

	root := h.Data[0]
	lastIndex := size - 1

	// Move last element to root and shrink slice
	h.Data[0] = h.Data[lastIndex]
	h.Data = h.Data[:lastIndex]

	if len(h.Data) > 0 {
		h.ShiftDown(0)
	}
	return root
}

func (h *MaxHeap) ShiftDown(index int) {
	size := len(h.Data)
	for {
		largest := index
		leftIdx := index*2 + 1
		rightIdx := index*2 + 2

		if leftIdx < size && h.Data[leftIdx] > h.Data[largest] {
			largest = leftIdx
		}

		if rightIdx < size && h.Data[rightIdx] > h.Data[largest] {
			largest = rightIdx
		}

		if largest != index {
			h.Data[index], h.Data[largest] = h.Data[largest], h.Data[index]
			index = largest
		} else {
			break
		}
	}
}

func main() {
	maxHeap := NewMaxHeap()
	// Inserting initial values one by one to maintain heap property
	for _, v := range []int{30, 20, 15, 5, 10, 12, 6} {
		maxHeap.Insert(v)
	}

	fmt.Println("Before Insert 40:", maxHeap.Data)
	maxHeap.Insert(40)
	fmt.Println("After Insert 40:", maxHeap.Data)

	popped := maxHeap.Pop()
	fmt.Printf("Popped: %d, Current Heap: %v\n", popped, maxHeap.Data)
}
