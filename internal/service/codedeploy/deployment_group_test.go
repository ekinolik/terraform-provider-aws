package codedeploy_test

import (
	"fmt"
	"reflect"
	"regexp"
	"sort"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/codedeploy"
	sdkacctest "github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	tfcodedeploy "github.com/hashicorp/terraform-provider-aws/internal/service/codedeploy"
)

func TestAccCodeDeployDeploymentGroup_basic(t *testing.T) {
	var group codedeploy.DeploymentGroupInfo
	resourceName := "aws_codedeploy_deployment_group.test"

	rName := sdkacctest.RandString(5)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(t) },
		ErrorCheck:   acctest.ErrorCheck(t, codedeploy.EndpointsID),
		Providers:    acctest.Providers,
		CheckDestroy: testAccCheckDeploymentGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDeploymentGroup(rName, false),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDeploymentGroupExists(resourceName, &group),
					acctest.CheckResourceAttrRegionalARN(resourceName, "arn", "codedeploy", fmt.Sprintf(`deploymentgroup:%s/%s`, "tf-acc-test-"+rName, "tf-acc-test-"+rName)),
					resource.TestCheckResourceAttr(resourceName, "app_name", "tf-acc-test-"+rName),
					resource.TestCheckResourceAttr(resourceName, "deployment_group_name", "tf-acc-test-"+rName),
					resource.TestCheckResourceAttr(resourceName, "deployment_config_name", "CodeDeployDefault.OneAtATime"),
					resource.TestCheckResourceAttrPair(resourceName, "service_role_arn", "aws_iam_role.test", "arn"),

					resource.TestCheckResourceAttr(resourceName, "blue_green_deployment_config.#", "0"),

					resource.TestCheckResourceAttr(resourceName, "ec2_tag_set.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "ec2_tag_filter.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "ec2_tag_filter.*", map[string]string{
						"key":   "filterkey",
						"type":  "KEY_AND_VALUE",
						"value": "filtervalue",
					}),

					resource.TestCheckResourceAttr(resourceName, "alarm_configuration.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "auto_rollback_configuration.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "trigger_configuration.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "compute_platform", "Server"),
					resource.TestCheckResourceAttrSet(resourceName, "deployment_group_id"),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "0"),
				),
			},
			{
				Config: testAccDeploymentGroupModified(rName, false),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDeploymentGroupExists(resourceName, &group),
					acctest.CheckResourceAttrRegionalARN(resourceName, "arn", "codedeploy", fmt.Sprintf(`deploymentgroup:%s/%s`, "tf-acc-test-"+rName, "tf-acc-test-updated-"+rName)),
					resource.TestCheckResourceAttr(resourceName, "app_name", "tf-acc-test-"+rName),
					resource.TestCheckResourceAttr(resourceName, "deployment_group_name", "tf-acc-test-updated-"+rName),
					resource.TestCheckResourceAttr(resourceName, "deployment_config_name", "CodeDeployDefault.OneAtATime"),
					resource.TestCheckResourceAttrPair(resourceName, "service_role_arn", "aws_iam_role.test_updated", "arn"),

					resource.TestCheckResourceAttr(resourceName, "ec2_tag_set.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "ec2_tag_filter.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "ec2_tag_filter.*", map[string]string{
						"key":   "filterkey",
						"type":  "KEY_AND_VALUE",
						"value": "anotherfiltervalue",
					}),

					resource.TestCheckResourceAttr(resourceName, "alarm_configuration.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "auto_rollback_configuration.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "trigger_configuration.#", "0"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateIdFunc: testAccDeploymentGroupImportStateIdFunc(resourceName),
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccCodeDeployDeploymentGroup_Basic_tagSet(t *testing.T) {
	var group codedeploy.DeploymentGroupInfo
	resourceName := "aws_codedeploy_deployment_group.test"

	rName := sdkacctest.RandString(5)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(t) },
		ErrorCheck:   acctest.ErrorCheck(t, codedeploy.EndpointsID),
		Providers:    acctest.Providers,
		CheckDestroy: testAccCheckDeploymentGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDeploymentGroup(rName, true),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDeploymentGroupExists(resourceName, &group),
					resource.TestCheckResourceAttr(resourceName, "app_name", "tf-acc-test-"+rName),
					resource.TestCheckResourceAttr(resourceName, "deployment_group_name", "tf-acc-test-"+rName),
					resource.TestCheckResourceAttr(resourceName, "deployment_config_name", "CodeDeployDefault.OneAtATime"),
					resource.TestCheckResourceAttrPair(resourceName, "service_role_arn", "aws_iam_role.test", "arn"),

					resource.TestCheckResourceAttr(resourceName, "ec2_tag_set.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "ec2_tag_set.*", map[string]string{
						"ec2_tag_filter.#": "1",
					}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "ec2_tag_set.*.ec2_tag_filter.*", map[string]string{
						"key":   "filterkey",
						"type":  "KEY_AND_VALUE",
						"value": "filtervalue",
					}),
					resource.TestCheckResourceAttr(resourceName, "ec2_tag_filter.#", "0"),

					resource.TestCheckResourceAttr(resourceName, "alarm_configuration.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "auto_rollback_configuration.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "trigger_configuration.#", "0"),
				),
			},
			{
				Config: testAccDeploymentGroupModified(rName, true),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDeploymentGroupExists(resourceName, &group),
					resource.TestCheckResourceAttr(resourceName, "app_name", "tf-acc-test-"+rName),
					resource.TestCheckResourceAttr(resourceName, "deployment_group_name", "tf-acc-test-updated-"+rName),
					resource.TestCheckResourceAttr(resourceName, "deployment_config_name", "CodeDeployDefault.OneAtATime"),
					resource.TestCheckResourceAttrPair(resourceName, "service_role_arn", "aws_iam_role.test_updated", "arn"),

					resource.TestCheckResourceAttr(resourceName, "ec2_tag_set.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "ec2_tag_set.*", map[string]string{
						"ec2_tag_filter.#": "1",
					}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "ec2_tag_set.*.ec2_tag_filter.*", map[string]string{
						"key":   "filterkey",
						"type":  "KEY_AND_VALUE",
						"value": "anotherfiltervalue",
					}),
					resource.TestCheckResourceAttr(resourceName, "ec2_tag_filter.#", "0"),

					resource.TestCheckResourceAttr(resourceName, "alarm_configuration.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "auto_rollback_configuration.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "trigger_configuration.#", "0"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateIdFunc: testAccDeploymentGroupImportStateIdFunc(resourceName),
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccCodeDeployDeploymentGroup_onPremiseTag(t *testing.T) {
	var group codedeploy.DeploymentGroupInfo
	resourceName := "aws_codedeploy_deployment_group.test"
	rName := sdkacctest.RandString(5)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(t) },
		ErrorCheck:   acctest.ErrorCheck(t, codedeploy.EndpointsID),
		Providers:    acctest.Providers,
		CheckDestroy: testAccCheckDeploymentGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDeploymentGroupOnPremiseTags(rName, true),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDeploymentGroupExists(resourceName, &group),
					resource.TestCheckResourceAttr(resourceName, "app_name", "tf-acc-test-"+rName),
					resource.TestCheckResourceAttr(resourceName, "deployment_group_name", "tf-acc-test-"+rName),
					resource.TestCheckResourceAttr(resourceName, "deployment_config_name", "CodeDeployDefault.OneAtATime"),
					resource.TestCheckResourceAttrPair(resourceName, "service_role_arn", "aws_iam_role.test", "arn"),

					resource.TestCheckResourceAttr(resourceName, "on_premises_instance_tag_set.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "on_premises_instance_tag_set.*", map[string]string{
						"on_premises_instance_tag_filter.#": "1",
					}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "on_premises_instance_tag_set.*.on_premises_instance_tag_filter.*", map[string]string{
						"key":   "filterkey",
						"type":  "KEY_AND_VALUE",
						"value": "filtervalue",
					}),
					resource.TestCheckResourceAttr(resourceName, "on_premises_instance_tag_filter.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "alarm_configuration.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "auto_rollback_configuration.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "trigger_configuration.#", "0"),
				),
			},
			{
				Config: testAccDeploymentGroupOnPremiseTags(rName, true),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDeploymentGroupExists(resourceName, &group),
					resource.TestCheckResourceAttr(resourceName, "app_name", "tf-acc-test-"+rName),
					resource.TestCheckResourceAttr(resourceName, "deployment_group_name", "tf-acc-test-"+rName),
					resource.TestCheckResourceAttr(resourceName, "deployment_config_name", "CodeDeployDefault.OneAtATime"),
					resource.TestCheckResourceAttrPair(resourceName, "service_role_arn", "aws_iam_role.test", "arn"),

					resource.TestCheckResourceAttr(resourceName, "on_premises_instance_tag_set.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "on_premises_instance_tag_set.*", map[string]string{
						"on_premises_instance_tag_filter.#": "1",
					}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "on_premises_instance_tag_set.*.on_premises_instance_tag_filter.*", map[string]string{
						"key":   "filterkey",
						"type":  "KEY_AND_VALUE",
						"value": "filtervalue",
					}),
					resource.TestCheckResourceAttr(resourceName, "on_premises_instance_tag_filter.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "alarm_configuration.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "auto_rollback_configuration.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "trigger_configuration.#", "0"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateIdFunc: testAccDeploymentGroupImportStateIdFunc(resourceName),
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccCodeDeployDeploymentGroup_disappears(t *testing.T) {
	var group codedeploy.DeploymentGroupInfo
	resourceName := "aws_codedeploy_deployment_group.test"
	rName := sdkacctest.RandString(5)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(t) },
		ErrorCheck:   acctest.ErrorCheck(t, codedeploy.EndpointsID),
		Providers:    acctest.Providers,
		CheckDestroy: testAccCheckDeploymentGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDeploymentGroup(rName, false),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDeploymentGroupExists(resourceName, &group),
					acctest.CheckResourceDisappears(acctest.Provider, tfcodedeploy.ResourceDeploymentGroup(), resourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccCodeDeployDeploymentGroup_Disappears_app(t *testing.T) {
	var group codedeploy.DeploymentGroupInfo
	resourceName := "aws_codedeploy_deployment_group.test"
	rName := sdkacctest.RandString(5)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(t) },
		ErrorCheck:   acctest.ErrorCheck(t, codedeploy.EndpointsID),
		Providers:    acctest.Providers,
		CheckDestroy: testAccCheckDeploymentGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDeploymentGroup(rName, false),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDeploymentGroupExists(resourceName, &group),
					acctest.CheckResourceDisappears(acctest.Provider, tfcodedeploy.ResourceApp(), "aws_codedeploy_app.test"),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccCodeDeployDeploymentGroup_tags(t *testing.T) {
	var group codedeploy.DeploymentGroupInfo
	resourceName := "aws_codedeploy_deployment_group.test"

	rName := sdkacctest.RandString(5)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(t) },
		ErrorCheck:   acctest.ErrorCheck(t, codedeploy.EndpointsID),
		Providers:    acctest.Providers,
		CheckDestroy: testAccCheckDeploymentGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDeploymentGroupTags1Config(rName, "key1", "value1"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDeploymentGroupExists(resourceName, &group),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.key1", "value1"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateIdFunc: testAccDeploymentGroupImportStateIdFunc(resourceName),
				ImportStateVerify: true,
			},
			{
				Config: testAccDeploymentGroupTags2Config(rName, "key1", "value1updated", "key2", "value2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDeploymentGroupExists(resourceName, &group),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "2"),
					resource.TestCheckResourceAttr(resourceName, "tags.key1", "value1updated"),
					resource.TestCheckResourceAttr(resourceName, "tags.key2", "value2"),
				),
			},
			{
				Config: testAccDeploymentGroupTags1Config(rName, "key2", "value2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDeploymentGroupExists(resourceName, &group),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.key2", "value2"),
				),
			},
		},
	})
}

