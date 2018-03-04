// This file contains utilities for working at the top-most level of this library
package terratest

import (
	"fmt"
	"github.com/gruntwork-io/terratest/aws"
	"github.com/gruntwork-io/terratest/util"
)

// A RandomResourceCollection is a typed holder for resources we need as we do a Terraform run.
// Some of these resources are dynamically generated (e.g. KeyPair) and others are randomly selected (e.g. AwsRegion).
type RandomResourceCollection struct {
	UniqueId    string      // A short unique id effective for namespacing resource names
	AwsRegion   string      // The AWS Region
	KeyPair     *Ec2Keypair // The EC2 KeyPair created in AWS
	AmiId       string      // A random AMI ID valid for the AwsRegion
	AccountId   string      // The AWS account ID
	SnsTopicArn string      // The ARN of the SNS Topic created in AWS
}

// Represents an EC2 KeyPair created in AWS
type Ec2Keypair struct {
	Name       string // The name assigned in AWS to the EC2 KeyPair
	PublicKey  string // The public key
	PrivateKey string // The private key in .pem format
}

// RandomResourecCollectionOpts represents the options passed when creating a new RandomResourceCollection
type RandomResourceCollectionOpts struct {
	ForbiddenRegions []string // A list of AWS Regions in which a test should never run
	ApprovedRegions  []string // The only AWS Regions in which a test should be run
}

func NewRandomResourceCollectionOptions() *RandomResourceCollectionOpts {
	return &RandomResourceCollectionOpts{}
}

// Create an instance of all properties in a RandomResourceCollection
func CreateRandomResourceCollection(ro *RandomResourceCollectionOpts) (*RandomResourceCollection, error) {
	r := &RandomResourceCollection{}

	r.AwsRegion = aws.GetRandomRegion(ro.ApprovedRegions, ro.ForbiddenRegions)

	r.UniqueId = util.UniqueId()

	// Fetch a random AMI ID
	r.AmiId = aws.GetUbuntu1604Ami(r.AwsRegion)

	// Generate a key pair and create it in AWS as an EC2 KeyPair
	keyPair, err := util.GenerateRSAKeyPair(2048)
	if err != nil {
		return r, fmt.Errorf("Failed to generate random key pair: %s\n", err.Error())
	}

	err = aws.CreateEC2KeyPair(r.AwsRegion, r.UniqueId, keyPair.PublicKey)
	if err != nil {
		return r, fmt.Errorf("Failed to create EC2 KeyPair: %s\n", err.Error())
	}

	ec2KeyPair := &Ec2Keypair{}
	ec2KeyPair.Name = r.UniqueId
	ec2KeyPair.PublicKey = keyPair.PublicKey
	ec2KeyPair.PrivateKey = keyPair.PrivateKey

	r.KeyPair = ec2KeyPair

	// Create an SNS Topic
	snsTopicArn, err := aws.CreateSnsTopic(r.AwsRegion, r.UniqueId)
	if err != nil {
		return r, fmt.Errorf("Failed to create SNS Topic: %s\n", err.Error())
	}
	r.SnsTopicArn = snsTopicArn

	// Get the AWS Account ID
	r.AccountId, err = aws.GetAccountId()
	if err != nil {
		return r, fmt.Errorf("Failed to get AWS Account Id: %s\n", err.Error())
	}

	return r, nil
}

// Destroy any persistent resources referenced in the given RandomResourceCollection.
func (r *RandomResourceCollection) DestroyResources() error {
	var err error

	if r != nil && r.AwsRegion != "" {
		if r.KeyPair.Name != "" {
			err = aws.DeleteEC2KeyPair(r.AwsRegion, r.KeyPair.Name)
		}
		if err != nil {
			return err
		}

		if r.SnsTopicArn != "" {
			err = aws.DeleteSNSTopic(r.AwsRegion, r.SnsTopicArn)
		}
		if err != nil {
			return err
		}
	}

	return nil
}

// Return the AWS Availability Zones for a given AWS region
func (r *RandomResourceCollection) FetchAwsAvailabilityZones() []string {
	if r != nil && r.AwsRegion != "" {
		return aws.GetAvailabilityZones(r.AwsRegion)
	}
	return nil
}

func (r *RandomResourceCollection) GetRandomPrivateCidrBlock(prefix int) string {
	return util.GetRandomPrivateCidrBlock(prefix)
}

func (r *RandomResourceCollection) GetDefaultVpc() (aws.Vpc, error) {
	return aws.GetDefaultVpc(r.AwsRegion)
}
