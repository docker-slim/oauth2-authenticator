package browser

import (
	"os"
	"fmt"
	"time"
	"os/exec"
)

var Names = map[string]string{
	Chrome: "google chrome",
	Firefox: "firefox",
	Brave: "brave browser",
	Chromium: "chromium",
	Safari: "safari",
}

var locations = map[string]string{
	Chrome: "/Applications/Google Chrome.app",
	Firefox: "/Applications/Firefox.app",
	Brave: "/Applications/Brave Browser.app",
	Chromium: "/Applications/Chromium.app",
	Safari: "/Applications/Safari.app",
}

func IsInstalled(browser string) bool {
	location, found := locations[browser]
	if !found {
		return false
	}

	if info, err := os.Stat(location); err == nil && info.IsDir() {
		return true
	}

	return false
}

func OpenURL(url, browser string, wait int) error {
	if url == "" {
		return ErrBadURL
	}

	if browser != Default && !IsInstalled(browser) {
		return ErrBrowserNotInstalled
	}

	var args []string
	if browser != Default { 
		args = append(args, "-a")
		arg := fmt.Sprintf("%s", locations[browser])
		args = append(args, arg)
	}

	args = append(args, fmt.Sprintf("%s",url))

	cmd := exec.Command("open", args...)
	err := cmd.Start()
	if err != nil {
		return err
	}

	if wait > 0 {
		errorCh := make(chan error, 1)
		go func() {
			errorCh <- cmd.Wait()
		}()

		select {
		case err := <- errorCh:
			return err
		case <- time.After(time.Duration(wait) * time.Second):
			return nil
		}
	}

	return nil
}
