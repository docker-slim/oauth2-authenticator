package browser

import (
	"time"
	"os/exec"
)

var Names = map[string]string{
	Chrome: "google-chrome",
	Firefox: "firefox",
	Brave: "brave-browser",
	Chromium: "chromium-browser",
}

func IsInstalled(browser string) bool {
	name, found := Names[browser]
	if !found {
		return false
	}

	_, err := exec.LookPath(name)
	if err == nil {
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
	
	cmdName := "xdg-open"
	if browser != Default {
		cmdName = Names[browser]
	}

	cmd := exec.Command(cmdName, url)
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
