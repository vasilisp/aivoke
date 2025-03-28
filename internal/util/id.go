package util

import (
	"fmt"
	"regexp"
)

var regex = regexp.MustCompile(`^[a-zA-Z0-9]+(?:-[a-zA-Z0-9]+)*$`)

func ValidateID(id string) error {
	if !regex.MatchString(id) {
		return fmt.Errorf("invalid ID: %v", id)
	}
	return nil
}
