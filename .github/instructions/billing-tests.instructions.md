---
applyTo: "internal/billing/*_test.go"
---

# Billing tests

`TestSuggestedCreditUsesV2Cap` fails on a clean tree on purpose. The client constant `SUGGESTED_CREDIT_API_VERSION` ships as v1. Do not edit that test to make the suite green.

Dispute `dsp_1043` claims 40000 cents against a Scale invoice of 24900 cents. v2 caps a suggestion at the plan price. Keep both API routes.
