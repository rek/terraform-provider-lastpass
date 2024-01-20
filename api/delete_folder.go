package api

import (
	"bytes"
	"errors"
	"os/exec"
	"strings"
)

// Delete a user on a shared folder
func (c *Client) DeleteFolder(folder string, email string) error {
	err := c.login()
	if err != nil {
		return err
	}
	var errbuf bytes.Buffer
	cmd := exec.Command("lpass", "share", "userdel", folder, email)
	cmd.Stderr = &errbuf
	err = cmd.Run()
	if err != nil {
		// Make sure the secret is not removed manually.
		if strings.Contains(errbuf.String(), "Could not find specified account") {
			return nil
		}
		var err = errors.New(errbuf.String())
		return err
	}
	return nil
}
