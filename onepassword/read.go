package onepassword

import (
	"bytes"
	"encoding/json"
	"errors"

	exec "golang.org/x/sys/execabs"
)

const (
	usernameLabel = "username"
	passwordLabel = "password"
)

// GetCredentials loads credentials from 1password.
func (c *Client) GetCredentials(host, path string) (*Credentials, error) { //nolint:funlen,gocognit,gocyclo,lll // TODO: refactor
	var stdout bytes.Buffer

	var stderr bytes.Buffer

	title := host

	if path != "" {
		title += "/" + path
	}

	args := []string{"--cache", "--session", c.token, "item", "get", "--format", "json"}

	if c.Vault != "" {
		args = append(args, "--vault", c.Vault)
	}

	args = append(args, title)

	// TODO: handle session expired error
	cmd := exec.Command("op", args...)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return nil, errors.New(stderr.String()) //nolint:goerr113 // TODO: refactor
	}

	var username string

	var password string

	var res response
	if err := json.Unmarshal(stdout.Bytes(), &res); err != nil {
		return nil, err
	}

	for _, field := range res.Fields {
		if field.Name == usernameLabel {
			username = field.Value
		}

		if field.Name == passwordLabel {
			password = field.Value
		}
	}

	return &Credentials{
		Username: username,
		Password: password,
	}, nil
}