func TestAccCodeDeployDeploymentGroup_Trigger_basic(t *testing.T) {
	var group codedeploy.DeploymentGroupInfo
	resourceName := "aws_codedeploy_deployment_group.test"

	rName := sdkacctest.RandString(5)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(t) },
		ErrorCheck:   acctest.ErrorCheck(t, codedeploy.EndpointsID),
		Providers:    acctest.Providers,
		CheckDestroy: testAccCheckDeploymentGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDeploymentGroup_triggerConfiguration_create(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDeploymentGroupExists(resourceName, &group),
					resource.TestCheckResourceAttr(
						resourceName, "app_name", "tf-acc-test-"+rName),
					resource.TestCheckResourceAttr(
						resourceName, "deployment_group_name", "tf-acc-test-"+rName),
					testAccCheckCodeDeployDeploymentGroupTriggerEvents(&group, "test-trigger", []string{
						"DeploymentFailure",
					}),
				),
			},
			{
				Config: testAccDeploymentGroup_triggerConfiguration_update(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDeploymentGroupExists(resourceName, &group),
					resource.TestCheckResourceAttr(
						resourceName, "app_name", "tf-acc-test-"+rName),
					resource.TestCheckResourceAttr(
						resourceName, "deployment_group_name", "tf-acc-test-"+rName),
					testAccCheckCodeDeployDeploymentGroupTriggerEvents(&group, "test-trigger", []string{
						"DeploymentFailure",
						"DeploymentSuccess",
					}),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateIdFunc: testAccDeploymentGroupImportStateIdFunc(resourceName),
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccCodeDeployDeploymentGroup_Trigger_multiple(t *testing.T) {
	var group codedeploy.DeploymentGroupInfo
	resourceName := "aws_codedeploy_deployment_group.test"

	rName := sdkacctest.RandString(5)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(t) },
		ErrorCheck:   acctest.ErrorCheck(t, codedeploy.EndpointsID),
		Providers:    acctest.Providers,
		CheckDestroy: testAccCheckDeploymentGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDeploymentGroup_triggerConfiguration_createMultiple(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDeploymentGroupExists(resourceName, &group),
					resource.TestCheckResourceAttr(
						resourceName, "app_name", "tf-acc-test-"+rName),
					resource.TestCheckResourceAttr(
						resourceName, "deployment_group_name", "tf-acc-test-"+rName),
					testAccCheckCodeDeployDeploymentGroupTriggerEvents(&group, "test-trigger-1", []string{
						"DeploymentFailure",
					}),
					testAccCheckCodeDeployDeploymentGroupTriggerEvents(&group, "test-trigger-2", []string{
						"InstanceFailure",
					}),
					testAccCheckCodeDeployDeploymentGroupTriggerTargetArn(&group, "test-trigger-2",
						regexp.MustCompile(fmt.Sprintf("^arn:%s:sns:[^:]+:[0-9]{12}:tf-acc-test-2-%s$", acctest.Partition(), rName))),
				),
			},
			{
				Config: testAccDeploymentGroup_triggerConfiguration_updateMultiple(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDeploymentGroupExists(resourceName, &group),
					resource.TestCheckResourceAttr(
						resourceName, "app_name", "tf-acc-test-"+rName),
					resource.TestCheckResourceAttr(
						resourceName, "deployment_group_name", "tf-acc-test-"+rName),
					testAccCheckCodeDeployDeploymentGroupTriggerEvents(&group, "test-trigger-1", []string{
						"DeploymentFailure",
						"DeploymentStart",
						"DeploymentStop",
						"DeploymentSuccess",
					}),
					testAccCheckCodeDeployDeploymentGroupTriggerEvents(&group, "test-trigger-2", []string{
						"InstanceFailure",
					}),
					testAccCheckCodeDeployDeploymentGroupTriggerTargetArn(&group, "test-trigger-2",
						regexp.MustCompile(fmt.Sprintf("^arn:%s:sns:[^:]+:[0-9]{12}:tf-acc-test-3-%s$", acctest.Partition(), rName))),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateIdFunc: testAccDeploymentGroupImportStateIdFunc(resourceName),
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccCodeDeployDeploymentGroup_AutoRollback_create(t *testing.T) {
	var group codedeploy.DeploymentGroupInfo
	resourceName := "aws_codedeploy_deployment_group.test"

	rName := sdkacctest.RandString(5)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(t) },
		ErrorCheck:   acctest.ErrorCheck(t, codedeploy.EndpointsID),
		Providers:    acctest.Providers,
		CheckDestroy: testAccCheckDeploymentGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: test_config_auto_rollback_configuration_create(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDeploymentGroupExists(resourceName, &group),
					resource.TestCheckResourceAttr(
						resourceName, "auto_rollback_configuration.#", "1"),
					resource.TestCheckResourceAttr(
						resourceName, "auto_rollback_configuration.0.enabled", "true"),
					resource.TestCheckResourceAttr(
						resourceName, "auto_rollback_configuration.0.events.#", "1"),
					resource.TestCheckTypeSetElemAttr(resourceName, "auto_rollback_configuration.0.events.*", "DEPLOYMENT_FAILURE"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateIdFunc: testAccDeploymentGroupImportStateIdFunc(resourceName),
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccCodeDeployDeploymentGroup_AutoRollback_update(t *testing.T) {
	var group codedeploy.DeploymentGroupInfo
	resourceName := "aws_codedeploy_deployment_group.test"

	rName := sdkacctest.RandString(5)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(t) },
		ErrorCheck:   acctest.ErrorCheck(t, codedeploy.EndpointsID),
		Providers:    acctest.Providers,
		CheckDestroy: testAccCheckDeploymentGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: test_config_auto_rollback_configuration_create(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDeploymentGroupExists(resourceName, &group),
					resource.TestCheckResourceAttr(
						resourceName, "auto_rollback_configuration.#", "1"),
					resource.TestCheckResourceAttr(
						resourceName, "auto_rollback_configuration.0.enabled", "true"),
					resource.TestCheckResourceAttr(
						resourceName, "auto_rollback_configuration.0.events.#", "1"),
					resource.TestCheckTypeSetElemAttr(resourceName, "auto_rollback_configuration.0.events.*", "DEPLOYMENT_FAILURE"),
				),
			},
			{
				Config: test_config_auto_rollback_configuration_update(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDeploymentGroupExists(resourceName, &group),
					resource.TestCheckResourceAttr(
						resourceName, "auto_rollback_configuration.#", "1"),
					resource.TestCheckResourceAttr(
						resourceName, "auto_rollback_configuration.0.enabled", "true"),
					resource.TestCheckResourceAttr(
						resourceName, "auto_rollback_configuration.0.events.#", "2"),
					resource.TestCheckTypeSetElemAttr(resourceName, "auto_rollback_configuration.0.events.*", "DEPLOYMENT_STOP_ON_ALARM"),
					resource.TestCheckTypeSetElemAttr(resourceName, "auto_rollback_configuration.0.events.*", "DEPLOYMENT_FAILURE"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateIdFunc: testAccDeploymentGroupImportStateIdFunc(resourceName),
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccCodeDeployDeploymentGroup_AutoRollback_delete(t *testing.T) {
	var group codedeploy.DeploymentGroupInfo
	resourceName := "aws_codedeploy_deployment_group.test"

	rName := sdkacctest.RandString(5)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(t) },
		ErrorCheck:   acctest.ErrorCheck(t, codedeploy.EndpointsID),
		Providers:    acctest.Providers,
		CheckDestroy: testAccCheckDeploymentGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: test_config_auto_rollback_configuration_create(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDeploymentGroupExists(resourceName, &group),
					resource.TestCheckResourceAttr(
						resourceName, "auto_rollback_configuration.#", "1"),
					resource.TestCheckResourceAttr(
						resourceName, "auto_rollback_configuration.0.enabled", "true"),
					resource.TestCheckResourceAttr(
						resourceName, "auto_rollback_configuration.0.events.#", "1"),
					resource.TestCheckTypeSetElemAttr(resourceName, "auto_rollback_configuration.0.events.*", "DEPLOYMENT_FAILURE"),
				),
			},
			{
				Config: test_config_auto_rollback_configuration_none(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDeploymentGroupExists(resourceName, &group),
					resource.TestCheckResourceAttr(
						resourceName, "auto_rollback_configuration.#", "0"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateIdFunc: testAccDeploymentGroupImportStateIdFunc(resourceName),
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccCodeDeployDeploymentGroup_AutoRollback_disable(t *testing.T) {
	var group codedeploy.DeploymentGroupInfo
	resourceName := "aws_codedeploy_deployment_group.test"

	rName := sdkacctest.RandString(5)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(t) },
		ErrorCheck:   acctest.ErrorCheck(t, codedeploy.EndpointsID),
		Providers:    acctest.Providers,
		CheckDestroy: testAccCheckDeploymentGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: test_config_auto_rollback_configuration_create(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDeploymentGroupExists(resourceName, &group),
					resource.TestCheckResourceAttr(
						resourceName, "auto_rollback_configuration.#", "1"),
					resource.TestCheckResourceAttr(
						resourceName, "auto_rollback_configuration.0.enabled", "true"),
					resource.TestCheckResourceAttr(
						resourceName, "auto_rollback_configuration.0.events.#", "1"),
					resource.TestCheckTypeSetElemAttr(resourceName, "auto_rollback_configuration.0.events.*", "DEPLOYMENT_FAILURE"),
				),
			},
			{
				Config: test_config_auto_rollback_configuration_disable(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDeploymentGroupExists(resourceName, &group),
					resource.TestCheckResourceAttr(
						resourceName, "auto_rollback_configuration.#", "1"),
					resource.TestCheckResourceAttr(
						resourceName, "auto_rollback_configuration.0.enabled", "false"),
					resource.TestCheckResourceAttr(
						resourceName, "auto_rollback_configuration.0.events.#", "1"),
					resource.TestCheckTypeSetElemAttr(resourceName, "auto_rollback_configuration.0.events.*", "DEPLOYMENT_FAILURE"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateIdFunc: testAccDeploymentGroupImportStateIdFunc(resourceName),
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccCodeDeployDeploymentGroup_Alarm_create(t *testing.T) {
	var group codedeploy.DeploymentGroupInfo
	resourceName := "aws_codedeploy_deployment_group.test"

	rName := sdkacctest.RandString(5)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(t) },
		ErrorCheck:   acctest.ErrorCheck(t, codedeploy.EndpointsID),
		Providers:    acctest.Providers,
		CheckDestroy: testAccCheckDeploymentGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: test_config_alarm_configuration_create(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDeploymentGroupExists(resourceName, &group),
					resource.TestCheckResourceAttr(
						resourceName, "alarm_configuration.#", "1"),
					resource.TestCheckResourceAttr(
						resourceName, "alarm_configuration.0.enabled", "true"),
					resource.TestCheckResourceAttr(
						resourceName, "alarm_configuration.0.alarms.#", "1"),
					resource.TestCheckTypeSetElemAttr(resourceName, "alarm_configuration.0.alarms.*", "test-alarm"),
					resource.TestCheckResourceAttr(
						resourceName, "alarm_configuration.0.ignore_poll_alarm_failure", "false"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateIdFunc: testAccDeploymentGroupImportStateIdFunc(resourceName),
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccCodeDeployDeploymentGroup_Alarm_update(t *testing.T) {
	var group codedeploy.DeploymentGroupInfo
	resourceName := "aws_codedeploy_deployment_group.test"

	rName := sdkacctest.RandString(5)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(t) },
		ErrorCheck:   acctest.ErrorCheck(t, codedeploy.EndpointsID),
		Providers:    acctest.Providers,
		CheckDestroy: testAccCheckDeploymentGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: test_config_alarm_configuration_create(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDeploymentGroupExists(resourceName, &group),
					resource.TestCheckResourceAttr(
						resourceName, "alarm_configuration.#", "1"),
					resource.TestCheckResourceAttr(
						resourceName, "alarm_configuration.0.enabled", "true"),
					resource.TestCheckResourceAttr(
						resourceName, "alarm_configuration.0.alarms.#", "1"),
					resource.TestCheckTypeSetElemAttr(resourceName, "alarm_configuration.0.alarms.*", "test-alarm"),
					resource.TestCheckResourceAttr(
						resourceName, "alarm_configuration.0.ignore_poll_alarm_failure", "false"),
				),
			},
			{
				Config: test_config_alarm_configuration_update(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDeploymentGroupExists(resourceName, &group),
					resource.TestCheckResourceAttr(
						resourceName, "alarm_configuration.#", "1"),
					resource.TestCheckResourceAttr(
						resourceName, "alarm_configuration.0.enabled", "true"),
					resource.TestCheckResourceAttr(
						resourceName, "alarm_configuration.0.alarms.#", "2"),
					resource.TestCheckTypeSetElemAttr(resourceName, "alarm_configuration.0.alarms.*", "test-alarm"),
					resource.TestCheckTypeSetElemAttr(resourceName, "alarm_configuration.0.alarms.*", "test-alarm-2"),
					resource.TestCheckResourceAttr(
						resourceName, "alarm_configuration.0.ignore_poll_alarm_failure", "true"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateIdFunc: testAccDeploymentGroupImportStateIdFunc(resourceName),
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccCodeDeployDeploymentGroup_Alarm_delete(t *testing.T) {
	var group codedeploy.DeploymentGroupInfo
	resourceName := "aws_codedeploy_deployment_group.test"

	rName := sdkacctest.RandString(5)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(t) },
		ErrorCheck:   acctest.ErrorCheck(t, codedeploy.EndpointsID),
		Providers:    acctest.Providers,
		CheckDestroy: testAccCheckDeploymentGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: test_config_alarm_configuration_create(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDeploymentGroupExists(resourceName, &group),
					resource.TestCheckResourceAttr(
						resourceName, "alarm_configuration.#", "1"),
					resource.TestCheckResourceAttr(
						resourceName, "alarm_configuration.0.enabled", "true"),
					resource.TestCheckResourceAttr(
						resourceName, "alarm_configuration.0.alarms.#", "1"),
					resource.TestCheckTypeSetElemAttr(resourceName, "alarm_configuration.0.alarms.*", "test-alarm"),
					resource.TestCheckResourceAttr(
						resourceName, "alarm_configuration.0.ignore_poll_alarm_failure", "false"),
				),
			},
			{
				Config: test_config_alarm_configuration_none(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDeploymentGroupExists(resourceName, &group),
					resource.TestCheckResourceAttr(
						resourceName, "alarm_configuration.#", "0"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateIdFunc: testAccDeploymentGroupImportStateIdFunc(resourceName),
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccCodeDeployDeploymentGroup_Alarm_disable(t *testing.T) {
	var group codedeploy.DeploymentGroupInfo
	resourceName := "aws_codedeploy_deployment_group.test"

	rName := sdkacctest.RandString(5)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(t) },
		ErrorCheck:   acctest.ErrorCheck(t, codedeploy.EndpointsID),
		Providers:    acctest.Providers,
		CheckDestroy: testAccCheckDeploymentGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: test_config_alarm_configuration_create(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDeploymentGroupExists(resourceName, &group),
					resource.TestCheckResourceAttr(
						resourceName, "alarm_configuration.#", "1"),
					resource.TestCheckResourceAttr(
						resourceName, "alarm_configuration.0.enabled", "true"),
					resource.TestCheckResourceAttr(
						resourceName, "alarm_configuration.0.alarms.#", "1"),
					resource.TestCheckTypeSetElemAttr(resourceName, "alarm_configuration.0.alarms.*", "test-alarm"),
					resource.TestCheckResourceAttr(
						resourceName, "alarm_configuration.0.ignore_poll_alarm_failure", "false"),
				),
			},
			{
				Config: test_config_alarm_configuration_disable(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDeploymentGroupExists(resourceName, &group),
					resource.TestCheckResourceAttr(
						resourceName, "alarm_configuration.#", "1"),
					resource.TestCheckResourceAttr(
						resourceName, "alarm_configuration.0.enabled", "false"),
					resource.TestCheckResourceAttr(
						resourceName, "alarm_configuration.0.alarms.#", "1"),
					resource.TestCheckTypeSetElemAttr(resourceName, "alarm_configuration.0.alarms.*", "test-alarm"),
					resource.TestCheckResourceAttr(
						resourceName, "alarm_configuration.0.ignore_poll_alarm_failure", "false"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateIdFunc: testAccDeploymentGroupImportStateIdFunc(resourceName),
				ImportStateVerify: true,
			},
		},
	})
}

// When no configuration is provided, a deploymentStyle object with default values is computed
func TestAccCodeDeployDeploymentGroup_DeploymentStyle_default(t *testing.T) {
	var group codedeploy.DeploymentGroupInfo
	resourceName := "aws_codedeploy_deployment_group.test"

	rName := sdkacctest.RandString(5)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(t) },
		ErrorCheck:   acctest.ErrorCheck(t, codedeploy.EndpointsID),
		Providers:    acctest.Providers,
		CheckDestroy: testAccCheckDeploymentGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: test_config_deployment_style_default(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDeploymentGroupExists(resourceName, &group),
					resource.TestCheckResourceAttr(
						resourceName, "deployment_style.#", "1"),
					resource.TestCheckResourceAttrSet(
						resourceName, "deployment_style.0.deployment_option"),
					resource.TestCheckResourceAttrSet(
						resourceName, "deployment_style.0.deployment_type"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateIdFunc: testAccDeploymentGroupImportStateIdFunc(resourceName),
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccCodeDeployDeploymentGroup_DeploymentStyle_create(t *testing.T) {
	var group codedeploy.DeploymentGroupInfo
	resourceName := "aws_codedeploy_deployment_group.test"

	rName := sdkacctest.RandString(5)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(t) },
		ErrorCheck:   acctest.ErrorCheck(t, codedeploy.EndpointsID),
		Providers:    acctest.Providers,
		CheckDestroy: testAccCheckDeploymentGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: test_config_deployment_style_create(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDeploymentGroupExists(resourceName, &group),
					resource.TestCheckResourceAttr(
						resourceName, "deployment_style.#", "1"),
					resource.TestCheckResourceAttr(
						resourceName, "deployment_style.0.deployment_option", "WITH_TRAFFIC_CONTROL"),
					resource.TestCheckResourceAttr(
						resourceName, "deployment_style.0.deployment_type", "BLUE_GREEN"),

					resource.TestCheckResourceAttr(
						resourceName, "load_balancer_info.#", "1"),
					resource.TestCheckResourceAttr(
						resourceName, "load_balancer_info.0.elb_info.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "load_balancer_info.0.elb_info.*", map[string]string{
						"name": "acc-test-codedeploy-dep-group",
					}),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateIdFunc: testAccDeploymentGroupImportStateIdFunc(resourceName),
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccCodeDeployDeploymentGroup_DeploymentStyle_update(t *testing.T) {
	var group codedeploy.DeploymentGroupInfo
	resourceName := "aws_codedeploy_deployment_group.test"

	rName := sdkacctest.RandString(5)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(t) },
		ErrorCheck:   acctest.ErrorCheck(t, codedeploy.EndpointsID),
		Providers:    acctest.Providers,
		CheckDestroy: testAccCheckDeploymentGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: test_config_deployment_style_create(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDeploymentGroupExists(resourceName, &group),
					resource.TestCheckResourceAttr(
						resourceName, "deployment_style.#", "1"),
					resource.TestCheckResourceAttr(
						resourceName, "deployment_style.0.deployment_option", "WITH_TRAFFIC_CONTROL"),
					resource.TestCheckResourceAttr(
						resourceName, "deployment_style.0.deployment_type", "BLUE_GREEN"),
				),
			},
			{
				Config: test_config_deployment_style_update(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDeploymentGroupExists(resourceName, &group),
					resource.TestCheckResourceAttr(
						resourceName, "deployment_style.#", "1"),
					resource.TestCheckResourceAttr(
						resourceName, "deployment_style.0.deployment_option", "WITHOUT_TRAFFIC_CONTROL"),
					resource.TestCheckResourceAttr(
						resourceName, "deployment_style.0.deployment_type", "IN_PLACE"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateIdFunc: testAccDeploymentGroupImportStateIdFunc(resourceName),
				ImportStateVerify: true,
			},
		},
	})
}

// Delete reverts to default configuration. It does not remove the deployment_style block
func TestAccCodeDeployDeploymentGroup_DeploymentStyle_delete(t *testing.T) {
	var group codedeploy.DeploymentGroupInfo
	resourceName := "aws_codedeploy_deployment_group.test"

	rName := sdkacctest.RandString(5)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(t) },
		ErrorCheck:   acctest.ErrorCheck(t, codedeploy.EndpointsID),
		Providers:    acctest.Providers,
		CheckDestroy: testAccCheckDeploymentGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: test_config_deployment_style_create(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDeploymentGroupExists(resourceName, &group),
					resource.TestCheckResourceAttr(
						resourceName, "deployment_style.#", "1"),
					resource.TestCheckResourceAttr(
						resourceName, "deployment_style.0.deployment_option", "WITH_TRAFFIC_CONTROL"),
					resource.TestCheckResourceAttr(
						resourceName, "deployment_style.0.deployment_type", "BLUE_GREEN"),
				),
			},
			{
				Config: test_config_deployment_style_default(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDeploymentGroupExists(resourceName, &group),
					resource.TestCheckResourceAttr(
						resourceName, "deployment_style.#", "1"),
					resource.TestCheckResourceAttr(
						resourceName, "deployment_style.0.deployment_option", "WITHOUT_TRAFFIC_CONTROL"),
					resource.TestCheckResourceAttr(
						resourceName, "deployment_style.0.deployment_type", "IN_PLACE"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateIdFunc: testAccDeploymentGroupImportStateIdFunc(resourceName),
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccCodeDeployDeploymentGroup_LoadBalancerInfo_create(t *testing.T) {
	var group codedeploy.DeploymentGroupInfo
	resourceName := "aws_codedeploy_deployment_group.test"

	rName := sdkacctest.RandString(5)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(t) },
		ErrorCheck:   acctest.ErrorCheck(t, codedeploy.EndpointsID),
		Providers:    acctest.Providers,
		CheckDestroy: testAccCheckDeploymentGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: test_config_load_balancer_info_create(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDeploymentGroupExists(resourceName, &group),
					resource.TestCheckResourceAttr(
						resourceName, "load_balancer_info.#", "1"),
					resource.TestCheckResourceAttr(
						resourceName, "load_balancer_info.0.elb_info.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "load_balancer_info.0.elb_info.*", map[string]string{
						"name": "acc-test-codedeploy-dep-group",
					}),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateIdFunc: testAccDeploymentGroupImportStateIdFunc(resourceName),
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccCodeDeployDeploymentGroup_LoadBalancerInfo_update(t *testing.T) {
	var group codedeploy.DeploymentGroupInfo
	resourceName := "aws_codedeploy_deployment_group.test"

	rName := sdkacctest.RandString(5)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(t) },
		ErrorCheck:   acctest.ErrorCheck(t, codedeploy.EndpointsID),
		Providers:    acctest.Providers,
		CheckDestroy: testAccCheckDeploymentGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: test_config_load_balancer_info_create(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDeploymentGroupExists(resourceName, &group),
					resource.TestCheckResourceAttr(
						resourceName, "load_balancer_info.#", "1"),
					resource.TestCheckResourceAttr(
						resourceName, "load_balancer_info.0.elb_info.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "load_balancer_info.0.elb_info.*", map[string]string{
						"name": "acc-test-codedeploy-dep-group",
					}),
				),
			},
			{
				Config: test_config_load_balancer_info_update(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDeploymentGroupExists(resourceName, &group),
					resource.TestCheckResourceAttr(
						resourceName, "load_balancer_info.#", "1"),
					resource.TestCheckResourceAttr(
						resourceName, "load_balancer_info.0.elb_info.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "load_balancer_info.0.elb_info.*", map[string]string{
						"name": "acc-test-codedeploy-dep-group-2",
					}),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateIdFunc: testAccDeploymentGroupImportStateIdFunc(resourceName),
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccCodeDeployDeploymentGroup_LoadBalancerInfo_delete(t *testing.T) {
	var group codedeploy.DeploymentGroupInfo
	resourceName := "aws_codedeploy_deployment_group.test"

	rName := sdkacctest.RandString(5)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(t) },
		ErrorCheck:   acctest.ErrorCheck(t, codedeploy.EndpointsID),
		Providers:    acctest.Providers,
		CheckDestroy: testAccCheckDeploymentGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: test_config_load_balancer_info_create(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDeploymentGroupExists(resourceName, &group),
					resource.TestCheckResourceAttr(
						resourceName, "load_balancer_info.#", "1"),
					resource.TestCheckResourceAttr(
						resourceName, "load_balancer_info.0.elb_info.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "load_balancer_info.0.elb_info.*", map[string]string{
						"name": "acc-test-codedeploy-dep-group",
					}),
				),
			},
			{
				Config: test_config_load_balancer_info_none(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDeploymentGroupExists(resourceName, &group),
					resource.TestCheckResourceAttr(
						resourceName, "load_balancer_info.#", "0"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateIdFunc: testAccDeploymentGroupImportStateIdFunc(resourceName),
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccCodeDeployDeploymentGroup_LoadBalancerInfoTargetGroupInfo_create(t *testing.T) {
	var group codedeploy.DeploymentGroupInfo
	resourceName := "aws_codedeploy_deployment_group.test"

	rName := sdkacctest.RandString(5)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(t) },
		ErrorCheck:   acctest.ErrorCheck(t, codedeploy.EndpointsID),
		Providers:    acctest.Providers,
		CheckDestroy: testAccCheckDeploymentGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: test_config_load_balancer_info_target_group_info_create(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDeploymentGroupExists(resourceName, &group),
					resource.TestCheckResourceAttr(
						resourceName, "load_balancer_info.#", "1"),

					resource.TestCheckResourceAttr(
						resourceName, "load_balancer_info.0.target_group_info.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "load_balancer_info.0.target_group_info.*", map[string]string{
						"name": "acc-test-codedeploy-dep-group",
					}),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateIdFunc: testAccDeploymentGroupImportStateIdFunc(resourceName),
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccCodeDeployDeploymentGroup_LoadBalancerInfoTargetGroupInfo_update(t *testing.T) {
	var group codedeploy.DeploymentGroupInfo
	resourceName := "aws_codedeploy_deployment_group.test"

	rName := sdkacctest.RandString(5)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(t) },
		ErrorCheck:   acctest.ErrorCheck(t, codedeploy.EndpointsID),
		Providers:    acctest.Providers,
		CheckDestroy: testAccCheckDeploymentGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: test_config_load_balancer_info_target_group_info_create(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDeploymentGroupExists(resourceName, &group),
					resource.TestCheckResourceAttr(
						resourceName, "load_balancer_info.#", "1"),
					resource.TestCheckResourceAttr(
						resourceName, "load_balancer_info.0.target_group_info.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "load_balancer_info.0.target_group_info.*", map[string]string{
						"name": "acc-test-codedeploy-dep-group",
					}),
				),
			},
			{
				Config: test_config_load_balancer_info_target_group_info_update(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDeploymentGroupExists(resourceName, &group),
					resource.TestCheckResourceAttr(
						resourceName, "load_balancer_info.#", "1"),
					resource.TestCheckResourceAttr(
						resourceName, "load_balancer_info.0.target_group_info.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "load_balancer_info.0.target_group_info.*", map[string]string{
						"name": "acc-test-codedeploy-dep-group-2",
					}),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateIdFunc: testAccDeploymentGroupImportStateIdFunc(resourceName),
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccCodeDeployDeploymentGroup_LoadBalancerInfoTargetGroupInfo_delete(t *testing.T) {
	var group codedeploy.DeploymentGroupInfo
	resourceName := "aws_codedeploy_deployment_group.test"

	rName := sdkacctest.RandString(5)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(t) },
		ErrorCheck:   acctest.ErrorCheck(t, codedeploy.EndpointsID),
		Providers:    acctest.Providers,
		CheckDestroy: testAccCheckDeploymentGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: test_config_load_balancer_info_target_group_info_create(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDeploymentGroupExists(resourceName, &group),
					resource.TestCheckResourceAttr(
						resourceName, "load_balancer_info.#", "1"),
					resource.TestCheckResourceAttr(
						resourceName, "load_balancer_info.0.target_group_info.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "load_balancer_info.0.target_group_info.*", map[string]string{
						"name": "acc-test-codedeploy-dep-group",
					}),
				),
			},
			{
				Config: test_config_load_balancer_info_target_group_info_delete(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDeploymentGroupExists(resourceName, &group),
					resource.TestCheckResourceAttr(
						resourceName, "load_balancer_info.#", "0"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateIdFunc: testAccDeploymentGroupImportStateIdFunc(resourceName),
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccCodeDeployDeploymentGroup_InPlaceDeploymentWithTrafficControl_create(t *testing.T) {
	var group codedeploy.DeploymentGroupInfo
	resourceName := "aws_codedeploy_deployment_group.test"

	rName := sdkacctest.RandString(5)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(t) },
		ErrorCheck:   acctest.ErrorCheck(t, codedeploy.EndpointsID),
		Providers:    acctest.Providers,
		CheckDestroy: testAccCheckDeploymentGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: test_config_in_place_deployment_with_traffic_control_create(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDeploymentGroupExists(resourceName, &group),
					resource.TestCheckResourceAttr(
						resourceName, "deployment_style.#", "1"),
					resource.TestCheckResourceAttr(
						resourceName, "deployment_style.0.deployment_option", "WITH_TRAFFIC_CONTROL"),
					resource.TestCheckResourceAttr(
						resourceName, "deployment_style.0.deployment_type", "IN_PLACE"),

					resource.TestCheckResourceAttr(
						resourceName, "load_balancer_info.#", "1"),
					resource.TestCheckResourceAttr(
						resourceName, "load_balancer_info.0.elb_info.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "load_balancer_info.0.elb_info.*", map[string]string{
						"name": "acc-test-codedeploy-dep-group",
					}),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateIdFunc: testAccDeploymentGroupImportStateIdFunc(resourceName),
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccCodeDeployDeploymentGroup_InPlaceDeploymentWithTrafficControl_update(t *testing.T) {
	var group codedeploy.DeploymentGroupInfo
	resourceName := "aws_codedeploy_deployment_group.test"

	rName := sdkacctest.RandString(5)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(t) },
		ErrorCheck:   acctest.ErrorCheck(t, codedeploy.EndpointsID),
		Providers:    acctest.Providers,
		CheckDestroy: testAccCheckDeploymentGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: test_config_in_place_deployment_with_traffic_control_create(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDeploymentGroupExists(resourceName, &group),
					resource.TestCheckResourceAttr(
						resourceName, "deployment_style.#", "1"),
					resource.TestCheckResourceAttr(
						resourceName, "deployment_style.0.deployment_option", "WITH_TRAFFIC_CONTROL"),
					resource.TestCheckResourceAttr(
						resourceName, "deployment_style.0.deployment_type", "IN_PLACE"),

					resource.TestCheckResourceAttr(
						resourceName, "load_balancer_info.#", "1"),
					resource.TestCheckResourceAttr(
						resourceName, "load_balancer_info.0.elb_info.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "load_balancer_info.0.elb_info.*", map[string]string{
						"name": "acc-test-codedeploy-dep-group",
					}),
				),
			},
			{
				Config: test_config_in_place_deployment_with_traffic_control_update(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDeploymentGroupExists(resourceName, &group),
					resource.TestCheckResourceAttr(
						resourceName, "deployment_style.#", "1"),
					resource.TestCheckResourceAttr(
						resourceName, "deployment_style.0.deployment_option", "WITH_TRAFFIC_CONTROL"),
					resource.TestCheckResourceAttr(
						resourceName, "deployment_style.0.deployment_type", "BLUE_GREEN"),

					resource.TestCheckResourceAttr(
						resourceName, "load_balancer_info.#", "1"),
					resource.TestCheckResourceAttr(
						resourceName, "load_balancer_info.0.elb_info.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "load_balancer_info.0.elb_info.*", map[string]string{
						"name": "acc-test-codedeploy-dep-group",
					}),
					resource.TestCheckResourceAttr(
						resourceName, "blue_green_deployment_config.0.deployment_ready_option.#", "1"),
					resource.TestCheckResourceAttr(
						resourceName, "blue_green_deployment_config.0.deployment_ready_option.0.action_on_timeout", "CONTINUE_DEPLOYMENT"),

					resource.TestCheckResourceAttr(
						resourceName, "blue_green_deployment_config.0.green_fleet_provisioning_option.#", "1"),
					resource.TestCheckResourceAttr(
						resourceName, "blue_green_deployment_config.0.green_fleet_provisioning_option.0.action", "DISCOVER_EXISTING"),

					resource.TestCheckResourceAttr(
						resourceName, "blue_green_deployment_config.0.terminate_blue_instances_on_deployment_success.#", "1"),
					resource.TestCheckResourceAttr(
						resourceName, "blue_green_deployment_config.0.terminate_blue_instances_on_deployment_success.0.action", "KEEP_ALIVE"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateIdFunc: testAccDeploymentGroupImportStateIdFunc(resourceName),
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccCodeDeployDeploymentGroup_BlueGreenDeployment_create(t *testing.T) {
	var group codedeploy.DeploymentGroupInfo
	resourceName := "aws_codedeploy_deployment_group.test"

	rName := sdkacctest.RandString(5)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(t) },
		ErrorCheck:   acctest.ErrorCheck(t, codedeploy.EndpointsID),
		Providers:    acctest.Providers,
		CheckDestroy: testAccCheckDeploymentGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: test_config_blue_green_deployment_config_create_with_asg(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDeploymentGroupExists(resourceName, &group),
					resource.TestCheckResourceAttr(
						resourceName, "deployment_style.#", "1"),
					resource.TestCheckResourceAttr(
						resourceName, "deployment_style.0.deployment_option", "WITH_TRAFFIC_CONTROL"),
					resource.TestCheckResourceAttr(
						resourceName, "deployment_style.0.deployment_type", "BLUE_GREEN"),

					resource.TestCheckResourceAttr(
						resourceName, "blue_green_deployment_config.#", "1"),

					resource.TestCheckResourceAttr(
						resourceName, "blue_green_deployment_config.0.deployment_ready_option.#", "1"),
					resource.TestCheckResourceAttr(
						resourceName, "blue_green_deployment_config.0.deployment_ready_option.0.action_on_timeout", "STOP_DEPLOYMENT"),
					resource.TestCheckResourceAttr(
						resourceName, "blue_green_deployment_config.0.deployment_ready_option.0.wait_time_in_minutes", "60"),

					resource.TestCheckResourceAttr(
						resourceName, "blue_green_deployment_config.0.green_fleet_provisioning_option.#", "1"),
					resource.TestCheckResourceAttr(
						resourceName, "blue_green_deployment_config.0.green_fleet_provisioning_option.0.action", "COPY_AUTO_SCALING_GROUP"),

					resource.TestCheckResourceAttr(
						resourceName, "blue_green_deployment_config.0.terminate_blue_instances_on_deployment_success.#", "1"),
					resource.TestCheckResourceAttr(
						resourceName, "blue_green_deployment_config.0.terminate_blue_instances_on_deployment_success.0.action", "TERMINATE"),
					resource.TestCheckResourceAttr(
						resourceName, "blue_green_deployment_config.0.terminate_blue_instances_on_deployment_success.0.termination_wait_time_in_minutes", "120"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateIdFunc: testAccDeploymentGroupImportStateIdFunc(resourceName),
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccCodeDeployDeploymentGroup_BlueGreenDeployment_updateWithASG(t *testing.T) {
	var group codedeploy.DeploymentGroupInfo
	resourceName := "aws_codedeploy_deployment_group.test"

	rName := sdkacctest.RandString(5)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(t) },
		ErrorCheck:   acctest.ErrorCheck(t, codedeploy.EndpointsID),
		Providers:    acctest.Providers,
		CheckDestroy: testAccCheckDeploymentGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: test_config_blue_green_deployment_config_create_with_asg(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDeploymentGroupExists(resourceName, &group),
					resource.TestCheckResourceAttr(
						resourceName, "blue_green_deployment_config.#", "1"),

					resource.TestCheckResourceAttr(
						resourceName, "blue_green_deployment_config.0.deployment_ready_option.#", "1"),
					resource.TestCheckResourceAttr(
						resourceName, "blue_green_deployment_config.0.deployment_ready_option.0.action_on_timeout", "STOP_DEPLOYMENT"),
					resource.TestCheckResourceAttr(
						resourceName, "blue_green_deployment_config.0.deployment_ready_option.0.wait_time_in_minutes", "60"),

					resource.TestCheckResourceAttr(
						resourceName, "blue_green_deployment_config.0.green_fleet_provisioning_option.#", "1"),
					resource.TestCheckResourceAttr(
						resourceName, "blue_green_deployment_config.0.green_fleet_provisioning_option.0.action", "COPY_AUTO_SCALING_GROUP"),

					resource.TestCheckResourceAttr(
						resourceName, "blue_green_deployment_config.0.terminate_blue_instances_on_deployment_success.#", "1"),
					resource.TestCheckResourceAttr(
						resourceName, "blue_green_deployment_config.0.terminate_blue_instances_on_deployment_success.0.action", "TERMINATE"),
					resource.TestCheckResourceAttr(
						resourceName, "blue_green_deployment_config.0.terminate_blue_instances_on_deployment_success.0.termination_wait_time_in_minutes", "120"),
				),
			},
			{
				Config: test_config_blue_green_deployment_config_update_with_asg(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDeploymentGroupExists(resourceName, &group),
					resource.TestCheckResourceAttr(
						resourceName, "blue_green_deployment_config.#", "1"),

					resource.TestCheckResourceAttr(
						resourceName, "blue_green_deployment_config.0.deployment_ready_option.#", "1"),
					resource.TestCheckResourceAttr(
						resourceName, "blue_green_deployment_config.0.deployment_ready_option.0.action_on_timeout", "STOP_DEPLOYMENT"),
					resource.TestCheckResourceAttr(
						resourceName, "blue_green_deployment_config.0.deployment_ready_option.0.wait_time_in_minutes", "60"),

					resource.TestCheckResourceAttr(
						resourceName, "blue_green_deployment_config.0.green_fleet_provisioning_option.#", "1"),
					resource.TestCheckResourceAttr(
						resourceName, "blue_green_deployment_config.0.green_fleet_provisioning_option.0.action", "COPY_AUTO_SCALING_GROUP"),

					resource.TestCheckResourceAttr(
						resourceName, "blue_green_deployment_config.0.terminate_blue_instances_on_deployment_success.#", "1"),
					resource.TestCheckResourceAttr(
						resourceName, "blue_green_deployment_config.0.terminate_blue_instances_on_deployment_success.0.action", "KEEP_ALIVE"),
				),
			},
		},
	})
}

func TestAccCodeDeployDeploymentGroup_BlueGreenDeployment_update(t *testing.T) {
	var group codedeploy.DeploymentGroupInfo
	resourceName := "aws_codedeploy_deployment_group.test"

	rName := sdkacctest.RandString(5)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(t) },
		ErrorCheck:   acctest.ErrorCheck(t, codedeploy.EndpointsID),
		Providers:    acctest.Providers,
		CheckDestroy: testAccCheckDeploymentGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: test_config_blue_green_deployment_config_create_no_asg(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDeploymentGroupExists(resourceName, &group),
					resource.TestCheckResourceAttr(
						resourceName, "deployment_style.#", "1"),
					resource.TestCheckResourceAttr(
						resourceName, "deployment_style.0.deployment_option", "WITH_TRAFFIC_CONTROL"),
					resource.TestCheckResourceAttr(
						resourceName, "deployment_style.0.deployment_type", "BLUE_GREEN"),

					resource.TestCheckResourceAttr(
						resourceName, "blue_green_deployment_config.#", "1"),

					resource.TestCheckResourceAttr(
						resourceName, "blue_green_deployment_config.0.deployment_ready_option.#", "1"),
					resource.TestCheckResourceAttr(
						resourceName, "blue_green_deployment_config.0.deployment_ready_option.0.action_on_timeout", "STOP_DEPLOYMENT"),
					resource.TestCheckResourceAttr(
						resourceName, "blue_green_deployment_config.0.deployment_ready_option.0.wait_time_in_minutes", "60"),

					resource.TestCheckResourceAttr(
						resourceName, "blue_green_deployment_config.0.green_fleet_provisioning_option.#", "1"),
					resource.TestCheckResourceAttr(
						resourceName, "blue_green_deployment_config.0.green_fleet_provisioning_option.0.action", "DISCOVER_EXISTING"),

					resource.TestCheckResourceAttr(
						resourceName, "blue_green_deployment_config.0.terminate_blue_instances_on_deployment_success.#", "1"),
					resource.TestCheckResourceAttr(
						resourceName, "blue_green_deployment_config.0.terminate_blue_instances_on_deployment_success.0.action", "TERMINATE"),
					resource.TestCheckResourceAttr(
						resourceName, "blue_green_deployment_config.0.terminate_blue_instances_on_deployment_success.0.termination_wait_time_in_minutes", "120"),
				),
			},
			{
				Config: test_config_blue_green_deployment_config_update_no_asg(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDeploymentGroupExists(resourceName, &group),
					resource.TestCheckResourceAttr(
						resourceName, "deployment_style.#", "1"),
					resource.TestCheckResourceAttr(
						resourceName, "deployment_style.0.deployment_option", "WITH_TRAFFIC_CONTROL"),
					resource.TestCheckResourceAttr(
						resourceName, "deployment_style.0.deployment_type", "BLUE_GREEN"),

					resource.TestCheckResourceAttr(
						resourceName, "blue_green_deployment_config.#", "1"),

					resource.TestCheckResourceAttr(
						resourceName, "blue_green_deployment_config.0.deployment_ready_option.#", "1"),
					resource.TestCheckResourceAttr(
						resourceName, "blue_green_deployment_config.0.deployment_ready_option.0.action_on_timeout", "CONTINUE_DEPLOYMENT"),

					resource.TestCheckResourceAttr(
						resourceName, "blue_green_deployment_config.0.green_fleet_provisioning_option.#", "1"),
					resource.TestCheckResourceAttr(
						resourceName, "blue_green_deployment_config.0.green_fleet_provisioning_option.0.action", "DISCOVER_EXISTING"),

					resource.TestCheckResourceAttr(
						resourceName, "blue_green_deployment_config.0.terminate_blue_instances_on_deployment_success.#", "1"),
					resource.TestCheckResourceAttr(
						resourceName, "blue_green_deployment_config.0.terminate_blue_instances_on_deployment_success.0.action", "KEEP_ALIVE"),
				),
			},
		},
	})
}

// Without "Computed: true" on blue_green_deployment_config, removing the resource
// from configuration causes an error, because the remote resource still exists.
func TestAccCodeDeployDeploymentGroup_BlueGreenDeployment_delete(t *testing.T) {
	var group codedeploy.DeploymentGroupInfo
	resourceName := "aws_codedeploy_deployment_group.test"

	rName := sdkacctest.RandString(5)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(t) },
		ErrorCheck:   acctest.ErrorCheck(t, codedeploy.EndpointsID),
		Providers:    acctest.Providers,
		CheckDestroy: testAccCheckDeploymentGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: test_config_blue_green_deployment_config_create_no_asg(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDeploymentGroupExists(resourceName, &group),
					resource.TestCheckResourceAttr(
						resourceName, "deployment_style.#", "1"),
					resource.TestCheckResourceAttr(
						resourceName, "deployment_style.0.deployment_type", "BLUE_GREEN"),

					resource.TestCheckResourceAttr(
						resourceName, "blue_green_deployment_config.#", "1"),

					resource.TestCheckResourceAttr(
						resourceName, "blue_green_deployment_config.0.deployment_ready_option.#", "1"),
					resource.TestCheckResourceAttr(
						resourceName, "blue_green_deployment_config.0.deployment_ready_option.0.action_on_timeout", "STOP_DEPLOYMENT"),
					resource.TestCheckResourceAttr(
						resourceName, "blue_green_deployment_config.0.deployment_ready_option.0.wait_time_in_minutes", "60"),

					resource.TestCheckResourceAttr(
						resourceName, "blue_green_deployment_config.0.green_fleet_provisioning_option.#", "1"),
					resource.TestCheckResourceAttr(
						resourceName, "blue_green_deployment_config.0.green_fleet_provisioning_option.0.action", "DISCOVER_EXISTING"),

					resource.TestCheckResourceAttr(
						resourceName, "blue_green_deployment_config.0.terminate_blue_instances_on_deployment_success.#", "1"),
					resource.TestCheckResourceAttr(
						resourceName, "blue_green_deployment_config.0.terminate_blue_instances_on_deployment_success.0.action", "TERMINATE"),
					resource.TestCheckResourceAttr(
						resourceName, "blue_green_deployment_config.0.terminate_blue_instances_on_deployment_success.0.termination_wait_time_in_minutes", "120"),
				),
			},
			{
				Config: test_config_blue_green_deployment_config_delete(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDeploymentGroupExists(resourceName, &group),
					resource.TestCheckResourceAttr(
						resourceName, "deployment_style.#", "1"),
					resource.TestCheckResourceAttr(
						resourceName, "deployment_style.0.deployment_type", "IN_PLACE"),

					// The state is preserved, but AWS ignores it
					resource.TestCheckResourceAttr(
						resourceName, "blue_green_deployment_config.#", "1"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateIdFunc: testAccDeploymentGroupImportStateIdFunc(resourceName),
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccCodeDeployDeploymentGroup_BlueGreenDeployment_complete(t *testing.T) {
	var group codedeploy.DeploymentGroupInfo
	resourceName := "aws_codedeploy_deployment_group.test"

	rName := sdkacctest.RandString(5)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(t) },
		ErrorCheck:   acctest.ErrorCheck(t, codedeploy.EndpointsID),
		Providers:    acctest.Providers,
		CheckDestroy: testAccCheckDeploymentGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: test_config_blue_green_deployment_complete(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDeploymentGroupExists(resourceName, &group),

					resource.TestCheckResourceAttr(
						resourceName, "deployment_style.#", "1"),
					resource.TestCheckResourceAttr(
						resourceName, "deployment_style.0.deployment_option", "WITH_TRAFFIC_CONTROL"),
					resource.TestCheckResourceAttr(
						resourceName, "deployment_style.0.deployment_type", "BLUE_GREEN"),

					resource.TestCheckResourceAttr(
						resourceName, "load_balancer_info.#", "1"),
					resource.TestCheckResourceAttr(
						resourceName, "load_balancer_info.0.elb_info.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "load_balancer_info.0.elb_info.*", map[string]string{
						"name": "acc-test-codedeploy-dep-group",
					}),
					resource.TestCheckResourceAttr(
						resourceName, "blue_green_deployment_config.#", "1"),

					resource.TestCheckResourceAttr(
						resourceName, "blue_green_deployment_config.0.deployment_ready_option.#", "1"),
					resource.TestCheckResourceAttr(
						resourceName, "blue_green_deployment_config.0.deployment_ready_option.0.action_on_timeout", "STOP_DEPLOYMENT"),
					resource.TestCheckResourceAttr(
						resourceName, "blue_green_deployment_config.0.deployment_ready_option.0.wait_time_in_minutes", "60"),

					resource.TestCheckResourceAttr(
						resourceName, "blue_green_deployment_config.0.green_fleet_provisioning_option.#", "1"),
					resource.TestCheckResourceAttr(
						resourceName, "blue_green_deployment_config.0.green_fleet_provisioning_option.0.action", "DISCOVER_EXISTING"),

					resource.TestCheckResourceAttr(
						resourceName, "blue_green_deployment_config.0.terminate_blue_instances_on_deployment_success.#", "1"),
					resource.TestCheckResourceAttr(
						resourceName, "blue_green_deployment_config.0.terminate_blue_instances_on_deployment_success.0.action", "KEEP_ALIVE"),
					resource.TestCheckResourceAttr(
						resourceName, "blue_green_deployment_config.0.terminate_blue_instances_on_deployment_success.0.termination_wait_time_in_minutes", "0"),
				),
			},
			{
				Config: test_config_blue_green_deployment_complete_updated(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDeploymentGroupExists(resourceName, &group),

					resource.TestCheckResourceAttr(
						resourceName, "deployment_style.#", "1"),
					resource.TestCheckResourceAttr(
						resourceName, "deployment_style.0.deployment_option", "WITH_TRAFFIC_CONTROL"),
					resource.TestCheckResourceAttr(
						resourceName, "deployment_style.0.deployment_type", "BLUE_GREEN"),

					resource.TestCheckResourceAttr(
						resourceName, "load_balancer_info.#", "1"),
					resource.TestCheckResourceAttr(
						resourceName, "load_balancer_info.0.elb_info.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "load_balancer_info.0.elb_info.*", map[string]string{
						"name": "acc-test-codedeploy-dep-group",
					}),
					resource.TestCheckResourceAttr(
						resourceName, "blue_green_deployment_config.#", "1"),

					resource.TestCheckResourceAttr(
						resourceName, "blue_green_deployment_config.0.deployment_ready_option.#", "1"),
					resource.TestCheckResourceAttr(
						resourceName, "blue_green_deployment_config.0.deployment_ready_option.0.action_on_timeout", "CONTINUE_DEPLOYMENT"),

					resource.TestCheckResourceAttr(
						resourceName, "blue_green_deployment_config.0.green_fleet_provisioning_option.#", "1"),
					resource.TestCheckResourceAttr(
						resourceName, "blue_green_deployment_config.0.green_fleet_provisioning_option.0.action", "DISCOVER_EXISTING"),

					resource.TestCheckResourceAttr(
						resourceName, "blue_green_deployment_config.0.terminate_blue_instances_on_deployment_success.#", "1"),
					resource.TestCheckResourceAttr(
						resourceName, "blue_green_deployment_config.0.terminate_blue_instances_on_deployment_success.0.action", "KEEP_ALIVE"),
					resource.TestCheckResourceAttr(
						resourceName, "blue_green_deployment_config.0.terminate_blue_instances_on_deployment_success.0.termination_wait_time_in_minutes", "0"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateIdFunc: testAccDeploymentGroupImportStateIdFunc(resourceName),
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccCodeDeployDeploymentGroup_ECS_blueGreen(t *testing.T) {
	var group codedeploy.DeploymentGroupInfo
	rName := fmt.Sprintf("tf-acc-test-%s", sdkacctest.RandString(5))
	ecsClusterResourceName := "aws_ecs_cluster.test"
	ecsServiceResourceName := "aws_ecs_service.test"
	lbTargetGroupBlueResourceName := "aws_lb_target_group.blue"
	lbTargetGroupGreenResourceName := "aws_lb_target_group.green"
	resourceName := "aws_codedeploy_deployment_group.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(t) },
		ErrorCheck:   acctest.ErrorCheck(t, codedeploy.EndpointsID),
		Providers:    acctest.Providers,
		CheckDestroy: testAccCheckDeploymentGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDeploymentGroupECSBlueGreenConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDeploymentGroupExists(resourceName, &group),
					resource.TestCheckResourceAttr(resourceName, "ecs_service.#", "1"),
					resource.TestCheckResourceAttrPair(resourceName, "ecs_service.0.cluster_name", ecsClusterResourceName, "name"),
					resource.TestCheckResourceAttrPair(resourceName, "ecs_service.0.service_name", ecsServiceResourceName, "name"),
					resource.TestCheckResourceAttr(resourceName, "load_balancer_info.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "load_balancer_info.0.target_group_pair_info.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "load_balancer_info.0.target_group_pair_info.0.prod_traffic_route.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "load_balancer_info.0.target_group_pair_info.0.prod_traffic_route.0.listener_arns.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "load_balancer_info.0.target_group_pair_info.0.target_group.#", "2"),
					resource.TestCheckResourceAttrPair(resourceName, "load_balancer_info.0.target_group_pair_info.0.target_group.0.name", lbTargetGroupBlueResourceName, "name"),
					resource.TestCheckResourceAttrPair(resourceName, "load_balancer_info.0.target_group_pair_info.0.target_group.1.name", lbTargetGroupGreenResourceName, "name"),
					resource.TestCheckResourceAttr(resourceName, "load_balancer_info.0.target_group_pair_info.0.test_traffic_route.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "blue_green_deployment_config.0.deployment_ready_option.0.action_on_timeout", "CONTINUE_DEPLOYMENT"),
					resource.TestCheckResourceAttr(resourceName, "blue_green_deployment_config.0.terminate_blue_instances_on_deployment_success.0.action", "TERMINATE"),
					resource.TestCheckResourceAttr(resourceName, "blue_green_deployment_config.0.terminate_blue_instances_on_deployment_success.0.termination_wait_time_in_minutes", "5"),
				),
			},
			{
				Config: testAccDeploymentGroupECSBlueGreenUpdateConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDeploymentGroupExists(resourceName, &group),
					resource.TestCheckResourceAttr(resourceName, "ecs_service.#", "1"),
					resource.TestCheckResourceAttrPair(resourceName, "ecs_service.0.cluster_name", ecsClusterResourceName, "name"),
					resource.TestCheckResourceAttrPair(resourceName, "ecs_service.0.service_name", ecsServiceResourceName, "name"),
					resource.TestCheckResourceAttr(resourceName, "load_balancer_info.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "load_balancer_info.0.target_group_pair_info.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "load_balancer_info.0.target_group_pair_info.0.prod_traffic_route.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "load_balancer_info.0.target_group_pair_info.0.prod_traffic_route.0.listener_arns.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "load_balancer_info.0.target_group_pair_info.0.target_group.#", "2"),
					resource.TestCheckResourceAttrPair(resourceName, "load_balancer_info.0.target_group_pair_info.0.target_group.0.name", lbTargetGroupBlueResourceName, "name"),
					resource.TestCheckResourceAttrPair(resourceName, "load_balancer_info.0.target_group_pair_info.0.target_group.1.name", lbTargetGroupGreenResourceName, "name"),
					resource.TestCheckResourceAttr(resourceName, "load_balancer_info.0.target_group_pair_info.0.test_traffic_route.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "blue_green_deployment_config.0.deployment_ready_option.0.action_on_timeout", "STOP_DEPLOYMENT"),
					resource.TestCheckResourceAttr(resourceName, "blue_green_deployment_config.0.deployment_ready_option.0.wait_time_in_minutes", "30"),
					resource.TestCheckResourceAttr(resourceName, "blue_green_deployment_config.0.terminate_blue_instances_on_deployment_success.0.action", "TERMINATE"),
					resource.TestCheckResourceAttr(resourceName, "blue_green_deployment_config.0.terminate_blue_instances_on_deployment_success.0.termination_wait_time_in_minutes", "60"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateIdFunc: testAccDeploymentGroupImportStateIdFunc(resourceName),
				ImportStateVerify: true,
			},
		},
	})
}

func TestDeploymentGroup_buildTriggerConfigs(t *testing.T) {
	input := []interface{}{
		map[string]interface{}{
			"trigger_events": schema.NewSet(schema.HashString, []interface{}{
				"DeploymentFailure",
			}),
			"trigger_name":       "test-trigger",
			"trigger_target_arn": "arn:aws:sns:us-west-2:123456789012:test-topic", // lintignore:AWSAT003,AWSAT005 // unit test
		},
	}

	expected := []*codedeploy.TriggerConfig{
		{
			TriggerEvents: []*string{
				aws.String("DeploymentFailure"),
			},
			TriggerName:      aws.String("test-trigger"),
			TriggerTargetArn: aws.String("arn:aws:sns:us-west-2:123456789012:test-topic"), // lintignore:AWSAT003,AWSAT005 // unit test
		},
	}

	actual := tfcodedeploy.BuildTriggerConfigs(input)

	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("tfcodedeploy.BuildTriggerConfigs output is not correct.\nGot:\n%#v\nExpected:\n%#v\n",
			actual, expected)
	}
}

func TestDeploymentGroup_triggerConfigsToMap(t *testing.T) {
	input := []*codedeploy.TriggerConfig{
		{
			TriggerEvents: []*string{
				aws.String("DeploymentFailure"),
				aws.String("InstanceFailure"),
			},
			TriggerName:      aws.String("test-trigger-2"),
			TriggerTargetArn: aws.String("arn:aws:sns:us-west-2:123456789012:test-topic-2"), // lintignore:AWSAT003,AWSAT005 // unit test
		},
	}

	expected := map[string]interface{}{
		"trigger_events": schema.NewSet(schema.HashString, []interface{}{
			"DeploymentFailure",
			"InstanceFailure",
		}),
		"trigger_name":       "test-trigger-2",
		"trigger_target_arn": "arn:aws:sns:us-west-2:123456789012:test-topic-2", // lintignore:AWSAT003,AWSAT005 // unit test
	}

	actual := tfcodedeploy.TriggerConfigsToMap(input)[0]

	fatal := false

	if actual["trigger_name"] != expected["trigger_name"] {
		fatal = true
	}

	if actual["trigger_target_arn"] != expected["trigger_target_arn"] {
		fatal = true
	}

	actualEvents := actual["trigger_events"].(*schema.Set)
	expectedEvents := expected["trigger_events"].(*schema.Set)
	if !actualEvents.Equal(expectedEvents) {
		fatal = true
	}

	if fatal {
		t.Fatalf("tfcodedeploy.TriggerConfigsToMap output is not correct.\nGot:\n%#v\nExpected:\n%#v\n",
			actual, expected)
	}
}

func TestDeploymentGroup_buildAutoRollbackConfig(t *testing.T) {
	input := []interface{}{
		map[string]interface{}{
			"events": schema.NewSet(schema.HashString, []interface{}{
				"DEPLOYMENT_FAILURE",
			}),
			"enabled": true,
		},
	}

	expected := &codedeploy.AutoRollbackConfiguration{
		Events: []*string{
			aws.String("DEPLOYMENT_FAILURE"),
		},
		Enabled: aws.Bool(true),
	}

	actual := tfcodedeploy.BuildAutoRollbackConfig(input)

	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("tfcodedeploy.BuildAutoRollbackConfig output is not correct.\nGot:\n%#v\nExpected:\n%#v\n",
			actual, expected)
	}
}

func TestDeploymentGroup_autoRollbackConfigToMap(t *testing.T) {
	input := &codedeploy.AutoRollbackConfiguration{
		Events: []*string{
			aws.String("DEPLOYMENT_FAILURE"),
			aws.String("DEPLOYMENT_STOP_ON_ALARM"),
		},
		Enabled: aws.Bool(false),
	}

	expected := map[string]interface{}{
		"events": schema.NewSet(schema.HashString, []interface{}{
			"DEPLOYMENT_FAILURE",
			"DEPLOYMENT_STOP_ON_ALARM",
		}),
		"enabled": false,
	}

	actual := tfcodedeploy.AutoRollbackConfigToMap(input)[0]

	fatal := false

	if actual["enabled"] != expected["enabled"] {
		fatal = true
	}

	actualEvents := actual["events"].(*schema.Set)
	expectedEvents := expected["events"].(*schema.Set)
	if !actualEvents.Equal(expectedEvents) {
		fatal = true
	}

	if fatal {
		t.Fatalf("tfcodedeploy.AutoRollbackConfigToMap output is not correct.\nGot:\n%#v\nExpected:\n%#v\n",
			actual, expected)
	}
}

func TestDeploymentGroup_expandDeploymentStyle(t *testing.T) {
	input := []interface{}{
		map[string]interface{}{
			"deployment_option": "WITH_TRAFFIC_CONTROL",
			"deployment_type":   "BLUE_GREEN",
		},
	}

	expected := &codedeploy.DeploymentStyle{
		DeploymentOption: aws.String("WITH_TRAFFIC_CONTROL"),
		DeploymentType:   aws.String("BLUE_GREEN"),
	}

	actual := tfcodedeploy.ExpandDeploymentStyle(input)

	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("tfcodedeploy.ExpandDeploymentStyle output is not correct.\nGot:\n%#v\nExpected:\n%#v\n",
			actual, expected)
	}
}

func TestDeploymentGroup_flattenDeploymentStyle(t *testing.T) {
	expected := map[string]interface{}{
		"deployment_option": "WITHOUT_TRAFFIC_CONTROL",
		"deployment_type":   "IN_PLACE",
	}

	input := &codedeploy.DeploymentStyle{
		DeploymentOption: aws.String("WITHOUT_TRAFFIC_CONTROL"),
		DeploymentType:   aws.String("IN_PLACE"),
	}

	actual := tfcodedeploy.FlattenDeploymentStyle(input)[0]

	fatal := false

	if actual["deployment_option"] != expected["deployment_option"] {
		fatal = true
	}

	if actual["deployment_type"] != expected["deployment_type"] {
		fatal = true
	}

	if fatal {
		t.Fatalf("tfcodedeploy.FlattenDeploymentStyle output is not correct.\nGot:\n%#v\nExpected:\n%#v\n",
			actual, expected)
	}
}

func TestDeploymentGroup_expandLoadBalancerInfo(t *testing.T) {
	testCases := []struct {
		Input    []interface{}
		Expected *codedeploy.LoadBalancerInfo
	}{
		{
			Input:    nil,
			Expected: &codedeploy.LoadBalancerInfo{},
		},
		{
			Input: []interface{}{
				map[string]interface{}{
					"elb_info": schema.NewSet(tfcodedeploy.LoadBalancerInfoHash, []interface{}{
						map[string]interface{}{
							"name": "acc-test-codedeploy-dep-group",
						},
						map[string]interface{}{
							"name": "acc-test-codedeploy-dep-group-2",
						},
					}),
				},
			},
			Expected: &codedeploy.LoadBalancerInfo{
				ElbInfoList: []*codedeploy.ELBInfo{
					{
						Name: aws.String("acc-test-codedeploy-dep-group"),
					},
					{
						Name: aws.String("acc-test-codedeploy-dep-group-2"),
					},
				},
			},
		},
	}

	for _, tc := range testCases {
		actual := tfcodedeploy.ExpandLoadBalancerInfo(tc.Input)
		if !reflect.DeepEqual(actual, tc.Expected) {
			t.Fatalf("tfcodedeploy.ExpandLoadBalancerInfo output is not correct.\nGot:\n%#v\nExpected:\n%#v\n",
				actual, tc.Expected)
		}
	}
}

func TestDeploymentGroup_flattenLoadBalancerInfo(t *testing.T) {
	input := &codedeploy.LoadBalancerInfo{
		TargetGroupInfoList: []*codedeploy.TargetGroupInfo{
			{
				Name: aws.String("abc-tg"),
			},
			{
				Name: aws.String("xyz-tg"),
			},
		},
	}

	expected := map[string]interface{}{
		"target_group_info": schema.NewSet(tfcodedeploy.LoadBalancerInfoHash, []interface{}{
			map[string]interface{}{
				"name": "abc-tg",
			},
			map[string]interface{}{
				"name": "xyz-tg",
			},
		}),
	}

	actual := tfcodedeploy.FlattenLoadBalancerInfo(input)[0].(map[string]interface{})

	fatal := false

	a := actual["target_group_info"].(*schema.Set)
	e := expected["target_group_info"].(*schema.Set)
	if !a.Equal(e) {
		fatal = true
	}

	if fatal {
		t.Fatalf("tfcodedeploy.FlattenLoadBalancerInfo output is not correct.\nGot:\n%#v\nExpected:\n%#v\n",
			actual, expected)
	}
}

func TestDeploymentGroup_expandBlueGreenDeploymentConfig(t *testing.T) {
	input := []interface{}{
		map[string]interface{}{
			"deployment_ready_option": []interface{}{
				map[string]interface{}{
					"action_on_timeout":    "CONTINUE_DEPLOYMENT",
					"wait_time_in_minutes": 60,
				},
			},

			"green_fleet_provisioning_option": []interface{}{
				map[string]interface{}{
					"action": "COPY_AUTO_SCALING_GROUP",
				},
			},

			"terminate_blue_instances_on_deployment_success": []interface{}{
				map[string]interface{}{
					"action":                           "TERMINATE",
					"termination_wait_time_in_minutes": 90,
				},
			},
		},
	}

	expected := &codedeploy.BlueGreenDeploymentConfiguration{
		DeploymentReadyOption: &codedeploy.DeploymentReadyOption{
			ActionOnTimeout:   aws.String("CONTINUE_DEPLOYMENT"),
			WaitTimeInMinutes: aws.Int64(60),
		},

		GreenFleetProvisioningOption: &codedeploy.GreenFleetProvisioningOption{
			Action: aws.String("COPY_AUTO_SCALING_GROUP"),
		},

		TerminateBlueInstancesOnDeploymentSuccess: &codedeploy.BlueInstanceTerminationOption{
			Action:                       aws.String("TERMINATE"),
			TerminationWaitTimeInMinutes: aws.Int64(90),
		},
	}

	actual := tfcodedeploy.ExpandBlueGreenDeploymentConfig(input)

	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("tfcodedeploy.ExpandBlueGreenDeploymentConfig output is not correct.\nGot:\n%#v\nExpected:\n%#v\n",
			actual, expected)
	}
}

func TestDeploymentGroup_flattenBlueGreenDeploymentConfig(t *testing.T) {
	input := &codedeploy.BlueGreenDeploymentConfiguration{
		DeploymentReadyOption: &codedeploy.DeploymentReadyOption{
			ActionOnTimeout:   aws.String("STOP_DEPLOYMENT"),
			WaitTimeInMinutes: aws.Int64(120),
		},

		GreenFleetProvisioningOption: &codedeploy.GreenFleetProvisioningOption{
			Action: aws.String("DISCOVER_EXISTING"),
		},

		TerminateBlueInstancesOnDeploymentSuccess: &codedeploy.BlueInstanceTerminationOption{
			Action:                       aws.String("KEEP_ALIVE"),
			TerminationWaitTimeInMinutes: aws.Int64(90),
		},
	}

	expected := map[string]interface{}{
		"deployment_ready_option": []map[string]interface{}{
			{
				"action_on_timeout":    "STOP_DEPLOYMENT",
				"wait_time_in_minutes": 120,
			},
		},

		"green_fleet_provisioning_option": []map[string]interface{}{
			{
				"action": "DISCOVER_EXISTING",
			},
		},

		"terminate_blue_instances_on_deployment_success": []map[string]interface{}{
			{
				"action":                           "KEEP_ALIVE",
				"termination_wait_time_in_minutes": 90,
			},
		},
	}

	actual := tfcodedeploy.FlattenBlueGreenDeploymentConfig(input)[0]

	fatal := false

	a := actual["deployment_ready_option"].([]map[string]interface{})[0]
	if a["action_on_timeout"].(string) != "STOP_DEPLOYMENT" {
		fatal = true
	}

	if a["wait_time_in_minutes"].(int64) != 120 {
		fatal = true
	}

	b := actual["green_fleet_provisioning_option"].([]map[string]interface{})[0]
	if b["action"].(string) != "DISCOVER_EXISTING" {
		fatal = true
	}

	c := actual["terminate_blue_instances_on_deployment_success"].([]map[string]interface{})[0]
	if c["action"].(string) != "KEEP_ALIVE" {
		fatal = true
	}

	if c["termination_wait_time_in_minutes"].(int64) != 90 {
		fatal = true
	}

	if fatal {
		t.Fatalf("tfcodedeploy.FlattenBlueGreenDeploymentConfig output is not correct.\nGot:\n%#v\nExpected:\n%#v\n",
			actual, expected)
	}
}

func TestDeploymentGroup_buildAlarmConfig(t *testing.T) {
	input := []interface{}{
		map[string]interface{}{
			"alarms": schema.NewSet(schema.HashString, []interface{}{
				"test-alarm",
			}),
			"enabled":                   true,
			"ignore_poll_alarm_failure": false,
		},
	}

	expected := &codedeploy.AlarmConfiguration{
		Alarms: []*codedeploy.Alarm{
			{
				Name: aws.String("test-alarm"),
			},
		},
		Enabled:                aws.Bool(true),
		IgnorePollAlarmFailure: aws.Bool(false),
	}

	actual := tfcodedeploy.BuildAlarmConfig(input)

	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("tfcodedeploy.BuildAlarmConfig output is not correct.\nGot:\n%#v\nExpected:\n%#v\n",
			actual, expected)
	}
}

func TestDeploymentGroup_alarmConfigToMap(t *testing.T) {
	input := &codedeploy.AlarmConfiguration{
		Alarms: []*codedeploy.Alarm{
			{
				Name: aws.String("test-alarm-2"),
			},
			{
				Name: aws.String("test-alarm"),
			},
		},
		Enabled:                aws.Bool(false),
		IgnorePollAlarmFailure: aws.Bool(true),
	}

	expected := map[string]interface{}{
		"alarms": schema.NewSet(schema.HashString, []interface{}{
			"test-alarm-2",
			"test-alarm",
		}),
		"enabled":                   false,
		"ignore_poll_alarm_failure": true,
	}

	actual := tfcodedeploy.AlarmConfigToMap(input)[0]

	fatal := false

	if actual["enabled"] != expected["enabled"] {
		fatal = true
	}

	if actual["ignore_poll_alarm_failure"] != expected["ignore_poll_alarm_failure"] {
		fatal = true
	}

	actualAlarms := actual["alarms"].(*schema.Set)
	expectedAlarms := expected["alarms"].(*schema.Set)
	if !actualAlarms.Equal(expectedAlarms) {
		fatal = true
	}

	if fatal {
		t.Fatalf("tfcodedeploy.AlarmConfigToMap output is not correct.\nGot:\n%#v\nExpected:\n%#v\n",
			actual, expected)
	}
}

func testAccCheckCodeDeployDeploymentGroupTriggerEvents(group *codedeploy.DeploymentGroupInfo, triggerName string, expectedEvents []string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		found := false
		for _, actual := range group.TriggerConfigurations {
			if *actual.TriggerName == triggerName {
				found = true

				numberOfEvents := len(actual.TriggerEvents)
				if numberOfEvents != len(expectedEvents) {
					return fmt.Errorf("Trigger events do not match. Expected: %d. Got: %d.",
						len(expectedEvents), numberOfEvents)
				}

				actualEvents := make([]string, 0, numberOfEvents)
				for _, event := range actual.TriggerEvents {
					actualEvents = append(actualEvents, *event)
				}
				sort.Strings(actualEvents)

				if !reflect.DeepEqual(actualEvents, expectedEvents) {
					return fmt.Errorf("Trigger events do not match.\nExpected: %v\nGot: %v\n",
						expectedEvents, actualEvents)
				}
				break
			}
		}
		if found {
			return nil
		} else {
			return fmt.Errorf("trigger configuration %q not found", triggerName)
		}
	}
}

func testAccCheckCodeDeployDeploymentGroupTriggerTargetArn(group *codedeploy.DeploymentGroupInfo, triggerName string, r *regexp.Regexp) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		found := false
		for _, actual := range group.TriggerConfigurations {
			if *actual.TriggerName == triggerName {
				found = true
				if !r.MatchString(*actual.TriggerTargetArn) {
					return fmt.Errorf("Trigger target arn does not match regular expression.\nRegex: %v\nTriggerTargetArn: %v\n",
						r, *actual.TriggerTargetArn)
				}
				break
			}
		}
		if found {
			return nil
		} else {
			return fmt.Errorf("trigger configuration %q not found", triggerName)
		}
	}
}

