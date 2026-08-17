package api

import (
	"net/http"

	journeyapplication "github.com/Nil3s1/go-ic-wallet/internal/journey/application"
	walletapplication "github.com/Nil3s1/go-ic-wallet/internal/wallet/application"
)

type Dependencies struct {
	CreateCard           *walletapplication.CreateCardCommandHandler
	AddBalance           *walletapplication.AddBalanceCommandHandler
	ApplyPayment         *walletapplication.ApplyPaymentCommandHandler
	HasSufficientBalance *walletapplication.HasSufficientBalanceQueryHandler
	StartJourney         *journeyapplication.StartJourneyCommandHandler
	EndJourney           *journeyapplication.EndJourneyCommandHandler
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
