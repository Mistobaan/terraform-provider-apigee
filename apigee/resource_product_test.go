package apigee

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/zambien/go-apigee-edge"
	"log"
	"strings"
	"testing"
)

func TestAccProduct_Updated(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckProductDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckProductConfigRequired,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckProductExists("apigee_product.foo_product", "foo_product"),
					resource.TestCheckResourceAttr(
						"apigee_product.foo_product", "name", "foo_product"),
					resource.TestCheckResourceAttr(
						"apigee_product.foo_product", "display_name", "foo_product"),
					resource.TestCheckResourceAttr(
						"apigee_product.foo_product", "approval_type", "manual"),
				),
			},

			resource.TestStep{
				Config: testAccCheckProductConfigUpdated,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckProductExists("apigee_product.foo_product", "foo_product_updated"),
					resource.TestCheckResourceAttr(
						"apigee_product.foo_product", "name", "foo_product_updated"),
					resource.TestCheckResourceAttr(
						"apigee_product.foo_product", "display_name", "foo_product_updated_different"),
					resource.TestCheckResourceAttr(
						"apigee_product.foo_product", "description", "no one ever fills this out"),
					resource.TestCheckResourceAttr(
						"apigee_product.foo_product", "approval_type", "auto"),
					resource.TestCheckResourceAttr(
						"apigee_product.foo_product", "api_resources.0", "/**"),
					resource.TestCheckResourceAttr(
						"apigee_product.foo_product", "proxies.0", "helloworld"),
					resource.TestCheckResourceAttr(
						"apigee_product.foo_product", "quota", "1000"),
					resource.TestCheckResourceAttr(
						"apigee_product.foo_product", "quota_interval", "2"),
					resource.TestCheckResourceAttr(
						"apigee_product.foo_product", "scopes.0", "READ"),
					resource.TestCheckResourceAttr(
						"apigee_product.foo_product", "quota_time_unit", "minute"),
					resource.TestCheckResourceAttr(
						"apigee_product.foo_product", "attributes.access", "public"),
					resource.TestCheckResourceAttr(
						"apigee_product.foo_product", "attributes.custom1", "customval1"),
					resource.TestCheckResourceAttr(
						"apigee_product.foo_product", "attributes.custom2", "customval2"),
				),
			},
		},
	})
}

func testAccCheckProductDestroy(s *terraform.State) error {

	client := testAccProvider.Meta().(*apigee.EdgeClient)

	if err := productDestroyHelper(s, client); err != nil {
		return err
	}
	return nil
}

func testAccCheckProductExists(n string, name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*apigee.EdgeClient)
		if err := productExistsHelper(s, client, name); err != nil {
			log.Print("Error in testAccCheckProductExists: %s", err)
			return err
		}
		return nil
	}
}

const testAccCheckProductConfigRequired = `
resource "apigee_product" "foo_product" {
   name = "foo_product"
   approval_type = "manual"
}
`

const testAccCheckProductConfigUpdated = `
resource "apigee_product" "foo_product" {
   name = "foo_product_updated"
   display_name = "foo_product_updated_different"
   description = "no one ever fills this out"
   approval_type = "auto"

   api_resources = ["/**"]
   proxies = ["helloworld"]

   quota = "1000"
   quota_interval = "2"
   quota_time_unit = "minute"

   scopes = ["READ"]

   attributes {
      access = "public"

      custom1 = "customval1"
      custom2 = "customval2"
   }
}
`

func productDestroyHelper(s *terraform.State, client *apigee.EdgeClient) error {

	for _, r := range s.RootModule().Resources {
		id := r.Primary.ID

		if id == "" {
			return fmt.Errorf("No product ID is set")
		}

		_, _, err := client.Products.Get("foo_product")

		if err != nil {
			if strings.Contains(err.Error(), "404 ") {
				return nil
			}
			return fmt.Errorf("Received an error retrieving product  %+v\n", err)
		}
	}

	return fmt.Errorf("Product still exists")
}

func productExistsHelper(s *terraform.State, client *apigee.EdgeClient, name string) error {

	for _, r := range s.RootModule().Resources {
		id := r.Primary.ID

		if id == "" {
			return fmt.Errorf("No product ID is set")
		}

		if productData, _, err := client.Products.Get(name); err != nil {
			return fmt.Errorf("Received an error retrieving product  %+v\n", productData)
		} else {
			log.Print("Created product name: %s", productData.Name)
		}

	}
	return nil
}
