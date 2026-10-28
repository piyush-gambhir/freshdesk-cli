// Package config stores profile metadata separately from keychain credentials.
package config

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
	"unicode"

	"github.com/gofrs/flock"
	"go.yaml.in/yaml/v3"
)

type Profile struct {
	Domain     string `yaml:"domain" json:"domain"`
	TokenStore string `yaml:"token_store" json:"token_store"`
	SecretID   string `yaml:"secret_id" json:"-"`
	ReadOnly   bool   `yaml:"read_only,omitempty" json:"read_only"`
}

type Config struct {
	CurrentProfile string             `yaml:"current_profile,omitempty"`
	Profiles       map[string]Profile `yaml:"profiles"`
}

func DefaultPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".config", "freshdesk-cli", "config.yaml"), nil
}

func Path() (string, error) {
	if value := os.Getenv("FRESHDESK_CONFIG"); value != "" {
		return value, nil
	}
	if base := os.Getenv("XDG_CONFIG_HOME"); base != "" {
		return filepath.Join(base, "freshdesk-cli", "config.yaml"), nil
	}
	return DefaultPath()
}

func Load(path string) (*Config, error) {
	c := &Config{Profiles: map[string]Profile{}}
	b, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return c, nil
	}
	if err != nil {
		return nil, fmt.Errorf("read config: %w", err)
	}
	if err := yaml.Unmarshal(b, c); err != nil {
		return nil, fmt.Errorf("invalid YAML config at %s", path)
	}
	if c.Profiles == nil {
		c.Profiles = map[string]Profile{}
	}
	return c, nil
}

// Update serializes profile and credential transactions across processes.
func Update(ctx context.Context, path string, mutate func(*Config) error) error {
	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		return err
	}
	lock := flock.New(path + ".lock")
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	ok, err := lock.TryLockContext(ctx, 25*time.Millisecond)
	if err != nil {
		return err
	}
	if !ok {
		return fmt.Errorf("configuration is locked by another process")
	}
	defer lock.Unlock()
	if err := os.Chmod(path+".lock", 0600); err != nil {
		return err
	}
	c, err := Load(path)
	if err != nil {
		return err
	}
	if err := mutate(c); err != nil {
		return err
	}
	b, err := yaml.Marshal(c)
	if err != nil {
		return err
	}
	return atomicWrite(path, b)
}

func (c *Config) ProfileName(flag string) string {
	if flag != "" {
		return flag
	}
	if value := os.Getenv("FRESHDESK_PROFILE"); value != "" {
		return value
	}
	return c.CurrentProfile
}

func ValidName(name string) error {
	if name == "" || len(name) > 64 || strings.TrimSpace(name) != name || strings.ContainsAny(name, ":/\\") || strings.ContainsFunc(name, unicode.IsControl) {
		return fmt.Errorf("profile name must be 1-64 characters without control characters, slashes, colons, or whitespace at the ends")
	}
	return nil
}
