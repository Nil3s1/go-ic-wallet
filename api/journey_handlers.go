package api

import (
	"encoding/json"
	"net/http"

	journeyapplication "github.com/Nil3s1/go-ic-wallet/internal/journey/application"
)

type startJourneyRequest struct {
	MediaId      string `json:"mediaId"`
	StartStation string `json:"startStation"`
}

type endJourneyRequest struct {
	MediaId    string `json:"mediaId"`
	EndStation string `json:"endStation"`
}

type startJourneyResponse struct {
	MediaId     string `json:"mediaId"`
	IsOnJourney bool   `json:"isOnJourney"`
}

type endJourneyResponse struct {
	MediaId           string `json:"mediaId"`
	IsOnJourney       bool   `json:"isOnJourney"`
	DistanceTravelled uint   `json:"distanceTravelled"`
	Fare              uint   `json:"fare"`
}

func startJourneyHandler(h *journeyapplication.StartJourneyCommandHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req startJourneyRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}

		cmd := journeyapplication.StartJourneyCommand{
			MediaId:      req.MediaId,
			StartStation: req.StartStation,
		}

		res, err := h.Handle(r.Context(), cmd)
		if err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}

		writeJSON(w, http.StatusOK, startJourneyResponse{
			MediaId:     res.MediaId,
			IsOnJourney: res.IsOnJourney,
		})
	}
}

func endJourneyHandler(h *journeyapplication.EndJourneyCommandHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req endJourneyRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}

		cmd := journeyapplication.EndJourneyCommand{
			MediaId:    req.MediaId,
			EndStation: req.EndStation,
		}

		res, err := h.Handle(r.Context(), cmd)
		if err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}

		writeJSON(w, http.StatusOK, endJourneyResponse{
			MediaId:           res.MediaId,
			IsOnJourney:       res.IsOnJourney,
			DistanceTravelled: res.DistanceTravelled,
			Fare:              res.Fare,
		})
	}
}
