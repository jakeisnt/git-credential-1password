package onepassword

import (
	"bytes"
	"errors"
	"fmt"

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

	usernameField := fmt.Sprintf("username=%s", username)
	passwordField := fmt.Sprintf("password=%s", password)
	protoField := fmt.Sprintf("protocol[text]=%s", protocol)

	args := []string{"--cache", "--session", c.token, "item", "create", "--format", "json", "--category", "login",
		usernameField, passwordField, protoField, "--title", title, "--tags", "git-credential-1password"}

	cmd := exec.Command("op", args...)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return errors.New(stderr.String()) // nolint:goerr113 // TODO: refactor
	}

	return nil
}
