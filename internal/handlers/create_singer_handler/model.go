package create_singer_handler

import (
	"github.com/go-playground/validator/v10"
)

type CreateSingersRequest struct {
	SingerId string `json:"singer_id" validate:"required"`
	Name     string `json:"name" validate:"required"`
	LastName string `json:"last_name" validate:"required"`
}

func (req *CreateSingersRequest) Validate(validate *validator.Validate) error {
	return validate.Struct(req)
}