func testAccCheckDeploymentGroupDestroy(s *terraform.State) error {
	conn := acctest.Provider.Meta().(*conns.AWSClient).CodeDeployConn

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aws_codedeploy_deployment_group" {
			continue
		}

		resp, err := conn.GetDeploymentGroup(&codedeploy.GetDeploymentGroupInput{
			ApplicationName:     aws.String(rs.Primary.Attributes["app_name"]),
			DeploymentGroupName: aws.String(rs.Primary.Attributes["deployment_group_name"]),
		})

		if ae, ok := err.(awserr.Error); ok && ae.Code() == "ApplicationDoesNotExistException" {
			continue
		}

		if err == nil {
			if resp.DeploymentGroupInfo.DeploymentGroupName != nil {
				return fmt.Errorf("CodeDeploy deployment group still exists:\n%#v", *resp.DeploymentGroupInfo.DeploymentGroupName)
			}
		}

		return err
	}

	return nil
}

func testAccCheckDeploymentGroupExists(name string, group *codedeploy.DeploymentGroupInfo) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		conn := acctest.Provider.Meta().(*conns.AWSClient).CodeDeployConn

		resp, err := conn.GetDeploymentGroup(&codedeploy.GetDeploymentGroupInput{
			ApplicationName:     aws.String(rs.Primary.Attributes["app_name"]),
			DeploymentGroupName: aws.String(rs.Primary.Attributes["deployment_group_name"]),
		})

		if err != nil {
			return err
		}

		*group = *resp.DeploymentGroupInfo

		return nil
	}
}

func testAccDeploymentGroupImportStateIdFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("Not found: %s", resourceName)
		}

		return fmt.Sprintf("%s:%s", rs.Primary.Attributes["app_name"], rs.Primary.Attributes["deployment_group_name"]), nil
	}
}

func testAccDeploymentGroupBaseConfig(rName string) string {
	return fmt.Sprintf(`
resource "aws_codedeploy_app" "test" {
  name = "tf-acc-test-%[1]s"
}

resource "aws_iam_role_policy" "test" {
  name = "tf-acc-test-%[1]s"
  role = aws_iam_role.test.id

  policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": [
        "autoscaling:CompleteLifecycleAction",
        "autoscaling:DeleteLifecycleHook",
        "autoscaling:DescribeAutoScalingGroups",
        "autoscaling:DescribeLifecycleHooks",
        "autoscaling:PutLifecycleHook",
        "autoscaling:RecordLifecycleActionHeartbeat",
        "ec2:DescribeInstances",
        "ec2:DescribeInstanceStatus",
        "tag:GetTags",
        "tag:GetResources"
      ],
      "Resource": "*"
    }
  ]
}
EOF
}

data "aws_partition" "current" {}

resource "aws_iam_role" "test" {
  name = "tf-acc-test-%[1]s"

  assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Sid": "",
      "Effect": "Allow",
      "Principal": {
        "Service": [
          "codedeploy.${data.aws_partition.current.dns_suffix}"
        ]
      },
      "Action": "sts:AssumeRole"
    }
  ]
}
EOF
}
`, rName)
}

