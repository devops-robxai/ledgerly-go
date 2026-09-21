package httpserver

import (
	"encoding/json"
	"net/http"
	"sort"
	"strings"

	"github.com/devops-robxai/ledgerly-go/internal/billing"
)

type pageData struct {
	Title    string
	Active   string
	Flash    string
	Customer *billing.Customer
	Invoice  *billing.Invoice
	Dispute  *billing.Dispute
	Invoices []*billing.Invoice
	Disputes []*billing.Dispute

	SuggestedCreditCents int
	SuggestedCreditUSD   string
	SuggestedAPIVersion  string
	SuggestedAPIPath     string

	Stats struct {
		InvoiceCount int
		OpenInvoices int
		DisputeCount int
		OverdueCount int
	}
}

func (s *Server) handleDashboard(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	invoices := s.Store.ListInvoices()
	disputes := s.Store.ListDisputes()
	data := pageData{Title: "Dashboard", Active: "dashboard", Invoices: invoices, Disputes: disputes}
	data.Stats.InvoiceCount = len(invoices)
	data.Stats.DisputeCount = len(disputes)
	for _, inv := range invoices {
		if inv.Status == "open" || inv.Status == "disputed" {
			data.Stats.OpenInvoices++
		}
		if inv.Status == "overdue" {
			data.Stats.OverdueCount++
		}
	}
	sortInvoices(data.Invoices)
	s.render(w, "dashboard.html", data)
}

func (s *Server) handleInvoices(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/invoices" {
		http.NotFound(w, r)
		return
	}
	invoices := s.Store.ListInvoices()
	sortInvoices(invoices)
	s.render(w, "invoices.html", pageData{
		Title: "Invoices", Active: "invoices", Invoices: invoices,
	})
}

func (s *Server) handleInvoiceDetail(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/invoices/")
	path = strings.Trim(path, "/")
	if path == "" {
		http.NotFound(w, r)
		return
	}

	// POST /invoices/{id}/email
	if strings.HasSuffix(path, "/email") {
		id := strings.TrimSuffix(path, "/email")
		id = strings.TrimSuffix(id, "/")
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		inv, ok := s.Store.GetInvoice(id)
		if !ok {
			http.NotFound(w, r)
			return
		}
		if err := r.ParseForm(); err != nil {
			http.Error(w, "bad form", http.StatusBadRequest)
			return
		}
		email := strings.TrimSpace(r.FormValue("email"))
		if email == "" {
			http.Error(w, "email required", http.StatusBadRequest)
			return
		}
		if err := s.Store.UpdateCustomerEmail(inv.CustomerID, email); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Redirect(w, r, "/invoices/"+id+"?saved=1", http.StatusSeeOther)
		return
	}

	if strings.Contains(path, "/") {
		http.NotFound(w, r)
		return
	}

	inv, ok := s.Store.GetInvoice(path)
	if !ok {
		http.NotFound(w, r)
		return
	}
	cust, _ := s.Store.GetCustomer(inv.CustomerID)
	flash := ""
	if r.URL.Query().Get("saved") == "1" {
		flash = "Customer email updated."
	}
	s.render(w, "invoice_detail.html", pageData{
		Title: "Invoice " + path, Active: "invoices",
		Invoice: inv, Customer: cust, Flash: flash,
	})
}

func (s *Server) handleDisputes(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/disputes" {
		http.NotFound(w, r)
		return
	}
	disputes := s.Store.ListDisputes()
	sortDisputes(disputes)
	s.render(w, "disputes.html", pageData{
		Title: "Disputes", Active: "disputes", Disputes: disputes,
	})
}

func (s *Server) handleDisputeDetail(w http.ResponseWriter, r *http.Request) {
	id := trimID("/disputes/", r.URL.Path)
	if id == "" || strings.Contains(id, "/") {
		http.NotFound(w, r)
		return
	}
	d, ok := s.Store.GetDispute(id)
	if !ok {
		http.NotFound(w, r)
		return
	}
	inv, _ := s.Store.GetInvoice(d.InvoiceID)
	var cust *billing.Customer
	if inv != nil {
		cust, _ = s.Store.GetCustomer(inv.CustomerID)
	}

	// PLANTED SEAM: follows SUGGESTED_CREDIT_API_VERSION (v1) → $400 for dsp_1043.
	credit, err := s.Store.ResolveSuggestedCredit(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	s.render(w, "dispute_detail.html", pageData{
		Title: "Dispute " + id, Active: "disputes",
		Dispute: d, Invoice: inv, Customer: cust,
		SuggestedCreditCents: credit.SuggestedCreditCents,
		SuggestedCreditUSD:   billing.FormatUSD(credit.SuggestedCreditCents),
		SuggestedAPIVersion:  credit.APIVersion,
		SuggestedAPIPath:     billing.SuggestedCreditPath(id),
	})
}

func (s *Server) handleSuggestedCreditV1(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	id, ok := disputeAPIID(r)
	if !ok {
		http.NotFound(w, r)
		return
	}
	resp, err := s.Store.SuggestedCreditV1(id)
	if err != nil {
		http.Error(w, `{"error":"Dispute not found"}`, http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Deprecation", "true")
	w.Header().Set("Link", `</api/v2/disputes/`+id+`/suggested-credit>; rel="successor-version"`)
	_ = json.NewEncoder(w).Encode(resp)
}

func (s *Server) handleSuggestedCreditV2(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	id, ok := disputeAPIID(r)
	if !ok {
		http.NotFound(w, r)
		return
	}
	resp, err := s.Store.SuggestedCreditV2(id)
	if err != nil {
		http.Error(w, `{"error":"Dispute not found"}`, http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)
}

func sortInvoices(invoices []*billing.Invoice) {
	sort.Slice(invoices, func(i, j int) bool { return invoices[i].ID < invoices[j].ID })
}

func sortDisputes(disputes []*billing.Dispute) {
	sort.Slice(disputes, func(i, j int) bool { return disputes[i].ID < disputes[j].ID })
}
