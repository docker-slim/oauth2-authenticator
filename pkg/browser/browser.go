package browser

import (
	"errors"
)

//Browser types
const (
	Default = ""
	Chrome = "chrome"
	Firefox = "firefox"
	Brave = "brave"
	Chromium = "chromium"
	Safari = "safari"
)

//Common errors
var ErrBrowserNotInstalled = errors.New("browser not installed")
var ErrBadURL = errors.New("bad URL")

func IsValid(browser string) bool {
	if browser == Default {
		return true
	}

	if _, found := Names[browser]; found {
		return true
	}

	return false
}

