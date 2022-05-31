package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"strings"
	"log"
)

type Text struct {
	Input    string
	Output   string
	ErrorNum int
	ErrorMes string
}

//handles error messages
func errorHandler(w http.ResponseWriter, r *http.Request, status int) {
	w.WriteHeader(status)
	if status == http.StatusNotFound {
		t, err := template.ParseFiles("errorPage.html")
		if err != nil {
			errorHandler(w, r, http.StatusInternalServerError)
			return
		}
		em := "HTTP status 404: Page Not Found"
		p := Text{ErrorNum: status, ErrorMes: em}
		t.Execute(w, p)
	}
	if status == http.StatusInternalServerError {
		t, err := template.ParseFiles("errorPage.html")
		if err!=nil{
			fmt.Fprint(w, "HTTP status 500: Internal Server Error -missing errorPage.html file")
		}
		em := "HTTP status 500: Internal Server Error"
		p := Text{ErrorNum: status, ErrorMes: em}
		t.Execute(w, p)
	}
	if status == http.StatusBadRequest {
		t, err := template.ParseFiles("errorPage.html")
		if err != nil {
			errorHandler(w, r, http.StatusInternalServerError)
			return
		}
		print:=""
		for i := 32; i <= 126; i++ {
			print += string(byte(i))
		}
		em := ("HTTP status 400: Bad Request."+ "\n"+"Please pick a banner, use printable characters or enter text."+"\n"+"Printable Characters:"+"\n"+print)
		p := Text{ErrorNum: status, ErrorMes: em}
		t.Execute(w, p)

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
	
	http.HandleFunc("/", homePage)
	http.HandleFunc("/ascii-art", asciiPage)
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("css"))))
	
	fmt.Printf("Starting server at port 8080\n")
	fmt.Println("Go on http://localhost:8080") // Prints the link of the website on the command prompt
	fmt.Printf("\nTo shutdown the server and exit the code hit \"crtl+C\"\n")
	if err := http.ListenAndServe(":8080", nil); err != nil { // Launches the server on port 8080 if port 8080 is not already busy, else quit
		log.Fatal(err)
		 }
}
