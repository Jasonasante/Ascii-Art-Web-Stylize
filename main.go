package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"strings"
)

type Text struct {
	Input  string
	Output string
}

func errorHandler(w http.ResponseWriter, r *http.Request, status int) {
	w.WriteHeader(status)
	if status == http.StatusNotFound {
		fmt.Fprint(w, "<h1>HTTP status 404: Page Not Found</h1>")
	}
	if status == http.StatusInternalServerError {
		fmt.Fprint(w, "<h1>HTTP status 500: Internal Server Error</h1>")
	}
	if status == http.StatusBadRequest {
		fmt.Fprintln(w, "<h1>HTTP status 400: Bad Request."+"\n"+"Please pick a banner, use printable characters or enter text.</h1>")
		fmt.Fprintln(w, "Printable Characters:")
		for i := 32; i <= 126; i++ {
			fmt.Fprint(w, string(byte(i)))
		}

	}
}

// this handles the homepage
func homePage(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		errorHandler(w, r, http.StatusNotFound)
		return
	}
	t, err := template.ParseFiles("template.html")
	if err != nil {
		errorHandler(w, r, http.StatusInternalServerError)
		return
	}
	t.Execute(w, Text{})
}

// this handles the ascii-art output page
func asciiPage(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/ascii-art" {
		errorHandler(w, r, http.StatusNotFound)
		return
	}
	r.ParseForm()
	input := r.Form["input"][0]
	if input == "" {
		errorHandler(w, r, http.StatusBadRequest)
		return
	}
	for _, ele := range input {
		if (ele != 13) && (ele != 10) && (ele < 32 || ele > 126) {
			errorHandler(w, r, http.StatusBadRequest)
			return
		}
	}
	if r.FormValue("banner") == "" {
		errorHandler(w, r, http.StatusBadRequest)
		return
	}
	banner := r.Form["banner"][0]
	_, err := os.ReadFile(banner + ".txt")
	if err != nil {
		errorHandler(w, r, http.StatusInternalServerError)
		return
	}
	output := asciiArt(input, banner)
	p := Text{Input: input, Output: output}
	t, err := template.ParseFiles("result.html")
	if err != nil {
		errorHandler(w, r, http.StatusInternalServerError)
		return
	}
	t.Execute(w, p)
	// t.Execute(w, Text{input, output}) this is another way of executing
}

// this handles the conversion from the a string to a GUI
func asciiArt(s string, b string) string {
	var emptyString string
	var inputString []string
	Content, _ := os.ReadFile(b + ".txt")
	asciiSlice2 := make([][]string, 95)
	s = strings.Replace(s, "\r\n", "\\n", -1)
	inputString = strings.Split(s, "\\n")
	for i := 0; i < len(asciiSlice2); i++ {
		asciiSlice2[i] = make([]string, 9)
	}
	var bubbleCount int
	count := 0
	for i := 1; i < len(Content); i++ {
		if Content[i] == '\n' && bubbleCount <= 94 {
			asciiSlice2[bubbleCount][count] = emptyString
			emptyString = ""
			count++
		}
		if count == 9 {
			count = 0
			bubbleCount++
		} else {
			if Content[i] != '\n' && Content[i] != '\r' {
				emptyString += string(Content[i])
			}
		}
	}
	var outputStr string
	var tempOutput [][]string
	for _, str := range inputString {
		for _, aRune := range str {
			tempOutput = append(tempOutput, asciiSlice2[aRune-rune(32)])
		}
		for i := range tempOutput[0] {
			for _, char := range tempOutput {
				outputStr += char[i]
			}
			outputStr += "\n"
		}
		tempOutput = nil
	}
	return outputStr
}

// this opens the server
func main() {
	fmt.Println("Starting Server at Port 8080")
	fmt.Println("now open a broswer and enter: localhost:8080 into the URL")
	http.HandleFunc("/", homePage)
	http.HandleFunc("/ascii-art", asciiPage)
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("css"))))
	http.ListenAndServe(":8080", nil)
}
