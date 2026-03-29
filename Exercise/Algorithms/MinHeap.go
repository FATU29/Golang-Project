package main

import "fmt"

type MinHeap struct {
	Data []int
}

func (h *MinHeap) Insert(val int) {
	h.Data = append(h.Data, val)
	h.ShiftUp(len(h.Data) - 1)
}

func (h *MinHeap) ShiftUp(index int) {
	for index > 0 {
		parent := (index - 1) / 2
		// For MinHeap, swap if child is SMALLER than parent
		if h.Data[index] < h.Data[parent] {
			h.Data[index], h.Data[parent] = h.Data[parent], h.Data[index]
			index = parent
		} else {
			break
		}
	}
}

func (h *MinHeap) Pop() int {
	n := len(h.Data)
	if n == 0 {
		return -1
	}

	val := h.Data[0]
	lastIdx := n - 1

	// Move last element to root
	h.Data[0] = h.Data[lastIdx]
	h.Data = h.Data[:lastIdx]

	if len(h.Data) > 0 {
		h.ShiftDown(0)
	}
	return val
}

func (h *MinHeap) ShiftDown(index int) {
	size := len(h.Data)
	for {
		smallest := index // Reset smallest to current node at each level
		left := index*2 + 1
		right := index*2 + 2

		if left < size && h.Data[left] < h.Data[smallest] {
			smallest = left
		}

		if right < size && h.Data[right] < h.Data[smallest] {
			smallest = right
		}

		if smallest != index {
			h.Data[index], h.Data[smallest] = h.Data[smallest], h.Data[index]
			index = smallest
		} else {
			break
		}
	}
}

func main() {
	// Better to start empty or Heapify an existing slice
	minHeap := &MinHeap{Data: []int{}}
	values := []int{10, 40, 15, 100, 50, 30, 40}

	for _, v := range values {
		minHeap.Insert(v)
	}

	minHeap.Insert(12)
	fmt.Println("Heap after all inserts:", minHeap.Data)

	val := minHeap.Pop()
	fmt.Println("Popped (Min):", val)
	fmt.Println("Heap after Pop:", minHeap.Data)
}
