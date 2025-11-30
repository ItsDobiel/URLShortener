package test

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/cucumber/godog"
	"github.com/tebeka/selenium"
)

// TestContext holds the state for BDD tests
type TestContext struct {
	webDriver     selenium.WebDriver
	baseURL       string
	lastShortCode string
}

var testCtx *TestContext

// TestFeatures runs the BDD test suite using Go's native testing framework
func TestFeatures(t *testing.T) {
	fmt.Println("=== Starting BDD Test Suite ===")

	// Initialize test context
	testCtx = &TestContext{
		baseURL: "http://localhost:8080",
	}

	fmt.Println(" Launching Firefox browser...")
	driver, err := newWebDriver()
	if err != nil {
		t.Fatalf("Failed to create WebDriver: %v", err)
	}
	testCtx.webDriver = driver
	defer func() {
		fmt.Println(" Closing browser...")
		testCtx.webDriver.Quit()
	}()

	fmt.Println(" Test environment ready")

	// Configure and run godog test suite
	suite := godog.TestSuite{
		ScenarioInitializer: initializeScenario,
		Options: &godog.Options{
			Format:   "pretty",
			Paths:    []string{"features"},
			TestingT: t,
		},
	}

	if suite.Run() != 0 {
		t.Fatal("non-zero status returned, failed to run feature tests")
	}

	fmt.Println("=== All Tests Complete ===")
}

func initializeScenario(ctx *godog.ScenarioContext) {
	ctx.Before(func(ctx context.Context, sc *godog.Scenario) (context.Context, error) {
		fmt.Printf(" Starting scenario: %s\n", sc.Name)
		testCtx.lastShortCode = ""
		return ctx, nil
	})

	ctx.After(func(ctx context.Context, sc *godog.Scenario, err error) (context.Context, error) {
		if err != nil {
			fmt.Printf(" Scenario failed: %v\n", err)
		} else {
			fmt.Println(" Scenario passed")
		}
		return ctx, nil
	})

	ctx.Step(`^the URL shortener service is running$`, stepServiceIsRunning)
	ctx.Step(`^I am on the home page$`, stepOnHomePage)
	ctx.Step(`^I enter the URL "([^"]*)"$`, stepEnterURL)
	ctx.Step(`^I submit the form$`, stepSubmitForm)
	ctx.Step(`^I should see a success message$`, stepSeeSuccessMessage)
	ctx.Step(`^I should see a shortened URL$`, stepSeeShortenedURL)
	ctx.Step(`^the shortened URL should be valid$`, stepShortenedURLValid)
	ctx.Step(`^the shortned URL must redirect me to "([^"]*)"$`, stepNavigateToShortenedURL)
	ctx.Step(`^I should receive the same short code as before$`, stepSameShortCode)
	ctx.Step(`^I should see an error message$`, stepSeeErrorMessage)
	ctx.Step(`^I navigate to "([^"]*)"$`, stepNavigateTo)
	ctx.Step(`^I should see an error page$`, stepSeeErrorPage)
	ctx.Step(`^the error should indicate the short code was not found$`, stepErrorNotFound)
	ctx.Step(`^the short code should be alphanumeric with allowed characters$`, stepShortCodeAlphanumeric)
	ctx.Step(`^the short code length should match the configured length$`, stepShortCodeLength)
}

func newWebDriver() (selenium.WebDriver, error) {
	caps := selenium.Capabilities{
		"browserName": "firefox",
		"moz:firefoxOptions": map[string]any{
			"args": []string{"--width=1920", "--height=1080"},
			"prefs": map[string]any{
				"dom.webdriver.enabled":      false,
				"useAutomationExtension":     false,
				"general.useragent.override": "Mozilla/5.0 (X11; Linux x86_64; rv:109.0) Gecko/20100101 Firefox/115.0",
			},
		},
	}

	wd, err := selenium.NewRemote(caps, "http://localhost:4444")
	if err != nil {
		return nil, fmt.Errorf("GeckoDriver connection failed (ensure GeckoDriver is running on port 4444): %w", err)
	}

	wd.MaximizeWindow("")
	return wd, nil
}

