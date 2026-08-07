package utils

import (
	"fmt"

	"github.com/nyaruka/phonenumbers"
)

func IsValidPhoneNumber(phone string) bool {
	num, err := phonenumbers.Parse(phone, "")
	if err != nil {
		return false
	}

	return phonenumbers.IsValidNumber(num)
}

func GetPhoneCountry(phone string) (string, error) {
	num, err := phonenumbers.Parse(phone, "")
	if err != nil {
		return "", err
	}

	if !phonenumbers.IsValidNumber(num) {
		return "", fmt.Errorf("invalid phone number")
	}

	return phonenumbers.GetRegionCodeForNumber(num), nil
}
