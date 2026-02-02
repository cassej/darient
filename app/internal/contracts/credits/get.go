package clients

import "api/internal/contracts"

var Get = contracts.Contract{
	Method: "GET",
	URI:    "/credits/{id}",
	Required: map[string]contracts.FieldSpec{
		"id": {
			Type: "int",
			Min:  1,
		},
	},
}