func testAccDeploymentGroup(rName string, tagGroup bool) string {
	var tagGroupOrFilter string
	if tagGroup {
		tagGroupOrFilter = `
ec2_tag_set {
  ec2_tag_filter {
    key   = "filterkey"
    type  = "KEY_AND_VALUE"
    value = "filtervalue"
  }
}
`

	} else {
		tagGroupOrFilter = `
ec2_tag_filter {
  key   = "filterkey"
  type  = "KEY_AND_VALUE"
  value = "filtervalue"
}
`

	}

	return testAccDeploymentGroupBaseConfig(rName) + fmt.Sprintf(`
resource "aws_codedeploy_deployment_group" "test" {
  app_name              = aws_codedeploy_app.test.name
  deployment_group_name = "tf-acc-test-%[1]s"
  service_role_arn      = aws_iam_role.test.arn
  %[2]s
}
`, rName, tagGroupOrFilter)
}

func testAccDeploymentGroupModified(rName string, tagGroup bool) string {
	var tagGroupOrFilter string
	if tagGroup {
		tagGroupOrFilter = `
ec2_tag_set {
  ec2_tag_filter {
    key   = "filterkey"
    type  = "KEY_AND_VALUE"
    value = "anotherfiltervalue"
  }
}
`

	} else {
		tagGroupOrFilter = `
ec2_tag_filter {
  key   = "filterkey"
  type  = "KEY_AND_VALUE"
  value = "anotherfiltervalue"
}
`

	}

	return testAccDeploymentGroupBaseConfig(rName) + fmt.Sprintf(`
resource "aws_codedeploy_deployment_group" "test" {
  app_name              = aws_codedeploy_app.test.name
  deployment_group_name = "tf-acc-test-updated-%[1]s"
  service_role_arn      = aws_iam_role.test_updated.arn
  %[2]s
}

resource "aws_iam_role_policy" "test_updated" {
  name = "tf-acc-test-%[1]s"
  role = aws_iam_role.test_updated.id

  policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": [
        "autoscaling:CompleteLifecycleAction",
        "autoscaling:DeleteLifecycleHook",
        "autoscaling:DescribeAutoScalingGroups",
        "autoscaling:DescribeLifecycleHooks",
        "autoscaling:PutLifecycleHook",
        "autoscaling:RecordLifecycleActionHeartbeat",
        "ec2:DescribeInstances",
        "ec2:DescribeInstanceStatus",
        "tag:GetTags",
        "tag:GetResources"
      ],
      "Resource": "*"
    }
  ]
}
EOF
}

resource "aws_iam_role" "test_updated" {
  name = "tf-acc-test-updated-%[1]s"

  assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Sid": "",
      "Effect": "Allow",
      "Principal": {
        "Service": [
          "codedeploy.${data.aws_partition.current.dns_suffix}"
        ]
      },
      "Action": "sts:AssumeRole"
    }
  ]
}
EOF
}
`, rName, tagGroupOrFilter)
}

