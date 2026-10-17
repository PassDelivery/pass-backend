package core_http_request

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	core_errors "github.com/PassDelivery/pass-backend/internal/core/errors"
	"github.com/go-playground/validator/v10"
)

var requestValidator = validator.New()

type validatble interface {
	Validate() error
}

func DecodeAndValidateRequest(r *http.Request, dest any) error {
	if err := json.NewDecoder(r.Body).Decode(dest); err != nil {
		return fmt.Errorf("decode json: %v: %w", err, core_errors.ErrInvalidArgument)
	}

	var err error

	v, ok := dest.(validatble)
	if ok {
		err = v.Validate()
	}else {
		err = requestValidator.Struct(dest)
	}
	if err != nil && !errors.Is(err, core_errors.ErrInvalidArgument) {
		return fmt.Errorf("request validation: %v: %w", err, core_errors.ErrInvalidArgument)
	}
	if err != nil {
		return fmt.Errorf("request validation: %w", err)
	}

	return nil
}