package api

import (
	"encoding/json"
	"net/http"

	"github.com/Nil3s1/go-ic-wallet/internal/journey"
)

type startJourneyRequest struct {
	MediaId      string `json:"mediaId"`
	StartStation string `json:"startStation"`
}

type endJourneyRequest struct {
	MediaId    string `json:"mediaId"`
	EndStation string `json:"endStation"`
}

func startJourneyHandler(h *journey.StartJourneyCommandHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req startJourneyRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}

		cmd := journey.StartJourneyCommand{
			MediaId:      req.MediaId,
			StartStation: req.StartStation,
		}

		if err := h.Handle(r.Context(), cmd); err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

func endJourneyHandler(h *journey.EndJourneyCommandHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req endJourneyRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}

		cmd := journey.EndJourneyCommand{
			MediaId:    req.MediaId,
			EndStation: req.EndStation,
		}

		if err := h.Handle(r.Context(), cmd); err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