func testAccDeploymentGroupOnPremiseTags(rName string, tagGroup bool) string {
	var tagGroupOrFilter string
	if tagGroup {
		tagGroupOrFilter = `
on_premises_instance_tag_set {
  on_premises_instance_tag_filter {
    key   = "filterkey"
    type  = "KEY_AND_VALUE"
    value = "filtervalue"
  }
}
`
	} else {
		tagGroupOrFilter = `
on_premises_instance_tag_filter {
  key   = "filterkey"
  type  = "KEY_AND_VALUE"
  value = "filtervalue"
}
`
	}
	return testAccDeploymentGroupBaseConfig(rName) + fmt.Sprintf(`
resource "aws_codedeploy_deployment_group" "test" {
  app_name               = aws_codedeploy_app.test.name
  deployment_group_name  = "tf-acc-test-%[1]s"
  deployment_config_name = "CodeDeployDefault.OneAtATime"
  service_role_arn       = aws_iam_role.test.arn

  %[2]s
}
`, rName, tagGroupOrFilter)
}

func baseCodeDeployConfig(rName string) string {
	return fmt.Sprintf(`
resource "aws_codedeploy_app" "test" {
  name = "tf-acc-test-%[1]s"
}

resource "aws_iam_role_policy" "test" {
  name = "tf-acc-test-%[1]s"
  role = aws_iam_role.test.id

  policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": [
        "autoscaling:CompleteLifecycleAction",
        "autoscaling:DeleteLifecycleHook",
        "autoscaling:DescribeAutoScalingGroups",
        "autoscaling:DescribeLifecycleHooks",
        "autoscaling:PutLifecycleHook",
        "autoscaling:RecordLifecycleActionHeartbeat",
        "codedeploy:*",
        "ec2:DescribeInstances",
        "ec2:DescribeInstanceStatus",
        "tag:GetTags",
        "tag:GetResources",
        "sns:Publish"
      ],
      "Resource": "*"
    }
  ]
}
EOF
}

data "aws_partition" "current" {}

resource "aws_iam_role" "test" {
  name = "tf-acc-test-%[1]s"

  assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Sid": "",
      "Effect": "Allow",
      "Principal": {
        "Service": "codedeploy.${data.aws_partition.current.dns_suffix}"
      },
      "Action": "sts:AssumeRole"
    }
  ]
}
EOF
}

resource "aws_sns_topic" "test" {
  name = "tf-acc-test-%[1]s"
}
`, rName)
}

