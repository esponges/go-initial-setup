package sample_post_request

import (
	"testing"

	"github.com/go-playground/validator/v10"
)

func TestSamplePostRequestValidate(t *testing.T) {
	tests := []struct {
		name string
		req  SamplePostRequest
		want bool
	}{
		{
			name: "valid request",
			req: SamplePostRequest{
				Name: "John Doe",
			},
			want: true,
		},
		{
			name: "missing name",
			req: SamplePostRequest{
				Name: "",
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			validate := validator.New()
			err := tt.req.Validate(validate)
			if tt.want {
				if err != nil {
					t.Errorf("Validate() error = %v, want nil", err)
				}
			} else {
				if err == nil {
					t.Errorf("Validate() error = nil, want non-nil")
				} else if _, ok := err.(validator.ValidationErrors); !ok {
					t.Errorf("Validate() error = %v, want validator.ValidationErrors", err)
				}
			}
		})
	}
}
