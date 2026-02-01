package credits

import "api/internal/contracts"

var Create = contracts.Contract{
	Required: map[string]contracts.FieldSpec{
		"client_id": {
			Type: "uuid",
		},
		"bank_id": {
			Type: "uuid",
		},
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
	},
}