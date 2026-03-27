package domain

import (
	"errors"
)

type gameService struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &gameService{repo: repo}
}

func (s *gameService) Validate(oldBoard, newBoard Board) error {
	diffs := 0
	for r := 0; r < 3; r++ {
		for c := 0; c < 3; c++ {
			if oldBoard[r][c] != 0 && oldBoard[r][c] != newBoard[r][c] {
				return errors.New("invalid move: you changed an existing mark")
			}
			if oldBoard[r][c] == 0 && newBoard[r][c] != 0 {
				if newBoard[r][c] != 1 {
					return errors.New("invalid move: only user (1) can move")
				}
				diffs++
			}
		}
	}
	if diffs != 1 {
		return errors.New("invalid move: must place exactly one mark")
	}
	return nil
}

func (s *gameService) PlayNextMove(game Game) (Game, error) {
	bestScore := -1000
	var bestMove [2]int
	found := false

	for r := 0; r < 3; r++ {
		for c := 0; c < 3; c++ {
			if game.Board[r][c] == 0 {
				game.Board[r][c] = 2
				score := s.minimax(game.Board, 0, false)
				game.Board[r][c] = 0

				if score > bestScore {
					bestScore = score
					bestMove = [2]int{r, c}
					found = true
				}
			}
		}
	}

	if found {
		game.Board[bestMove[0]][bestMove[1]] = 2
	}

	err := s.repo.Save(game)
	return game, err
}

func (s *gameService) minimax(board Board, depth int, isMaximizing bool) int {
	over, winner := s.IsGameOver(board)
	if over {
		if winner == 2 {
			return 10 - depth
		}
		if winner == 1 {
			return depth - 10
		}
		return 0
	}

	if isMaximizing {
		bestScore := -1000
		for r := 0; r < 3; r++ {
			for c := 0; c < 3; c++ {
				if board[r][c] == 0 {
					board[r][c] = 2
					score := s.minimax(board, depth+1, false)
					board[r][c] = 0
					if score > bestScore {
						bestScore = score
					}
				}
			}
		}
		return bestScore
	} else {
		bestScore := 1000
		for r := 0; r < 3; r++ {
			for c := 0; c < 3; c++ {
				if board[r][c] == 0 {
					board[r][c] = 1
					score := s.minimax(board, depth+1, true)
					board[r][c] = 0
					if score < bestScore {
						bestScore = score
					}
				}
			}
		}
		return bestScore
	}
}

func (s *gameService) IsGameOver(b Board) (bool, int) {
	for i := 0; i < 3; i++ {
		if b[i][0] != 0 && b[i][0] == b[i][1] && b[i][1] == b[i][2] {
			return true, b[i][0]
		}
		if b[0][i] != 0 && b[0][i] == b[1][i] && b[1][i] == b[2][i] {
			return true, b[0][i]
		}
	}
	if b[0][0] != 0 && b[0][0] == b[1][1] && b[1][1] == b[2][2] {
		return true, b[0][0]
	}
	if b[0][2] != 0 && b[0][2] == b[1][1] && b[1][1] == b[2][0] {
		return true, b[0][2]
	}
	for r := 0; r < 3; r++ {
		for c := 0; c < 3; c++ {
			if b[r][c] == 0 {
				return false, 0
			}
		}
	}
	return true, 0
}
