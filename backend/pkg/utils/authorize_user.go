package utils

import "errors"

type ContextKey string

func AuthorizeUser(role string, allowedRoles ...string) (bool, error) {
	for _, allowedRole := range allowedRoles {
		if allowedRole == role {
			return true, nil
		}
	}

	return false, errors.New("user not authorized")
}
