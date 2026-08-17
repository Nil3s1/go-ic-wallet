package application

import "time"

type CreateCardResult struct {
	CardNo  string
	ValidTo time.Time
}

type BalanceResult struct {
	CardNo     string
	OldBalance uint
	NewBalance uint
}
