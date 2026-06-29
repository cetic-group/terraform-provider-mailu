package provider

import (
	"context"
	"testing"

	"github.com/cetic-group/terraform-provider-mailu/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	datasourceschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	resourceschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

func TestResourcesAndDataSourcesRegistered(t *testing.T) {
	t.Parallel()

	p := New("test")()

	resourceNames := make(map[string]bool)
	for _, factory := range p.Resources(context.Background()) {
		r := factory()
		var resp resource.MetadataResponse
		r.Metadata(context.Background(), resource.MetadataRequest{ProviderTypeName: "mailu"}, &resp)
		resourceNames[resp.TypeName] = true
	}

	for _, name := range []string{
		"mailu_domain",
		"mailu_user",
		"mailu_alias",
		"mailu_alternative_domain",
		"mailu_domain_manager",
		"mailu_relay",
		"mailu_token",
	} {
		if !resourceNames[name] {
			t.Fatalf("resource %q is not registered", name)
		}
	}

	dataSourceNames := make(map[string]bool)
	for _, factory := range p.DataSources(context.Background()) {
		ds := factory()
		var resp datasource.MetadataResponse
		ds.Metadata(context.Background(), datasource.MetadataRequest{ProviderTypeName: "mailu"}, &resp)
		dataSourceNames[resp.TypeName] = true
	}

	for _, name := range []string{"mailu_domain", "mailu_user", "mailu_dkim"} {
		if !dataSourceNames[name] {
			t.Fatalf("data source %q is not registered", name)
		}
	}
}

func TestMVPResourceSensitiveAttributes(t *testing.T) {
	t.Parallel()

	userSchema := resourceSchema(t, NewUserResource())
	rawPassword, ok := userSchema.Attributes["raw_password"].(resourceschema.StringAttribute)
	if !ok {
		t.Fatalf("raw_password has type %T, want resource string attribute", userSchema.Attributes["raw_password"])
	}
	if !rawPassword.Sensitive {
		t.Fatal("raw_password must be sensitive")
	}

	replyBody, ok := userSchema.Attributes["reply_body"].(resourceschema.StringAttribute)
	if !ok {
		t.Fatalf("reply_body has type %T, want resource string attribute", userSchema.Attributes["reply_body"])
	}
	if !replyBody.Sensitive {
		t.Fatal("reply_body must be sensitive")
	}

	userDataSourceSchema := dataSourceSchema(t, NewUserDataSource())
	dataSourceReplyBody, ok := userDataSourceSchema.Attributes["reply_body"].(datasourceschema.StringAttribute)
	if !ok {
		t.Fatalf("data source reply_body has type %T, want data source string attribute", userDataSourceSchema.Attributes["reply_body"])
	}
	if !dataSourceReplyBody.Sensitive {
		t.Fatal("data source reply_body must be sensitive")
	}

	tokenSchema := resourceSchema(t, NewTokenResource())
	generatedToken, ok := tokenSchema.Attributes["token"].(resourceschema.StringAttribute)
	if !ok {
		t.Fatalf("token has type %T, want resource string attribute", tokenSchema.Attributes["token"])
	}
	if !generatedToken.Sensitive {
		t.Fatal("generated token must be sensitive")
	}
}

func TestMVPImportIDsAreNormalized(t *testing.T) {
	t.Parallel()

	if got, want := normalizeDomain(" EXAMPLE.COM "), "example.com"; got != want {
		t.Fatalf("domain import id = %q, want %q", got, want)
	}
	if got, want := normalizeEmail(" ADMIN@EXAMPLE.COM "), "admin@example.com"; got != want {
		t.Fatalf("email import id = %q, want %q", got, want)
	}
}

func TestImportIDValidation(t *testing.T) {
	t.Parallel()

	if got, err := validateDomainImportID(" EXAMPLE.COM "); err != nil || got != "example.com" {
		t.Fatalf("domain import = %q, %v; want example.com, nil", got, err)
	}
	for _, input := range []string{"", " ", "example.com/admin", "admin@example.com"} {
		if _, err := validateDomainImportID(input); err == nil {
			t.Fatalf("domain import %q should fail", input)
		}
	}

	if got, err := validateEmailImportID(" ADMIN@EXAMPLE.COM "); err != nil || got != "admin@example.com" {
		t.Fatalf("email import = %q, %v; want admin@example.com, nil", got, err)
	}
	for _, input := range []string{"", "admin", "admin@", "@example.com", "admin/example.com"} {
		if _, err := validateEmailImportID(input); err == nil {
			t.Fatalf("email import %q should fail", input)
		}
	}

	if got, err := validateTokenImportID(" 42 "); err != nil || got != "42" {
		t.Fatalf("token import = %q, %v; want 42, nil", got, err)
	}
	for _, input := range []string{"", " ", "token/42", "token 42"} {
		if _, err := validateTokenImportID(input); err == nil {
			t.Fatalf("token import %q should fail", input)
		}
	}
}

func TestRelaySMTPRejectsCredentials(t *testing.T) {
	t.Parallel()

	var diags diag.Diagnostics
	validateRelaySMTP("smtp://user:pass@mail.example.com:587", &diags)
	if !diags.HasError() {
		t.Fatal("expected diagnostics for relay SMTP credentials")
	}

	var clean diag.Diagnostics
	validateRelaySMTP("mail.example.com:587", &clean)
	if clean.HasError() {
		t.Fatalf("unexpected diagnostics for credential-free SMTP: %v", clean)
	}
}

func TestTokenApplyAPIDoesNotPersistGeneratedToken(t *testing.T) {
	t.Parallel()

	var model tokenModel
	var diags diag.Diagnostics
	model.applyAPI(context.Background(), &client.Token{
		ID:    client.FlexibleString("42"),
		Email: "ADMIN@EXAMPLE.COM",
		Token: "generated-secret",
	}, false, &diags)
	if diags.HasError() {
		t.Fatalf("unexpected diagnostics: %v", diags)
	}
	if !model.Token.IsNull() {
		t.Fatalf("generated token should not be stored in state: %#v", model.Token)
	}
}

func TestMVPRequiredAttributesForceReplacement(t *testing.T) {
	t.Parallel()

	for name, r := range map[string]resource.Resource{
		"domain": NewDomainResource(),
		"user":   NewUserResource(),
		"alias":  NewAliasResource(),
	} {
		schema := resourceSchema(t, r)
		attrName := map[string]string{
			"domain": "name",
			"user":   "email",
			"alias":  "email",
		}[name]

		attr, ok := schema.Attributes[attrName].(resourceschema.StringAttribute)
		if !ok {
			t.Fatalf("%s.%s has type %T, want resource string attribute", name, attrName, schema.Attributes[attrName])
		}
		if len(attr.PlanModifiers) == 0 {
			t.Fatalf("%s.%s must require replacement", name, attrName)
		}
	}
}

func resourceSchema(t *testing.T, r resource.Resource) resourceschema.Schema {
	t.Helper()

	var resp resource.SchemaResponse
	r.Schema(context.Background(), resource.SchemaRequest{}, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("resource schema diagnostics: %v", resp.Diagnostics)
	}

	return resp.Schema
}

func dataSourceSchema(t *testing.T, ds datasource.DataSource) datasourceschema.Schema {
	t.Helper()

	var resp datasource.SchemaResponse
	ds.Schema(context.Background(), datasource.SchemaRequest{}, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("data source schema diagnostics: %v", resp.Diagnostics)
	}

	return resp.Schema
}

func TestClientErrorNotFoundDetection(t *testing.T) {
	t.Parallel()

	if !isNotFound(&client.APIError{StatusCode: 404}) {
		t.Fatal("404 API error should be treated as not found")
	}
	if isNotFound(&client.APIError{StatusCode: 500}) {
		t.Fatal("500 API error should not be treated as not found")
	}
}