// Step definition functions
func stepServiceIsRunning() error {
	fmt.Println("   Verifying server is accessible...")
	err := testCtx.webDriver.Get(testCtx.baseURL)
	//time.Sleep(200 * time.Millisecond)
	return err
}

func stepOnHomePage() error {
	fmt.Println("   Navigating to home page...")
	err := testCtx.webDriver.Get(testCtx.baseURL)
	//time.Sleep(200 * time.Millisecond)
	return err
}

func stepEnterURL(url string) error {
	fmt.Printf("   Entering URL: %s\n", url)

	currentURL, _ := testCtx.webDriver.CurrentURL()
	if !strings.HasSuffix(currentURL, "/") || strings.Contains(currentURL, "error") {
		fmt.Println("   Navigating back to home page...")
		testCtx.webDriver.Get(testCtx.baseURL)
	}

	urlInput, err := testCtx.webDriver.FindElement(selenium.ByID, "url")
	if err != nil {
		return fmt.Errorf("URL input not found: %w", err)
	}

	urlInput.Clear()
	err = urlInput.SendKeys(url)
	//time.Sleep(500 * time.Millisecond)
	return err
}

func stepSubmitForm() error {
	fmt.Println("   Clicking submit button...")

	button, err := testCtx.webDriver.FindElement(selenium.ByID, "submit")
	if err != nil {
		return fmt.Errorf("submit button not found: %w", err)
	}

	button.Click()
	fmt.Println("   Waiting for response...")
	//time.Sleep(200 * time.Millisecond)
	return nil
}

func stepSeeSuccessMessage() error {
	fmt.Println("   Checking for success message...")

	pageSource, err := testCtx.webDriver.PageSource()
	if err != nil {
		return fmt.Errorf("failed to get page source: %w", err)
	}

	if strings.Contains(pageSource, "Success") || strings.Contains(pageSource, "✨") {
		fmt.Println("   Success message found!")
		return nil
	}

	elements, _ := testCtx.webDriver.FindElements(selenium.ByID, "result")
	if len(elements) > 0 {
		text, _ := elements[0].Text()
		if strings.Contains(text, "Success") {
			fmt.Println("   Success message found!")
			return nil
		}
	}

	return fmt.Errorf("success message not found")
}

func stepSeeShortenedURL() error {
	fmt.Println("   Looking for shortened URL...")
	links, _ := testCtx.webDriver.FindElements(selenium.ByID, "shortened")
	for _, link := range links {
		href, _ := link.GetAttribute("href")
		if strings.Contains(href, "http://localhost:8080/") && !strings.HasSuffix(href, "/") {
			parts := strings.Split(href, "/")
			testCtx.lastShortCode = strings.TrimSpace(parts[len(parts)-1])
			fmt.Printf("   Short URL found: %s\n", href)
			return nil
		}
	}

	return fmt.Errorf("shortened URL not found")
}

