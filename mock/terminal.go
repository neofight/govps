package mock

import (
	"errors"
)

type Terminal struct {
	Password string
}

func (mock Terminal) ReadPassword(fd int) ([]byte, error) {

	if mock.Password == "" {
		return []byte(""), errors.New("Error reading password")
	}

	return []byte(mock.Password), nil
}
