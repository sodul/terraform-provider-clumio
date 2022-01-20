// Copyright 2021. Clumio, Inc.

package clumio

const (
	awsAccountIDRegexpInternalPattern = `(aws|\d{12})`
	awsPartitionRegexpInternalPattern = `aws(-[a-z]+)*`
	awsRegionRegexpInternalPattern    = `[a-z]{2}(-[a-z]+)+-\d`
	awsAccountIDRegexpPattern         = "^" + awsAccountIDRegexpInternalPattern + "$"
	awsPartitionRegexpPattern         = "^" + awsPartitionRegexpInternalPattern + "$"
	awsRegionRegexpPattern            = "^" + awsRegionRegexpInternalPattern + "$"
	awsProfile                        = "AWS_PROFILE"
	awsSharedCredsFile                = "AWS_SHARED_CREDENTIALS_FILE"
	kMaxRetries                       = 8
)
