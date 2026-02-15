package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// EnvironmentSecrets holds the secrets for a single environment.
// Each environment has its own secrets file: secrets.<env>.yaml (git-crypt encrypted in repo).
type EnvironmentSecrets struct {
	Repo RepoSecrets `yaml:"repo"`
}

// RepoSecrets holds git repository credentials.
type RepoSecrets struct {
	URL            string `yaml:"url"`
	TargetRevision string `yaml:"targetRevision"`
	SSHPrivateKey  string `yaml:"sshPrivateKey"`
}

// SecretsFileName returns the secrets file name for the given environment.
func SecretsFileName(env string) string {
	return fmt.Sprintf("secrets.%s.yaml", env)
}

// LoadSecrets reads and parses a per-environment secrets file.
// With git-crypt, unlock the repo first so these files are decrypted on disk.
func LoadSecrets(filePath string) (*EnvironmentSecrets, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read secrets: %w", err)
	}

	var secrets EnvironmentSecrets
	if err := yaml.Unmarshal(data, &secrets); err != nil {
		return nil, fmt.Errorf("failed to parse secrets: %w", err)
	}

	return &secrets, nil
}
