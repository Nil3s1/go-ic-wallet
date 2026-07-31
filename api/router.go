package api

import (
	"net/http"

	"github.com/Nil3s1/go-ic-wallet/internal/journey"
	"github.com/Nil3s1/go-ic-wallet/internal/wallet"
)

type Dependencies struct {
	CreateCard           *wallet.CreateCardCommandHandler
	AddBalance           *wallet.AddBalanceCommandHandler
	ApplyPayment         *wallet.ApplyPaymentCommandHandler
	HasSufficientBalance *wallet.HasSufficientBalanceQueryHandler
	StartJourney         *journey.StartJourneyCommandHandler
	EndJourney           *journey.EndJourneyCommandHandler
}

func NewRouter(deps Dependencies) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /wallet/cards/new", createCardHandler(deps.CreateCard))
	mux.HandleFunc("POST /wallet/cards/{cardNo}/addBalance", addBalanceHandler(deps.AddBalance))
	mux.HandleFunc("POST /wallet/cards/{cardNo}/applyPayment", applyPaymentHandler(deps.ApplyPayment))

	mux.HandleFunc("POST /journey/start", startJourneyHandler(deps.StartJourney))
	mux.HandleFunc("POST /journey/end", endJourneyHandler(deps.EndJourney))

	return mux
}
