package outscale

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccOutscaleOAPIVMDataSource_basic(t *testing.T) {
	omi := os.Getenv("OUTSCALE_IMAGEID")
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccOAPIVMDataSourceConfig(omi, "c4.large"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"outscale_vm.outscale_vm", "image_id", omi),
					resource.TestCheckResourceAttr(
						"outscale_vm.outscale_vm", "vm_type", "c4.large"),
				),
			},
		},
	})
}

func testAccOAPIVMDataSourceConfig(omi, vmType string) string {
	return fmt.Sprintf(`
		resource "outscale_net" "outscale_net" {
			ip_range = "10.0.0.0/16"

			tags {
				key = "Name"
				value = "testacc-vm-ds"
			}
		}	
		 
 		resource "outscale_subnet" "outscale_subnet" {
			net_id         = "${outscale_net.outscale_net.net_id}"
			ip_range       = "10.0.0.0/24"
			subregion_name = "eu-west-2a"
		}
		 
 		resource "outscale_vm" "outscale_vm" {
			image_id     = "%s"
			vm_type      = "%s"
			keypair_name = "terraform-basic"
			subnet_id    = "${outscale_subnet.outscale_subnet.subnet_id}"
		}
		 
    data "outscale_vm" "basic_web" {
		 filter {
				name   = "vm_ids"
				values = ["${outscale_vm.outscale_vm.vm_id}"]
		  }
		}
	`, omi, vmType)
}
