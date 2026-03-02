package main

import (
	"fmt"
	"html/template"
	"net/http"
)

func createCompanyAccessCode() string {
	flcl := getFlaarumClient()

	for {
		rs := "c" + untestedRandomString(15)
		count, err := flcl.CountRows(fmt.Sprintf(`
			table: seed_companies
			where:
				access_code = %s
			`, rs))
		if err != nil {
			fmt.Println(err)
		}

		if count == 0 {
			return rs
		}
	}
}

func registerCompanyHandler(w http.ResponseWriter, r *http.Request) {
	userId, role, err := getCurrentUserData(r)
	if err != nil {
		fmt.Println(err)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	if userId == 0 {
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	if role != "admin" {
		errorPage(w, fmt.Errorf("You must be an admin to view this page."))
		return
	}

	if r.Method != http.MethodPost {
		tmpl := template.Must(template.ParseFiles("templates/base.html", "templates/register_company.html"))
		tmpl.Execute(w, nil)
	} else {

		newCode := createCompanyAccessCode()

		toWrite := map[string]string{
			"name":        r.FormValue("company_name"),
			"cacno":       r.FormValue("cacno"),
			"access_code": newCode,
		}

		flcl := getFlaarumClient()

		retId, err := flcl.InsertRowStr("seed_companies", toWrite)
		if err != nil {
			errorPage(w, err)
			return
		}

		http.Redirect(w, r, fmt.Sprintf("/scoya/%d", retId), http.StatusTemporaryRedirect)

	}

}

func seedCompanyAccessHandler(w http.ResponseWriter, r *http.Request) {
	userId, role, err := getCurrentUserData(r)
	if err != nil {
		fmt.Println(err)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	if userId == 0 {
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	if role != "admin" {
		errorPage(w, fmt.Errorf("You must be an admin to view this page."))
		return
	}

	flcl := getFlaarumClient()

	sCoyId := r.PathValue("id")
	aRow, err := flcl.SearchForOne(fmt.Sprintf(`
		table: seed_companies
		where:
		  id = %s
		`, sCoyId))
	if err != nil {
		errorPage(w, err)
		return
	}

	tmpl := template.Must(template.ParseFiles("templates/base.html", "templates/view_company_accesscode.html"))
	tmpl.Execute(w, aRow)
}
