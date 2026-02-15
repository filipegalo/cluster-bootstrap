package config

import (
	"os"
	"path/filepath"
	"strings"
)

// GitCryptAttributesPattern is the .gitattributes line for git-crypt to encrypt secrets files.
const GitCryptAttributesPattern = "secrets.*.yaml filter=git-crypt diff=git-crypt"

// EnsureGitCryptAttributes ensures .gitattributes in outputDir contains the git-crypt pattern
// for secrets.*.yaml. It appends the line if missing.
func EnsureGitCryptAttributes(outputDir string) error {
	path := filepath.Join(outputDir, ".gitattributes")
	data, err := os.ReadFile(path)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	content := string(data)
	if strings.Contains(content, GitCryptAttributesPattern) {
		return nil // already configured
	}

	var newContent string
	if len(content) > 0 && !strings.HasSuffix(content, "\n") {
		newContent = content + "\n" + GitCryptAttributesPattern + "\n"
	} else {
		newContent = content + GitCryptAttributesPattern + "\n"
	}

	return os.WriteFile(path, []byte(newContent), 0644)
}
