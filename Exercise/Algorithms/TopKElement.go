package main

import "fmt"

type Node struct {
	Number int
	Count  int
}

type MinHeapTopK struct {
	Values []Node
}

func (h *MinHeapTopK) Insert(val *Node) {
	h.Values = append(h.Values, *val)
	h.ShiftUp()
}

func (h *MinHeapTopK) ShiftUp() {
	curr := len(h.Values) - 1
	for curr > 0 {
		parent := (curr - 1) / 2
		if h.Values[curr].Count < h.Values[parent].Count {
			h.Values[curr], h.Values[parent] = h.Values[parent], h.Values[curr]
			curr = parent
		} else {
			break
		}
	}
}

func (h *MinHeapTopK) Pop() *Node {
	size := len(h.Values)

	if size <= 0 {
		return nil
	}
	firstEL := h.Values[0]

	h.Values[0] = h.Values[size-1]
	h.Values = h.Values[:size-1]
	h.ShiftDown()

	return &firstEL
}

func (h *MinHeapTopK) ShiftDown() {
	size := len(h.Values)

	if size <= 0 {
		return
	}

	n := 0
	smallest := 0
	var idx1 int
	var idx2 int

	for {
		idx1 = n*2 + 1
		idx2 = n*2 + 2
		smallest = n

		if size > idx1 && h.Values[smallest].Count > h.Values[idx1].Count {
			smallest = idx1
		}

		if size > idx2 && h.Values[smallest].Count > h.Values[idx2].Count {
			smallest = idx2
		}

		if smallest != n {
			h.Values[smallest], h.Values[n] = h.Values[n], h.Values[smallest]
			n = smallest
		} else {
			break
		}
	}

}

func topKFrequent(nums []int, k int) []int {
	sequences := make(map[int]int)

	for _, num := range nums {
		sequences[num]++
	}

	minHeap := &MinHeapTopK{
		Values: make([]Node, 0),
	}

	for key, val := range sequences {
		minHeap.Insert(&Node{
			Number: key,
			Count:  val,
		})

		if len(minHeap.Values) > k {
			minHeap.Pop()
		}
	}

	result := make([]int, 0, k)

	for _, val := range minHeap.Values {
		result = append(result, val.Number)
	}

	return result
}

func main() {

	values := []int{1, 1, 1, 2, 2, 3}
	k := 2

	result := topKFrequent(values, k)
	fmt.Println(result)

}
