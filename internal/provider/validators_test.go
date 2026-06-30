package provider

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	resourceschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestIsEmailAddress(t *testing.T) {
	t.Parallel()

	valid := []string{"admin@example.com", "john.doe+tag@sub.example.co.uk"}
	for _, v := range valid {
		if !isEmailAddress(v) {
			t.Errorf("isEmailAddress(%q) = false, want true", v)
		}
	}

	invalid := []string{"", "admin", "admin@", "@example.com", "admin@localhost", "a b@example.com", "admin@@example.com"}
	for _, v := range invalid {
		if isEmailAddress(v) {
			t.Errorf("isEmailAddress(%q) = true, want false", v)
		}
	}
}

func TestIsDomainName(t *testing.T) {
	t.Parallel()

	valid := []string{"example.com", "sub.example.co.uk", "a-b.example.org"}
	for _, v := range valid {
		if !isDomainName(v) {
			t.Errorf("isDomainName(%q) = false, want true", v)
		}
	}

	invalid := []string{"", "example", "exa mple.com", "example_.com", "-example.com", "example.com/", "user@example.com"}
	for _, v := range invalid {
		if isDomainName(v) {
			t.Errorf("isDomainName(%q) = true, want false", v)
		}
	}
}

func TestIsIPOrCIDR(t *testing.T) {
	t.Parallel()

	valid := []string{"192.0.2.1", "203.0.113.0/24", "2001:db8::1", "2001:db8::/32"}
	for _, v := range valid {
		if !isIPOrCIDR(v) {
			t.Errorf("isIPOrCIDR(%q) = false, want true", v)
		}
	}

	invalid := []string{"", "999.0.0.1", "192.0.2.1/40", "not-an-ip", "192.0.2.0/24/8"}
	for _, v := range invalid {
		if isIPOrCIDR(v) {
			t.Errorf("isIPOrCIDR(%q) = true, want false", v)
		}
	}
}

func TestStringValidatorRejectsBadValue(t *testing.T) {
	t.Parallel()

	v := emailValidator()

	var bad validator.StringResponse
	v.ValidateString(context.Background(), validator.StringRequest{
		ConfigValue: types.StringValue("not-an-email"),
	}, &bad)
	if !bad.Diagnostics.HasError() {
		t.Fatal("expected error for invalid email")
	}

	var good validator.StringResponse
	v.ValidateString(context.Background(), validator.StringRequest{
		ConfigValue: types.StringValue("admin@example.com"),
	}, &good)
	if good.Diagnostics.HasError() {
		t.Fatalf("unexpected error for valid email: %v", good.Diagnostics)
	}

	// null/unknown values must be skipped (other layers handle requiredness).
	var null validator.StringResponse
	v.ValidateString(context.Background(), validator.StringRequest{
		ConfigValue: types.StringNull(),
	}, &null)
	if null.Diagnostics.HasError() {
		t.Fatalf("null value must not error: %v", null.Diagnostics)
	}
}

func TestSetValidatorChecksEachElement(t *testing.T) {
	t.Parallel()

	v := stringSetValidator("email address", isEmailAddress)

	set, diags := types.SetValue(types.StringType, []attr.Value{
		types.StringValue("good@example.com"),
		types.StringValue("bad-one"),
	})
	if diags.HasError() {
		t.Fatalf("set build failed: %v", diags)
	}

	var resp validator.SetResponse
	v.ValidateSet(context.Background(), validator.SetRequest{ConfigValue: set}, &resp)
	if !resp.Diagnostics.HasError() {
		t.Fatal("expected error for set containing an invalid element")
	}
}

func TestInt64BetweenValidator(t *testing.T) {
	t.Parallel()

	v := int64Between(0, 100)

	for _, in := range []int64{0, 50, 100} {
		var resp validator.Int64Response
		v.ValidateInt64(context.Background(), validator.Int64Request{ConfigValue: types.Int64Value(in)}, &resp)
		if resp.Diagnostics.HasError() {
			t.Errorf("value %d should be valid: %v", in, resp.Diagnostics)
		}
	}

	for _, in := range []int64{-1, 101} {
		var resp validator.Int64Response
		v.ValidateInt64(context.Background(), validator.Int64Request{ConfigValue: types.Int64Value(in)}, &resp)
		if !resp.Diagnostics.HasError() {
			t.Errorf("value %d should be rejected", in)
		}
	}
}

func TestSchemaAttributesCarryValidators(t *testing.T) {
	t.Parallel()

	cases := []struct {
		resource resource.Resource
		attr     string
	}{
		{NewDomainResource(), "name"},
		{NewDomainResource(), "alternatives"},
		{NewUserResource(), "email"},
		{NewUserResource(), "forward_destination"},
		{NewUserResource(), "spam_threshold"},
		{NewAliasResource(), "email"},
		{NewAliasResource(), "destination"},
		{NewTokenResource(), "email"},
		{NewTokenResource(), "authorized_ips"},
	}

	for _, tc := range cases {
		schema := resourceSchema(t, tc.resource)
		if validatorCount(schema.Attributes[tc.attr]) == 0 {
			t.Errorf("attribute %q is missing a validator", tc.attr)
		}
	}
}

func validatorCount(attr resourceschema.Attribute) int {
	switch a := attr.(type) {
	case resourceschema.StringAttribute:
		return len(a.Validators)
	case resourceschema.Int64Attribute:
		return len(a.Validators)
	case resourceschema.SetAttribute:
		return len(a.Validators)
	default:
		return 0
	}
}
