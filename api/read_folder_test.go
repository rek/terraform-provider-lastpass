package api

import (
	"testing"
)

func TestParseUsers(t *testing.T) {
	input := `User                     RO Admin  Hide OutEnt Accept
One Guy <one.guy@cntxt.com>    _   _   _   x   x
Another Nice Dude <another.dude@cntxt.com>     _   x   _   _   x
`

	expectedUsers := []User{
		{Name: "One Guy", Email: "one.guy@cntxt.com", RO: false, Admin: false, Hide: false, OutEnt: true, Accept: true},
		{Name: "Another Nice Dude", Email: "another.dude@cntxt.com", RO: false, Admin: true, Hide: false, OutEnt: false, Accept: true},
	}

	users := parseUsers(input)

	if len(users) != len(expectedUsers) {
		t.Errorf("Expected %d users, got %d", len(expectedUsers), len(users))
	}

	for i := range users {
		if users[i] != expectedUsers[i] {
			t.Errorf("User %d mismatch: expected %+v, got %+v", i, expectedUsers[i], users[i])
		}
	}
}
