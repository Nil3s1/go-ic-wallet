package wallet

import "time"

type CardProjection struct {
	CardNo         string
	ValidTo        time.Time
	CurrentBalance int //Currency in cents
}
