package helpers_test

import (
	"errors"
	"testing"

	"github.com/ShukinDmitriy/GophKeeper/internal/helpers"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
)

func TestExtractErrors(t *testing.T) {
	validate := validator.New(validator.WithRequiredStructEnabled())

	type testStruct struct {
		Field        string `validate:"required"`
		AnotherField string `validate:"required"`
	}

	ts1 := testStruct{
		Field:        "",
		AnotherField: "value",
	}

	ts2 := testStruct{
		Field:        "",
		AnotherField: "",
	}

	type args struct {
		err error
	}
	tests := []struct {
		name string
		args args
		want helpers.ValidationError
	}{
		{
			name: "nothing to extract",
			args: args{
				err: errors.New("some error"),
			},
			want: helpers.ValidationError{},
		},
		{
			name: "one validation error",
			args: args{
				err: validate.Struct(&ts1),
			},
			want: helpers.ValidationError{
				"field": map[string]bool{
					"required": true,
				},
			},
		},
		{
			name: "several validation error",
			args: args{
				err: validate.Struct(&ts2),
			},
			want: helpers.ValidationError{
				"field": map[string]bool{
					"required": true,
				},
				"anotherfield": map[string]bool{
					"required": true,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, helpers.ExtractErrors(tt.args.err), "ExtractErrors(%v)", tt.args.err)
		})
	}
}
