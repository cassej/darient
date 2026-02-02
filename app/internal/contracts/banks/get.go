package banks

import "api/internal/contracts"

var Get = contracts.Contract{
	Method: "GET",
	URI:    "/banks/{id}",
	Required: map[string]contracts.FieldSpec{
		"id": {
			Type: "int",
			Min:  1,
		},
	},
}