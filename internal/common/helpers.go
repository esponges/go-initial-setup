package common

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/go-playground/validator/v10"
)

type Validator interface {
	Validate(validate *validator.Validate) error
}

func UnmarshalAndValidateRequest(request *http.Request, reqContract Validator, validator *validator.Validate) ([]byte, error) {
	body, err := io.ReadAll(request.Body)
	if err != nil {
		return nil, err
	}

	log.Println("request body: " + string(body))

	if err := json.Unmarshal(body, &reqContract); err != nil {
		return nil, err
	}

	log.Println("unmarshaled request: " + string(body))

	if err := reqContract.Validate(validator); err != nil {
		return nil, err
	}

	return body, nil
}
