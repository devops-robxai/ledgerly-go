package billing

import "fmt"

// SUGGESTED_CREDIT_API_VERSION selects which suggested-credit API the UI/client uses.
//
// PLANTED SEAM (workshop): this ships as "v1" on purpose. v1 returns the raw
// disputed claim ($400 for dsp_1043). v2 caps at the Scale plan price ($249).
// The failing test in suggested_credit_test.go expects "v2". Demo fix: flip
// this constant (or the call site) to "v2".
const SUGGESTED_CREDIT_API_VERSION = "v1"

// SuggestedCreditPath returns the HTTP path the client should call.
func SuggestedCreditPath(disputeID string) string {
	return fmt.Sprintf("/api/%s/disputes/%s/suggested-credit", SUGGESTED_CREDIT_API_VERSION, disputeID)
}

// SuggestDisputeCredit applies the v2 catalog cap: never suggest more than the plan price.
func SuggestDisputeCredit(disputedAmountCents, planPriceCents int) (int, error) {
	if disputedAmountCents < 0 || planPriceCents < 0 {
		return 0, fmt.Errorf("credit inputs must be non-negative cents")
	}
	if disputedAmountCents < planPriceCents {
		return disputedAmountCents, nil
	}
	return planPriceCents, nil
}

// SuggestedCreditResponse is the JSON shape for both v1 and v2 endpoints.
type SuggestedCreditResponse struct {
	DisputeID            string `json:"disputeId"`
	SuggestedCreditCents int    `json:"suggestedCreditCents"`
	APIVersion           string `json:"apiVersion"`
}
