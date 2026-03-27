package web

type GameRequest struct {
	Board [][]int `json:"board"`
}

type GameResponse struct {
	ID      string  `json:"id"`
	Board   [][]int `json:"board"`
	Message string  `json:"message,omitempty"`
}
