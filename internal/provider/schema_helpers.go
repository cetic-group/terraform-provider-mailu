package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// The helpers below build Optional+Computed attributes that keep their prior
// state value when the configuration leaves them unset. Without
// UseStateForUnknown the framework marks every server-defaulted attribute as
// "(known after apply)" on unrelated changes, producing noisy or perpetual
// diffs.

func optionalComputedString() schema.StringAttribute {
	return schema.StringAttribute{
		Optional:      true,
		Computed:      true,
		PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
	}
}

func optionalComputedBool() schema.BoolAttribute {
	return schema.BoolAttribute{
		Optional:      true,
		Computed:      true,
		PlanModifiers: []planmodifier.Bool{boolplanmodifier.UseStateForUnknown()},
	}
}

func optionalComputedInt64() schema.Int64Attribute {
	return schema.Int64Attribute{
		Optional:      true,
		Computed:      true,
		PlanModifiers: []planmodifier.Int64{int64planmodifier.UseStateForUnknown()},
	}
}

func optionalComputedStringSet() schema.SetAttribute {
	return schema.SetAttribute{
		Optional:      true,
		Computed:      true,
		ElementType:   types.StringType,
		PlanModifiers: []planmodifier.Set{setplanmodifier.UseStateForUnknown()},
	}
}
