package aws

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/stretchr/testify/require"
)

// GetParameter retrieves the latest version of SSM Parameter at keyName with decryption.
func GetParameter(t *testing.T, awsRegion string, keyName string) string {
	keyValue, err := GetParameterE(t, awsRegion, keyName)
	require.NoError(t, err)
	return keyValue
}

// GetParameterE retrieves the latest version of SSM Parameter at keyName with decryption.
func GetParameterE(t *testing.T, awsRegion string, keyName string) (string, error) {
	ssmClient, err := NewSsmClientE(t, awsRegion)
	if err != nil {
		return "", err
	}

	resp, err := ssmClient.GetParameter(&ssm.GetParameterInput{Name: aws.String(keyName), WithDecryption: aws.Bool(true)})
	if err != nil {
		return "", err
	}

	parameter := *resp.Parameter
	return *parameter.Value, nil
}

// PutParameter creates new version of SSM Parameter at keyName with keyValue as SecureString.
func PutParameter(t *testing.T, awsRegion string, keyName string, keyDescription string, keyValue string) int64 {
	version, err := PutParameterE(t, awsRegion, keyName, keyDescription, keyValue)
	require.NoError(t, err)
	return version
}

// PutParameterE creates new version of SSM Parameter at keyName with keyValue as SecureString.
func PutParameterE(t *testing.T, awsRegion string, keyName string, keyDescription string, keyValue string) (int64, error) {
	ssmClient, err := NewSsmClientE(t, awsRegion)
	if err != nil {
		return 0, err
	}

	resp, err := ssmClient.PutParameter(&ssm.PutParameterInput{Name: aws.String(keyName), Description: aws.String(keyDescription), Value: aws.String(keyValue), Type: aws.String("SecureString")})
	if err != nil {
		return 0, err
	}

	return *resp.Version, nil
}

// NewSsmClient creates a SSM client.
func NewSsmClient(t *testing.T, region string) *ssm.SSM {
	client, err := NewSsmClientE(t, region)
	require.NoError(t, err)
	return client
}

// NewSsmClientE creates an SSM client.
func NewSsmClientE(t *testing.T, region string) (*ssm.SSM, error) {
	sess, err := NewAuthenticatedSession(region)
	if err != nil {
		return nil, err
	}

	return ssm.New(sess), nil
}
