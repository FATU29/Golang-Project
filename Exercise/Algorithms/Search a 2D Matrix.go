package main

import "fmt"

func searchMatrix(matrix [][]int, target int) bool {
	rows, cols := len(matrix), len(matrix[0])

	top, bot := 0, len(matrix)-1

	for top <= bot {
		row := (top + bot) / 2

		if matrix[row][1] == target {

		}

	}

}

func main() {
	fmt.Println("Hello World")
	matrix := [][]int{{1, 2, 4, 8}, {10, 11, 12, 13}, {14, 20, 30, 40}}
	target := 10

	fmt.Println(matrix, target)
}
