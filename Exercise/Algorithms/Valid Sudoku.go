package main

import "fmt"

func isValidSudoku(board [][]byte) bool {
	rows := len(board)
	cols := len(board[0])

	var hashRows, hashCols, hashBoxes [9][9]bool

	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			if board[i][j] == '.' {
				continue
			}

			num := board[i][j] - '1'
			boxIdx := (i/3)*3 + (j / 3)

			if hashRows[i][num] || hashCols[j][num] || hashBoxes[boxIdx][num] {
				return false
			}

			hashRows[i][num] = true
			hashCols[j][num] = true
			hashBoxes[boxIdx][num] = true
		}
	}
	return true
}

func main() {
	fmt.Println("Hello World")

	board := [][]byte{
		{'5', '3', '.', '.', '7', '.', '.', '.', '.'},
		{'6', '.', '.', '1', '9', '5', '.', '.', '.'},
		{'.', '9', '8', '.', '.', '.', '.', '6', '.'},
		{'8', '.', '.', '.', '6', '.', '.', '.', '3'},
		{'4', '.', '.', '8', '.', '3', '.', '.', '1'},
		{'7', '.', '.', '.', '2', '.', '.', '.', '6'},
		{'.', '6', '.', '.', '.', '.', '2', '8', '.'},
		{'.', '.', '.', '4', '1', '9', '.', '.', '5'},
		{'.', '.', '.', '.', '8', '.', '.', '7', '9'},
	}

	fmt.Println(isValidSudoku(board))
}
