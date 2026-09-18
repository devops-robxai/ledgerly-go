package billing

import "time"

type Customer struct {
	ID    string
	Name  string
	Email string
}

type Invoice struct {
	ID          string
	CustomerID  string
	Plan        PlanID
	AmountCents int
	Status      string // paid | open | overdue | disputed
	IssuedAt    time.Time
	DueAt       time.Time
	Memo        string
}

type Dispute struct {
	ID                   string
	InvoiceID            string
	Status               string // open | needs_review | resolved
	Reason               string
	DisputedAmountCents  int
	OpenedAt             time.Time
}
