package credits

import "api/internal/contracts"

var Delete = contracts.Contract{
	Method: "DELETE",
	URI:    "/credits/{id}",
	Required: map[string]contracts.FieldSpec{
		"id": {
			Type: "int",
			Min:  1,
		},
	},
}