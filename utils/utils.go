package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type Envelope map[string]interface{}

func WriteJSON(w http.ResponseWriter, status int, data Envelope) error {
	js, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		http.Error(w, "Failed to marshal JSON", http.StatusInternalServerError)
		return err
	}
	js = append(js, '\n') // Add a newline at the end
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, err = w.Write(js)
	if err != nil {
		http.Error(w, "Failed to write response", http.StatusInternalServerError)
		return err
	}
	return nil
}

func ReadIDParam(r *http.Request) (int64, error) {
	idStr := chi.URLParam(r, "id")
	if idStr == "" {
		return 0, fmt.Errorf("missing ID parameter")
	}
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid ID parameter: %w", err)
	}
	return id, nil
}

func MatchRegex(pattern, str string) bool {
	pattern = regexp.MustCompile(pattern).String()
	matched, err := regexp.MatchString(pattern, str)
	if err != nil {
		return false
	}
	return matched
}