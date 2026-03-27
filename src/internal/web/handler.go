package web

import (
	"encoding/json"
	"net/http"
	"sort"
	"tictactoe/internal/domain"
)

type Handler struct {
	svc  domain.Service
	repo domain.Repository
}

func NewHandler(svc domain.Service, repo domain.Repository) *Handler {
	return &Handler{svc: svc, repo: repo}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /game/{id}", h.handlePlay)
	mux.HandleFunc("GET /game/{id}", h.handleGet)
	mux.HandleFunc("GET /games", h.handleGetAllAndSort)
}

func (h *Handler) handlePlay(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	var req GameRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	newBoard := toDomainBoard(req.Board)

	game, err := h.repo.Get(id)
	if err != nil {
		game = domain.Game{ID: id, Board: domain.Board{}}
	}

	var over bool
	over, _ = h.svc.IsGameOver(game.Board)
	if over {
		http.Error(w, "game is already over", http.StatusBadRequest)
		return
	}

	if err = h.svc.Validate(game.Board, newBoard); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	game.Board = newBoard

	over, _ = h.svc.IsGameOver(game.Board)
	if over {
		h.repo.Save(game)
		json.NewEncoder(w).Encode(GameResponse{ID: game.ID, Board: fromDomainBoard(game.Board), Message: "User wins!"})
		return
	}

	game, err = h.svc.PlayNextMove(game)
	if err != nil {
		http.Error(w, "server error computing move", http.StatusInternalServerError)
		return
	}

	over, _ = h.svc.IsGameOver(game.Board)
	if over {
		h.repo.Save(game)
		json.NewEncoder(w).Encode(GameResponse{ID: game.ID, Board: fromDomainBoard(game.Board), Message: "AI wins!"})
		return
	}

	json.NewEncoder(w).Encode(GameResponse{
		ID:    game.ID,
		Board: fromDomainBoard(game.Board),
	})
}

func (h *Handler) handleGet(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	game, err := h.repo.Get(id)
	if err != nil {
		http.Error(w, "game not found", http.StatusNotFound)
		return
	}

	over, winner := h.svc.IsGameOver(game.Board)
	if over {
		var message string
		switch winner {
		case 1:
			message = "User won!"
		case 2:
			message = "AI won!"
		default:
			message = "Game ended with a draw."
		}
		json.NewEncoder(w).Encode(GameResponse{ID: game.ID, Board: fromDomainBoard(game.Board), Message: message})
		return
	}

	json.NewEncoder(w).Encode(GameResponse{
		ID:    game.ID,
		Board: fromDomainBoard(game.Board),
	})
}

func (h *Handler) handleGetAllAndSort(w http.ResponseWriter, r *http.Request) {
	games, err := h.repo.GetAllAndSort()

	if err != nil {
		http.Error(w, "no games found", http.StatusNotFound)
		return
	}

	sort.Slice(games, func(i, j int) bool {
		return games[i].ID < games[j].ID
	})

	var over bool
	var winner int

	for _, g := range games {
		over, winner = h.svc.IsGameOver(g.Board)
		if over {
			var message string
			switch winner {
			case 1:
				message = "User won!"
			case 2:
				message = "AI won!"
			default:
				message = "Game ended with a draw."
			}
			json.NewEncoder(w).Encode(GameResponse{ID: g.ID, Board: fromDomainBoard(g.Board), Message: message})
			continue
		}

		json.NewEncoder(w).Encode(GameResponse{
			ID:    g.ID,
			Board: fromDomainBoard(g.Board),
		})
	}
}
