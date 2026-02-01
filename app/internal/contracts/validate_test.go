package contracts

import (
	"errors"
	"fmt"
	"reflect"
	"testing"
)

func TestValidate(t *testing.T) {
	tests := []struct {
		name     string
		input    map[string]any
		contract Contract
		wantErr  error
		want     map[string]any
	}{
		{
			name:  "valid required fields",
			input: map[string]any{"name": "Test", "type": "PRIVATE"},
			contract: Contract{
				Required: map[string]FieldSpec{
					"name": {Type: "string", Min: 2, Max: 100},
					"type": {Type: "enum", Options: []string{"PRIVATE", "GOVERNMENT"}},
				},
			},
			wantErr: nil,
			want:    map[string]any{"name": "Test", "type": "PRIVATE"},
		},
		{
			name:  "missing required field",
			input: map[string]any{"type": "PRIVATE"},
			contract: Contract{
				Required: map[string]FieldSpec{
					"name": {Type: "string"},
				},
			},
			wantErr: fmt.Errorf("%w: name", ErrRequired),
			want:    nil,
		},
		{
			name:  "invalid optional field",
			input: map[string]any{"extra": "invalid"},
			contract: Contract{
				Optional: map[string]FieldSpec{
					"optional": {Type: "string"},
				},
			},
			wantErr: fmt.Errorf("%w: extra", ErrUnexpectedField),
			want:    nil,
		},
		{
			name:  "nil input for optional",
			input: map[string]any{"optional": nil},
			contract: Contract{
				Optional: map[string]FieldSpec{
					"optional": {Type: "string"},
				},
			},
			wantErr: nil,
			want:    map[string]any{},
		},
		{
			name:  "unexpected field",
			input: map[string]any{"name": "Test", "unexpected": 123},
			contract: Contract{
				Required: map[string]FieldSpec{
					"name": {Type: "string"},
				},
			},
			wantErr: fmt.Errorf("%w: unexpected", ErrUnexpectedField),
			want:    nil,
		},
		{
			name:  "empty input for empty contract",
			input: map[string]any{},
			contract: Contract{
				Required: map[string]FieldSpec{},
				Optional: map[string]FieldSpec{},
			},
			wantErr: nil,
			want:    map[string]any{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Validate(tt.input, tt.contract)
			if !errors.Is(err, tt.wantErr) && (tt.wantErr == nil || err.Error() != tt.wantErr.Error()) {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Validate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValidateField(t *testing.T) {
	tests := []struct {
		name    string
		field   string
		value   any
		spec    FieldSpec
		wantErr error
	}{
		// string
		{
			name:  "valid string",
			field: "name",
			value: "Test",
			spec:  FieldSpec{Type: "string", Min: 2, Max: 100},
			wantErr: nil,
		},
		{
			name:  "string too short",
			field: "name",
			value: "T",
			spec:  FieldSpec{Type: "string", Min: 2},
			wantErr: ErrTooShort,
		},
		{
			name:  "string too long",
			field: "name",
			value: strings.Repeat("a", 101),
			spec:  FieldSpec{Type: "string", Max: 100},
			wantErr: ErrTooLong,
		},
		{
			name:  "string invalid type",
			field: "name",
			value: 123,
			spec:  FieldSpec{Type: "string"},
			wantErr: ErrInvalidType,
		},

		// email
		{
			name:  "valid email",
			field: "email",
			value: "test@example.com",
			spec:  FieldSpec{Type: "email"},
			wantErr: nil,
		},
		{
			name:  "invalid email",
			field: "email",
			value: "invalid",
			spec:  FieldSpec{Type: "email"},
			wantErr: ErrInvalidEmail,
		},
		{
			name:  "email invalid type",
			field: "email",
			value: 123,
			spec:  FieldSpec{Type: "email"},
			wantErr: ErrInvalidType,
		},

		// date
		{
			name:  "valid date",
			field: "date",
			value: "2026-02-01",
			spec:  FieldSpec{Type: "date"},
			wantErr: nil,
		},
		{
			name:  "invalid date",
			field: "date",
			value: "2026-13-01",
			spec:  FieldSpec{Type: "date"},
			wantErr: ErrInvalidDate,
		},
		{
			name:  "date invalid type",
			field: "date",
			value: 123,
			spec:  FieldSpec{Type: "date"},
			wantErr: ErrInvalidType,
		},

		// int
		{
			name:  "valid int",
			field: "age",
			value: 30,
			spec:  FieldSpec{Type: "int", Min: 18, Max: 99},
			wantErr: nil,
		},
		{
			name:  "int too small",
			field: "age",
			value: 17,
			spec:  FieldSpec{Type: "int", Min: 18},
			wantErr: ErrTooSmall,
		},
		{
			name:  "int too big",
			field: "age",
			value: 100,
			spec:  FieldSpec{Type: "int", Max: 99},
			wantErr: ErrTooBig,
		},
		{
			name:  "int from float without fraction",
			field: "age",
			value: 30.0,
			spec:  FieldSpec{Type: "int"},
			wantErr: nil,
		},
		{
			name:  "int from float with fraction",
			field: "age",
			value: 30.5,
			spec:  FieldSpec{Type: "int"},
			wantErr: ErrInvalidType,
		},
		{
			name:  "int invalid type",
			field: "age",
			value: "thirty",
			spec:  FieldSpec{Type: "int"},
			wantErr: ErrInvalidType,
		},

		// number
		{
			name:  "valid number",
			field: "price",
			value: 99.99,
			spec:  FieldSpec{Type: "number", MinVal: 0.01, MaxVal: 1000.00},
			wantErr: nil,
		},
		{
			name:  "number too small",
			field: "price",
			value: 0.00,
			spec:  FieldSpec{Type: "number", MinVal: 0.01},
			wantErr: ErrTooSmall,
		},
		{
			name:  "number too big",
			field: "price",
			value: 1000.01,
			spec:  FieldSpec{Type: "number", MaxVal: 1000.00},
			wantErr: ErrTooBig,
		},
		{
			name:  "number from int",
			field: "price",
			value: 99,
			spec:  FieldSpec{Type: "number"},
			wantErr: nil,
		},
		{
			name:  "number invalid type",
			field: "price",
			value: "99.99",
			spec:  FieldSpec{Type: "number"},
			wantErr: ErrInvalidType,
		},

		// enum
		{
			name:  "valid enum",
			field: "type",
			value: "private",
			spec:  FieldSpec{Type: "enum", Options: []string{"PRIVATE", "GOVERNMENT"}},
			wantErr: nil,
		},
		{
			name:  "invalid enum",
			field: "type",
			value: "corp",
			spec:  FieldSpec{Type: "enum", Options: []string{"PRIVATE", "GOVERNMENT"}},
			wantErr: ErrInvalidEnum,
		},
		{
			name:  "enum invalid type",
			field: "type",
			value: 123,
			spec:  FieldSpec{Type: "enum", Options: []string{"A", "B"}},
			wantErr: ErrInvalidType,
		},

		// uuid
		{
			name:  "valid uuid",
			field: "id",
			value: "123e4567-e89b-12d3-a456-426614174000",
			spec:  FieldSpec{Type: "uuid"},
			wantErr: nil,
		},
		{
			name:  "invalid uuid",
			field: "id",
			value: "invalid-uuid",
			spec:  FieldSpec{Type: "uuid"},
			wantErr: ErrInvalidUUID,
		},
		{
			name:  "uuid invalid type",
			field: "id",
			value: 123,
			spec:  FieldSpec{Type: "uuid"},
			wantErr: ErrInvalidType,
		},

		// unsupported
		{
			name:  "unsupported type",
			field: "unknown",
			value: "value",
			spec:  FieldSpec{Type: "unknown"},
			wantErr: ErrUnsupportedType,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateField(tt.field, tt.value, tt.spec)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("ValidateField() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNormalize(t *testing.T) {
	tests := []struct {
		name  string
		value any
		spec  FieldSpec
		want  any
	}{
		{
			name:  "normalize string",
			value: "  Test  ",
			spec:  FieldSpec{Type: "string"},
			want:  "Test",
		},
		{
			name:  "normalize email",
			value: "  TEST@EXAMPLE.COM  ",
			spec:  FieldSpec{Type: "email"},
			want:  "test@example.com",
		},
		{
			name:  "normalize uuid",
			value: "  123E4567-E89B-12D3-A456-426614174000  ",
			spec:  FieldSpec{Type: "uuid"},
			want:  "123e4567-e89b-12d3-a456-426614174000",
		},
		{
			name:  "normalize date",
			value: "  2026-02-01  ",
			spec:  FieldSpec{Type: "date"},
			want:  "2026-02-01",
		},
		{
			name:  "normalize int from float",
			value: 30.0,
			spec:  FieldSpec{Type: "int"},
			want:  30,
		},
		{
			name:  "normalize int from int",
			value: 30,
			spec:  FieldSpec{Type: "int"},
			want:  30,
		},
		{
			name:  "normalize number from int",
			value: 99,
			spec:  FieldSpec{Type: "number"},
			want:  99.0,
		},
		{
			name:  "normalize number from float",
			value: 99.99,
			spec:  FieldSpec{Type: "number"},
			want:  99.99,
		},
		{
			name:  "normalize enum",
			value: "private",
			spec:  FieldSpec{Type: "enum"},
			want:  "PRIVATE",
		},
		{
			name:  "normalize unknown type",
			value: "value",
			spec:  FieldSpec{Type: "unknown"},
			want:  "value",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Normalize(tt.value, tt.spec)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Normalize() = %v, want %v", got, tt.want)
			}
		})
	}
}