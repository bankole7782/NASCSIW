package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

func verifyAccessCode(w http.ResponseWriter, r *http.Request) {
	flcl := getFlaarumClient()
	w.Header().Set("Content-Type", "application/json")

	accessCode := r.PathValue("code")
	aRow, err := flcl.SearchForOne(fmt.Sprintf(`
		table: seed_companies
		where:
		  access_code = %s
		`, accessCode))
	if err != nil {
		fmt.Println(err)

		w.WriteHeader(http.StatusBadRequest)
		data := map[string]string{"error": "true", "msg": "not_found"}
		jsonStr, _ := json.Marshal(data)
		w.Write(jsonStr)

		return
	}

	data := map[string]string{
		"error":        "false",
		"company_name": (*aRow)["name"].(string),
		"cacno":        (*aRow)["cacno"].(string),
	}
	jsonStr, _ := json.Marshal(data)
	w.WriteHeader(http.StatusOK)
	w.Write(jsonStr)
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

		w.WriteHeader(http.StatusBadRequest)
		data := map[string]string{"error": "true", "msg": "not_found"}
		jsonStr, _ := json.Marshal(data)
		w.Write(jsonStr)
		return
	}

	photoPath := getPhotosStorePath()
	tmpPath := filepath.Join(photoPath, ".tmp")
	os.MkdirAll(tmpPath, 0777)

	r.ParseMultipartForm(10000 << 20)

	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		fmt.Println(err)

		w.WriteHeader(http.StatusBadRequest)
		data := map[string]string{"error": "true", "msg": err.Error()}
		jsonStr, _ := json.Marshal(data)
		w.Write(jsonStr)
		return
	}
	defer file.Close()

	rawFile, err := io.ReadAll(file)
	if err != nil {
		fmt.Println(err)

		w.WriteHeader(http.StatusBadRequest)
		data := map[string]string{"error": "true", "msg": err.Error()}
		jsonStr, _ := json.Marshal(data)
		w.Write(jsonStr)
		return
	}

	fileName := untestedRandomString(100) + filepath.Ext(fileHeader.Filename)
	inPath := filepath.Join(photoPath, fileName)
	os.WriteFile(inPath, rawFile, 0777)

	data := map[string]string{"error": "false", "photo_name": fileName}
	jsonStr, _ := json.Marshal(data)
	w.WriteHeader(http.StatusOK)
	w.Write(jsonStr)
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

		w.WriteHeader(http.StatusBadRequest)
		data := map[string]string{"error": "true", "msg": "not_found"}
		jsonStr, _ := json.Marshal(data)
		w.Write(jsonStr)
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
		fmt.Println(err)

		w.WriteHeader(http.StatusBadRequest)
		data := map[string]string{"error": "true", "msg": err.Error()}
		jsonStr, _ := json.Marshal(data)
		w.Write(jsonStr)
		return
	}

	data := map[string]string{"error": "false", "field_id": strconv.FormatInt(retId, 10)}
	jsonStr, _ := json.Marshal(data)
	w.WriteHeader(http.StatusOK)
	w.Write(jsonStr)
}
