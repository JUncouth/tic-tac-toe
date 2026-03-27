package datasource

import (
	"errors"
	"tictactoe/internal/domain"
)

type memoryRepository struct {
	db *Storage
}

func NewRepository(db *Storage) domain.Repository {
	return &memoryRepository{db: db}
}

func (r *memoryRepository) Save(game domain.Game) error {
	data := fromDomain(game)
	r.db.m.Store(data.ID, data)
	return nil
}

func (r *memoryRepository) Get(id string) (domain.Game, error) {
	val, ok := r.db.m.Load(id)
	if !ok {
		return domain.Game{}, errors.New("game not found")
	}
	data := val.(GameData)
	return toDomain(data), nil
}

func (r *memoryRepository) GetAllAndSort() ([]domain.Game, error) {
	var games []domain.Game

	r.db.m.Range(func(key, value interface{}) bool {
		data := value.(GameData)
		games = append(games, toDomain(data))
		return true
	})

	if len(games) == 0 {
		return games, errors.New("no games found")
	}

	return games, nil
}
