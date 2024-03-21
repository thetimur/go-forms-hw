package main

import (
	"crypto/md5"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"strconv"
)

type FormData struct {
	Textarea  string
	Radio     string
	Select    string
	Hidden    string
	Email     string
	Password  string
	Checkbox  bool
	Number    int
	Date      string
	Color     string
	Range     int
	FileName  string
	FileSize  int64
	FileHash  string
	Phone     string
	OtherData map[string]string
}

var formTemplate = template.Must(template.New("form").Parse(`
<!DOCTYPE html>
<html>
<head>
    <title>Form Submission</title>
</head>
<body>
    <h2>Form</h2>
    <form enctype="multipart/form-data" action="/" method="post">
        <label for="textarea">Textarea:</label><br>
        <textarea id="textarea" name="textarea" required></textarea><br><br>

        <input type="radio" id="radio1" name="radio" value="Option1" required>
        <label for="radio1">Option 1</label><br>
        <input type="radio" id="radio2" name="radio" value="Option2">
        <label for="radio2">Option 2</label><br><br>

        <label for="select">Select:</label><br>
        <select id="select" name="select" required>
            <option value="Choice1">Choice 1</option>
            <option value="Choice2">Choice 2</option>
        </select><br><br>

        <input type="file" id="file" name="file" required><br><br>

        <input type="hidden" id="hidden" name="hidden" value="HiddenValue"><br>


        <label for="email">Email:</label><br>
        <input type="email" id="email" name="email" required><br><br>

        <label for="password">Password:</label><br>
        <input type="password" id="password" name="password" required><br><br>

        <input type="checkbox" id="checkbox" name="checkbox">
        <label for="checkbox">Checkbox</label><br><br>

        <label for="number">Number:</label><br>
        <input type="number" id="number" name="number" required><br><br>

        <label for="date">Date:</label><br>
        <input type="date" id="date" name="date" required><br><br>

        <label for="color">Color:</label><br>
        <input type="color" id="color" name="color" required><br><br>

        <label for="range">Range:</label><br>
        <input type="range" id="range" name="range" min="0" max="100" required><br><br>
		
		<label for="tel">Telephone:</label><br>
		<input type="tel" id="tel" name="tel" required><br>
		
		<br><input type="submit" value="Submit">
    </form>

    <h2>Submitted Data</h2>
    <p><strong>Textarea:</strong> {{.Textarea}}</p>
    <p><strong>Radio:</strong> {{.Radio}}</p>
    <p><strong>Select:</strong> {{.Select}}</p>
    <p><strong>Hidden:</strong> {{.Hidden}}</p>
	<p><strong>Email:</strong> {{.Email}}</p>
    <p><strong>Password:</strong> {{.Password}}</p>
    <p><strong>Checkbox:</strong> {{if .Checkbox}}Checked{{else}}Unchecked{{end}}</p>
    <p><strong>Number:</strong> {{.Number}}</p>
    <p><strong>Date:</strong> {{.Date}}</p>
    <p><strong>Color:</strong> {{.Color}}</p>
    <p><strong>Range:</strong> {{.Range}}</p>
    <p><strong>File Size:</strong> {{.FileSize}} bytes</p>
    <p><strong>File Hash (MD5):</strong> {{.FileHash}}</p>
	<p><strong>Phone:</strong> {{.Phone}}</p>
    {{range $key, $value := .OtherData}}
        <p><strong>{{$key}}:</strong> {{$value}}</p>
    {{end}}
</body>
</html>
`))

func main() {
	http.HandleFunc("/", formHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func formHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		if err := r.ParseMultipartForm(10 << 20); err != nil {
			http.Error(w, "Error parsing form", http.StatusInternalServerError)
			return
		}

		file, header, err := r.FormFile("file")
		if err != nil {
			http.Error(w, "Error retrieving the file", http.StatusInternalServerError)
			return
		}
		defer file.Close()

		fileSize := header.Size
		hash := md5.New()
		if _, err := io.Copy(hash, file); err != nil {
			http.Error(w, "Error calculating file hash", http.StatusInternalServerError)
			return
		}
		fileHash := fmt.Sprintf("%x", hash.Sum(nil))

		checkbox := r.FormValue("checkbox") == "on"
		number, _ := strconv.Atoi(r.FormValue("number"))
		rangeValue, _ := strconv.Atoi(r.FormValue("range"))

		formData := FormData{
			Textarea:  r.FormValue("textarea"),
			Radio:     r.FormValue("radio"),
			Select:    r.FormValue("select"),
			Hidden:    r.FormValue("hidden"),
			Email:     r.FormValue("email"),
			Password:  r.FormValue("password"),
			Checkbox:  checkbox,
			Number:    number,
			Date:      r.FormValue("date"),
			Color:     r.FormValue("color"),
			Range:     rangeValue,
			FileSize:  fileSize,
			FileHash:  fileHash,
			Phone:     r.FormValue("tel"),
			OtherData: make(map[string]string),
		}

		for key, values := range r.Form {
			if key != "textarea" && key != "radio" && key != "select" && key != "hidden" && key != "file" {
				formData.OtherData[key] = values[0]
			}
		}

		if err := formTemplate.Execute(w, formData); err != nil {
			http.Error(w, "Error rendering template",
				http.StatusInternalServerError)
			return
		}
	} else {
		if err := formTemplate.Execute(w, FormData{}); err != nil {
			http.Error(w, "Error rendering template", http.StatusInternalServerError)
		}
	}
}
