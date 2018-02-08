package aws

import (
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloud9"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccAWSCloud9EnvironmentEc2_basic(t *testing.T) {
	var conf cloud9.Environment

	rString := acctest.RandString(8)
	envName := fmt.Sprintf("tf_acc_env_basic_%s", rString)
	uEnvName := fmt.Sprintf("tf_acc_env_basic_updated_%s", rString)

	resourceName := "aws_cloud9_environment_ec2.test"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAWSCloud9EnvironmentEc2Destroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccAWSCloud9EnvironmentEc2Config(envName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSCloud9EnvironmentEc2Exists(resourceName, &conf),
					resource.TestCheckResourceAttr(resourceName, "instance_type", "t2.micro"),
					resource.TestCheckResourceAttr(resourceName, "name", envName),
					resource.TestCheckResourceAttrSet(resourceName, "arn"),
					resource.TestCheckResourceAttrSet(resourceName, "owner_arn"),
				),
			},
			resource.TestStep{
				Config: testAccAWSCloud9EnvironmentEc2Config(uEnvName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSCloud9EnvironmentEc2Exists(resourceName, &conf),
					resource.TestCheckResourceAttr(resourceName, "instance_type", "t2.micro"),
					resource.TestCheckResourceAttr(resourceName, "name", uEnvName),
					resource.TestCheckResourceAttrSet(resourceName, "arn"),
					resource.TestCheckResourceAttrSet(resourceName, "owner_arn"),
				),
			},
		},
	})
}

func TestAccAWSCloud9EnvironmentEc2_allFields(t *testing.T) {
	var conf cloud9.Environment

	rString := acctest.RandString(8)
	envName := fmt.Sprintf("tf_acc_env_basic_%s", rString)
	uEnvName := fmt.Sprintf("tf_acc_env_basic_updated_%s", rString)
	description := fmt.Sprintf("Tf Acc Test %s", rString)
	uDescription := fmt.Sprintf("Tf Acc Test Updated %s", rString)

	resourceName := "aws_cloud9_environment_ec2.test"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAWSCloud9EnvironmentEc2Destroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccAWSCloud9EnvironmentEc2AllFieldsConfig(envName, description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSCloud9EnvironmentEc2Exists(resourceName, &conf),
					resource.TestCheckResourceAttr(resourceName, "instance_type", "t2.micro"),
					resource.TestCheckResourceAttr(resourceName, "name", envName),
					resource.TestCheckResourceAttrSet(resourceName, "arn"),
					resource.TestCheckResourceAttrSet(resourceName, "owner_arn"),
					resource.TestCheckResourceAttrSet(resourceName, "type"),
				),
			},
			resource.TestStep{
				Config: testAccAWSCloud9EnvironmentEc2AllFieldsConfig(uEnvName, uDescription),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSCloud9EnvironmentEc2Exists(resourceName, &conf),
					resource.TestCheckResourceAttr(resourceName, "instance_type", "t2.micro"),
					resource.TestCheckResourceAttr(resourceName, "name", uEnvName),
					resource.TestCheckResourceAttrSet(resourceName, "arn"),
					resource.TestCheckResourceAttrSet(resourceName, "owner_arn"),
					resource.TestCheckResourceAttrSet(resourceName, "type"),
				),
			},
		},
	})
}

func TestAccAWSCloud9EnvironmentEc2_importBasic(t *testing.T) {
	var conf cloud9.Environment

	rString := acctest.RandString(8)
	name := fmt.Sprintf("tf_acc_api_doc_part_import_%s", rString)

	resourceName := "aws_cloud9_environment_ec2.test"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAWSCloud9EnvironmentEc2Destroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccAWSCloud9EnvironmentEc2Config(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSCloud9EnvironmentEc2Exists(resourceName, &conf),
					resource.TestCheckResourceAttr(resourceName, "instance_type", "t2.micro"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttrSet(resourceName, "arn"),
					resource.TestCheckResourceAttrSet(resourceName, "owner_arn"),
					resource.TestCheckResourceAttrSet(resourceName, "type"),
				),
			},
			resource.TestStep{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"instance_type"},
			},
		},
	})
}

func testAccCheckAWSCloud9EnvironmentEc2Exists(n string, res *cloud9.Environment) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Cloud9 Environment EC2 ID is set")
		}

		conn := testAccProvider.Meta().(*AWSClient).cloud9conn

		out, err := conn.DescribeEnvironments(&cloud9.DescribeEnvironmentsInput{
			EnvironmentIds: []*string{aws.String(rs.Primary.ID)},
		})
		if err != nil {
			if isAWSErr(err, cloud9.ErrCodeNotFoundException, "") {
				return fmt.Errorf("Cloud9 Environment EC2 (%q) not found", rs.Primary.ID)
			}
			return err
		}
		if len(out.Environments) == 0 {
			return fmt.Errorf("Cloud9 Environment EC2 (%q) not found", rs.Primary.ID)
		}
		env := out.Environments[0]

		*res = *env

		return nil
	}
}

func testAccCheckAWSCloud9EnvironmentEc2Destroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*AWSClient).cloud9conn

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aws_cloud9_environment_ec2" {
			continue
		}

		out, err := conn.DescribeEnvironments(&cloud9.DescribeEnvironmentsInput{
			EnvironmentIds: []*string{aws.String(rs.Primary.ID)},
		})
		if err != nil {
			if isAWSErr(err, cloud9.ErrCodeNotFoundException, "") {
				return nil
			}
			// :'-(
			if isAWSErr(err, "AccessDeniedException", "is not authorized to access this resource") {
				return nil
			}
			return err
		}
		if len(out.Environments) == 0 {
			return nil
		}

		return fmt.Errorf("Cloud9 Environment EC2 %q still exists.", rs.Primary.ID)
	}
	return nil
}

func testAccAWSCloud9EnvironmentEc2Config(name string) string {
	return fmt.Sprintf(`
resource "aws_cloud9_environment_ec2" "test" {
  instance_type = "t2.micro"
  name = "%s"
}
`, name)
}

func testAccAWSCloud9EnvironmentEc2AllFieldsConfig(name, description string) string {
	return fmt.Sprintf(`
resource "aws_cloud9_environment_ec2" "test" {
  instance_type = "t2.micro"
  name = "%s"
  description = "%s"
  automatic_stop_time_minutes = 60
  subnet_id = "${aws_subnet.test.id}"
  depends_on = ["aws_route_table_association.test"]
}

resource "aws_vpc" "test" {
  cidr_block = "10.10.0.0/16"
}

resource "aws_subnet" "test" {
  vpc_id = "${aws_vpc.test.id}"
  cidr_block = "10.10.0.0/19"
}

resource "aws_internet_gateway" "test" {
  vpc_id = "${aws_vpc.test.id}"
}

resource "aws_route_table" "test" {
  vpc_id = "${aws_vpc.test.id}"
  route {
    cidr_block = "0.0.0.0/0"
    gateway_id = "${aws_internet_gateway.test.id}"
  }
}

resource "aws_route_table_association" "test" {
  subnet_id      = "${aws_subnet.test.id}"
  route_table_id = "${aws_route_table.test.id}"
}
`, name, description)
}
