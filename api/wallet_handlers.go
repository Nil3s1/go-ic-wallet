package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Nil3s1/go-ic-wallet/internal/wallet"
)

type createCardRequest struct {
	CardNo         string `json:"cardNo"`
	InitialBalance int    `json:"initialBalance"`
}

type addBalanceRequest struct {
	Amount int `json:"amount"`
}

type applyPaymentRequest struct {
	Amount      int    `json:"amount"`
	ReferenceId string `json:"referenceId"`
}

type sufficientBalanceResponse struct {
	CardNo            string `json:"cardNo"`
	Amount            int    `json:"amount"`
	SufficientBalance bool   `json:"sufficientBalance"`
}

func createCardHandler(h *wallet.CreateCardCommandHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req createCardRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}

		cmd := wallet.CreateCardCommand{
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

func addBalanceHandler(h *wallet.AddBalanceCommandHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req addBalanceRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}

		cmd := wallet.AddBalanceCommand{
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

func applyPaymentHandler(h *wallet.ApplyPaymentCommandHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req applyPaymentRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}

		cmd := wallet.ApplyPaymentCommand{
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

func hasSufficientBalanceHandler(h *wallet.HasSufficientBalanceQueryHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cardNo := r.PathValue("cardNo")

		amount, err := strconv.Atoi(r.URL.Query().Get("amount"))
		if err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}

		sufficient, err := h.Handle(wallet.HasSufficientBalanceQuery{
			CardNo: cardNo,
			Amount: amount,
		})
		if err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}

		writeJSON(w, http.StatusOK, sufficientBalanceResponse{
			CardNo:            cardNo,
			Amount:            amount,
			SufficientBalance: sufficient,
		})
	}
}
