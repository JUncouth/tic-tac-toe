package domain

type Board [3][3]int

type Game struct {
	ID    string
	Board Board
}
