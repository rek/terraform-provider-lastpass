package api

import (
	"bytes"
	"errors"
	"os/exec"
)

// Folder describes a Lastpass object.
type Folder struct {
	Name string `json:"fullname"`
}

// FolderShare describes a Lastpass folder share object.
type FolderShare struct {
	Id       string
	Folder   string
	Name     string
	Email    string
	ReadOnly bool
	Admin    bool
	Hide     bool
	OutEnt   bool
	Accept   bool
}

// Create is used to shared a folder with a user
func (c *Client) CreateFolder(folder_share FolderShare) (FolderShare, error) {
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
	cmd := exec.Command("lpass", "share", "useradd", "--read-only="+readOnly, "--hidden="+hidden, "--admin="+admin, folder_share.Folder, folder_share.Email)
	return c.createFolderShare(folder_share, cmd)
}

func (c *Client) createFolderShare(folder_share FolderShare, cmd *exec.Cmd) (FolderShare, error) {
	err := c.login()
	if err != nil {
		return folder_share, err
	}
	var outbuf, errbuf bytes.Buffer
	cmd.Stdin = &outbuf
	cmd.Stderr = &errbuf
	err = cmd.Run()
	if err != nil {
		var err = errors.New(errbuf.String())
		return folder_share, err
	}

	// folder_share.Id = folder_share.Folder + "/" + folder_share.Email

	return folder_share, err
}
