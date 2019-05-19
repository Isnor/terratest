package aws

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/stretchr/testify/require"
)

// GetDynamoDbTableTags fetches resource tags of a specified dynamoDB table.
func GetDynamoDbTableTags(t *testing.T, region string, tableName string) []*dynamodb.Tag {
	tags, err := GetDynamoDbTableTagsE(t, region, tableName)
	require.NoError(t, err)
	return tags
}

// GetDynamoDbTableTagsE fetches resource tags of a specified dynamoDB table.
func GetDynamoDbTableTagsE(t *testing.T, region string, tableName string) ([]*dynamodb.Tag, error) {
	table := GetDynamoDBTable(t, region, tableName)
	out, err := NewDynamoDBClient(t, region).ListTagsOfResource(&dynamodb.ListTagsOfResourceInput{
		ResourceArn: table.TableArn,
	})
	if err != nil {
		return nil, err
	}
	return out.Tags, err
}

// GetDynamoDBTableTimeToLive fetches information about TTL configuration of a specified dynamoDB table.
func GetDynamoDBTableTimeToLive(t *testing.T, region string, tableName string) *dynamodb.TimeToLiveDescription {
	ttl, err := GetDynamoDBTableTimeToLiveE(t, region, tableName)
	require.NoError(t, err)
	return ttl
}

// GetDynamoDBTableTimeToLiveE fetches information about TTL configuration of a specified dynamoDB table.
func GetDynamoDBTableTimeToLiveE(t *testing.T, region string, tableName string) (*dynamodb.TimeToLiveDescription, error) {
	out, err := NewDynamoDBClient(t, region).DescribeTimeToLive(&dynamodb.DescribeTimeToLiveInput{
		TableName: aws.String(tableName),
	})
	if err != nil {
		return nil, err
	}
	return out.TimeToLiveDescription, err
}

// GetDynamoDBTable fetches information about specified dynamoDB table.
func GetDynamoDBTable(t *testing.T, region string, tableName string) *dynamodb.TableDescription {
	table, err := GetDynamoDBTableE(t, region, tableName)
	require.NoError(t, err)
	return table
}

// GetDynamoDBTableE fetches information about specified dynamoDB table.
func GetDynamoDBTableE(t *testing.T, region string, tableName string) (*dynamodb.TableDescription, error) {
	out, err := NewDynamoDBClient(t, region).DescribeTable(&dynamodb.DescribeTableInput{
		TableName: aws.String(tableName),
	})
	if err != nil {
		return nil, err
	}
	return out.Table, err
}

// NewDynamoDBClient creates a DynamoDB client.
func NewDynamoDBClient(t *testing.T, region string) *dynamodb.DynamoDB {
	client, err := NewDynamoDBClientE(t, region)
	require.NoError(t, err)
	return client
}

// NewDynamoDBClientE creates a DynamoDB client.
func NewDynamoDBClientE(t *testing.T, region string) (*dynamodb.DynamoDB, error) {
	sess, err := NewAuthenticatedSession(region)
	if err != nil {
		return nil, err
	}
	return dynamodb.New(sess), nil
}
