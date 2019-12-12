package main

import (
	"fmt"

	"github.com/equinox-io/equinox"
)

// assigned when creating a new application in the dashboard
const appID = "app_ja6WuaZgwsF"

// public portion of signing key generated by `equinox genkey`
var publicKey = []byte(`
-----BEGIN ECDSA PUBLIC KEY-----
MHYwEAYHKoZIzj0CAQYFK4EEACIDYgAEDNdY0eUMfybJP52Wsy/Csy/MMX5Zua81
4aefCgDOiXP+weaRPl9JQTDEAhygBS0ksEW8G8vtFNJVVJycPdBB8VVk+mM/bqxU
gkLzMYCOBaVtEukjyxZCUutxwwf/8XfW
-----END ECDSA PUBLIC KEY-----
`)

func equinoxUpdate() error {
	var opts equinox.Options
	if err := opts.SetPublicKeyPEM(publicKey); err != nil {
		return err
	}

	// check for the update
	resp, err := equinox.Check(appID, opts)
	switch {
	case err == equinox.NotAvailableErr:
		fmt.Println("No update available, already at the latest version!")
		return nil
	case err != nil:
		fmt.Println("Update failed:", err)
		return err
	}

	// fetch the update and apply it
	err = resp.Apply()
	if err != nil {
		return err
	}

	fmt.Printf("Updated to new version: %s!\n", resp.ReleaseVersion)
	return nil
}