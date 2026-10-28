package cmd

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/piyush-gambhir/freshdesk-cli/cli-go/internal/build"
	"github.com/piyush-gambhir/freshdesk-cli/cli-go/internal/client"
	"github.com/piyush-gambhir/freshdesk-cli/cli-go/internal/config"
	"github.com/piyush-gambhir/freshdesk-cli/cli-go/internal/secrets"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

func (a *app) authCmd() *cobra.Command {
	c := &cobra.Command{Use: "auth", Short: "Manage API keys and named profiles"}
	c.AddCommand(a.loginCmd(), a.statusCmd(), a.profilesCmd(), a.useCmd(), a.logoutCmd())
	return c
}

func (a *app) configCmd() *cobra.Command {
	c := &cobra.Command{Use: "config", Short: "Inspect configuration and choose a profile"}
	status := a.statusCmd()
	status.Use = "show"
	list := a.profilesCmd()
	list.Use = "list-profiles"
	use := a.useCmd()
	use.Use = "use-profile NAME"
	c.AddCommand(status, list, use)
	return c
}

func (a *app) loginCmd() *cobra.Command {
	var insecure, noVerify, keyStdin, saveReadOnly bool
	c := &cobra.Command{Use: "login", Short: "Verify and save an API key in the OS keychain", Args: cobra.NoArgs, Annotations: effect("local-write"),
		Long: "Get your API key from Freshdesk > profile picture > Profile settings.\nLogin checks GET /api/v2/agents/me before saving. Use --no-verify to save without a request.\nThe default store is the OS keychain. --insecure-storage opts into a plaintext 0600 file.",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, path, err := a.load()
			if err != nil {
				return err
			}
			name := cfg.ProfileName(a.profile)
			if name == "" {
				name = "default"
			}
			if err := config.ValidName(name); err != nil {
				return err
			}
			origin, err := a.target()
			if err != nil {
				return err
			}
			key := a.apiKey
			if key == "" {
				key = os.Getenv("FRESHDESK_API_KEY")
			}
			if keyStdin {
				if a.terminal() {
					return fmt.Errorf("--key-stdin requires piped input")
				}
				b, err := readLimited(a.in, 4096)
				if err != nil {
					return err
				}
				key = strings.TrimSpace(string(b))
			} else if key == "" {
				if a.noInput || !a.terminal() {
					return fmt.Errorf("provide FRESHDESK_API_KEY or --key-stdin for non-interactive login")
				}
				fmt.Fprint(a.errOut, "Freshdesk API key (hidden): ")
				b, err := term.ReadPassword(int(a.in.(*os.File).Fd()))
				fmt.Fprintln(a.errOut)
				if err != nil {
					return err
				}
				key = strings.TrimSpace(string(b))
			}
			if err := validKey(key); err != nil {
				return err
			}
			a.resolvedKey = key
			if !noVerify {
				c := client.Client{Origin: origin, APIKey: key, Transport: a.transport, Timeout: a.timeout, Version: build.Version}
				if _, err := c.Do(cmd.Context(), client.Plan{Method: "GET", URL: origin + "/api/v2/agents/me", Effect: "read"}); err != nil {
					return err
				}
			}
			store := a.store(path)
			backend := secrets.Keychain
			if insecure {
				backend = secrets.File
			}
			var previous config.Profile
			var random [16]byte
			if _, err := rand.Read(random[:]); err != nil {
				return err
			}
			secretID := name + ":" + hex.EncodeToString(random[:])
			var saved secrets.Saved
			err = config.Update(cmd.Context(), path, func(cfg *config.Config) error {
				previous = cfg.Profiles[name]
				if previous.ReadOnly {
					return fmt.Errorf("profile is read-only")
				}
				var err error
				saved, err = store.Save(cmd.Context(), secretID, key, backend, false)
				if err != nil {
					return fmt.Errorf("%w; use --insecure-storage only to explicitly choose plaintext storage, or use FRESHDESK_API_KEY without saving", err)
				}
				cfg.Profiles[name] = config.Profile{Domain: origin, TokenStore: string(saved.Backend), SecretID: secretID, ReadOnly: saveReadOnly}
				cfg.CurrentProfile = name
				return nil
			})
			if err != nil {
				if saved.Backend != "" {
					_ = store.Delete(cmd.Context(), saved.Backend, secretID)
				}
				return err
			}
			if previous.SecretID != "" {
				if err := store.Delete(cmd.Context(), secrets.Backend(previous.TokenStore), previous.SecretID); err != nil {
					fmt.Fprintln(a.errOut, "Warning: previous credential cleanup failed; remove the old keychain/file entry manually")
				}
			}
			return a.print(map[string]any{"saved": true, "verified": !noVerify, "profile": name, "domain": origin, "token_store": saved.Backend, "read_only": saveReadOnly})
		}}
	c.Flags().BoolVar(&insecure, "insecure-storage", false, "Use an owner-only plaintext secrets file instead of the keychain")
	c.Flags().BoolVar(&noVerify, "no-verify", false, "Save without calling Freshdesk")
	c.Flags().BoolVar(&keyStdin, "key-stdin", false, "Read the API key from stdin")
	c.Flags().BoolVar(&saveReadOnly, "save-read-only", false, "Save a read-only restriction in this profile")
	return c
}

