package api

import (
	"log"
	"bytes"
	"errors"
	"os/exec"
)

// Update is called to update secret with upstream
func (c *Client) UpdateFolder(folder_share FolderShare) error {
	err := c.login()
	if err != nil {
		return err
	}
	var readOnly = "false"
	if folder_share.ReadOnly {
		readOnly = "true"
	}
	var hidden = "false"
	if folder_share.Hide {
		readOnly = "true"
	}
	var admin = "false"
	if folder_share.Admin {
		admin = "true"
	}

	log.Printf("[INFO] ================= 1--->" + "--read-only="+readOnly)
	cmd := exec.Command("lpass", "share", "usermod", "--read-only="+readOnly, "--hidden="+hidden, "--admin="+admin, folder_share.Folder, folder_share.Email)

	var outbuf, errbuf bytes.Buffer
	cmd.Stdin = &outbuf
	cmd.Stderr = &errbuf
	err = cmd.Run()
	if err != nil {
		var err = errors.New(errbuf.String())
		return err
	}

	return err
}
