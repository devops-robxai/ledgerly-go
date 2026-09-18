package billing

import "testing"

// TestSuggestedCreditUsesV2Cap is the planted red test for Cursor 101.
//
// It fails until SUGGESTED_CREDIT_API_VERSION (or the dispute-detail call site)
// uses v2. Do not edit this test to get green — flip the client to v2 instead.
func TestSuggestedCreditUsesV2Cap(t *testing.T) {
	path := SuggestedCreditPath("dsp_1043")
	want := "/api/v2/disputes/dsp_1043/suggested-credit"
	if path != want {
		t.Fatalf("suggested credit client still on deprecated API:\n  got  %s\n  want %s\n\nWorkshop fix: set SUGGESTED_CREDIT_API_VERSION to \"v2\" so dsp_1043 caps at $249 (Scale).", path, want)
	}

	store := NewStore()
	Seed(store)
	resp, err := store.ResolveSuggestedCredit("dsp_1043")
	if err != nil {
		t.Fatalf("ResolveSuggestedCredit: %v", err)
	}
	if resp.SuggestedCreditCents != 24900 {
		t.Fatalf("dsp_1043 suggested credit = %d cents (%s), want 24900 ($249.00 Scale cap)", resp.SuggestedCreditCents, FormatUSD(resp.SuggestedCreditCents))
	}
	if resp.APIVersion != "v2" {
		t.Fatalf("apiVersion = %q, want %q", resp.APIVersion, "v2")
	}
}

func TestSuggestDisputeCreditCapsAtPlan(t *testing.T) {
	got, err := SuggestDisputeCredit(40000, 24900)
	if err != nil {
		t.Fatal(err)
	}
	if got != 24900 {
		t.Fatalf("got %d, want 24900", got)
	}
}

func TestV1ReturnsRawClaim(t *testing.T) {
	store := NewStore()
	Seed(store)
	resp, err := store.SuggestedCreditV1("dsp_1043")
	if err != nil {
		t.Fatal(err)
	}
	if resp.SuggestedCreditCents != 40000 {
		t.Fatalf("v1 should return raw claim 40000, got %d", resp.SuggestedCreditCents)
	}
}

func TestV2CapsScale(t *testing.T) {
	store := NewStore()
	Seed(store)
	resp, err := store.SuggestedCreditV2("dsp_1043")
	if err != nil {
		t.Fatal(err)
	}
	if resp.SuggestedCreditCents != 24900 {
		t.Fatalf("v2 should cap at 24900, got %d", resp.SuggestedCreditCents)
	}
}
