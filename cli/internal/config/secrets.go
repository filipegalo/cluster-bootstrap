package config

import (
	"fmt"

	"gopkg.in/yaml.v3"

	"github.com/user-cube/cluster-bootstrap/cluster-bootstrap/internal/sops"
)

// EnvironmentSecrets holds the secrets for a single environment.
// Each environment has its own secrets file: secrets.<env>.enc.yaml
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
	return fmt.Sprintf("secrets.%s.enc.yaml", env)
}

// LoadSecrets decrypts and parses a per-environment SOPS-encrypted secrets file.
func LoadSecrets(filePath string, sopsOpts *sops.Options) (*EnvironmentSecrets, error) {
	plaintext, err := sops.Decrypt(filePath, sopsOpts)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt secrets: %w", err)
	}

	var secrets EnvironmentSecrets
	if err := yaml.Unmarshal(plaintext, &secrets); err != nil {
		return nil, fmt.Errorf("failed to parse secrets: %w", err)
	}

	return &secrets, nil
}
