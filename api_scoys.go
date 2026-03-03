package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

func verifyAccessCode(w http.ResponseWriter, r *http.Request) {
	flcl := getFlaarumClient()

	accessCode := r.PathValue("code")
	aRow, err := flcl.SearchForOne(fmt.Sprintf(`
		table: seed_companies
		where:
		  access_code = %s
		`, accessCode))
	if err != nil {
		fmt.Println(err)
		http.Error(w, "error: true\nmsg: not_found\n", 400)
		return
	}

	fmt.Fprintf(w, "error: false\ncompany_name: %s\ncacno: %s\n", (*aRow)["name"].(string),
		(*aRow)["cacno"].(string))
}

func submitPhoto(w http.ResponseWriter, r *http.Request) {
	flcl := getFlaarumClient()

	accessCode := r.PathValue("code")
	_, err := flcl.SearchForOne(fmt.Sprintf(`
		table: seed_companies
		where:
		  access_code = %s
		`, accessCode))
	if err != nil {
		fmt.Println(err)
		http.Error(w, "error: true\nmsg: not_found\n", 400)
		return
	}

	photoPath := getPhotosStorePath()
	tmpPath := filepath.Join(photoPath, ".tmp")
	os.MkdirAll(tmpPath, 0777)

	r.ParseMultipartForm(10000 << 20)

	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		http.Error(w, fmt.Sprintf("error: true\nmsg: %s\n", err), http.StatusBadRequest)
		fmt.Println(err)
		return
	}
	defer file.Close()

	rawFile, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, fmt.Sprintf("error: true\nmsg: %s\n", err), http.StatusBadRequest)
		fmt.Println(err)
		return
	}

	fileName := untestedRandomString(100) + filepath.Ext(fileHeader.Filename)
	inPath := filepath.Join(photoPath, fileName)
	os.WriteFile(inPath, rawFile, 0777)

	fmt.Fprintf(w, "error: false\nphoto_name: %s\n", fileName)
}

func submitProductionPlan(w http.ResponseWriter, r *http.Request) {
	flcl := getFlaarumClient()

	accessCode := r.PathValue("code")
	aRow, err := flcl.SearchForOne(fmt.Sprintf(`
		table: seed_companies
		where:
		  access_code = %s
		`, accessCode))
	if err != nil {
		fmt.Println(err)
		http.Error(w, "error: true\nmsg: not_found\n", 400)
		return
	}

	_ = aRow
	toInsert := make(map[string]string)
	for field := range r.PostForm {
		if field == "access_code" {
			continue
		}

		toInsert[field] = r.FormValue(field)
	}

	toInsert["company_id"] = strconv.FormatInt((*aRow)["id"].(int64), 10)

	retId, err := flcl.InsertRowStr("production_plans", toInsert)
	if err != nil {
		http.Error(w, fmt.Sprintf("error: true\nmsg: %s\n", err), http.StatusBadRequest)
		fmt.Println(err)
		return
	}

	fmt.Fprintf(w, "error: false\nfield_id: %d\n", retId)
}
