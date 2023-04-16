package main

import (
	"fmt"
	"net/mail"
	"net/url"

	"github.com/google/uuid"
)

func shouldBeURL(val any) error {
	str, ok := val.(string)
	if !ok {
		return fmt.Errorf("value is not a string")
	}

	if _, err := url.ParseRequestURI(str); err != nil {
		return fmt.Errorf("value is not a valid URL")
	}

	return nil
}

func shouldBeEmail(val any) error {
	str, ok := val.(string)
	if !ok {
		return fmt.Errorf("value is not a string")
	}

	if _, err := mail.ParseAddress(str); err != nil {
		return fmt.Errorf("value is not a valid email address")
	}

	return nil
}

func shouldBeGUID(val any) error {
	str, ok := val.(string)
	if !ok {
		return fmt.Errorf("value is not a string")
	}

	if _, err := uuid.Parse(str); err != nil {
		return fmt.Errorf("value is not a valid GUID")
	}

	return nil
}
