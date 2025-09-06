package common

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
)

// CheckBodyParam checks if the body parameter is provided in JSON or file format.
func CheckBodyParam(bodyJSON string, bodyFile string) (string, error) {
	var body []byte
	if bodyJSON != "" {
		bodyJSON = strings.Trim(bodyJSON, "'")
		body = []byte(bodyJSON)
	} else {
		if bodyFile != "" {
			data, err := os.ReadFile(bodyFile)
			if err != nil {
				return "", err
			}
			body = data
		}
	}

	if len(body) == 0 {
		return "", nil
	}

	ok, err := isValidJSON(body)
	if err != nil {
		return "", err
	}
	if !ok {
		return "", fmt.Errorf("invalid JSON format")
	}

	return string(body), nil
}

func isValidJSON(data []byte) (bool, error) {
	var js interface{}
	if err := json.Unmarshal(data, &js); err != nil {
		return false, err
	}
	return true, nil
}

// NewID generates a new ID for each request.
func NewID() int64 {
	ns := time.Now().UnixMilli() * 1000000
	return ns + rand.Int63n(1000000)
}
