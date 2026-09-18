package billing

import (
	"fmt"
	"sync"
)

// Store is an in-memory billing book for the workshop demo.
type Store struct {
	mu        sync.RWMutex
	customers map[string]*Customer
	invoices  map[string]*Invoice
	disputes  map[string]*Dispute
}

func NewStore() *Store {
	return &Store{
		customers: make(map[string]*Customer),
		invoices:  make(map[string]*Invoice),
		disputes:  make(map[string]*Dispute),
	}
}

func (s *Store) Reset() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.customers = make(map[string]*Customer)
	s.invoices = make(map[string]*Invoice)
	s.disputes = make(map[string]*Dispute)
	seedLocked(s)
}

func (s *Store) ListCustomers() []*Customer {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]*Customer, 0, len(s.customers))
	for _, c := range s.customers {
		cp := *c
		out = append(out, &cp)
	}
	return out
}

func (s *Store) GetCustomer(id string) (*Customer, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	c, ok := s.customers[id]
	if !ok {
		return nil, false
	}
	cp := *c
	return &cp, true
}

func (s *Store) UpdateCustomerEmail(id, email string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	c, ok := s.customers[id]
	if !ok {
		return fmt.Errorf("customer %s not found", id)
	}
	c.Email = email
	return nil
}

func (s *Store) ListInvoices() []*Invoice {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]*Invoice, 0, len(s.invoices))
	for _, inv := range s.invoices {
		cp := *inv
		out = append(out, &cp)
	}
	return out
}

func (s *Store) GetInvoice(id string) (*Invoice, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	inv, ok := s.invoices[id]
	if !ok {
		return nil, false
	}
	cp := *inv
	return &cp, true
}

func (s *Store) ListDisputes() []*Dispute {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]*Dispute, 0, len(s.disputes))
	for _, d := range s.disputes {
		cp := *d
		out = append(out, &cp)
	}
	return out
}

func (s *Store) GetDispute(id string) (*Dispute, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	d, ok := s.disputes[id]
	if !ok {
		return nil, false
	}
	cp := *d
	return &cp, true
}

// SuggestedCreditV1 returns the raw disputed claim (deprecated).
func (s *Store) SuggestedCreditV1(disputeID string) (*SuggestedCreditResponse, error) {
	d, ok := s.GetDispute(disputeID)
	if !ok {
		return nil, fmt.Errorf("dispute not found")
	}
	return &SuggestedCreditResponse{
		DisputeID:            d.ID,
		SuggestedCreditCents: d.DisputedAmountCents,
		APIVersion:           "v1",
	}, nil
}

// SuggestedCreditV2 caps the suggestion at the invoice plan price.
func (s *Store) SuggestedCreditV2(disputeID string) (*SuggestedCreditResponse, error) {
	d, ok := s.GetDispute(disputeID)
	if !ok {
		return nil, fmt.Errorf("dispute not found")
	}
	inv, ok := s.GetInvoice(d.InvoiceID)
	if !ok {
		return nil, fmt.Errorf("invoice not found")
	}
	price, err := PlanPriceCents(inv.Plan)
	if err != nil {
		return nil, err
	}
	cents, err := SuggestDisputeCredit(d.DisputedAmountCents, price)
	if err != nil {
		return nil, err
	}
	return &SuggestedCreditResponse{
		DisputeID:            d.ID,
		SuggestedCreditCents: cents,
		APIVersion:           "v2",
	}, nil
}

// ResolveSuggestedCredit follows the planted client version constant.
// Dispute HTML detail uses this so the page shows v1's $400 today.
func (s *Store) ResolveSuggestedCredit(disputeID string) (*SuggestedCreditResponse, error) {
	switch SUGGESTED_CREDIT_API_VERSION {
	case "v2":
		return s.SuggestedCreditV2(disputeID)
	default:
		return s.SuggestedCreditV1(disputeID)
	}
}