func (a *app) statusCmd() *cobra.Command {
	var verify bool
	c := &cobra.Command{Use: "status", Short: "Inspect credential source without revealing or loading the key", Args: cobra.NoArgs, Annotations: effect("read"), RunE: func(cmd *cobra.Command, args []string) error {
		cfg, path, err := a.load()
		if err != nil {
			return err
		}
		name := cfg.ProfileName(a.profile)
		p := cfg.Profiles[name]
		source := "none"
		configured := false
		switch {
		case a.apiKey != "":
			source = "flag"
			configured = true
		case os.Getenv("FRESHDESK_API_KEY") != "":
			source = "environment"
			configured = true
		case p.SecretID != "":
			source = p.TokenStore
			configured = true
		}
		domain := a.domain
		if domain == "" {
			domain = os.Getenv("FRESHDESK_DOMAIN")
		}
		if domain == "" {
			domain = p.Domain
		}
		if domain != "" {
			domain, err = client.Domain(domain)
			if err != nil {
				return err
			}
		} else {
			configured = false
		}
		result := map[string]any{"configured": configured, "profile": name, "domain": domain, "credential_source": source, "config_file": path, "read_only": a.readOnly || p.ReadOnly, "verified": false}
		if verify {
			c, err := a.client()
			if err != nil {
				return err
			}
			res, err := c.Do(cmd.Context(), client.Plan{Method: "GET", URL: c.Origin + "/api/v2/agents/me", Effect: "read"})
			if err != nil {
				return err
			}
			result["verified"] = true
			result["agent"] = res.Data
		}
		return a.print(result)
	}}
	c.Flags().BoolVar(&verify, "verify", false, "Verify access with one GET /agents/me request")
	return c
}

func (a *app) profilesCmd() *cobra.Command {
	return &cobra.Command{Use: "list", Short: "List saved profiles without credentials", Args: cobra.NoArgs, Annotations: effect("read"), RunE: func(cmd *cobra.Command, args []string) error {
		cfg, _, err := a.load()
		if err != nil {
			return err
		}
		names := make([]string, 0, len(cfg.Profiles))
		for n := range cfg.Profiles {
			names = append(names, n)
		}
		sort.Strings(names)
		rows := make([]map[string]any, 0, len(names))
		for _, n := range names {
			p := cfg.Profiles[n]
			rows = append(rows, map[string]any{"profile": n, "domain": p.Domain, "current": n == cfg.CurrentProfile, "token_store": p.TokenStore, "read_only": p.ReadOnly})
		}
		return a.print(rows)
	}}
}

func (a *app) useCmd() *cobra.Command {
	return &cobra.Command{Use: "use NAME", Short: "Select the default saved profile", Args: cobra.ExactArgs(1), Annotations: effect("local-write"), RunE: func(cmd *cobra.Command, args []string) error {
		path, err := config.Path()
		if err != nil {
			return err
		}
		if err := config.Update(cmd.Context(), path, func(cfg *config.Config) error {
			if _, ok := cfg.Profiles[args[0]]; !ok {
				return fmt.Errorf("profile not found")
			}
			cfg.CurrentProfile = args[0]
			return nil
		}); err != nil {
			return err
		}
		return a.print(map[string]string{"current_profile": args[0]})
	}}
}

func (a *app) logoutCmd() *cobra.Command {
	return &cobra.Command{Use: "logout", Short: "Remove the selected profile and saved key from this machine", Long: "Removes local credentials only. To revoke the key remotely, regenerate it in Freshdesk Profile settings.", Args: cobra.NoArgs, Annotations: effect("local-write"), RunE: func(cmd *cobra.Command, args []string) error {
		path, err := config.Path()
		if err != nil {
			return err
		}
		name := ""
		if err := config.Update(cmd.Context(), path, func(cfg *config.Config) error {
			name = cfg.ProfileName(a.profile)
			p, ok := cfg.Profiles[name]
			if !ok {
				return fmt.Errorf("profile not found")
			}
			if p.ReadOnly {
				return fmt.Errorf("profile is read-only")
			}
			if p.SecretID != "" {
				if err := a.store(path).Delete(cmd.Context(), secrets.Backend(p.TokenStore), p.SecretID); err != nil {
					return err
				}
			}
			delete(cfg.Profiles, name)
			if cfg.CurrentProfile == name {
				cfg.CurrentProfile = ""
			}
			return nil
		}); err != nil {
			return err
		}
		return a.print(map[string]any{"profile": name, "logged_out": true})
	}}
}
