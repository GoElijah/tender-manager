package app

import (
	"errors"
	"reflect"
)

func ValidateNotEmpty(reqs ...any) error {
	for _, v := range reqs {
		if reflect.ValueOf(v).IsZero() {
			return errors.New("body param is empty for")
		}
	}
	return nil

}

func ValidateUserHasAccess(userOrganization, tenderOrganization string) error {
	if userOrganization != tenderOrganization {
		return errors.New("Access denied")
	}
	return nil
}
