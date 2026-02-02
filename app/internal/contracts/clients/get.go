package clients

import "api/internal/contracts"

var Get = contracts.Contract{
	Method: "GET",
	URI:    "/clients/{id}",
	Required: map[string]contracts.FieldSpec{
		"id": {
			Type: "int",
			Min:  1,
		},
	},
}