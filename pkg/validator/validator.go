package validator

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"
)

func ValidateWebsite(website string) error {
	if website == "" {
		return fmt.Errorf("website name cannot be empty")
	}

	website = strings.TrimPrefix(website, "http://")
	website = strings.TrimPrefix(website, "https://")
	website = strings.TrimPrefix(website, "www.")

	if idx := strings.Index(website, "/"); idx >= 0 {
		website = website[:idx]
	}

	if err := validateDomainFormat(website); err != nil {
		return err
	}

	return nil
}

func validateDomainFormat(domain string) error {
	if len(domain) > 253 {
		return fmt.Errorf("domain name too long (max 253 characters)")
	}

	if len(domain) < 1 {
		return fmt.Errorf("domain name too short")
	}

	domainRegex := regexp.MustCompile(`^[a-zA-Z0-9]([a-zA-Z0-9\-]{0,61}[a-zA-Z0-9])?(\.[a-zA-Z0-9]([a-zA-Z0-9\-]{0,61}[a-zA-Z0-9])?)*$`)
	
	if !domainRegex.MatchString(domain) {
		return fmt.Errorf("invalid domain format: %s", domain)
	}

	if !strings.Contains(domain, ".") {
		return fmt.Errorf("domain must contain at least one dot: %s", domain)
	}

	testURL := "http://" + domain
	if _, err := url.Parse(testURL); err != nil {
		return fmt.Errorf("invalid domain format: %s", domain)
	}

	return nil
}

func CleanWebsite(website string) string {
	website = strings.TrimPrefix(website, "http://")
	website = strings.TrimPrefix(website, "https://")
	
	website = strings.TrimPrefix(website, "www.")
	
	if idx := strings.Index(website, "/"); idx >= 0 {
		website = website[:idx]
	}
	
	website = strings.ToLower(website)
	
	return strings.TrimSpace(website)
}
