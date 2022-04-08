package smoothstack_demo

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/gruntwork-io/terratest/modules/aws"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

type S3TagsJson struct {
	Deployment string `json:"deployment"`
	Env        string `json:"env"`
	Region     string `json:"region"`
}
type S3Json struct {
	Name       string     `json:"name"`
	Versioning string     `json:"versioning"`
	Tags       S3TagsJson `json:"tags"`
}

var deployment_passed bool
var ExpectedS3 S3Json

func init() {
	jsonFile, err := os.Open("./aws_s3_testdata.json")
	if err != nil {
		fmt.Println(err)
	}
	byteValue, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal(byteValue, &ExpectedS3)
}

func TestAWSS3BucketInput(t *testing.T) {
	// outFolder := "./"
	// planFilePath := filepath.Join(outFolder, "plan.out")
	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: "../",
		VarFiles:     []string{"input.tfvars"},
		// PlanFilePath: planFilePath,
		NoColor: true,
		EnvVars: map[string]string{
			"AWS_DEFAULT_REGION": ExpectedS3.Tags.Region,
		},
	})

	defer terraform.Destroy(t, terraformOptions)
	terraform.InitAndApply(t, terraformOptions)
	ActualS3 := S3Json{
		Name: terraform.Output(t, terraformOptions, "bucket_name"),
		Tags: S3TagsJson{
			Region:     terraform.Output(t, terraformOptions, "tag_region"),
			Deployment: terraform.Output(t, terraformOptions, "tag_deployment"),
			Env:        terraform.Output(t, terraformOptions, "tag_enviornment"),
		},
	}

	if assert.Equal(t, ExpectedS3.Name, ActualS3.Name) {
		deployment_passed = true
		t.Logf("PASS: The expected S3 Bucket name:%v matches the Actual S3 Bucket name:%v", ExpectedS3.Name, ActualS3.Name)
	} else {
		deployment_passed = false
		terraform.Destroy(t, terraformOptions)
		t.Fatalf("FAIL: Expected %v, but found %v", ExpectedS3.Name, ActualS3.Name)
	}

	if assert.Equal(t, ExpectedS3.Tags, ActualS3.Tags) {
		deployment_passed = true
		t.Logf("PASS: The expected S3 Bucket region:%v matches the Actual S3 Bucket tags:%v", ExpectedS3.Tags.Region, ActualS3.Tags.Region)
		t.Logf("PASS: The expected S3 Bucket environment:%v matches the Actual S3 Bucket tags:%v", ExpectedS3.Tags.Env, ActualS3.Tags.Env)
		t.Logf("PASS: The expected S3 Bucket deployment:%v matches the Actual S3 Bucket tags:%v", ExpectedS3.Tags.Deployment, ActualS3.Tags.Deployment)
	} else {
		deployment_passed = false
		terraform.Destroy(t, terraformOptions)
		t.Fatalf("FAIL: Expected %v, but found %v", ExpectedS3.Tags, ActualS3.Tags)
	}

	bucketID := terraform.Output(t, terraformOptions, "bucket_id")

	actualStatus := aws.GetS3BucketVersioning(t, ExpectedS3.Tags.Region, bucketID)

	assert.Equal(t, ExpectedS3.Versioning, actualStatus)
	time.Sleep(120 * time.Second)
	fmt.Println("Sleep Over.....")

}
