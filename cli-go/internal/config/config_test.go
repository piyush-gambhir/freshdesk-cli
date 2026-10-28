package config

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"testing"
)

func TestAtomicUpdatesPreserveConcurrentProfiles(t *testing.T) {
	path := filepath.Join(t.TempDir(), "nested", "config.yaml")
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			err := Update(context.Background(), path, func(c *Config) error {
				c.Profiles[fmt.Sprint(i)] = Profile{Domain: "https://acme.freshdesk.com"}
				return nil
			})
			if err != nil {
				t.Error(err)
			}
		}(i)
	}
	wg.Wait()
	c, err := Load(path)
	if err != nil || len(c.Profiles) != 10 {
		t.Fatalf("profiles lost: %+v %v", c, err)
	}
	if runtime.GOOS != "windows" {
		stat, err := os.Stat(path)
		if err != nil || stat.Mode().Perm() != 0600 {
			t.Fatalf("config permissions: %v %v", stat, err)
		}
	}
}

func TestFailedTransactionDoesNotReplaceConfig(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.yaml")
	_ = os.WriteFile(path, []byte("current_profile: before\nprofiles: {}\n"), 0600)
	err := Update(context.Background(), path, func(c *Config) error { c.CurrentProfile = "after"; return fmt.Errorf("abort") })
	if err == nil {
		t.Fatal("failed transaction succeeded")
	}
	c, err := Load(path)
	if err != nil || c.CurrentProfile != "before" {
		t.Fatalf("failed update changed config: %+v %v", c, err)
	}
}

func TestMalformedConfigDoesNotEchoSecrets(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.yaml")
	_ = os.WriteFile(path, []byte("secret-key: [oops"), 0600)
	if _, err := Load(path); err == nil || err.Error() != "invalid YAML config at "+path {
		t.Fatalf("unexpected config error %v", err)
	}
}
