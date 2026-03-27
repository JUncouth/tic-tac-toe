package datasource

import "tictactoe/internal/domain"

func toDomain(g GameData) domain.Game {
	return domain.Game{
		ID:    g.ID,
		Board: domain.Board(g.Board),
	}
}

func fromDomain(g domain.Game) GameData {
	return GameData{
		ID:    g.ID,
		Board: [3][3]int(g.Board),
	}
}