func testAccDeploymentGroup_triggerConfiguration_create(rName string) string {
	return fmt.Sprintf(`
resource "aws_codedeploy_deployment_group" "test" {
  app_name              = aws_codedeploy_app.test.name
  deployment_group_name = "tf-acc-test-%[1]s"
  service_role_arn      = aws_iam_role.test.arn

  trigger_configuration {
    trigger_events     = ["DeploymentFailure"]
    trigger_name       = "test-trigger"
    trigger_target_arn = aws_sns_topic.test.arn
  }
}
`, rName) + baseCodeDeployConfig(rName)
}

func testAccDeploymentGroup_triggerConfiguration_update(rName string) string {
	return fmt.Sprintf(`
resource "aws_codedeploy_deployment_group" "test" {
  app_name              = aws_codedeploy_app.test.name
  deployment_group_name = "tf-acc-test-%[1]s"
  service_role_arn      = aws_iam_role.test.arn

  trigger_configuration {
    trigger_events     = ["DeploymentSuccess", "DeploymentFailure"]
    trigger_name       = "test-trigger"
    trigger_target_arn = aws_sns_topic.test.arn
  }
}
`, rName) + baseCodeDeployConfig(rName)
}

func testAccDeploymentGroup_triggerConfiguration_createMultiple(rName string) string {
	return fmt.Sprintf(`
resource "aws_codedeploy_deployment_group" "test" {
  app_name              = aws_codedeploy_app.test.name
  deployment_group_name = "tf-acc-test-%[1]s"
  service_role_arn      = aws_iam_role.test.arn

  trigger_configuration {
    trigger_events     = ["DeploymentFailure"]
    trigger_name       = "test-trigger-1"
    trigger_target_arn = aws_sns_topic.test.arn
  }

  trigger_configuration {
    trigger_events     = ["InstanceFailure"]
    trigger_name       = "test-trigger-2"
    trigger_target_arn = aws_sns_topic.test_2.arn
  }
}

resource "aws_sns_topic" "test_2" {
  name = "tf-acc-test-2-%[1]s"
}
`, rName) + baseCodeDeployConfig(rName)
}

func testAccDeploymentGroup_triggerConfiguration_updateMultiple(rName string) string {
	return fmt.Sprintf(`
resource "aws_codedeploy_deployment_group" "test" {
  app_name              = aws_codedeploy_app.test.name
  deployment_group_name = "tf-acc-test-%[1]s"
  service_role_arn      = aws_iam_role.test.arn

  trigger_configuration {
    trigger_events     = ["DeploymentStart", "DeploymentSuccess", "DeploymentFailure", "DeploymentStop"]
    trigger_name       = "test-trigger-1"
    trigger_target_arn = aws_sns_topic.test.arn
  }

  trigger_configuration {
    trigger_events     = ["InstanceFailure"]
    trigger_name       = "test-trigger-2"
    trigger_target_arn = aws_sns_topic.test_3.arn
  }
}

resource "aws_sns_topic" "test_2" {
  name = "tf-acc-test-2-%[1]s"
}

resource "aws_sns_topic" "test_3" {
  name = "tf-acc-test-3-%[1]s"
}
`, rName) + baseCodeDeployConfig(rName)
}

func test_config_auto_rollback_configuration_create(rName string) string {
	return fmt.Sprintf(`
resource "aws_codedeploy_deployment_group" "test" {
  app_name              = aws_codedeploy_app.test.name
  deployment_group_name = "tf-acc-test-%[1]s"
  service_role_arn      = aws_iam_role.test.arn

  auto_rollback_configuration {
    enabled = true
    events  = ["DEPLOYMENT_FAILURE"]
  }
}
`, rName) + baseCodeDeployConfig(rName)
}

func test_config_auto_rollback_configuration_update(rName string) string {
	return fmt.Sprintf(`
resource "aws_codedeploy_deployment_group" "test" {
  app_name              = aws_codedeploy_app.test.name
  deployment_group_name = "tf-acc-test-%[1]s"
  service_role_arn      = aws_iam_role.test.arn

  auto_rollback_configuration {
    enabled = true
    events  = ["DEPLOYMENT_FAILURE", "DEPLOYMENT_STOP_ON_ALARM"]
  }
}
`, rName) + baseCodeDeployConfig(rName)
}

func test_config_auto_rollback_configuration_none(rName string) string {
	return fmt.Sprintf(`
resource "aws_codedeploy_deployment_group" "test" {
  app_name              = aws_codedeploy_app.test.name
  deployment_group_name = "tf-acc-test-%[1]s"
  service_role_arn      = aws_iam_role.test.arn
}
`, rName) + baseCodeDeployConfig(rName)
}

func test_config_auto_rollback_configuration_disable(rName string) string {
	return fmt.Sprintf(`
resource "aws_codedeploy_deployment_group" "test" {
  app_name              = aws_codedeploy_app.test.name
  deployment_group_name = "tf-acc-test-%[1]s"
  service_role_arn      = aws_iam_role.test.arn

  auto_rollback_configuration {
    enabled = false
    events  = ["DEPLOYMENT_FAILURE"]
  }
}
`, rName) + baseCodeDeployConfig(rName)
}

func test_config_alarm_configuration_create(rName string) string {
	return fmt.Sprintf(`
resource "aws_codedeploy_deployment_group" "test" {
  app_name              = aws_codedeploy_app.test.name
  deployment_group_name = "tf-acc-test-%[1]s"
  service_role_arn      = aws_iam_role.test.arn

  alarm_configuration {
    alarms  = ["test-alarm"]
    enabled = true
  }
}
`, rName) + baseCodeDeployConfig(rName)
}

func test_config_alarm_configuration_update(rName string) string {
	return fmt.Sprintf(`
resource "aws_codedeploy_deployment_group" "test" {
  app_name              = aws_codedeploy_app.test.name
  deployment_group_name = "tf-acc-test-%[1]s"
  service_role_arn      = aws_iam_role.test.arn

  alarm_configuration {
    alarms                    = ["test-alarm", "test-alarm-2"]
    enabled                   = true
    ignore_poll_alarm_failure = true
  }
}
`, rName) + baseCodeDeployConfig(rName)
}

func test_config_alarm_configuration_none(rName string) string {
	return fmt.Sprintf(`
resource "aws_codedeploy_deployment_group" "test" {
  app_name              = aws_codedeploy_app.test.name
  deployment_group_name = "tf-acc-test-%[1]s"
  service_role_arn      = aws_iam_role.test.arn
}
`, rName) + baseCodeDeployConfig(rName)
}

func test_config_alarm_configuration_disable(rName string) string {
	return fmt.Sprintf(`
resource "aws_codedeploy_deployment_group" "test" {
  app_name              = aws_codedeploy_app.test.name
  deployment_group_name = "tf-acc-test-%[1]s"
  service_role_arn      = aws_iam_role.test.arn

  alarm_configuration {
    alarms  = ["test-alarm"]
    enabled = false
  }
}
`, rName) + baseCodeDeployConfig(rName)
}

func test_config_deployment_style_default(rName string) string {
	return fmt.Sprintf(`
resource "aws_codedeploy_deployment_group" "test" {
  app_name              = aws_codedeploy_app.test.name
  deployment_group_name = "tf-acc-test-%[1]s"
  service_role_arn      = aws_iam_role.test.arn
}
`, rName) + baseCodeDeployConfig(rName)
}

func test_config_deployment_style_create(rName string) string {
	return fmt.Sprintf(`
resource "aws_codedeploy_deployment_group" "test" {
  app_name              = aws_codedeploy_app.test.name
  deployment_group_name = "tf-acc-test-%[1]s"
  service_role_arn      = aws_iam_role.test.arn

  deployment_style {
    deployment_option = "WITH_TRAFFIC_CONTROL"
    deployment_type   = "BLUE_GREEN"
  }

  load_balancer_info {
    elb_info {
      name = "acc-test-codedeploy-dep-group"
    }
  }
}
`, rName) + baseCodeDeployConfig(rName)
}

func test_config_deployment_style_update(rName string) string {
	return fmt.Sprintf(`
resource "aws_codedeploy_deployment_group" "test" {
  app_name              = aws_codedeploy_app.test.name
  deployment_group_name = "tf-acc-test-%[1]s"
  service_role_arn      = aws_iam_role.test.arn

  deployment_style {
    deployment_option = "WITHOUT_TRAFFIC_CONTROL"
    deployment_type   = "IN_PLACE"
  }
}
`, rName) + baseCodeDeployConfig(rName)
}

func test_config_load_balancer_info_none(rName string) string {
	return fmt.Sprintf(`
resource "aws_codedeploy_deployment_group" "test" {
  app_name              = aws_codedeploy_app.test.name
  deployment_group_name = "tf-acc-test-%[1]s"
  service_role_arn      = aws_iam_role.test.arn
}
`, rName) + baseCodeDeployConfig(rName)
}

func test_config_load_balancer_info_create(rName string) string {
	return fmt.Sprintf(`
resource "aws_codedeploy_deployment_group" "test" {
  app_name              = aws_codedeploy_app.test.name
  deployment_group_name = "tf-acc-test-%[1]s"
  service_role_arn      = aws_iam_role.test.arn

  load_balancer_info {
    elb_info {
      name = "acc-test-codedeploy-dep-group"
    }
  }
}
`, rName) + baseCodeDeployConfig(rName)
}

func test_config_load_balancer_info_update(rName string) string {
	return fmt.Sprintf(`
resource "aws_codedeploy_deployment_group" "test" {
  app_name              = aws_codedeploy_app.test.name
  deployment_group_name = "tf-acc-test-%[1]s"
  service_role_arn      = aws_iam_role.test.arn

  load_balancer_info {
    elb_info {
      name = "acc-test-codedeploy-dep-group-2"
    }
  }
}
`, rName) + baseCodeDeployConfig(rName)
}

func test_config_load_balancer_info_target_group_info_create(rName string) string {
	return fmt.Sprintf(`
resource "aws_codedeploy_deployment_group" "test" {
  app_name              = aws_codedeploy_app.test.name
  deployment_group_name = "tf-acc-test-%[1]s"
  service_role_arn      = aws_iam_role.test.arn

  load_balancer_info {
    target_group_info {
      name = "acc-test-codedeploy-dep-group"
    }
  }
}
`, rName) + baseCodeDeployConfig(rName)
}

func test_config_load_balancer_info_target_group_info_update(rName string) string {
	return fmt.Sprintf(`
resource "aws_codedeploy_deployment_group" "test" {
  app_name              = aws_codedeploy_app.test.name
  deployment_group_name = "tf-acc-test-%[1]s"
  service_role_arn      = aws_iam_role.test.arn

  load_balancer_info {
    target_group_info {
      name = "acc-test-codedeploy-dep-group-2"
    }
  }
}
`, rName) + baseCodeDeployConfig(rName)
}

func test_config_load_balancer_info_target_group_info_delete(rName string) string {
	return fmt.Sprintf(`
resource "aws_codedeploy_deployment_group" "test" {
  app_name              = aws_codedeploy_app.test.name
  deployment_group_name = "tf-acc-test-%[1]s"
  service_role_arn      = aws_iam_role.test.arn
}
`, rName) + baseCodeDeployConfig(rName)
}

func test_config_in_place_deployment_with_traffic_control_create(rName string) string {
	return fmt.Sprintf(`
resource "aws_codedeploy_deployment_group" "test" {
  app_name              = aws_codedeploy_app.test.name
  deployment_group_name = "tf-acc-test-%[1]s"
  service_role_arn      = aws_iam_role.test.arn

  deployment_style {
    deployment_option = "WITH_TRAFFIC_CONTROL"
    deployment_type   = "IN_PLACE"
  }

  load_balancer_info {
    elb_info {
      name = "acc-test-codedeploy-dep-group"
    }
  }
}
`, rName) + baseCodeDeployConfig(rName)
}

func test_config_in_place_deployment_with_traffic_control_update(rName string) string {
	return fmt.Sprintf(`
resource "aws_codedeploy_deployment_group" "test" {
  app_name              = aws_codedeploy_app.test.name
  deployment_group_name = "tf-acc-test-%[1]s"
  service_role_arn      = aws_iam_role.test.arn

  deployment_style {
    deployment_option = "WITH_TRAFFIC_CONTROL"
    deployment_type   = "BLUE_GREEN"
  }

  load_balancer_info {
    elb_info {
      name = "acc-test-codedeploy-dep-group"
    }
  }

  blue_green_deployment_config {
    deployment_ready_option {
      action_on_timeout = "CONTINUE_DEPLOYMENT"
    }

    green_fleet_provisioning_option {
      action = "DISCOVER_EXISTING"
    }

    terminate_blue_instances_on_deployment_success {
      action = "KEEP_ALIVE"
    }
  }
}
`, rName) + baseCodeDeployConfig(rName)
}

func test_config_blue_green_deployment_config_delete(rName string) string {
	return fmt.Sprintf(`
resource "aws_codedeploy_deployment_group" "test" {
  app_name              = aws_codedeploy_app.test.name
  deployment_group_name = "tf-acc-test-%[1]s"
  service_role_arn      = aws_iam_role.test.arn
}
`, rName) + baseCodeDeployConfig(rName)
}

func test_config_blue_green_deployment_config_create_with_asg(rName string) string {
	return fmt.Sprintf(`
resource "aws_codedeploy_deployment_group" "test" {
  app_name              = aws_codedeploy_app.test.name
  deployment_group_name = "tf-acc-test-%[1]s"
  service_role_arn      = aws_iam_role.test.arn

  deployment_style {
    deployment_option = "WITH_TRAFFIC_CONTROL"
    deployment_type   = "BLUE_GREEN"
  }

  autoscaling_groups = [aws_autoscaling_group.test.name]

  load_balancer_info {
    elb_info {
      name = "acc-test-codedeploy-dep-group"
    }
  }

  blue_green_deployment_config {
    deployment_ready_option {
      action_on_timeout    = "STOP_DEPLOYMENT"
      wait_time_in_minutes = 60
    }

    green_fleet_provisioning_option {
      action = "COPY_AUTO_SCALING_GROUP"
    }

    terminate_blue_instances_on_deployment_success {
      action                           = "TERMINATE"
      termination_wait_time_in_minutes = 120
    }
  }
}

data "aws_ami" "amzn-ami-minimal-hvm-ebs" {
  most_recent = true
  owners      = ["amazon"]

  filter {
    name   = "name"
    values = ["amzn-ami-minimal-hvm-*"]
  }

  filter {
    name   = "root-device-type"
    values = ["ebs"]
  }
}

data "aws_availability_zones" "available" {
  state = "available"

  filter {
    name   = "opt-in-status"
    values = ["opt-in-not-required"]
  }
}

data "aws_subnet" "test" {
  availability_zone = data.aws_availability_zones.available.names[0]
  default_for_az    = "true"
}

resource "aws_launch_configuration" "test" {
  image_id      = data.aws_ami.amzn-ami-minimal-hvm-ebs.id
  instance_type = "t2.micro"
  name_prefix   = "tf-acc-test-codedeploy-deployment-group-"

  lifecycle {
    create_before_destroy = true
  }
}

resource "aws_autoscaling_group" "test" {
  name             = "tf-acc-test-codedeploy-deployment-group-%[1]s"
  max_size         = 2
  min_size         = 0
  desired_capacity = 1

  vpc_zone_identifier = [data.aws_subnet.test.id]

  launch_configuration = aws_launch_configuration.test.name

  lifecycle {
    create_before_destroy = true
  }
}
`, rName) + baseCodeDeployConfig(rName)
}

