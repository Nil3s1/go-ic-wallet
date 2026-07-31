package wallet

import "time"

type CardModel struct {
	CardNo         string
	ValidTo        time.Time
	CurrentBalance int //Currency in cents
}
