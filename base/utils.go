package base

import (
	"fmt"
	"log"
	"strings"
)

// Panic prints a non-nil error and terminates the program
func Panic(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

// QuoteIdent formats identifiers according to PostgreSQL standards
func QuoteIdent(identifier string) string {
	if strings.ToLower(identifier) != identifier {
		return fmt.Sprintf("\"%s\"", identifier)
	}
	return identifier
}
