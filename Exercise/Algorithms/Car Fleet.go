package main

import (
	"fmt"
	"sort"
)

type Car struct {
	pos   int
	speed int
}

func carFleet(target int, position []int, speed []int) int {
	pair := make([]Car, 0, len(position))

	for idx, _ := range position {
		pair = append(pair, Car{
			pos:   position[idx],
			speed: speed[idx],
		})
	}

	sort.Slice(pair, func(i, j int) bool {
		return pair[i].pos > pair[j].pos
	})

	var stack []float64

	for _, val := range pair {

		time := float64(target-val.pos) / float64(val.speed)

		if len(stack) > 0 && stack[len(stack)-1] >= time {
			continue
		}

		stack = append(stack, time)
	}

	return len(stack)
}

func main() {
	fmt.Println("hello world")
	target := 10
	position := []int{1, 4}
	speed := []int{3, 2}

	fmt.Println(carFleet(target, position, speed))
}
