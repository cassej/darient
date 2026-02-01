package contracts

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"
	"unicode/utf8"
)

var (
	ErrRequired        = errors.New("required field missing")
	ErrInvalidType     = errors.New("invalid type")
	ErrTooShort        = errors.New("value too short")
	ErrTooLong         = errors.New("value too long")
	ErrInvalidEnum     = errors.New("invalid enum value")
	ErrUnexpectedField = errors.New("unexpected field")
	ErrUnsupportedType = errors.New("unsupported field type in contract")
	ErrInvalidEmail    = errors.New("invalid email format")
	ErrInvalidDate     = errors.New("invalid date format")
	ErrTooSmall        = errors.New("value too small")
	ErrTooBig          = errors.New("value too big")
	ErrInvalidUUID     = errors.New("invalid UUID format")
)

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
var uuidRegex = regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`)

type Contract struct {
	Required map[string]FieldSpec
	Optional map[string]FieldSpec
}

type FieldSpec struct {
	Type    string
	Min     int
	Max     int
	MinVal  float64
	MaxVal  float64
	Options []string
}

func Validate(input map[string]any, c Contract) (map[string]any, error) {
	result := make(map[string]any)

	// required
	for field, spec := range c.Required {
		val, exists := input[field]

		if !exists || val == nil {
			return nil, fmt.Errorf("%w: %s", ErrRequired, field)
		}

		if err := ValidateField(field, val, spec); err != nil {
			return nil, fmt.Errorf("%s: %w", field, err)
		}

		result[field] = Normalize(val, spec)
	}

	// optional
	for field, spec := range c.Optional {
		val, exists := input[field]

		if !exists || val == nil {
			continue
		}

		if err := ValidateField(field, val, spec); err != nil {
			return nil, fmt.Errorf("%s: %w", field, err)
		}

		result[field] = Normalize(val, spec)
	}

	for field := range input {
		if _, isReq := c.Required[field]; isReq {
			continue
		}

		if _, isOpt := c.Optional[field]; isOpt {
			continue
		}

		return nil, fmt.Errorf("%w: %s", ErrUnexpectedField, field)
	}

	return result, nil
}

func ValidateField(field string, value any, spec FieldSpec) error {
	switch spec.Type {
        case "string":
            s, ok := value.(string)

            if !ok {
		    	return fmt.Errorf("%w: expected string, got %T", ErrInvalidType, value)
            }

            s = strings.TrimSpace(s)
            length := utf8.RuneCountInString(s)

            if spec.Min > 0 && length < spec.Min {
			    return fmt.Errorf("%w: min length %d, got %d", ErrTooShort, spec.Min, length)
            }

            if spec.Max > 0 && length > spec.Max {
			    return fmt.Errorf("%w: max length %d, got %d", ErrTooLong, spec.Max, length)
            }

            return nil

        case "email":
            s, ok := value.(string)

            if !ok {
		    	return fmt.Errorf("%w: expected string for email, got %T", ErrInvalidType, value)
            }

            s = strings.TrimSpace(s)

            if !emailRegex.MatchString(s) {
			    return fmt.Errorf("%w: %s", ErrInvalidEmail, s)
            }

            return nil

        case "uuid":
            s, ok := value.(string)

            if !ok {
                return fmt.Errorf("%w: expected string for email, got %T", ErrInvalidType, value)
            }

            s = strings.TrimSpace(s)

            if !uuidRegex.MatchString(s) {
                return fmt.Errorf("%w: %s", ErrInvalidEmail, s)
            }

            return nil

        case "date":
            s, ok := value.(string)

            if !ok {
		    	return fmt.Errorf("%w: expected string for date, got %T", ErrInvalidType, value)
            }

            s = strings.TrimSpace(s)

            _, err := time.Parse("2006-01-02", s)
            if err != nil {
			    return fmt.Errorf("%w: expected YYYY-MM-DD format", ErrInvalidDate)
            }

            return nil

        case "int":
            var n int

            switch v := value.(type) {
                case float64:
                    n = int(v)

                case int:
                    n = v

                default:
                    return fmt.Errorf("%w: expected number for int, got %T", ErrInvalidType, value)
            }

            if spec.Min > 0 && n < spec.Min {
			    return fmt.Errorf("%w: min value %d, got %d", ErrTooSmall, spec.Min, n)
            }

            if spec.Max > 0 && n > spec.Max {
			    return fmt.Errorf("%w: max value %d, got %d", ErrTooBig, spec.Max, n)
            }

            return nil

        case "number":
            var n float64

            switch v := value.(type) {
                case float64:
                    n = v

                case int:
                    n = float64(v)

                default:
		    	    return fmt.Errorf("%w: expected number, got %T", ErrInvalidType, value)
            }

            if spec.MinVal > 0 && n < spec.MinVal {
			    return fmt.Errorf("%w: min value %.2f, got %.2f", ErrTooSmall, spec.MinVal, n)
            }

            if spec.MaxVal > 0 && n > spec.MaxVal {
			    return fmt.Errorf("%w: max value %.2f, got %.2f", ErrTooBig, spec.MaxVal, n)
            }

            return nil

        case "enum":
            s, ok := value.(string)

            if !ok {
			    return fmt.Errorf("%w: expected string for enum, got %T", ErrInvalidType, value)
            }

            s = strings.TrimSpace(strings.ToUpper(s))

            for _, opt := range spec.Options {
                if s == opt {
                    return nil
                }
            }

		    return fmt.Errorf("%w: %s (allowed: %s)", ErrInvalidEnum, s, strings.Join(spec.Options, ", "))

        default:
		    return fmt.Errorf("%w: %q for field %s", ErrUnsupportedType, spec.Type, field)
	}
}

func Normalize(value any, spec FieldSpec) any {
	switch spec.Type {
        case "string":
		    return strings.TrimSpace(value.(string))

        case "email":
		    return strings.TrimSpace(strings.ToLower(value.(string)))

        case "uuid":
            return strings.TrimSpace(strings.ToLower(value.(string)))

        case "date":
		    return strings.TrimSpace(value.(string))

        case "int":
            switch v := value.(type) {
                case float64:
                    return int(v)

                case int:
                    return v

                default:
                    return value
            }

        case "number":
            switch v := value.(type) {
                case float64:
                    return v

                case int:
                    return float64(v)

                default:
                    return value
            }

        case "enum":
		    return strings.TrimSpace(strings.ToUpper(value.(string)))

        default:
            return value
	}
}