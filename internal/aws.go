package internal

import (
	"fmt"
	"github.com/spf13/viper"
)

type AwsCredentials struct {
	AccessKey    string
	AccessSecret string
}

func (a *AwsCredentials) ToEnvs() []string {
	access := fmt.Sprintf("AWS_ACCESS_KEY_ID=%s", a.AccessKey)
	secret := fmt.Sprintf("AWS_SECRET_ACCESS_KEY=%s", a.AccessSecret)

	return []string{access, secret}
}

func GetAwsCred() (AwsCredentials, error) {
	awsCred := AwsCredentials{}

	accessKey := viper.Get("aws_access_key_id")
	if accessKey == nil {
		return awsCred, fmt.Errorf("environment variable AWS_ACCESS_KEY_ID is required")
	}
	awsCred.AccessKey = accessKey.(string)

	accessSecret := viper.Get("aws_secret_access_key")
	if accessSecret == nil {
		return awsCred, fmt.Errorf("environment AWS_SECRET_ACCESS_KEY variable is required")
	}
	awsCred.AccessSecret = accessSecret.(string)

	return awsCred, nil
}
