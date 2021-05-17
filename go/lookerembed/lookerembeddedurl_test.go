package lookerembed

import (
	"fmt"
	"testing"
	"time"
)

func TestRandomStringDoesNotRepeat(t *testing.T) {
	randomStrings := map[string]bool{}
	for i := 0; i < 100; i++ {
		randomString, err := randomString()
		if err != nil {
			t.Fatalf("Error while getting random string")
		}
		if _, exists := randomStrings[randomString]; exists {
			t.Errorf("Random string already in the set")
			return
		}
		if len(randomString) < 16 {
			t.Errorf("Random string too short")
			return
		}
		randomStrings[randomString] = true
	}
}

func TestStringificationWorksAsExpected(t *testing.T) {
	type paramsWithStringifiedPair struct {
		name string
		inp  *URLParams
		out  *stringifiedURLParams
		err  error
	}

	inputsOutputs := []paramsWithStringifiedPair{
		paramsWithStringifiedPair{
			"empty",
			&URLParams{
				Models:           []string{},
				SessionLength:    0,
				Permissions:      []string{},
				ForceLogoutLogin: true,
				UserAttributes:   map[string]string{},
			},
			&stringifiedURLParams{
				Models:           "[]",
				SessionLength:    "0",
				Permissions:      "[]",
				ForceLogoutLogin: "true",
				UserAttributes:   "{}",
			},
			nil,
		},
		paramsWithStringifiedPair{
			"some lists and map",
			&URLParams{
				Models:           []string{"foo", "bar"},
				SessionLength:    120,
				Permissions:      []string{"see_user_dashboards", "access_data"},
				ForceLogoutLogin: false,
				UserAttributes:   map[string]string{"foo": "bar"},
			},
			&stringifiedURLParams{
				Models:           "[\"foo\",\"bar\"]",
				SessionLength:    "120",
				Permissions:      "[\"see_user_dashboards\",\"access_data\"]",
				ForceLogoutLogin: "false",
				UserAttributes:   "{\"foo\": \"bar\"}",
			},
			nil,
		},
	}

	for _, test := range inputsOutputs {
		t.Run(test.name, func(t *testing.T) {
			input := test.inp
			expectedOutput := test.out
			output, err := newStringifiedURLParams(input)
			if err != test.err {
				t.Errorf("test: %s %q", test.name, err)
			}
			if output.Models != expectedOutput.Models {
				t.Errorf("test: %s models i %s; want %s", test.name, output.Models, expectedOutput.Models)
			}
			if output.SessionLength != expectedOutput.SessionLength {
				t.Errorf("test: %s sessionLength is %s; want %s", test.name, output.SessionLength, expectedOutput.SessionLength)
			}
			if output.Permissions != expectedOutput.Permissions {
				t.Errorf("test: %s permission is %s; want %s", test.name, output.Permissions, expectedOutput.Permissions)
			}
			if output.ForceLogoutLogin != expectedOutput.ForceLogoutLogin {
				t.Errorf("test: %s forcelogin is %s; want %s", test.name, output.ForceLogoutLogin, expectedOutput.ForceLogoutLogin)
			}
		})
	}

}

func TestSigning(t *testing.T) {
	type paramsStringPair struct {
		name string
		inp  *URLParams
		out  string
		err  error
	}
	inputsOutputs := []paramsStringPair{
		paramsStringPair{
			name: "somedata",
			inp: &URLParams{
				Host:             "chroniclepbldev.cloud.looker.com",
				Path:             "/embed/dashboards/5",
				ExternalUserID:   "username@acmeinc",
				ExternalGroupID:  "acmeinc",
				Models:           []string{"acmeinc"},
				GroupIDs:         []int64{3},
				SessionLength:    24 * 60 * 60,
				Permissions:      []string{"access_data", "see_user_dashboards", "see_lookml_dashboards", "see_looks"},
				UserAttributes:   map[string]string{},
				ForceLogoutLogin: true,
			},
			out: "YUff/8shh7e0ELnUi2JAfp4Q/8I=",
			err: nil,
		},
	}

	for _, test := range inputsOutputs {
		t.Run(test.name, func(t *testing.T) {
			input := test.inp
			expectedOutput := test.out

			nonce := "\"2d1a092d0bcb7ce292dea24b32f1acc2\""
			strParams, err := newStringifiedURLParams(input)
			if err != test.err {
				t.Errorf("%s got %q", test.name, err)
			}
			output := strParams.sign("adfagadg3egad", nonce, "1610834763")
			if output != expectedOutput {
				t.Errorf("%s got %s; want %s", test.name, output, expectedOutput)
			}
		})
	}
}

// Use this test to generate a URL with a key and
// then manually test it against the Looker Admin UI
// sandbox. 
// 
// You will need to set the correct host & SECRET
func TestLast(t *testing.T) {
	foo := &URLParams{
		Host:             "acmeincinstance.cloud.looker.com",
		Path:             "/embed/dashboards-next/47",
		ExternalUserID:   "username@acmeinc",
		ExternalGroupID:  "acmeinc",
		Models:           []string{"acmeinc"},
		GroupIDs:         []int64{70},
		SessionLength:    24 * 60 * 60,
		Permissions:      []string{"access_data", "see_user_dashboards", "see_lookml_dashboards", "see_looks"},
		UserAttributes:   map[string]string{},
		ForceLogoutLogin: true,
	}

	res, err := foo.CreateLookerSSOEmbeddedHostnameAndPath("SECRET", time.Minute)
	if err != nil {
		t.Errorf("failed to create Embedded URL %q", err)
	}

	// This is the URL to try at Looker console
	fmt.Printf("https://%s\n", res)
}
