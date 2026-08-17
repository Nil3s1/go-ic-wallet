package api

import (
	"encoding/json"
	"net/http"

	walletapplication "github.com/Nil3s1/go-ic-wallet/internal/wallet/application"
)

type createCardRequest struct {
	CardNo         string `json:"cardNo"`
	InitialBalance uint   `json:"initialBalance"`
}

type addBalanceRequest struct {
	Amount uint `json:"amount"`
}

type applyPaymentRequest struct {
	Amount      uint   `json:"amount"`
	ReferenceId string `json:"referenceId"`
}

func createCardHandler(h *walletapplication.CreateCardCommandHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req createCardRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}

		cmd := walletapplication.CreateCardCommand{
			CardNo:         req.CardNo,
			InitialBalance: req.InitialBalance,
		}

		if err := h.Handle(r.Context(), cmd); err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}

		w.WriteHeader(http.StatusCreated)
	}
}

func addBalanceHandler(h *walletapplication.AddBalanceCommandHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req addBalanceRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}

		cmd := walletapplication.AddBalanceCommand{
			CardNo: r.PathValue("cardNo"),
			Amount: req.Amount,
		}

		if err := h.Handle(r.Context(), cmd); err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

func applyPaymentHandler(h *walletapplication.ApplyPaymentCommandHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req applyPaymentRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}

		cmd := walletapplication.ApplyPaymentCommand{
			CardNo:      r.PathValue("cardNo"),
			Amount:      req.Amount,
			ReferenceId: req.ReferenceId,
		}

		if err := h.Handle(r.Context(), cmd); err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
