package cloudtrail

import (
	"testing"

	"github.com/aquasecurity/defsec/internal/types"

	"github.com/aquasecurity/defsec/pkg/state"

	"github.com/aquasecurity/defsec/pkg/providers/aws/cloudtrail"
	"github.com/aquasecurity/defsec/pkg/providers/aws/s3"
	"github.com/aquasecurity/defsec/pkg/scan"

	"github.com/stretchr/testify/assert"
)

func TestCheckBucketAccessLoggingRequired(t *testing.T) {
	tests := []struct {
		name     string
		inputCT  cloudtrail.CloudTrail
		inputS3  s3.S3
		expected bool
	}{
		{
			name: "Trail has bucket with logging enabled",
			inputCT: cloudtrail.CloudTrail{
				Trails: []cloudtrail.Trail{
					{
						Metadata:   types.NewTestMetadata(),
						BucketName: types.String("my-bucket", types.NewTestMetadata()),
					},
				},
			},
			inputS3: s3.S3{
				Buckets: []s3.Bucket{
					{
						Metadata: types.NewTestMetadata(),
						Name:     types.String("my-bucket", types.NewTestMetadata()),
						Logging: s3.Logging{
							Metadata: types.NewTestMetadata(),
							Enabled:  types.Bool(true, types.NewTestMetadata()),
						},
					},
				},
			},
			expected: false,
		},
		{
			name: "Trail has bucket without logging enabled",
			inputCT: cloudtrail.CloudTrail{
				Trails: []cloudtrail.Trail{
					{
						Metadata:   types.NewTestMetadata(),
						BucketName: types.String("my-bucket", types.NewTestMetadata()),
					},
				},
			},
			inputS3: s3.S3{
				Buckets: []s3.Bucket{
					{
						Metadata: types.NewTestMetadata(),
						Name:     types.String("my-bucket", types.NewTestMetadata()),
						Logging: s3.Logging{
							Metadata: types.NewTestMetadata(),
							Enabled:  types.Bool(false, types.NewTestMetadata()),
						},
					},
				},
			},
			expected: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var testState state.State
			testState.AWS.CloudTrail = test.inputCT
			testState.AWS.S3 = test.inputS3
			results := checkBucketAccessLoggingRequired.Evaluate(&testState)
			var found bool
			for _, result := range results {
				if result.Status() == scan.StatusFailed && result.Rule().LongID() == checkBucketAccessLoggingRequired.Rule().LongID() {
					found = true
				}
			}
			if test.expected {
				assert.True(t, found, "Rule should have been found")
			} else {
				assert.False(t, found, "Rule should not have been found")
			}
		})
	}
}