package apiserver

import (
	"testing"
	"time"

	"github.com/tebeka/selenium"
)

func TestMainPage(t *testing.T) {
	caps := selenium.Capabilities{"browserName": "chrome"}
	wd, err := selenium.NewRemote(caps, "http://localhost:4444/wd/hub")
	if err != nil {
		t.Fatalf("Failed to connect to Selenium server: %v", err)
	}
	defer wd.Quit()
	url := "http://localhost:8080"
	if err := wd.Get(url); err != nil {
		t.Fatalf("Failed to get page, get: %v", err)
	}
	// Rest of the test.../.
}

func TestSupport(t *testing.T) {
	caps := selenium.Capabilities{"browserName": "chrome"}
	wd, err := selenium.NewRemote(caps, "http://localhost:4444/wd/hub")
	if err != nil {
		t.Fatalf("Failed to connect to Selenium server: %v", err)
	}
	defer wd.Quit()

	if err := wd.Get("http://localhost:8080/support"); err != nil {
		t.Fatalf("Failed to open site, err: %v", err)
	}

	time.Sleep(2 * time.Second)

	problemTextField, err := wd.FindElement(selenium.ByCSSSelector, "#problemText")
	if err != nil {
		t.Fatalf("Error, cant put text: %v", err)
	}
	problemTextField.SendKeys("selenium test")

	submitButton, err := wd.FindElement(selenium.ByCSSSelector, "button[type='submit']")
	if err != nil {
		t.Fatalf("cannot send mail, err: %v", err)
	}

	time.Sleep(2 * time.Second)

	_, err = wd.FindElement(selenium.ByCSSSelector, "#result")
	if err != nil {
		t.Fatalf("Cannot find result, err: %v", err)
	}

	submitButton.Click()
}
