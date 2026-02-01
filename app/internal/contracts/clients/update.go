package clients

import "api/internal/contracts"

var Update = contracts.Contract{
	Optional: map[string]contracts.FieldSpec{
		"full_name": {
			Type: "string",
			Min:  2,
			Max:  255,
		},
		"email": {
			Type: "email",
		},
		"birth_date": {
			Type: "date",
		},
		"country": {
			Type: "string",
			Min:  2,
			Max:  100,
		},
	},
}