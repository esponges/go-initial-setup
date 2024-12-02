package sample_post_request

import "github.com/go-playground/validator/v10"

type SamplePostRequest struct {
	Name string `json:"name" validate:"required"`
}

type SamplePostResponse struct {
	Message string `json:"message"`
}

func (req *SamplePostRequest) Validate(validate *validator.Validate) error {
	return validate.Struct(req)
}
