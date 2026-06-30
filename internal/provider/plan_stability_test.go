package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	resourceschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

// optionalComputedAttributes that are not user-set must keep their prior state
// value across unrelated changes. Without a UseStateForUnknown plan modifier the
// framework marks them "(known after apply)" on every plan, producing noisy and
// potentially perpetual diffs. Each such attribute must carry a plan modifier.
func TestOptionalComputedAttributesUseStateForUnknown(t *testing.T) {
	t.Parallel()

	resources := map[string]resource.Resource{
		"domain": NewDomainResource(),
		"user":   NewUserResource(),
		"alias":  NewAliasResource(),
		"relay":  NewRelayResource(),
		"token":  NewTokenResource(),
	}

	for name, r := range resources {
		schema := resourceSchema(t, r)
		for attrName, attr := range schema.Attributes {
			if !isOptionalComputed(attr) {
				continue
			}
			if planModifierCount(attr) == 0 {
				t.Errorf("%s.%s is Optional+Computed but has no plan modifier (expected UseStateForUnknown)", name, attrName)
			}
		}
	}
}

func isOptionalComputed(attr resourceschema.Attribute) bool {
	return attr.IsOptional() && attr.IsComputed()
}

func planModifierCount(attr resourceschema.Attribute) int {
	switch a := attr.(type) {
	case resourceschema.StringAttribute:
		return len(a.PlanModifiers)
	case resourceschema.Int64Attribute:
		return len(a.PlanModifiers)
	case resourceschema.BoolAttribute:
		return len(a.PlanModifiers)
	case resourceschema.SetAttribute:
		return len(a.PlanModifiers)
	case resourceschema.ListAttribute:
		return len(a.PlanModifiers)
	default:
		return 0
	}
}
