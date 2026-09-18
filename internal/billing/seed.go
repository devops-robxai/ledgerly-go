package billing

import "time"

// DemoClock is frozen so overdue math never drifts during a workshop.
var DemoClock = time.Date(2026, 8, 23, 12, 0, 0, 0, time.UTC)

func Seed(s *Store) {
	s.mu.Lock()
	defer s.mu.Unlock()
	seedLocked(s)
}

func seedLocked(s *Store) {
	s.customers = map[string]*Customer{
		"cus_acme": {
			ID: "cus_acme", Name: "Acme North", Email: "billing@acme-north.example",
		},
		"cus_river": {
			ID: "cus_river", Name: "Riverstone Labs", Email: "ap@riverstone.example",
		},
		"cus_cobalt": {
			ID: "cus_cobalt", Name: "Cobalt Goods", Email: "finance@cobalt-goods.example",
		},
		"cus_bright": {
			ID: "cus_bright", Name: "Brightwell Labs", Email: "accounts@brightwell.example",
		},
		"cus_oak": {
			ID: "cus_oak", Name: "Oakiron Supply", Email: "payables@oakiron.example",
		},
	}

	s.invoices = map[string]*Invoice{
		"inv_1001": {
			ID: "inv_1001", CustomerID: "cus_acme", Plan: PlanStarter,
			AmountCents: 4900, Status: "paid",
			IssuedAt: DemoClock.AddDate(0, -2, 0), DueAt: DemoClock.AddDate(0, -1, -15),
			Memo: "Starter plan — June cycle.",
		},
		"inv_1022": {
			ID: "inv_1022", CustomerID: "cus_river", Plan: PlanGrowth,
			AmountCents: 9900, Status: "overdue",
			IssuedAt: DemoClock.AddDate(0, -1, -10), DueAt: DemoClock.AddDate(0, 0, -20),
			Memo: "Growth plan. Reminder sent.",
		},
		"inv_1035": {
			ID: "inv_1035", CustomerID: "cus_cobalt", Plan: PlanScale,
			AmountCents: 24900, Status: "open",
			IssuedAt: DemoClock.AddDate(0, 0, -5), DueAt: DemoClock.AddDate(0, 0, 10),
			Memo: "Scale plan — July cycle.",
		},
		"inv_1048": {
			ID: "inv_1048", CustomerID: "cus_bright", Plan: PlanScale,
			AmountCents: 24900, Status: "disputed",
			IssuedAt: DemoClock.AddDate(0, 0, -12), DueAt: DemoClock.AddDate(0, 0, -2),
			Memo: "Scale plan. Service-window dispute needs review.",
		},
		"inv_1051": {
			ID: "inv_1051", CustomerID: "cus_oak", Plan: PlanGrowth,
			AmountCents: 9900, Status: "paid",
			IssuedAt: DemoClock.AddDate(0, -1, 0), DueAt: DemoClock.AddDate(0, 0, -15),
			Memo: "Growth plan — paid mid-August.",
		},
	}

	s.disputes = map[string]*Dispute{
		// PLANTED: claims $400 against a $249 Scale invoice. v1 returns 40000; v2 caps at 24900.
		"dsp_1043": {
			ID: "dsp_1043", InvoiceID: "inv_1048", Status: "needs_review",
			Reason: "Billed on Scale. The signed order is Growth. Customer is claiming back more than this invoice charges.",
			DisputedAmountCents: 40000,
			OpenedAt:            DemoClock.AddDate(0, 0, -3),
		},
		"dsp_1020": {
			ID: "dsp_1020", InvoiceID: "inv_1022", Status: "open",
			Reason: "Duplicate charge on Growth renewal.",
			DisputedAmountCents: 9900,
			OpenedAt:            DemoClock.AddDate(0, 0, -8),
		},
		"dsp_0991": {
			ID: "dsp_0991", InvoiceID: "inv_1035", Status: "open",
			Reason: "Partial outage credit request.",
			DisputedAmountCents: 24900,
			OpenedAt:            DemoClock.AddDate(0, 0, -1),
		},
	}
}
