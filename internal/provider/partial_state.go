package provider

import (
	"context"

	"github.com/cetic-group/terraform-provider-mailu/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/diag"
)

func addPartialCreateWarning(diags *diag.Diagnostics, resourceName string, err error) {
	diags.AddWarning(
		"Mailu "+resourceName+" Created With Incomplete Read",
		"The object was created in Mailu, but the provider could not read it back immediately. Terraform stored the known identity to avoid creating a duplicate on the next apply. Run terraform refresh or terraform plan after the Mailu API is reachable. Error: "+clientSafeError(err),
	)
}

func clientSafeError(err error) string {
	if err == nil {
		return ""
	}

	return client.Redact(err.Error())
}

func setKnownCreateState(ctx context.Context, setter func(context.Context) diag.Diagnostics, diags *diag.Diagnostics, resourceName string, err error) {
	diags.Append(setter(ctx)...)
	if !diags.HasError() {
		addPartialCreateWarning(diags, resourceName, err)
	}
}
