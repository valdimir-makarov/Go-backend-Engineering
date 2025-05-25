package pkg

import (
	"errors"
	"strings"

	"github.com/sirupsen/logrus"
)

type CodeSanitizer struct {
	MaxLength int
}

func NewCodeSanitizer(maxLen int) *CodeSanitizer {
	return &CodeSanitizer{MaxLength: maxLen}
}

func (s *CodeSanitizer) SanitizeCode(code, lang string) error {
	if len(code) > s.MaxLength {
		return errors.New("code too long")
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
			logrus.Warnf("Disallowed pattern found in %s code: %s", lang, bad)
			return errors.New("disallowed code pattern: " + bad)
		}
	}
	return nil
}