func test_config_blue_green_deployment_config_update_with_asg(rName string) string {
	return fmt.Sprintf(`
resource "aws_codedeploy_deployment_group" "test" {
  app_name              = aws_codedeploy_app.test.name
  deployment_group_name = "tf-acc-test-%[1]s"
  service_role_arn      = aws_iam_role.test.arn

  deployment_style {
    deployment_option = "WITH_TRAFFIC_CONTROL"
    deployment_type   = "BLUE_GREEN"
  }

  autoscaling_groups = [aws_autoscaling_group.test.name]

  load_balancer_info {
    elb_info {
      name = "acc-test-codedeploy-dep-group"
    }
  }

  blue_green_deployment_config {
    deployment_ready_option {
      action_on_timeout    = "STOP_DEPLOYMENT"
      wait_time_in_minutes = 60
    }

    green_fleet_provisioning_option {
      action = "COPY_AUTO_SCALING_GROUP"
    }

    terminate_blue_instances_on_deployment_success {
      action = "KEEP_ALIVE"
    }
  }
}

data "aws_ami" "amzn-ami-minimal-hvm-ebs" {
  most_recent = true
  owners      = ["amazon"]

  filter {
    name   = "name"
    values = ["amzn-ami-minimal-hvm-*"]
  }

  filter {
    name   = "root-device-type"
    values = ["ebs"]
  }
}

data "aws_availability_zones" "available" {
  state = "available"

  filter {
    name   = "opt-in-status"
    values = ["opt-in-not-required"]
  }
}

data "aws_subnet" "test" {
  availability_zone = data.aws_availability_zones.available.names[0]
  default_for_az    = "true"
}

resource "aws_launch_configuration" "test" {
  image_id      = data.aws_ami.amzn-ami-minimal-hvm-ebs.id
  instance_type = "t2.micro"
  name_prefix   = "tf-acc-test-codedeploy-deployment-group-"

  lifecycle {
    create_before_destroy = true
  }
}

resource "aws_autoscaling_group" "test" {
  name             = "tf-acc-test-codedeploy-deployment-group-%[1]s"
  max_size         = 2
  min_size         = 0
  desired_capacity = 1

  vpc_zone_identifier = [data.aws_subnet.test.id]

  launch_configuration = aws_launch_configuration.test.name

  lifecycle {
    create_before_destroy = true
  }
}
`, rName) + baseCodeDeployConfig(rName)
}

func test_config_blue_green_deployment_config_create_no_asg(rName string) string {
	return fmt.Sprintf(`
resource "aws_codedeploy_deployment_group" "test" {
  app_name              = aws_codedeploy_app.test.name
  deployment_group_name = "tf-acc-test-%[1]s"
  service_role_arn      = aws_iam_role.test.arn

  deployment_style {
    deployment_option = "WITH_TRAFFIC_CONTROL"
    deployment_type   = "BLUE_GREEN"
  }

  load_balancer_info {
    elb_info {
      name = "acc-test-codedeploy-dep-group"
    }
  }

  blue_green_deployment_config {
    deployment_ready_option {
      action_on_timeout    = "STOP_DEPLOYMENT"
      wait_time_in_minutes = 60
    }

    green_fleet_provisioning_option {
      action = "DISCOVER_EXISTING"
    }

    terminate_blue_instances_on_deployment_success {
      action                           = "TERMINATE"
      termination_wait_time_in_minutes = 120
    }
  }
}
`, rName) + baseCodeDeployConfig(rName)
}

func test_config_blue_green_deployment_config_update_no_asg(rName string) string {
	return fmt.Sprintf(`
resource "aws_codedeploy_deployment_group" "test" {
  app_name              = aws_codedeploy_app.test.name
  deployment_group_name = "tf-acc-test-%[1]s"
  service_role_arn      = aws_iam_role.test.arn

  deployment_style {
    deployment_option = "WITH_TRAFFIC_CONTROL"
    deployment_type   = "BLUE_GREEN"
  }

  load_balancer_info {
    elb_info {
      name = "acc-test-codedeploy-dep-group"
    }
  }

  blue_green_deployment_config {
    deployment_ready_option {
      action_on_timeout = "CONTINUE_DEPLOYMENT"
    }

    green_fleet_provisioning_option {
      action = "DISCOVER_EXISTING"
    }

    terminate_blue_instances_on_deployment_success {
      action = "KEEP_ALIVE"
    }
  }
}
`, rName) + baseCodeDeployConfig(rName)
}

func test_config_blue_green_deployment_complete(rName string) string {
	return fmt.Sprintf(`
resource "aws_codedeploy_deployment_group" "test" {
  app_name              = aws_codedeploy_app.test.name
  deployment_group_name = "tf-acc-test-%[1]s"
  service_role_arn      = aws_iam_role.test.arn

  deployment_style {
    deployment_option = "WITH_TRAFFIC_CONTROL"
    deployment_type   = "BLUE_GREEN"
  }

  load_balancer_info {
    elb_info {
      name = "acc-test-codedeploy-dep-group"
    }
  }

  blue_green_deployment_config {
    deployment_ready_option {
      action_on_timeout    = "STOP_DEPLOYMENT"
      wait_time_in_minutes = 60
    }

    green_fleet_provisioning_option {
      action = "DISCOVER_EXISTING"
    }

    terminate_blue_instances_on_deployment_success {
      action = "KEEP_ALIVE"
    }
  }
}
`, rName) + baseCodeDeployConfig(rName)
}

func test_config_blue_green_deployment_complete_updated(rName string) string {
	return fmt.Sprintf(`
resource "aws_codedeploy_deployment_group" "test" {
  app_name              = aws_codedeploy_app.test.name
  deployment_group_name = "tf-acc-test-%[1]s"
  service_role_arn      = aws_iam_role.test.arn

  deployment_style {
    deployment_option = "WITH_TRAFFIC_CONTROL"
    deployment_type   = "BLUE_GREEN"
  }

  load_balancer_info {
    elb_info {
      name = "acc-test-codedeploy-dep-group"
    }
  }

  blue_green_deployment_config {
    deployment_ready_option {
      action_on_timeout = "CONTINUE_DEPLOYMENT"
    }

    green_fleet_provisioning_option {
      action = "DISCOVER_EXISTING"
    }

    terminate_blue_instances_on_deployment_success {
      action = "KEEP_ALIVE"
    }
  }
}
`, rName) + baseCodeDeployConfig(rName)
}

func testAccDeploymentGroupECSBaseConfig(rName string) string {
	return fmt.Sprintf(`
data "aws_availability_zones" "available" {
  state = "available"

  filter {
    name   = "opt-in-status"
    values = ["opt-in-not-required"]
  }
}

resource "aws_vpc" "test" {
  cidr_block = "10.0.0.0/16"

  tags = {
    Name = "tf-acc-test-codedeploy-deployment-group-ecs"
  }
}

resource "aws_subnet" "test" {
  count = 2

  availability_zone = data.aws_availability_zones.available.names[count.index]
  cidr_block        = cidrsubnet(aws_vpc.test.cidr_block, 8, count.index)
  vpc_id            = aws_vpc.test.id

  tags = {
    Name = "tf-acc-test-codedeploy-deployment-group-ecs"
  }
}

resource "aws_security_group" "test" {
  name   = %[1]q
  vpc_id = aws_vpc.test.id

  ingress {
    cidr_blocks = [aws_vpc.test.cidr_block]
    from_port   = 80
    protocol    = "6"
    to_port     = 8000
  }
}

resource "aws_lb_target_group" "blue" {
  name        = "${aws_lb.test.name}-blue"
  port        = 80
  protocol    = "HTTP"
  target_type = "ip"
  vpc_id      = aws_vpc.test.id
}

resource "aws_lb_target_group" "green" {
  name        = "${aws_lb.test.name}-green"
  port        = 80
  protocol    = "HTTP"
  target_type = "ip"
  vpc_id      = aws_vpc.test.id
}

resource "aws_lb" "test" {
  internal = true
  name     = %[1]q
  subnets  = aws_subnet.test[*].id
}

resource "aws_lb_listener" "test" {
  load_balancer_arn = aws_lb.test.arn
  port              = "80"
  protocol          = "HTTP"

  default_action {
    target_group_arn = aws_lb_target_group.blue.arn
    type             = "forward"
  }
}

resource "aws_ecs_cluster" "test" {
  name = %[1]q
}

resource "aws_ecs_task_definition" "test" {
  cpu                      = "256"
  family                   = %[1]q
  memory                   = "512"
  network_mode             = "awsvpc"
  requires_compatibilities = ["FARGATE"]

  container_definitions = <<DEFINITION
[
  {
    "cpu": 256,
    "essential": true,
    "image": "mongo:latest",
    "memory": 512,
    "name": "test",
    "networkMode": "awsvpc",
    "portMappings": [
      {
        "containerPort": 80,
        "hostPort": 80
      }
    ]
  }
]
DEFINITION
}

resource "aws_ecs_service" "test" {
  cluster         = aws_ecs_cluster.test.id
  desired_count   = 1
  launch_type     = "FARGATE"
  name            = %[1]q
  task_definition = aws_ecs_task_definition.test.arn

  deployment_controller {
    type = "CODE_DEPLOY"
  }

  load_balancer {
    container_name   = "test"
    container_port   = "80"
    target_group_arn = aws_lb_target_group.blue.id
  }

  network_configuration {
    assign_public_ip = true
    security_groups  = [aws_security_group.test.id]
    subnets          = aws_subnet.test[*].id
  }
}

resource "aws_codedeploy_app" "test" {
  compute_platform = "ECS"
  name             = %[1]q
}

data "aws_partition" "current" {}

resource "aws_iam_role" "test" {
  name = %[1]q

  assume_role_policy = <<POLICY
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Sid": "",
      "Effect": "Allow",
      "Action": "sts:AssumeRole",
      "Principal": {
        "Service": [
          "codedeploy.${data.aws_partition.current.dns_suffix}"
        ]
      }
    }
  ]
}
POLICY
}

resource "aws_iam_role_policy" "test" {
  role = aws_iam_role.test.name

  policy = <<POLICY
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": [
        "cloudwatch:DescribeAlarms",
        "ecs:CreateTaskSet",
        "ecs:DeleteTaskSet",
        "ecs:DescribeServices",
        "ecs:UpdateServicePrimaryTaskSet",
        "elasticloadbalancing:DescribeListeners",
        "elasticloadbalancing:DescribeRules",
        "elasticloadbalancing:DescribeTargetGroups",
        "elasticloadbalancing:ModifyListener",
        "elasticloadbalancing:ModifyRule",
        "lambda:InvokeFunction",
        "s3:GetObject",
        "s3:GetObjectMetadata",
        "s3:GetObjectVersion",
        "sns:Publish"
      ],
      "Effect": "Allow",
      "Resource": "*"
    }
  ]
}
POLICY
}
`, rName)
}

func testAccDeploymentGroupECSBlueGreenConfig(rName string) string {
	return testAccDeploymentGroupECSBaseConfig(rName) + fmt.Sprintf(`
resource "aws_codedeploy_deployment_group" "test" {
  app_name               = aws_codedeploy_app.test.name
  deployment_config_name = "CodeDeployDefault.ECSAllAtOnce"
  deployment_group_name  = %q
  service_role_arn       = aws_iam_role.test.arn

  auto_rollback_configuration {
    enabled = true
    events  = ["DEPLOYMENT_FAILURE"]
  }

  blue_green_deployment_config {
    deployment_ready_option {
      action_on_timeout = "CONTINUE_DEPLOYMENT"
    }

    terminate_blue_instances_on_deployment_success {
      action                           = "TERMINATE"
      termination_wait_time_in_minutes = 5
    }
  }

  deployment_style {
    deployment_option = "WITH_TRAFFIC_CONTROL"
    deployment_type   = "BLUE_GREEN"
  }

  ecs_service {
    cluster_name = aws_ecs_cluster.test.name
    service_name = aws_ecs_service.test.name
  }

  load_balancer_info {
    target_group_pair_info {
      prod_traffic_route {
        listener_arns = [aws_lb_listener.test.arn]
      }

      target_group {
        name = aws_lb_target_group.blue.name
      }

      target_group {
        name = aws_lb_target_group.green.name
      }
    }
  }
}
`, rName)
}

func testAccDeploymentGroupECSBlueGreenUpdateConfig(rName string) string {
	return testAccDeploymentGroupECSBaseConfig(rName) + fmt.Sprintf(`
resource "aws_codedeploy_deployment_group" "test" {
  app_name               = aws_codedeploy_app.test.name
  deployment_config_name = "CodeDeployDefault.ECSAllAtOnce"
  deployment_group_name  = %q
  service_role_arn       = aws_iam_role.test.arn

  auto_rollback_configuration {
    enabled = true
    events  = ["DEPLOYMENT_FAILURE"]
  }

  blue_green_deployment_config {
    deployment_ready_option {
      action_on_timeout    = "STOP_DEPLOYMENT"
      wait_time_in_minutes = 30
    }

    terminate_blue_instances_on_deployment_success {
      action                           = "TERMINATE"
      termination_wait_time_in_minutes = 60
    }
  }

  deployment_style {
    deployment_option = "WITH_TRAFFIC_CONTROL"
    deployment_type   = "BLUE_GREEN"
  }

  ecs_service {
    cluster_name = aws_ecs_cluster.test.name
    service_name = aws_ecs_service.test.name
  }

  load_balancer_info {
    target_group_pair_info {
      prod_traffic_route {
        listener_arns = [aws_lb_listener.test.arn]
      }

      target_group {
        name = aws_lb_target_group.blue.name
      }

      target_group {
        name = aws_lb_target_group.green.name
      }
    }
  }
}
`, rName)
}

func testAccDeploymentGroupTags1Config(rName, tagKey1, tagValue1 string) string {
	return testAccDeploymentGroupBaseConfig(rName) + fmt.Sprintf(`
resource "aws_codedeploy_deployment_group" "test" {
  app_name              = aws_codedeploy_app.test.name
  deployment_group_name = "tf-acc-test-%[1]s"
  service_role_arn      = aws_iam_role.test.arn

  tags = {
    %[2]q = %[3]q
  }
}
`, rName, tagKey1, tagValue1)
}

func testAccDeploymentGroupTags2Config(rName, tagKey1, tagValue1, tagKey2, tagValue2 string) string {
	return testAccDeploymentGroupBaseConfig(rName) + fmt.Sprintf(`
resource "aws_codedeploy_deployment_group" "test" {
  app_name              = aws_codedeploy_app.test.name
  deployment_group_name = "tf-acc-test-%[1]s"
  service_role_arn      = aws_iam_role.test.arn

  tags = {
    %[2]q = %[3]q
    %[4]q = %[5]q
  }
}
`, rName, tagKey1, tagValue1, tagKey2, tagValue2)
}
