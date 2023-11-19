package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
)

// EnvVars represents a map of environment variables
type EnvVars map[string]string

const Version = "0.1.0"

func main() {
	var listen = "0.0.0.0:8080"
	fmt.Printf("Starting server version %s on %s\n", Version, listen)
	http.HandleFunc("/env", envHandler)
	http.ListenAndServe(listen, nil)
}

// envHandler handles the HTTP request, checks for basic auth, and returns environment variables as JSON
func envHandler(w http.ResponseWriter, r *http.Request) {
	user, pass, ok := r.BasicAuth()
	if !ok || !checkCredentials(user, pass) {
		w.Header().Set("WWW-Authenticate", `Basic realm="restricted"`)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	envVars := getEnvVars("EXPOSER_")
	jsonData, err := json.Marshal(envVars)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

// checkCredentials verifies the provided username and password against environment variables
func checkCredentials(username, password string) bool {
	expectedUser := os.Getenv("AUTH_USER")
	expectedPass := os.Getenv("AUTH_PASS")

	return username == expectedUser && password == expectedPass
}

// getEnvVars reads all the environment variables with the given prefix and returns them as a map
func getEnvVars(prefix string) EnvVars {
	envVars := make(EnvVars)
	for _, env := range os.Environ() {
		keyVal := splitString(env, '=')
		if strings.HasPrefix(keyVal[0], prefix) {
			envVars[keyVal[0]] = keyVal[1]
		}
	}
	return envVars
}

// splitString splits a string into two parts based on the first occurrence of the separator
func splitString(str string, sep rune) [2]string {
	var split [2]string
	index := 0
	for i, char := range str {
		if char == sep {
			index = i
			break
		}
	}
	split[0] = str[:index]
	split[1] = str[index+1:]
	return split
}
