package banks

import "api/internal/contracts"

var Update = contracts.Contract{
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