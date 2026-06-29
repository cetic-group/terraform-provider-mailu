package provider

import (
	"context"
	"errors"
	"net/http"

	"github.com/cetic-group/terraform-provider-mailu/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func stringValue(value string) types.String {
	if value == "" {
		return types.StringNull()
	}

	return types.StringValue(value)
}

func int64Pointer(value types.Int64) *int64 {
	if value.IsNull() || value.IsUnknown() {
		return nil
	}

	raw := value.ValueInt64()
	return &raw
}

func boolPointer(value types.Bool) *bool {
	if value.IsNull() || value.IsUnknown() {
		return nil
	}

	raw := value.ValueBool()
	return &raw
}

func stringSliceFromSet(ctx context.Context, value types.Set, normalizer func(string) string) ([]string, diag.Diagnostics) {
	if value.IsNull() || value.IsUnknown() {
		return nil, nil
	}

	var raw []string
	diags := value.ElementsAs(ctx, &raw, false)
	if diags.HasError() {
		return nil, diags
	}

	return normalizeStrings(raw, normalizer), nil
}

func setFromStrings(ctx context.Context, values []string) (types.Set, diag.Diagnostics) {
	if values == nil {
		return types.SetNull(types.StringType), nil
	}

	return types.SetValueFrom(ctx, types.StringType, values)
}

func listFromStrings(ctx context.Context, values []string) (types.List, diag.Diagnostics) {
	if values == nil {
		return types.ListNull(types.StringType), nil
	}

	return types.ListValueFrom(ctx, types.StringType, values)
}

func isNotFound(err error) bool {
	var apiErr *client.APIError
	return errors.As(err, &apiErr) && apiErr.StatusCode == http.StatusNotFound
}

func addClientError(diags *diag.Diagnostics, summary string, err error) {
	diags.AddError(summary, client.Redact(err.Error()))
}
