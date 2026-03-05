package main

import (
	"fmt"
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
	email string required unique
	confirmed string
	role string required
	regdate string required
	password string required
::
`

	SessionsStmt = `
table: sessions
fields:
	session_code string required unique
	creation_dt datetime required
	user_id int required
::
foreign_keys:
    user_id users on_delete_delete
::
	`

	SeedCompanysStmt = `
table: seed_companies
fields:
	name string required unique
	cacno string required unique
	access_code string required unique
::
	`

	ProductionPlansStmt = `
table: production_plans
fields:
	username string required
	date date required
	company_id int required
	state string required
	address string required
	mobileNumber string required
	plantingArea string required
	latitude string required
	longitude string required
	crop string required
	seedClass string required
	variety string required
	productionYear int required
	plantingDate date required
	srcProductionYear int required
	srcQuantityProcured int required
	srcSeedCodexNumber string
	srcSeedClass string required
	srcSupplierName string required
	receipt_photo text
	field_photo_1 text
	field_photo_2 text
::	
foreign_keys:
  company_id seed_companies on_delete_delete
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

	stmts := []string{UsersFlaarumStmt, SessionsStmt, SeedCompanysStmt, ProductionPlansStmt}
	for _, stmt := range stmts {
		err = cl.CreateOrUpdateTable(stmt)
		if err != nil {
			panic(err)
		}
	}

	// load data
	accessCode := "cp"
	_, err = cl.SearchForOne(fmt.Sprintf(`
		table: seed_companies
		where:
		  access_code = %s
		`, accessCode))
	if err != nil {
		toWrite := map[string]string{
			"name":        "Test1",
			"cacno":       "RC001",
			"access_code": accessCode,
		}

		cl.InsertRowStr("seed_companies", toWrite)
	}

	return cl
}
