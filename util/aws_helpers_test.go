package util_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/YashdalfTheGray/statusinator/util"
)

func TestAWSHelpersGetSessionName(t *testing.T) {
	user, _ := os.LookupEnv("USER")
	testCases := []struct {
		desc, input, output string
		expectingErr        bool
	}{
		{
			desc:         "outputs correct session name given valid role ARN",
			input:        "arn:aws:iam::012345678901:role/test-role",
			output:       fmt.Sprintf("%s-test-role-statusinator", user),
			expectingErr: false,
		},
		{
			desc:         "outputs error given malformed role ARN",
			input:        "Theo is a good dog",
			output:       "",
			expectingErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			actual, err := util.GetSessionName(tc.input)
			if !tc.expectingErr && actual != tc.output {
				t.Errorf("Expected %s but got %s", tc.output, actual)
			}
			if tc.expectingErr && err == nil {
				t.Error("Expected error and got nil")
			}
		})
	}
}
