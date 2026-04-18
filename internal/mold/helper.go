package mold

import (
	"errors"

	"github.com/lib/pq"
)

// helper

func isUniqueViolation(err error) bool {
	var pqErr *pq.Error
	if errors.As(err, &pqErr) {
		return pqErr.Code == "23505"
	}

	return false
}
