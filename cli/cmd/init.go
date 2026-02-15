package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"

	"github.com/user-cube/cluster-bootstrap/cluster-bootstrap/internal/config"
)

var (
	outputDir string
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Interactive setup for git-crypt and per-environment secrets files",
	Long: `Interactively configure git-crypt and create per-environment secrets files.
Ensures .gitattributes includes secrets.*.yaml for git-crypt. Creates plaintext
secrets.<env>.yaml files that git-crypt will encrypt when committed.
Run 'git-crypt init' in the repo first if you have not already.`,
	RunE: runInit,
}

func init() {
	initCmd.Flags().StringVar(&outputDir, "output-dir", ".", "directory for secrets files and .gitattributes")

	rootCmd.AddCommand(initCmd)
}

func runInit(cmd *cobra.Command, args []string) error {
	// Step 0: Verify git-crypt is initialized
	gitCryptDir := filepath.Join(outputDir, ".git", "git-crypt")
	if _, err := os.Stat(gitCryptDir); os.IsNotExist(err) {
		return fmt.Errorf("git-crypt is not initialized in this repo — run 'git-crypt init' first")
	}

	// Step 1: Ensure .gitattributes has git-crypt pattern for secrets.*.yaml
	gitattributesPath := filepath.Join(outputDir, ".gitattributes")
	if err := config.EnsureGitCryptAttributes(outputDir); err != nil {
		return fmt.Errorf("failed to update .gitattributes: %w", err)
	}
	fmt.Printf("Updated %s (secrets.*.yaml → git-crypt)\n", gitattributesPath)

	// Step 2: Prompt for environment names and create per-environment secrets files
	created := 0
	for {
		var env string
		err := huh.NewInput().
			Title("Environment name (leave empty to finish)").
			Value(&env).
			Run()
		if err != nil {
			return fmt.Errorf("prompt failed: %w", err)
		}
		if env == "" {
			break
		}

		outputFile := filepath.Join(outputDir, config.SecretsFileName(env))
		if _, statErr := os.Stat(outputFile); statErr == nil {
			var overwrite bool
			err := huh.NewConfirm().
				Title(fmt.Sprintf("%s already exists. Overwrite?", config.SecretsFileName(env))).
				Value(&overwrite).
				Run()
			if err != nil {
				return fmt.Errorf("prompt failed: %w", err)
			}
			if !overwrite {
				continue
			}
		}

		envSecrets, err := promptEnvironmentSecrets(env)
		if err != nil {
			return err
		}

		data, err := yaml.Marshal(envSecrets)
		if err != nil {
			return fmt.Errorf("failed to marshal secrets: %w", err)
		}

		if err := os.WriteFile(outputFile, data, 0600); err != nil {
			return fmt.Errorf("failed to write %s: %w", outputFile, err)
		}

		fmt.Printf("Created %s (unlock repo with git-crypt to decrypt; files are encrypted in Git)\n", outputFile)
		created++
	}

	if created == 0 {
		return fmt.Errorf("no environments configured")
	}

	fmt.Println("\nYou can now run: cluster-bootstrap bootstrap <environment>")

	return nil
}

func promptEnvironmentSecrets(env string) (*config.EnvironmentSecrets, error) {
	fmt.Printf("\n--- Environment: %s ---\n", env)

	repoURL := "git@github.com:user-cube/cluster-bootstrap.git"
	targetRevision := "main"
	var sshKeyPath string

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Repository SSH URL").
				Value(&repoURL).
				Validate(requiredValidator("repository URL is required")),
			huh.NewInput().
				Title("Target revision (branch/tag)").
				Value(&targetRevision).
				Validate(requiredValidator("target revision is required")),
			huh.NewInput().
				Title("Path to SSH private key file").
				Value(&sshKeyPath).
				Validate(requiredValidator("SSH key path is required")),
		),
	)

	if err := form.Run(); err != nil {
		return nil, fmt.Errorf("prompt failed: %w", err)
	}

	// Read SSH key from filesystem
	sshKeyData, err := os.ReadFile(sshKeyPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read SSH key at %s: %w", sshKeyPath, err)
	}

	envSecrets := &config.EnvironmentSecrets{
		Repo: config.RepoSecrets{
			URL:            repoURL,
			TargetRevision: targetRevision,
			SSHPrivateKey:  string(sshKeyData),
		},
	}

	return envSecrets, nil
}

func requiredValidator(msg string) func(s string) error {
	return func(s string) error {
		if s == "" {
			return fmt.Errorf(msg)
		}
		return nil
	}
}
