package main

import (
	"errors"
	"log"
	"strings"
)

type CodeSanitizer struct {
	Maxlength int
}

func NewCodeSanitizer(maxLen int) *CodeSanitizer {

	return &CodeSanitizer{
		Maxlength: maxLen,
	}
}
func (s *CodeSanitizer) Sanitizer(code, lang string) error {
	if len(code) > s.Maxlength {
		log.Fatal("the Code is too long")
	}
	disallowed := []string{}
	switch lang {

	case "python":
		disallowed = []string{"import os", "import sys", "open(", "__import__", "eval(", "exec("}
	case "go":
		disallowed = []string{"os.", "syscall", "exec.", "unsafe"}
	case "nodejs":
		disallowed = []string{"require(", "child_process", "eval(", "fs."}
	default:
		return errors.New("unsupported language for sanitization")

	}

	for _, bad := range disallowed {
		if strings.Contains(code, bad) {
			log.Fatal("the Code is Faulty")
			return errors.New("disallowed pattern ")
		}
	}
	return nil

}
