package azurerm

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMServiceBusQueueAuthorizationRule_listen(t *testing.T) {
	testAccAzureRMServiceBusQueueAuthorizationRule(t, true, false, false)
}

func TestAccAzureRMServiceBusQueueAuthorizationRule_send(t *testing.T) {
	testAccAzureRMServiceBusQueueAuthorizationRule(t, false, true, false)
}

func TestAccAzureRMServiceBusQueueAuthorizationRule_listensend(t *testing.T) {
	testAccAzureRMServiceBusQueueAuthorizationRule(t, true, true, false)
}

func TestAccAzureRMServiceBusQueueAuthorizationRule_manage(t *testing.T) {
	testAccAzureRMServiceBusQueueAuthorizationRule(t, true, true, true)
}

func testAccAzureRMServiceBusQueueAuthorizationRule(t *testing.T, listen, send, manage bool) {
	resourceName := "azurerm_servicebus_queue_authorization_rule.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMServiceBusQueueAuthorizationRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMServiceBusQueueAuthorizationRule_base(tf.AccRandTimeInt(), testLocation(), listen, send, manage),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMServiceBusQueueAuthorizationRuleExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "name"),
					resource.TestCheckResourceAttrSet(resourceName, "namespace_name"),
					resource.TestCheckResourceAttrSet(resourceName, "queue_name"),
					resource.TestCheckResourceAttrSet(resourceName, "primary_key"),
					resource.TestCheckResourceAttrSet(resourceName, "secondary_key"),
					resource.TestCheckResourceAttrSet(resourceName, "primary_connection_string"),
					resource.TestCheckResourceAttrSet(resourceName, "secondary_connection_string"),
					resource.TestCheckResourceAttr(resourceName, "listen", strconv.FormatBool(listen)),
					resource.TestCheckResourceAttr(resourceName, "send", strconv.FormatBool(send)),
					resource.TestCheckResourceAttr(resourceName, "manage", strconv.FormatBool(manage)),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMServiceBusQueueAuthorizationRule_rightsUpdate(t *testing.T) {
	resourceName := "azurerm_servicebus_queue_authorization_rule.test"
	ri := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMServiceBusQueueAuthorizationRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMServiceBusQueueAuthorizationRule_base(ri, testLocation(), true, false, false),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMServiceBusQueueAuthorizationRuleExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "listen", "true"),
					resource.TestCheckResourceAttr(resourceName, "send", "false"),
					resource.TestCheckResourceAttr(resourceName, "manage", "false"),
				),
			},
			{
				Config: testAccAzureRMServiceBusQueueAuthorizationRule_base(ri, testLocation(), true, true, true),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMServiceBusQueueAuthorizationRuleExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "name"),
					resource.TestCheckResourceAttrSet(resourceName, "namespace_name"),
					resource.TestCheckResourceAttrSet(resourceName, "primary_key"),
					resource.TestCheckResourceAttrSet(resourceName, "secondary_key"),
					resource.TestCheckResourceAttrSet(resourceName, "primary_connection_string"),
					resource.TestCheckResourceAttrSet(resourceName, "secondary_connection_string"),
					resource.TestCheckResourceAttr(resourceName, "listen", "true"),
					resource.TestCheckResourceAttr(resourceName, "send", "true"),
					resource.TestCheckResourceAttr(resourceName, "manage", "true"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
func TestAccAzureRMServiceBusQueueAuthorizationRule_requiresImport(t *testing.T) {
	if !features.ShouldResourcesBeImported() {
		t.Skip("Skipping since resources aren't required to be imported")
		return
	}
	resourceName := "azurerm_servicebus_queue_authorization_rule.test"
	ri := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMServiceBusQueueAuthorizationRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMServiceBusQueueAuthorizationRule_base(ri, testLocation(), true, false, false),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMServiceBusQueueAuthorizationRuleExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "listen", "true"),
					resource.TestCheckResourceAttr(resourceName, "send", "false"),
					resource.TestCheckResourceAttr(resourceName, "manage", "false"),
				),
			},
			{
				Config:      testAccAzureRMServiceBusQueueAuthorizationRule_requiresImport(ri, testLocation(), true, false, false),
				ExpectError: testRequiresImportError("azurerm_servicebus_queue_authorization_rule"),
			},
		},
	})
}

func testCheckAzureRMServiceBusQueueAuthorizationRuleDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*ArmClient).servicebus.QueuesClient
	ctx := testAccProvider.Meta().(*ArmClient).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_servicebus_queue_authorization_rule" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		namespaceName := rs.Primary.Attributes["namespace_name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]
		queueName := rs.Primary.Attributes["queue_name"]

		resp, err := conn.GetAuthorizationRule(ctx, resourceGroup, namespaceName, queueName, name)
		if err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return err
			}
		}
	}

	return nil
}

func testCheckAzureRMServiceBusQueueAuthorizationRuleExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		name := rs.Primary.Attributes["name"]
		namespaceName := rs.Primary.Attributes["namespace_name"]
		queueName := rs.Primary.Attributes["queue_name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]

		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for ServiceBus Queue Authorization Rule: %s", name)
		}

		conn := testAccProvider.Meta().(*ArmClient).servicebus.QueuesClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext

		resp, err := conn.GetAuthorizationRule(ctx, resourceGroup, namespaceName, queueName, name)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: ServiceBus Queue Authorization Rule %q (Queue %q / Namespace %q / Resource Group: %q) does not exist", name, queueName, namespaceName, resourceGroup)
			}

			return fmt.Errorf("Bad: Get on ServiceBus Queue: %+v", err)
		}

		return nil
	}
}

func testAccAzureRMServiceBusQueueAuthorizationRule_base(rInt int, location string, listen, send, manage bool) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_servicebus_namespace" "test" {
  name                = "acctest-%[1]d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  sku                 = "Standard"
}

resource "azurerm_servicebus_queue" "test" {
  name                = "acctest-%[1]d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  namespace_name      = "${azurerm_servicebus_namespace.test.name}"

  enable_partitioning = true
}

resource "azurerm_servicebus_queue_authorization_rule" "test" {
  name                = "acctest-%[1]d"
  namespace_name      = "${azurerm_servicebus_namespace.test.name}"
  queue_name          = "${azurerm_servicebus_queue.test.name}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  listen = %[3]t
  send   = %[4]t
  manage = %[5]t
}
`, rInt, location, listen, send, manage)
}

func testAccAzureRMServiceBusQueueAuthorizationRule_requiresImport(rInt int, location string, listen, send, manage bool) string {
	return fmt.Sprintf(`
%s

resource "azurerm_servicebus_queue_authorization_rule" "import" {
  name                = "${azurerm_servicebus_queue_authorization_rule.test.name}"
  namespace_name      = "${azurerm_servicebus_queue_authorization_rule.test.namespace_name}"
  queue_name          = "${azurerm_servicebus_queue_authorization_rule.test.queue_name}"
  resource_group_name = "${azurerm_servicebus_queue_authorization_rule.test.resource_group_name}"

  listen = "${azurerm_servicebus_queue_authorization_rule.test.listen}"
  send   = "${azurerm_servicebus_queue_authorization_rule.test.send}"
  manage = "${azurerm_servicebus_queue_authorization_rule.test.manage}"
}
`, testAccAzureRMServiceBusQueueAuthorizationRule_base(rInt, location, listen, send, manage))
}
