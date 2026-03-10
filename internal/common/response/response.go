package response

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/lib/pq"
)

func JSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Context-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func Error(w http.ResponseWriter, status int, msg string) {
	JSON(w, status, map[string]string{"error": msg})
}

func BadRequest(w http.ResponseWriter, msg string) {
	Error(w, http.StatusBadRequest, msg)
}

func HandlePostgresError(err error) error {

	var pqErr *pq.Error

	if errors.As(err, &pqErr) {

		switch pqErr.Code {

		case "23505":
			return errors.New("duplicate record")

		case "23503":
			return errors.New("foreign key violation")

		case "23502":
			return errors.New("missing required field")

		}
	}

	return err
}
