package main

import (
	"fmt"
	"os"
	"time"

	"github.com/looker-open-source/sdk-codegen/go/lookerembed"
)

func main() {
	foo := &lookerembed.URLParams{
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
		println("failed to create Embedded URL %v", err)
		os.Exit(1)
	}

	// This is the URL to try at Looker console
	fmt.Printf("https://%s\n", res)
}