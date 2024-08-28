package utils

import (
	"errors"

	"github.com/google/uuid"
)

func ConvertStringToUUID(s string) (uuid.UUID, error) {
	u, err := uuid.Parse(s)
	if err != nil {
		return uuid.Nil, errors.New("invalid UUID string")
	}
	return u, nil
}
