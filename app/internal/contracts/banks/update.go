package banks

import "api/internal/contracts"

var Update = contracts.Contract{
    Method: "PUT",
	URI:    "/banks/{id}",
	Required: map[string]contracts.FieldSpec{
		"id": {
			Type: "int",
			Min:  1,
		},
	},
    Optional: map[string]contracts.FieldSpec{
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