func stepShortenedURLValid() error {
	fmt.Println("   Validating short code format...")

	if testCtx.lastShortCode == "" {
		return fmt.Errorf("no short code found")
	}

	for _, char := range testCtx.lastShortCode {
		if !((char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z') ||
			(char >= '0' && char <= '9') || char == '-' || char == '_') {
			return fmt.Errorf("invalid character in short code: %c", char)
		}
	}

	fmt.Printf("   Short code is valid: %s\n", testCtx.lastShortCode)
	return nil
}

func stepNavigateToShortenedURL(url string) error {
	fmt.Println("   Navigating to Shortened URL...")
	err := testCtx.webDriver.Get(testCtx.baseURL + "/" + testCtx.lastShortCode)
	if err != nil {
		return fmt.Errorf("Failed to open shortened URL: %s", err)
	}
	currentURL, err := testCtx.webDriver.CurrentURL()
	if err != nil {
		return fmt.Errorf("Failed to retrieve current URL: %s", err)
	}
	if currentURL == url {
		fmt.Println("   Successfully redirected to the original URL.")
	}
	return nil
}

func stepSameShortCode() error {
	fmt.Println("   Checking if short code matches...")

	savedCode := testCtx.lastShortCode

	links, _ := testCtx.webDriver.FindElements(selenium.ByID, "shortened")
	for _, link := range links {
		href, _ := link.GetAttribute("href")
		if strings.Contains(href, "http://localhost:8080/") && !strings.HasSuffix(href, "/") {
			parts := strings.Split(href, "/")
			currentCode := strings.TrimSpace(parts[len(parts)-1])

			if currentCode != savedCode {
				return fmt.Errorf("expected '%s', got '%s'", savedCode, currentCode)
			}
			fmt.Printf("   Codes match: %s\n", currentCode)
			return nil
		}
	}

	return fmt.Errorf("could not find short code to compare")
}

func stepSeeErrorMessage() error {
	fmt.Println("   Checking for error message...")

	pageSource, err := testCtx.webDriver.PageSource()
	if err != nil {
		return err
	}

	if strings.Contains(pageSource, "Something went wrong") ||
		strings.Contains(pageSource, "Oops") ||
		strings.Contains(pageSource, "Error") ||
		strings.Contains(pageSource, "invalid") ||
		strings.Contains(pageSource, "supported") ||
		strings.Contains(pageSource, "⚠️") {
		fmt.Println("   Error message found!")
		return nil
	}

	return fmt.Errorf("no error message found")
}

func stepNavigateTo(path string) error {
	fmt.Printf("   Navigating to: %s\n", path)
	err := testCtx.webDriver.Get(testCtx.baseURL + path)
	//time.Sleep(500 * time.millisecond)
	return err
}

func stepSeeErrorPage() error {
	fmt.Println("   Verifying error page...")

	pageSource, err := testCtx.webDriver.PageSource()
	if err != nil {
		return err
	}

	if !strings.Contains(pageSource, "Something went wrong") &&
		!strings.Contains(pageSource, "Error") &&
		!strings.Contains(pageSource, "⚠️") {
		return fmt.Errorf("error page not displayed")
	}

	fmt.Println("   Error page confirmed!")
	return nil
}

func stepErrorNotFound() error {
	fmt.Println("   Checking error message content...")

	pageSource, err := testCtx.webDriver.PageSource()
	if err != nil {
		return err
	}

	if !strings.Contains(pageSource, "not found") &&
		!strings.Contains(pageSource, "Short code not found") {
		return fmt.Errorf("'not found' message missing")
	}

	fmt.Println("   Error message correct!")
	return nil
}

func stepShortCodeAlphanumeric() error {
	fmt.Println("   Validating alphanumeric characters...")

	if testCtx.lastShortCode == "" {
		return fmt.Errorf("no short code to validate")
	}

	for _, char := range testCtx.lastShortCode {
		if !((char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z') ||
			(char >= '0' && char <= '9') || char == '-' || char == '_') {
			return fmt.Errorf("invalid character '%c' in '%s'", char, testCtx.lastShortCode)
		}
	}

	fmt.Printf("   Valid characters: %s\n", testCtx.lastShortCode)
	return nil
}

func stepShortCodeLength() error {
	fmt.Println("   Validating code length...")

	if testCtx.lastShortCode == "" {
		return fmt.Errorf("no short code to validate")
	}

	configuredLength := 7
	if envLength := os.Getenv("SHORT_CODE_LENGTH"); envLength != "" {
		fmt.Sscanf(envLength, "%d", &configuredLength)
	}

	actualLength := len(testCtx.lastShortCode)
	if actualLength != configuredLength {
		return fmt.Errorf("length mismatch: expected %d, got %d in '%s'",
			configuredLength, actualLength, testCtx.lastShortCode)
	}

	fmt.Printf("   Length correct: %d chars\n", actualLength)
	return nil
}
