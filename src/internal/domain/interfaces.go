package domain

type Repository interface {
	Save(game Game) error
	Get(id string) (Game, error)
	GetAllAndSort() ([]Game, error)
}

type Service interface {
	PlayNextMove(game Game) (Game, error)
	Validate(oldBoard, newBoard Board) error
	IsGameOver(board Board) (bool, int)
}
