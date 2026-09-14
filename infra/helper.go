package infra

import (
	"encoding/json"
	"io"
	"net/http"
)

func DecodeAndValidateDto(r io.Reader, dto Validatable) error {
	if err := json.NewDecoder(r).Decode(dto); err != nil {
		return err
	}

	if err := dto.Validate(); err != nil {
		return err
	}
	return nil
}

func CreateResponseWithError(w http.ResponseWriter, errStr string, code int) {
	CreateResponseWithJsonBody(w, NewErrorDto(errStr), code)
}

func CreateResponseWithJsonBody(w http.ResponseWriter, jsonBody any, code int) {
	b, err := json.MarshalIndent(jsonBody, "", "    ")
	if err != nil {
		panic(err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(b)
}
