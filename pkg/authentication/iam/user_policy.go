package iam

import (
	"fmt"
)

func defaultPolicy(bucket string) string {
	return fmt.Sprintf(`{
		"Version": "2012-10-17",
		"Statement": [
			{
				"Effect": "Allow",
				"Action": [
					"s3:GetObject",
					"s3:GetObjectVersion",
					"s3:PutObject",
					"s3:DeleteObject",
					"s3:DeleteObjectVersion",
					"s3:ListBucket"
				],
				"Resource": [
					"arn:aws:s3:::%s",
					"arn:aws:s3:::%s/*"
				]
			}
		]
	}`, bucket, bucket)
}
