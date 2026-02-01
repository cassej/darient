package credits

import "api/internal/contracts"

var Update = contracts.Contract{
	Optional: map[string]contracts.FieldSpec{
		"min_payment": {
			Type:   "number",
			MinVal: 0,
		},
		"max_payment": {
			Type:   "number",
			MinVal: 0,
		},
		"term_months": {
			Type: "int",
			Min:  1,
			Max:  360,
		},
		"credit_type": {
			Type:    "enum",
			Options: []string{"AUTO", "MORTGAGE", "COMMERCIAL"},
		},
		"status": {
			Type:    "enum",
			Options: []string{"PENDING", "APPROVED", "REJECTED"},
		},
	},
}