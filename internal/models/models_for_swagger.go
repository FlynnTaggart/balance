package models

// Payloads for correct swagger generation

type PayloadId struct {
	ID uint64 `json:"id"`
}

type PayloadAddBalance struct {
	ID     uint64  `json:"id"`
	Amount float32 `json:"amount"`
}

type PayloadDate struct {
	Year  int `params:"year" json:"year"`
	Month int `params:"month" json:"month"`
}

type PayloadErr struct {
	Message string `json:"message"`
}

type PayloadBalance struct {
	Balance float32 `json:"balance"`
}

type PayloadReserve struct {
	UserID    uint64  `json:"user_id"`
	ServiceID uint64  `json:"service_id"`
	OrderID   uint64  `json:"order_id"`
	Amount    float32 `json:"amount,omitempty"`
}

type PayloadLink struct {
	Balance float32 `json:"report_link"`
}
