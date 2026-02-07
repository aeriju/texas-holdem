package api

import (
	"encoding/json"
	"errors"
	"net/http"

	"texas-holdem/internal/poker"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

type BestHandRequest struct {
	Hole      []string `json:"hole"`
	Community []string `json:"community"`
}

type BestHandResponse struct {
	BestHand []string `json:"bestHand"`
	Category string   `json:"category"`
	Tiebreak []int    `json:"tiebreak"`
}

type HeadsUpRequest struct {
	Hand1 BestHandRequest `json:"hand1"`
	Hand2 BestHandRequest `json:"hand2"`
}

type HeadsUpResponse struct {
	Hand1   BestHandResponse `json:"hand1"`
	Hand2   BestHandResponse `json:"hand2"`
	Winner  string           `json:"winner"`
	Outcome string           `json:"outcome"`
}

type OddsRequest struct {
	Hole        []string `json:"hole"`
	Community   []string `json:"community"`
	Players     int      `json:"players"`
	Simulations int      `json:"simulations"`
}

type OddsResponse struct {
	WinProbability float64 `json:"winProbability"`
}

func RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/healthz", healthHandler)
	mux.HandleFunc("/api/v1/best-hand", bestHandHandler)
	mux.HandleFunc("/api/v1/heads-up", headsUpHandler)
	mux.HandleFunc("/api/v1/odds", oddsHandler)
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("ok"))
}

func bestHandHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	var req BestHandRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	resp, err := computeBest(req)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, resp)
}

func headsUpHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	var req HeadsUpRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	resp1, err := computeBest(req.Hand1)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	resp2, err := computeBest(req.Hand2)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	h1Cards, _ := poker.ParseCards(append(req.Hand1.Hole, req.Hand1.Community...))
	h2Cards, _ := poker.ParseCards(append(req.Hand2.Hole, req.Hand2.Community...))
	h1Rank, _ := poker.Evaluate7(h1Cards)
	h2Rank, _ := poker.Evaluate7(h2Cards)
	cmp := poker.Compare(h1Rank, h2Rank)
	winner := "tie"
	outcome := "tie"
	if cmp > 0 {
		winner = "hand1"
		outcome = "hand1 wins"
	} else if cmp < 0 {
		winner = "hand2"
		outcome = "hand2 wins"
	}

	writeJSON(w, http.StatusOK, HeadsUpResponse{
		Hand1:   resp1,
		Hand2:   resp2,
		Winner:  winner,
		Outcome: outcome,
	})
}

func oddsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	var req OddsRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	hole, err := poker.ParseCards(req.Hole)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	community, err := poker.ParseCards(req.Community)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	prob, err := poker.MonteCarlo(hole, community, req.Players, req.Simulations)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, OddsResponse{WinProbability: prob})
}

func computeBest(req BestHandRequest) (BestHandResponse, error) {
	if len(req.Hole) != 2 {
		return BestHandResponse{}, errors.New("hole must have 2 cards")
	}
	if len(req.Community) != 5 {
		return BestHandResponse{}, errors.New("community must have 5 cards")
	}
	allCards, err := poker.ParseCards(append(req.Hole, req.Community...))
	if err != nil {
		return BestHandResponse{}, err
	}
	if len(allCards) != 7 {
		return BestHandResponse{}, errors.New("invalid total cards")
	}
	rank, err := poker.Evaluate7(allCards)
	if err != nil {
		return BestHandResponse{}, err
	}
	best := make([]string, 0, len(rank.Best5))
	for _, c := range rank.Best5 {
		best = append(best, c.String())
	}
	return BestHandResponse{
		BestHand: best,
		Category: rank.Name(),
		Tiebreak: rank.Tiebreak,
	}, nil
}

func decodeJSON(r *http.Request, v any) error {
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(v); err != nil {
		return err
	}
	return nil
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, ErrorResponse{Error: msg})
}
