package onepassword

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	exec "golang.org/x/sys/execabs"
)

// StoreCredentials saves new credentials to 1password.
func (c *Client) StoreCredentials(protocol, host, path, username, password string) error {
	creds, err := c.GetCredentials(host, path)

	if creds != nil {
		return nil
	}

	var stdout bytes.Buffer

	var stderr bytes.Buffer

	// TODO: handle session expired error
	l := login{
		Notes: fmt.Sprintf("Protocol: %s", protocol),
		Fields: []field{
			{
				Name:        "username",
				Value:       username,
				Type:        "T",
				Designation: "username",
			},
			{
				Name:        "password",
				Value:       password,
				Type:        "P",
				Designation: "password",
			},
		},
	}

	data, err := json.Marshal(l)

	if err != nil {
		return err
	}

	// base64url encode data
	encData := base64.StdEncoding.EncodeToString(data)
	encData = strings.ReplaceAll(encData, "+", "-") // 62nd char of encoding
	encData = strings.ReplaceAll(encData, "/", "_") // 63rd char of encoding
	encData = strings.ReplaceAll(encData, "=", "")  // Remove any trailing '='s

	title := host

	if path != "" {
		title += "/" + path
	}

	cmd := exec.Command("op", "--cache", "--session", c.token, // nolint:gosec // TODO: validate
		"create", "item", "Login", encData, "--title", title, "--tags", "git-credential-1password")
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err = cmd.Run(); err != nil {
		return errors.New(stderr.String()) // nolint:goerr113 // TODO: refactor
	}

	return nil
}
