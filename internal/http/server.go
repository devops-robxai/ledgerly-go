package httpserver

import (
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"strings"

	"github.com/devops-robxai/ledgerly-go/internal/billing"
	"github.com/devops-robxai/ledgerly-go/internal/runbooks"
)

type Server struct {
	Store    *billing.Store
	tmpl     *template.Template
	mux      *http.ServeMux
	runbooks fs.FS
}

// NewFromFS parses HTML templates from templatesFS (*.html) and serves staticRoot at /static/.
// runbooksFS is the directory of workshop markdown (101.md); may be nil (routes 404).
func NewFromFS(store *billing.Store, templatesFS fs.FS, staticRoot fs.FS, runbooksFS fs.FS) (*Server, error) {
	funcs := template.FuncMap{
		"usd": billing.FormatUSD,
		"planLabel": func(p billing.PlanID) string {
			return billing.PlanLabel(p)
		},
		"md": runbooks.InlineHTML,
	}
	tmpl, err := template.New("").Funcs(funcs).ParseFS(templatesFS, "*.html")
	if err != nil {
		return nil, err
	}
	s := &Server{Store: store, tmpl: tmpl, mux: http.NewServeMux(), runbooks: runbooksFS}
	s.routes(staticRoot)
	return s, nil
}

func (s *Server) Handler() http.Handler { return s.mux }

func (s *Server) routes(staticRoot fs.FS) {
	if staticRoot != nil {
		s.mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.FS(staticRoot))))
	}

	s.mux.HandleFunc("/", s.handleDashboard)
	s.mux.HandleFunc("/invoices", s.handleInvoices)
	s.mux.HandleFunc("/invoices/", s.handleInvoiceDetail)
	s.mux.HandleFunc("/disputes", s.handleDisputes)
	s.mux.HandleFunc("/disputes/", s.handleDisputeDetail)

	s.mux.HandleFunc("/api/v1/disputes/", s.handleSuggestedCreditV1)
	s.mux.HandleFunc("/api/v2/disputes/", s.handleSuggestedCreditV2)

	s.mux.HandleFunc("/runbooks", s.handleRunbooks)
	s.mux.HandleFunc("/runbooks/", s.handleRunbooks)
	s.mux.HandleFunc("/workflows", s.redirectToRunbooks101)
	s.mux.HandleFunc("/analysis", s.redirectToRunbooks101)

	s.mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})
}

func (s *Server) render(w http.ResponseWriter, name string, data any) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := s.tmpl.ExecuteTemplate(w, name, data); err != nil {
		log.Printf("render %s: %v", name, err)
		http.Error(w, "template error", http.StatusInternalServerError)
	}
}

func trimID(prefix, p string) string {
	p = strings.TrimPrefix(p, prefix)
	p = strings.Trim(p, "/")
	if i := strings.IndexByte(p, '/'); i >= 0 {
		return p[:i]
	}
	return p
}

func disputeAPIID(r *http.Request) (id string, ok bool) {
	rest := r.URL.Path
	for _, prefix := range []string{"/api/v1/disputes/", "/api/v2/disputes/"} {
		if strings.HasPrefix(rest, prefix) {
			rest = strings.TrimPrefix(rest, prefix)
			break
		}
	}
	parts := strings.Split(strings.Trim(rest, "/"), "/")
	if len(parts) < 2 || parts[1] != "suggested-credit" {
		return "", false
	}
	return parts[0], true
}
