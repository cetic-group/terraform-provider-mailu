package provider

import (
	"fmt"
	"strings"
)

func validateDomainImportID(value string) (string, error) {
	id := normalizeDomain(value)
	if id == "" {
		return "", fmt.Errorf("import ID must not be empty")
	}
	if strings.ContainsAny(id, " \t\r\n/@") {
		return "", fmt.Errorf("import ID must be a domain name without spaces, slash, or @")
	}

	return id, nil
}

func validateEmailImportID(value string) (string, error) {
	id := normalizeEmail(value)
	if id == "" {
		return "", fmt.Errorf("import ID must not be empty")
	}
	if strings.ContainsAny(id, " \t\r\n/") {
		return "", fmt.Errorf("import ID must be an email address without spaces or slash")
	}

	local, domain, ok := strings.Cut(id, "@")
	if !ok || local == "" || domain == "" || strings.Contains(domain, "@") {
		return "", fmt.Errorf("import ID must use the format local@example.com")
	}

	return id, nil
}

func validateTokenImportID(value string) (string, error) {
	id := strings.TrimSpace(value)
	if id == "" {
		return "", fmt.Errorf("import ID must not be empty")
	}
	if strings.ContainsAny(id, " \t\r\n/") {
		return "", fmt.Errorf("import ID must not contain spaces or slash")
	}

	return id, nil
}
