package main

import (
	"fmt"
	"log"
	"net/mail"
	"net/url"
	"strings"

	"github.com/AlecAivazis/survey/v2/terminal"
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

func shouldBeSPSite(val any) error {
	str, ok := val.(string)
	if !ok {
		return fmt.Errorf("value is not a string")
	}
	if strings.Contains(str, ".aspx") {
		return fmt.Errorf("value should not be a page, but a site URL")
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

func shouldBeGUIDorEmpty(val any) error {
	str, _ := val.(string)
	if str == "" {
		return nil
	}
	return shouldBeGUID(val)
}

func interuptable(err error) {
	if err == terminal.InterruptErr {
		log.Fatal("interrupted")
	}
	if err != nil {
		fmt.Println(err)
	}
}
