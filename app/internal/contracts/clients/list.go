package clients

import "api/internal/contracts"

var List = contracts.Contract{
    Method: "GET",
	URI:    "/clients",
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