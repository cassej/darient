package banks

import "api/internal/contracts"

var Create = contracts.Contract{
	Required: map[string]contracts.FieldSpec{
		"name": {
			Type: "string",
			Min:  2,
			Max:  100,
		},
		"type": {
			Type:    "enum",
			Options: []string{"PRIVATE", "GOVERNMENT"},
		},
	},
}