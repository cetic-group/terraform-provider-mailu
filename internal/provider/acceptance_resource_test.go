package provider

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/cetic-group/terraform-provider-mailu/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// testAccProtoV6ProviderFactories wires the in-process provider for acceptance
// tests. The provider reads its endpoint and token from the environment, so no
// secrets are placed in Terraform configuration.
var testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"mailu": providerserver.NewProtocol6WithError(New("acc")()),
}

func testAccPreCheck(t *testing.T) {
	t.Helper()
	if !acceptanceEnabled() {
		t.Skip("set TF_ACC=1 to run acceptance tests")
	}
	if !getAcceptanceConfig().valid() {
		t.Fatal("MAILU_ENDPOINT, MAILU_API_TOKEN, and MAILU_ACC_DOMAIN are required when TF_ACC=1")
	}
}

func testAccClient(t *testing.T) *client.Client {
	t.Helper()
	cfg := getAcceptanceConfig()
	c, err := client.New(cfg.Endpoint, cfg.Token)
	if err != nil {
		t.Fatalf("acceptance client: %v", err)
	}
	return c
}

func TestAccDomainResource(t *testing.T) {
	testAccPreCheck(t)

	c := testAccClient(t)
	name := "tf-acc-domain." + getAcceptanceConfig().Domain

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckDomainDestroyed(c, name),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`resource "mailu_domain" "test" {
  name    = %[1]q
  comment = "managed by acceptance test"
}`, name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("mailu_domain.test", "name", name),
					resource.TestCheckResourceAttr("mailu_domain.test", "comment", "managed by acceptance test"),
				),
			},
			{
				// Update step exercises drift detection on a mutable field.
				Config: fmt.Sprintf(`resource "mailu_domain" "test" {
  name    = %[1]q
  comment = "updated by acceptance test"
}`, name),
				Check: resource.TestCheckResourceAttr("mailu_domain.test", "comment", "updated by acceptance test"),
			},
			{
				ResourceName:      "mailu_domain.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccUserResource(t *testing.T) {
	testAccPreCheck(t)

	c := testAccClient(t)
	domain := "tf-acc-user." + getAcceptanceConfig().Domain
	email := "tf-acc-user@" + domain

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy: resource.ComposeAggregateTestCheckFunc(
			testAccCheckUserDestroyed(c, email),
			testAccCheckDomainDestroyed(c, domain),
		),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`resource "mailu_domain" "test" {
  name = %[1]q
}

resource "mailu_user" "test" {
  email        = %[2]q
  raw_password = "Sup3r-Secret-Acc-Pass!"
  comment      = "managed by acceptance test"
  depends_on   = [mailu_domain.test]
}`, domain, email),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("mailu_user.test", "email", email),
					resource.TestCheckResourceAttrSet("mailu_user.test", "id"),
				),
			},
			{
				ResourceName:            "mailu_user.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"raw_password"},
			},
		},
	})
}

func TestAccAliasResource(t *testing.T) {
	testAccPreCheck(t)

	c := testAccClient(t)
	domain := "tf-acc-alias." + getAcceptanceConfig().Domain
	email := "tf-acc-alias@" + domain

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy: resource.ComposeAggregateTestCheckFunc(
			testAccCheckAliasDestroyed(c, email),
			testAccCheckDomainDestroyed(c, domain),
		),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`resource "mailu_domain" "test" {
  name = %[1]q
}

resource "mailu_alias" "test" {
  email       = %[2]q
  destination = ["tf-acc-dest@%[1]s"]
  depends_on  = [mailu_domain.test]
}`, domain, email),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("mailu_alias.test", "email", email),
					resource.TestCheckResourceAttr("mailu_alias.test", "destination.#", "1"),
				),
			},
			{
				ResourceName:      "mailu_alias.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckDomainDestroyed(c *client.Client, name string) resource.TestCheckFunc {
	return func(_ *terraform.State) error {
		return assertGone("domain", name, func() error {
			_, err := c.GetDomain(context.Background(), name)
			return err
		})
	}
}

func testAccCheckUserDestroyed(c *client.Client, email string) resource.TestCheckFunc {
	return func(_ *terraform.State) error {
		return assertGone("user", email, func() error {
			_, err := c.GetUser(context.Background(), email)
			return err
		})
	}
}

func testAccCheckAliasDestroyed(c *client.Client, email string) resource.TestCheckFunc {
	return func(_ *terraform.State) error {
		return assertGone("alias", email, func() error {
			_, err := c.GetAlias(context.Background(), email)
			return err
		})
	}
}

func assertGone(kind, id string, get func() error) error {
	err := get()
	if err == nil {
		return fmt.Errorf("%s %q still exists after destroy", kind, id)
	}
	var apiErr *client.APIError
	if errors.As(err, &apiErr) && apiErr.StatusCode == 404 {
		return nil
	}
	return fmt.Errorf("unexpected error checking destroyed %s %q: %w", kind, id, err)
}
