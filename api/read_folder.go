package api

import (
	"os/exec"
	"bytes"
	"errors"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
)

type User struct {
	Name     string
	Email    string
	ReadOnly bool
	Admin    bool
	Hide     bool
	OutEnt   bool
	Accept   bool
}

// id = folder_name/user@email.com
func (c *Client) ReadShare(id string) (FolderShare, error) {
	splitLine := strings.Split(id, "/")
	var folder_name = splitLine[0]
	var email = splitLine[1]
	data, err := c.ReadShares(folder_name)
	if err != nil {
		var folderShare FolderShare
		return folderShare, err
	}
	var user = findUser(data, email)

	return user, nil
}

// id = folder_name
func (c *Client) ReadShares(id string) ([]FolderShare, error) {
	var folderShares []FolderShare
	log.Printf("[INFO] Trying to read folder share '%s'", id)
	cmd := exec.Command("lpass", "share", "userls", id)
	var outbuf, errbuf bytes.Buffer
	cmd.Stdout = &outbuf
	cmd.Stderr = &errbuf
	err := cmd.Run()

	if err != nil {
		log.Printf("[ERROR] %v", err)
		var err = errors.New("Error reading shares")
		diag.FromErr(err)
		return folderShares, err
	}

	var shares = parseUsers(outbuf.String())

	// var debugInput = `User                     RO Admin  Hide OutEnt Accept
	// Faris Alfaris <faris.alfaris@cntxt.com>                  x   _   _   x   x
	// Another Nice Dude <another.dude@email.com>   _   x   _   _   x
	// `
	// 	var shares = parseUsers(debugInput)

	// log.Printf("[DEBUG] ReadShares data", shares)
	log.Printf("[INFO] ReadShares amount: %v", len(shares))

	if len(shares) == 0 {
		var err = errors.New("No shares found")
		diag.FromErr(err)
		return folderShares, err
	}

	for _, share := range shares {
		folderShare := FolderShare{
			Name:     share.Name,
			Email:    share.Email, // potentially not fully valid, it could have been trimmed by the cli
			ReadOnly: share.ReadOnly,
			Admin:    share.Admin,
			Hide:     share.Hide,
			OutEnt:   share.OutEnt,
			Accept:   share.Accept,
		}
		folderShares = append(folderShares, folderShare)
	}

	return folderShares, nil
}

func parseUsers(input string) []User {
	var users []User

	lines := strings.Split(input, "\n")

	// Skip the header line
	for _, line := range lines[1:] {
		// println("starting with line:", line)

		// Extract name as anything before '<'
		splitLine := strings.Split(line, "<")

		// ignore bad lines
		if len(splitLine) != 2 {
			continue
		}

		name := strings.TrimSpace(splitLine[0])
		fields := strings.Fields(splitLine[1])

		if len(fields) == 6 {
			user := User{
				Name:     name,
				Email:    strings.Trim(fields[0], " <>"),
				ReadOnly: fields[1] == "x",
				Admin:    fields[2] == "x",
				Hide:     fields[3] == "x",
				OutEnt:   fields[4] == "x",
				Accept:   fields[5] == "x",
			}
			users = append(users, user)
		}
	}

	return users
}

// this is duplicated, need to make it a util
func findUser(shares []FolderShare, email string) FolderShare {
	share_map := make(map[string]FolderShare)
	for _, share := range shares {
		share_map[share.Email] = share
	}

	// _, found := share_map[email]

	return share_map[email]
}
