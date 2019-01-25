package main

import (
	"fmt"
	"testing"
)

func TestError(t *testing.T) {
	err := fmt.Errorf("test")
	response := errorResponse(err)

	if response.StatusCode != 500 {
		t.Errorf("Incorrect status code: %d", response.StatusCode)
	}

	if response.Body != "{\"error\": \"test\"}" {
		t.Errorf("Incorrect message format: %s", response.Body)
	}
}
