package main

import (
	"math/rand"
	"net/http"
	"strconv"
	"text/template"
	"time"
)

var (
	lowerCharSet   = "abcdefghijklmnopqrstuvwxyz"
	upperCharSet   = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	specialCharSet = "!@#$%^&*()_"
	numberCharSet  = "1234567890"
	minSpecialChar = 2
	minUpperChar   = 2
	minNumberChar  = 2
	passwordLength = 10
)

type PasswordResponse struct {
	Passwords []string
}

func main() {
	rand.Seed(time.Now().UnixNano()) // Ensure randomness

	// Serve static files (like your HTML file) from the "static" folder
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// Route to serve the HTML file
	http.HandleFunc("/", serveHTML)

	// Route to handle password generation
	http.HandleFunc("/generate", handleGeneratePasswords)

	// Start the server
	http.ListenAndServe(":8080", nil)
}

func generatePassword() string {

	password := ""

	for i := 0; i < minSpecialChar; i++ {
		random := rand.Intn(len(specialCharSet))
		password = password + string(specialCharSet[random])
	}

	for i := 0; i < minUpperChar; i++ {
		random := rand.Intn(len(upperCharSet))
		password = password + string(upperCharSet[random])
	}

	for i := 0; i < minNumberChar; i++ {
		random := rand.Intn(len(numberCharSet))
		password = password + string(numberCharSet[random])
	}

	totalCharLenWithoutLowerChar := minUpperChar + minSpecialChar + minNumberChar

	remainingCharLen := passwordLength - totalCharLenWithoutLowerChar

	for i := 0; i < remainingCharLen; i++ {
		random := rand.Intn(len(lowerCharSet))
		password = password + string(lowerCharSet[random])
	}

	passwordRune := []rune(password)
	rand.Shuffle(len(passwordRune), func(i, j int) {
		passwordRune[i], passwordRune[j] = passwordRune[j], passwordRune[i]
	})

	password = string(passwordRune)
	return password
}

func serveHTML(w http.ResponseWriter, r *http.Request) {
	// Serve the static HTML page
	tmpl := "static/index.html"
	t, err := template.ParseFiles(tmpl)
	if err != nil {
		http.Error(w, "Error parsing template", http.StatusInternalServerError)
		return
	}
	t.Execute(w, nil)
}

func handleGeneratePasswords(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		// Read the number of passwords from the form
		r.ParseForm()
		numPasswords := r.FormValue("numPasswords")
		numPasswordsInt, err := strconv.Atoi(numPasswords)
		if err != nil || numPasswordsInt <= 0 {
			http.Error(w, "Invalid number of passwords", http.StatusBadRequest)
			return
		}

		// Generate the passwords
		passwords := make([]string, numPasswordsInt)
		for i := 0; i < numPasswordsInt; i++ {
			passwords[i] = generatePassword()
		}

		// Return the passwords back to the HTML page
		tmpl := "static/index.html"
		t, err := template.ParseFiles(tmpl)
		if err != nil {
			http.Error(w, "Error parsing template", http.StatusInternalServerError)
			return
		}

		// Create a struct to hold passwords for rendering
		resp := PasswordResponse{
			Passwords: passwords,
		}
		t.Execute(w, resp)
	}
}
