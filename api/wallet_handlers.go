package api

import (
	"encoding/json"
	"net/http"

	walletapplication "github.com/Nil3s1/go-ic-wallet/internal/wallet/application"
)

type createCardRequest struct {
	InitialBalance uint `json:"initialBalance"`
}

type addBalanceRequest struct {
	Amount uint `json:"amount"`
}

type applyPaymentRequest struct {
	Amount      uint   `json:"amount"`
	ReferenceId string `json:"referenceId"`
}

type createCardResponse struct {
	CardNo  string `json:"cardNo"`
	ValidTo string `json:"validTo"`
}

type balanceResponse struct {
	CardNo     string `json:"cardNo"`
	OldBalance uint   `json:"oldBalance"`
	NewBalance uint   `json:"newBalance"`
}

func createCardHandler(h *walletapplication.CreateCardCommandHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req createCardRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}

		cmd := walletapplication.CreateCardCommand{
			InitialBalance: req.InitialBalance,
		}
		result, err := h.Handle(r.Context(), cmd)

		if err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}

		writeJSON(w, http.StatusCreated, createCardResponse{
			CardNo:  result.CardNo,
			ValidTo: result.ValidTo.Format("2006-01-02"),
		})
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

		result, err := h.Handle(r.Context(), cmd)
		if err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}

		writeJSON(w, http.StatusOK, balanceResponse{
			CardNo:     result.CardNo,
			OldBalance: result.OldBalance,
			NewBalance: result.NewBalance,
		})
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

		result, err := h.Handle(r.Context(), cmd)
		if err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}

		writeJSON(w, http.StatusOK, balanceResponse{
			CardNo:     result.CardNo,
			OldBalance: result.OldBalance,
			NewBalance: result.NewBalance,
		})
	}
}
