package main

import (
	"os"

	"github.com/saenuma/flaarumlib"
)

const (
	// roles: seedcoy, officer, nasc_viewer, admin
	UsersFlaarumStmt = `
table: users
fields:
  firstname string required
	surname string required
	email string required
	confirmed string
	role string required
	regdate string required
	password string required
	organisation string
::
`

	SessionsStmt = `
table: sessions
fields:
	session_code string required
	creation_dt datetime required
	user_id int required
::
foreign_keys:
    user_id users on_delete_delete
::
	`
)

func getFlaarumClient() flaarumlib.Client {

	var cl flaarumlib.Client
	if os.Getenv("SAE_DEV") == "true" {
		cl = flaarumlib.NewClient("nascsiw")
	} else {
		cl = flaarumlib.NewClient("first_proj")
	}

	err := cl.Ping()
	if err != nil {
		panic(err)
	}

	stmts := []string{UsersFlaarumStmt, SessionsStmt}
	for _, stmt := range stmts {
		err = cl.CreateOrUpdateTable(stmt)
		if err != nil {
			panic(err)
		}
	}

	return cl
}
