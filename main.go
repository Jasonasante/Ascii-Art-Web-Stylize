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

// this handles the homepage
func homePage(w http.ResponseWriter, r *http.Request) {
	// if the handle is not / then it will return an error message
	//also the w.WriteHeader posts the status of the page on the network section when inspecting the page
	if r.URL.Path != "/" {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, "<h1>HTTP Status 404: Page Not Found</h1>")
		return
	}
	// t is the template file. ParseFiles opens up the template file and attempt to validate it.
	// If everything is correct there will be a nil error and a *template
	t, err := template.ParseFiles("template.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "<h1>HTTP Status 500: Internal Server Error(</h1>")
		fmt.Fprint(w, "<p>No banner Selected</p>")
		return
	}
	// t.Executes runs the data input (in this case Text{}) throught the template
	// which is then displayed on website via ResponseWriter
	t.Execute(w, Text{})
}

// this handles the ascii-art output page
func asciiPage(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/ascii-art" {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, "<h1>HTTP Status 404: Page Not Found</h1>")
		return
	}
	r.ParseForm()
	input := r.Form["input"][0]
	if input == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "<h1>HTTP Status 400: Bad Request</h1>")
		fmt.Fprint(w, "<p>Empty input string</p>")
		return
	}
	for _, ele := range input {
		if (ele != 13) && (ele != 10) && (ele < 32 || ele > 126) {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, "<h1>HTTP Status 400: Bad Request</h1>")
			fmt.Fprint(w, "<p>Incorrect character detected</p>")
			return
		}
	}
	if r.FormValue("banner") == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "<h1>HTTP Status 400: Bad Request</h1>")
		fmt.Fprint(w, "<p>No banner Selected</p>")
		return
	}
	banner := r.Form["banner"][0]
	_, err := os.ReadFile(banner + ".txt")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "<h1>HTTP Status 500: Internal Server Error(</h1>")
		fmt.Fprint(w, "<p>No banner file found</p>")
		return
	}
	output := asciiArt(input, banner)
	p := Text{Input: input, Output: output}
	t, err := template.ParseFiles("result.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "<h1>HTTP Status 500: Internal Server Error</h1>")
		fmt.Fprint(w, "<p>No template found</p>")
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
	// this handles the web request to the server with the path /
	// so it covers all paths that the user may visit on the website and it would be processed by handlerFunc
	// for example: t http://localhost:3000/some-other-path
	http.HandleFunc("/ascii-art", asciiPage)
	// starts up a web server listening on port 8080 using the default http handlers
	// so when we run the file, we open a browser and type: "http://localhost:8080/"
	// which is saying 'try to load a web page from this computer at port 8080'
	http.ListenAndServe(":8080", nil)
}
