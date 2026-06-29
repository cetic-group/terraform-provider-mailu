package provider

import (
	"context"
	"testing"

	"github.com/cetic-group/terraform-provider-mailu/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	datasourceschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	resourceschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

func TestMVPResourcesAndDataSourcesRegistered(t *testing.T) {
	t.Parallel()

	p := New("test")()

	resourceNames := make(map[string]bool)
	for _, factory := range p.Resources(context.Background()) {
		r := factory()
		var resp resource.MetadataResponse
		r.Metadata(context.Background(), resource.MetadataRequest{ProviderTypeName: "mailu"}, &resp)
		resourceNames[resp.TypeName] = true
	}

	for _, name := range []string{"mailu_domain", "mailu_user", "mailu_alias"} {
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

	for _, name := range []string{"mailu_domain", "mailu_user"} {
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
