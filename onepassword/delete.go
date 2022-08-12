package onepassword

import (
	"bytes"
	"errors"

	exec "golang.org/x/sys/execabs"
)

// DeleteCredentials deletes credentials from 1password.
func (c *Client) DeleteCredentials(_, host string, archive bool) error {
	var stdout bytes.Buffer

	var stderr bytes.Buffer

	args := []string{"--session", c.token, "item", "delete", "--format", "json", host}

	if archive {
		args = append(args, "--archive")
	}

	cmd := exec.Command("op", args...)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return errors.New(stderr.String()) // nolint:goerr113 // TODO: refactor
	}

	return nil
}
