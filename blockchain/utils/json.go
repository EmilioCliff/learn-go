package utils

import "fmt"

func JsonStatus(message string) string {
	return fmt.Sprintf("{\"message\": \"%s\"}", message)
}
