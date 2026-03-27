package web

import "tictactoe/internal/domain"

func toDomainBoard(b [][]int) domain.Board {
	var board domain.Board
	for r := 0; r < 3; r++ {
		for c := 0; c < 3; c++ {
			if r < len(b) && c < len(b[r]) {
				board[r][c] = b[r][c]
			}
		}
	}
	return board
}

func fromDomainBoard(b domain.Board) [][]int {
	res := make([][]int, 3)
	for r := 0; r < 3; r++ {
		res[r] = make([]int, 3)
		for c := 0; c < 3; c++ {
			res[r][c] = b[r][c]
		}
	}
	return res
}
