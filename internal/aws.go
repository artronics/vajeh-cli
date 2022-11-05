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

	accessKey := viper.Get("access_key")
	if accessKey == nil {
		return awsCred, fmt.Errorf("AWS access_key is required. You must pass it as AWS_ACCESS_KEY_ID")
	}
	awsCred.AccessKey = accessKey.(string)

	accessSecret := viper.Get("aws_secret_access_key")
	if accessSecret == nil {
		return awsCred, fmt.Errorf("AWS access_secret is required. You must pass it as AWS_SECRET_ACCESS_KEY")
	}
	awsCred.AccessSecret = accessSecret.(string)

	return awsCred, nil
}
