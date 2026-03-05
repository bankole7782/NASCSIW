package main

import (
	"fmt"
	"html/template"
	"math/rand/v2"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

func errorPage(w http.ResponseWriter, err error) {
	_, fn, line, _ := runtime.Caller(1)
	type Context struct {
		Message    template.HTML
		SourceFn   string
		SourceLine int
		DEVELOPER  bool
	}

	var ctx Context
	if os.Getenv("SAE_DEV") == "true" {
		msg := fmt.Sprintf("%+v", err)
		msg = strings.ReplaceAll(msg, "\n", "<br>")
		msg = strings.ReplaceAll(msg, "\t", "&nbsp;&nbsp;&nbsp;")
		ctx = Context{template.HTML(msg), fn, line, true}
	} else {
		msg := err.Error()
		ctx = Context{template.HTML(msg), fn, line, false}
	}
	tmpl := template.Must(template.ParseFiles("templates/base.html", "templates/error-page.html"))
	tmpl.Execute(w, ctx)
}

func messageUser(w http.ResponseWriter, r *http.Request, title string, msg template.HTML) {
	type Context struct {
		MsgTitle string
		Message  template.HTML
		Fullname string
	}

	ctx := Context{title, msg, ""}
	userId, _, _ := getCurrentUserData(r)
	tmpl := template.Must(template.ParseFiles("templates/base.html", "templates/msg.html"))
	if userId != 0 {
		fullname, _ := getUserFullName(userId)
		ctx = Context{title, msg, fullname}
		tmpl = template.Must(template.ParseFiles("templates/inside_base.html", "templates/msg.html"))
	}
	tmpl.Execute(w, ctx)
}

func randRange(min, max int) int {
	return rand.IntN(max-min) + min
}

func untestedRandomString(length int) string {
	const charset = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.IntN(len(charset))]
	}
	return string(b)
}

func getPhotosStorePath() string {
	if os.Getenv("SAE_DEV") == "true" {
		homeDir, _ := os.UserHomeDir()
		retPath := filepath.Join(homeDir, "op_nascsiw")
		return retPath
	} else {
		return "/opt/app/photos"
	}
}
