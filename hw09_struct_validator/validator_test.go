package hw09structvalidator

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

type (
	UserRole string
	Team     string
)

// Test the function on different structures and other types.
type (
	User struct {
		ID     string `json:"id" validate:"len:36"`
		Name   string
		Age    int      `validate:"min:18|max:50"`
		Email  string   `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role   UserRole `validate:"in:admin,stuff"`
		Member []Team   `validate:"in:otus,carte-noire,parlament,krk"`
		Phones []string `validate:"len:11"`
		meta   json.RawMessage
	}

	App struct {
		Version string `validate:"len:5"`
	}

	InvalidValidationRule struct {
		SomeValue int `validate:"len:10"`
	}

	IntMinMax struct {
		Value int `validate:"min:10|max:100"`
	}

	UnsupportedType struct {
		Field bool `validate:"min:10"`
	}

	Token struct {
		Header    []byte
		Payload   []byte
		Signature []byte
	}

	Response struct {
		Code int    `validate:"in:200,404,500"`
		Body string `json:"omitempty"`
	}

	Regexp struct {
		Name string `validate:"regexp:^\\w+@\\d+$"`
	}
)

func TestValidate(t *testing.T) {
	tests := []struct {
		in          any
		expectedErr error
	}{
		{UnsupportedType{true}, ErrUnsupportedType},
		{InvalidValidationRule{100}, ErrInvalidValidationRule},
		{[]string{"10", "20"}, ErrInvalidStructType},

		{App{"1"}, ErrValidationLength},
		{App{"123456"}, ErrValidationLength},
		{App{"12345"}, nil},
		{IntMinMax{9}, ErrValidationMinValue},
		{IntMinMax{101}, ErrValidationMaxValue},
		{IntMinMax{10}, nil},
		{IntMinMax{100}, nil},
		{Regexp{"hhh@555"}, nil},
		{Regexp{"555@hhh"}, ErrValidationRegexpValue},
		{Response{200, ""}, nil},
		{Response{201, ""}, ErrValidationInValue},
		{Token{}, nil},
		{User{
			strings.Repeat("1", 36), "pavel", 45,
			"pavel@bazhov.ru", "stuff",
			[]Team{"otus", "krk"},
			[]string{"66555-55-55"},
			json.RawMessage{},
		}, nil},
		{User{
			strings.Repeat("1", 36), "pavel", 45,
			"pavel@bazhov.ru", "stuff",
			[]Team{"otus", "donor"},
			[]string{"66555-55-55"},
			json.RawMessage{},
		}, ErrValidationInValue},
		{struct{ r []Response }{[]Response{{100, ""}}}, nil}, // do not support nested structures
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()

			err := Validate(tt.in)
			require.ErrorIs(t, err, tt.expectedErr)
		})
	}
}
