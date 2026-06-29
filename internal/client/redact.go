package client

import (
	"regexp"
	"strings"
)

var redactionPatterns = []*regexp.Regexp{
	regexp.MustCompile(`(?i)(authorization:\s*bearer\s+)[^\s]+`),
	regexp.MustCompile(`(?i)(authorization:\s*)[^\s]+`),
	regexp.MustCompile(`(?i)("?(?:token|raw_password|password|reply_body|smtp)"?\s*[:=]\s*")([^"]+)(")`),
	regexp.MustCompile(`(?i)((?:token|raw_password|password|reply_body|smtp)\s*=\s*)[^\s]+`),
	regexp.MustCompile(`(?i)(smtp(?:s)?://[^:/@\s]+:)[^@\s]+(@)`),
	regexp.MustCompile(`\$bcrypt-sha256\$[^\s",]+`),
}

func Redact(value string) string {
	redacted := value
	for _, pattern := range redactionPatterns {
		redacted = pattern.ReplaceAllString(redacted, `${1}<redacted>${3}`)
	}

	return strings.TrimSpace(redacted)
}
