package session

import "fmt"

func ErrPermissionDenied(which string) error {
	return fmt.Errorf("permissioned denied: " + which)
}

func ErrInvalidParameter(which string) error {
	return fmt.Errorf("invalid parameter: " + which)
}