package clients

import "api/internal/contracts"

var Update = contracts.Contract{
    Method: "GET",
	URI:    "/clients/{id}",
	Required: map[string]contracts.FieldSpec{
        "id": {
            Type: "int",
            Min:  1,
        },
    },
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