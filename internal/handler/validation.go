package handler

import (
	"fmt"
	"strings"
)

const (
	maxNameLength        = 150
	maxDescriptionLength = 1000
)

func normalizeAndValidatePayload(name string, description string) (string, string, error) {
	normalizedName := strings.TrimSpace(name)
	normalizedDescription := strings.TrimSpace(description)

	if normalizedName == "" {
		return "", "", fmt.Errorf("name must not be empty")
	}

	if len(normalizedName) > maxNameLength {
		return "", "", fmt.Errorf("name must be at most %d characters", maxNameLength)
	}

	if len(normalizedDescription) > maxDescriptionLength {
		return "", "", fmt.Errorf("description must be at most %d characters", maxDescriptionLength)
	}

	return normalizedName, normalizedDescription, nil
}
