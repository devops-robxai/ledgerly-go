package billing

import "fmt"

// Catalog prices are frozen for the demo. Never invent a fourth plan.
const (
	PlanStarter PlanID = "STARTER"
	PlanGrowth  PlanID = "GROWTH"
	PlanScale   PlanID = "SCALE"
)

type PlanID string

var planPriceCents = map[PlanID]int{
	PlanStarter: 4900,
	PlanGrowth:  9900,
	PlanScale:   24900, // $249 — v2 credit cap for Scale
}

var planLabel = map[PlanID]string{
	PlanStarter: "Starter",
	PlanGrowth:  "Growth",
	PlanScale:   "Scale",
}

func IsPlanID(s string) bool {
	_, ok := planPriceCents[PlanID(s)]
	return ok
}

func PlanPriceCents(plan PlanID) (int, error) {
	cents, ok := planPriceCents[plan]
	if !ok {
		return 0, fmt.Errorf("unknown plan %q (catalog: Starter $49, Growth $99, Scale $249)", plan)
	}
	return cents, nil
}

func PlanLabel(plan PlanID) string {
	if label, ok := planLabel[plan]; ok {
		return label
	}
	return string(plan)
}
