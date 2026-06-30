package provider

import (
	"context"
	"fmt"
	"net"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var domainNamePattern = regexp.MustCompile(`^[a-z0-9]([a-z0-9-]*[a-z0-9])?(\.[a-z0-9]([a-z0-9-]*[a-z0-9])?)+$`)

// isDomainName reports whether value is a structurally valid fully-qualified
// domain name. Validation is intentionally lenient (structural, case
// insensitive) so the provider rejects obvious mistakes at plan time without
// second-guessing what the Mailu API accepts.
func isDomainName(value string) bool {
	value = strings.ToLower(strings.TrimSpace(value))
	if value == "" || len(value) > 253 {
		return false
	}
	return domainNamePattern.MatchString(value)
}

// isEmailAddress reports whether value looks like local@domain with a valid
// domain part. It does not attempt full RFC 5322 validation.
func isEmailAddress(value string) bool {
	value = strings.TrimSpace(value)
	at := strings.IndexByte(value, '@')
	if at <= 0 {
		return false
	}
	local := value[:at]
	domain := value[at+1:]
	if local == "" || strings.ContainsAny(local, " \t\r\n@") {
		return false
	}
	return isDomainName(domain)
}

// isIPOrCIDR reports whether value is a single IP address or a CIDR range.
func isIPOrCIDR(value string) bool {
	value = strings.TrimSpace(value)
	if value == "" {
		return false
	}
	if net.ParseIP(value) != nil {
		return true
	}
	_, _, err := net.ParseCIDR(value)
	return err == nil
}

type stringPredicateValidator struct {
	label string
	check func(string) bool
}

func (v stringPredicateValidator) Description(_ context.Context) string {
	return "value must be a valid " + v.label
}

func (v stringPredicateValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v stringPredicateValidator) ValidateString(_ context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}
	value := req.ConfigValue.ValueString()
	if !v.check(value) {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Invalid "+v.label,
			fmt.Sprintf("%q is not a valid %s.", value, v.label),
		)
	}
}

func emailValidator() validator.String {
	return stringPredicateValidator{label: "email address", check: isEmailAddress}
}

func domainValidator() validator.String {
	return stringPredicateValidator{label: "domain name", check: isDomainName}
}

type setPredicateValidator struct {
	label string
	check func(string) bool
}

func (v setPredicateValidator) Description(_ context.Context) string {
	return "every element must be a valid " + v.label
}

func (v setPredicateValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v setPredicateValidator) ValidateSet(_ context.Context, req validator.SetRequest, resp *validator.SetResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}
	for _, element := range req.ConfigValue.Elements() {
		str, ok := element.(types.String)
		if !ok || str.IsNull() || str.IsUnknown() {
			continue
		}
		value := str.ValueString()
		if !v.check(value) {
			resp.Diagnostics.AddAttributeError(
				req.Path,
				"Invalid "+v.label,
				fmt.Sprintf("%q is not a valid %s.", value, v.label),
			)
		}
	}
}

func stringSetValidator(label string, check func(string) bool) validator.Set {
	return setPredicateValidator{label: label, check: check}
}

type int64BetweenValidator struct {
	min int64
	max int64
}

func (v int64BetweenValidator) Description(_ context.Context) string {
	return fmt.Sprintf("value must be between %d and %d inclusive", v.min, v.max)
}

func (v int64BetweenValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v int64BetweenValidator) ValidateInt64(_ context.Context, req validator.Int64Request, resp *validator.Int64Response) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}
	value := req.ConfigValue.ValueInt64()
	if value < v.min || value > v.max {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Value Out Of Range",
			fmt.Sprintf("must be between %d and %d inclusive, got %d.", v.min, v.max, value),
		)
	}
}

func int64Between(min, max int64) validator.Int64 {
	return int64BetweenValidator{min: min, max: max}
}
