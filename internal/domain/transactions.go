package domain

import "time"

type Transaction struct {
	ID         int
	SenderID   int
	ReceiverID int
	Type       string // Transaction type: "gift", "exchange", "lottery"
	Amount     int    // Points amount
	CreatedAt  time.Time
}
