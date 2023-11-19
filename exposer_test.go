package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

// TestCheckCredentials tests the checkCredentials function
func TestCheckCredentials(t *testing.T) {
	os.Setenv("AUTH_USER", "testuser")
	os.Setenv("AUTH_PASS", "testpass")

	tests := []struct {
		username string
		password string
		want     bool
	}{
		{"testuser", "testpass", true},
		{"wronguser", "testpass", false},
		{"testuser", "wrongpass", false},
		{"wronguser", "wrongpass", false},
	}

	for _, tt := range tests {
		t.Run(tt.username, func(t *testing.T) {
			if got := checkCredentials(tt.username, tt.password); got != tt.want {
				t.Errorf("checkCredentials() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestSplitString tests the splitString function
func TestSplitString(t *testing.T) {
	tests := []struct {
		str  string
		sep  rune
		want [2]string
	}{
		{"key=value", '=', [2]string{"key", "value"}},
		{"another_key=another_value", '=', [2]string{"another_key", "another_value"}},
		// Add more test cases as needed
	}

	for _, tt := range tests {
		t.Run(tt.str, func(t *testing.T) {
			if got := splitString(tt.str, tt.sep); got != tt.want {
				t.Errorf("splitString() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestEnvHandler tests the envHandler function
func TestEnvHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/env", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Setting up environment variables for testing
	os.Setenv("AUTH_USER", "testuser")
	os.Setenv("AUTH_PASS", "testpass")
	os.Setenv("EXPOSER_TEST", "testvalue")

	// Setting up basic auth
	req.SetBasicAuth("testuser", "testpass")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(envHandler)

	handler.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Check the response body
	expected := `{"EXPOSER_TEST":"testvalue"}`
	if strings.TrimSpace(rr.Body.String()) != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}
