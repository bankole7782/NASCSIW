package main

import (
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func getCurrentUserData(r *http.Request) (int64, string, error) {
	cookie, err := r.Cookie("things_things")
	if err != nil {
		return 0, "", nil
	}

	flcl := getFlaarumClient()

	count, err := flcl.CountRows(fmt.Sprintf(`
		table: sessions
		where:
		  session_code = %s
		`, cookie.Value))
	if err != nil {
		fmt.Println(err)
	}

	if count == 0 {
		return 0, "", errors.New("The cookie specified does not exist.")
	}

	sessionRow, err := flcl.SearchForOne(fmt.Sprintf(`
	table: sessions expand
	where:
	  session_code = %s`, cookie.Value))
	if err != nil {
		fmt.Println(err)
		return 0, "", err
	}

	return (*sessionRow)["user_id"].(int64), (*sessionRow)["user_id.role"].(string), nil
}

func createSessionCode() string {
	flcl := getFlaarumClient()

	for {
		rs := untestedRandomString(100)
		count, err := flcl.CountRows(fmt.Sprintf(`
			table: sessions
			where:
				session_code = %s
			`, rs))
		if err != nil {
			fmt.Println(err)
		}

		if count == 0 {
			return rs
		}
	}
}

func signinHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
	}

	if r.Method == http.MethodPost {
		email := r.FormValue("email")
		password := r.FormValue("password")

		flcl := getFlaarumClient()

		row, err := flcl.SearchForOne(fmt.Sprintf(`
		table: users
		where:
		  email = %s
		`, email))
		if err != nil {
			msg := template.HTML(fmt.Sprintf("email '%s' not found.", email))
			messageUser(w, r, "Invalid Credentails", msg)
			return
		}

		databasePassword := (*row)["password"].(string)
		err = bcrypt.CompareHashAndPassword([]byte(databasePassword), []byte(password))
		if err != nil {
			messageUser(w, r, "Invalid Credentails", "Email and Password combination not found.")
			return
		}

		now := time.Now()
		sessionCode := createSessionCode()

		toWrite := map[string]any{
			"session_code": sessionCode,
			"creation_dt":  time.Now(),
			"user_id":      (*row)["id"].(int64),
		}

		_, err = flcl.InsertRowAny("sessions", toWrite)
		if err != nil {
			errorPage(w, err)
			return
		}

		expires := now.Add(time.Hour * 24 * 30)
		cookie := &http.Cookie{
			Name:    "things_things",
			Value:   sessionCode,
			Path:    "/",
			Expires: expires,
		}
		http.SetCookie(w, cookie)

		http.Redirect(w, r, "/inside", http.StatusTemporaryRedirect)
	}
}

func signout(w http.ResponseWriter, r *http.Request) {
	cookie := &http.Cookie{
		Name:   "things_things",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	}
	http.SetCookie(w, cookie)
	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	trialCode := r.PathValue("code")
	if trialCode != REG_CODE {
		errorPage(w, fmt.Errorf("Invalid code '%s' supplied.", trialCode))
		return
	}
	userId, _, _ := getCurrentUserData(r)
	if userId > 0 {
		http.Redirect(w, r, "/inside", http.StatusTemporaryRedirect)
	}

	if r.Method != http.MethodPost {
		tmpl := template.Must(template.ParseFiles("templates/base.html", "templates/register.html"))
		tmpl.Execute(w, nil)
	} else {

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(r.FormValue("pw1")), bcrypt.DefaultCost)
		if err != nil {
			errorPage(w, err)
			return
		}

		toWrite := map[string]string{
			"email": r.FormValue("email"), "firstname": r.FormValue("firstname"),
			"surname": r.FormValue("surname"), "password": string(hashedPassword),
			"regdate": time.Now().Format("2006-01-02"), "role": r.FormValue("role"),
			"organisation": r.FormValue("organisation"),
		}

		flcl := getFlaarumClient()

		_, err = flcl.InsertRowStr("users", toWrite)
		if err != nil {
			errorPage(w, err)
			return
		}

		messageUser(w, r, "Succesful Registration", `You have sucessfully registered.
Go to your inbox to click mail confirmation
		`)

	}
}

func getUserFullName(userId int64) (string, error) {
	flcl := getFlaarumClient()

	userRow, err := flcl.SearchForOne(fmt.Sprintf(`
		table: users
		where:
		  id = %d
		`, userId))
	if err != nil {
		return "", err
	}

	fullName := (*userRow)["firstname"].(string) + " " + (*userRow)["surname"].(string)
	return fullName, nil
}
