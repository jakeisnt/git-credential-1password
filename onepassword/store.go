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
func (c *Client) StoreCredentials(protocol, host, path, username, password string) error { // nolint:funlen,lll // TODO: refactor
	creds, _ := c.GetCredentials(host, path) // nolint:errcheck // only check already existing

	if creds != nil {
		return nil
	}

	var stdout bytes.Buffer

	var stderr bytes.Buffer

	// TODO: handle session expired error

	title := host

	if path != "" {
		title += "/" + path
	}

	args := []string{"--cache", "--session", c.token}

	if isV2() {
		usernameField := fmt.Sprintf("username=%s", username)
		passwordField := fmt.Sprintf("password=%s", password)
		protoField := fmt.Sprintf("protocol[text]=%s", protocol)
		args = append(args, "item", "create", "--format", "json", "--category", "login",
			usernameField, passwordField, protoField)
	} else {
		data, err := json.Marshal(login{
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
		})

		if err != nil {
			return err
		}

		// base64url encode data
		encData := base64.StdEncoding.EncodeToString(data)
		encData = strings.ReplaceAll(encData, "+", "-") // 62nd char of encoding
		encData = strings.ReplaceAll(encData, "/", "_") // 63rd char of encoding
		encData = strings.ReplaceAll(encData, "=", "")  // Remove any trailing '='s

		args = append(args, "create", "item", "Login", encData)
	}

	args = append(args, "--title", title, "--tags", "git-credential-1password")

	cmd := exec.Command("op", args...)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return errors.New(stderr.String()) // nolint:goerr113 // TODO: refactor
	}

	return nil
}
