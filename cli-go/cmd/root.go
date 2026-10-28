package cmd

import (
	"bufio"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
	"unicode"

	"github.com/piyush-gambhir/freshdesk-cli/cli-go/internal/build"
	"github.com/piyush-gambhir/freshdesk-cli/cli-go/internal/client"
	"github.com/piyush-gambhir/freshdesk-cli/cli-go/internal/config"
	"github.com/piyush-gambhir/freshdesk-cli/cli-go/internal/output"
	"github.com/piyush-gambhir/freshdesk-cli/cli-go/internal/registry"
	"github.com/piyush-gambhir/freshdesk-cli/cli-go/internal/secrets"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"golang.org/x/term"
)

type app struct {
	format, profile, domain, apiKey, expectPlan    string
	timeout                                        time.Duration
	noInput, quiet, verbose, readOnly, dryRun, yes bool
	in                                             io.Reader
	out, errOut                                    io.Writer
	transport                                      http.RoundTripper
	resolvedKey                                    string
}

func envBool(key string) bool {
	value := os.Getenv(key)
	return value == "1" || strings.EqualFold(value, "true")
}

func NewRoot(in io.Reader, out, errOut io.Writer) *cobra.Command {
	return newRoot(&app{in: in, out: out, errOut: errOut})
}

func newRoot(a *app) *cobra.Command {
	root := &cobra.Command{Use: "freshdesk", Short: "Read and manage Freshdesk from your terminal", SilenceUsage: true, SilenceErrors: true,
		Long: "An independent Freshdesk REST API client for people and coding agents.\nData goes to stdout; diagnostics go to stderr. Writes support request previews and plan hashes.",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			a.readOnly = a.readOnly || envBool("FRESHDESK_READ_ONLY")
			a.noInput = a.noInput || envBool("FRESHDESK_NO_INPUT")
			if !output.Valid(a.format) {
				return fmt.Errorf("unsupported output %q; use table, json, yaml, csv", a.format)
			}
			if a.timeout <= 0 {
				return fmt.Errorf("--timeout must be positive")
			}
			if effect := cmd.Annotations["effect"]; cmd.Runnable() && effect == "" {
				return fmt.Errorf("command has no effect declaration")
			}
			if cmd.Annotations["effect"] == "local-write" && a.dryRun {
				return fmt.Errorf("--dry-run is supported for API operations, not local profile changes")
			}
			if cmd.Annotations["effect"] == "local-write" {
				ro, err := a.profileReadOnly()
				if err != nil {
					return err
				}
				if ro {
					return fmt.Errorf("local changes are blocked by read-only mode")
				}
			}
			return nil
		}}
	root.SetIn(a.in)
	root.SetOut(a.out)
	root.SetErr(a.errOut)
	f := root.PersistentFlags()
	f.StringVarP(&a.format, "output", "o", "table", "Output format: table, json, yaml, csv")
	f.StringVar(&a.profile, "profile", "", "Named profile (or FRESHDESK_PROFILE)")
	f.StringVar(&a.domain, "domain", "", "Account slug or HTTPS hostname (or FRESHDESK_DOMAIN)")
	f.StringVar(&a.apiKey, "api-key", "", "API key override (prefer FRESHDESK_API_KEY or auth login)")
	f.DurationVar(&a.timeout, "timeout", 30*time.Second, "HTTP request timeout")
	f.BoolVar(&a.noInput, "no-input", envBool("FRESHDESK_NO_INPUT"), "Never prompt; writes require --yes")
	f.BoolVarP(&a.quiet, "quiet", "q", envBool("FRESHDESK_QUIET"), "Suppress informational stderr output")
	f.BoolVarP(&a.verbose, "verbose", "v", envBool("FRESHDESK_VERBOSE"), "Log HTTP methods, paths, and status to stderr")
	f.BoolVar(&a.readOnly, "read-only", envBool("FRESHDESK_READ_ONLY"), "Block API writes, profile changes, and self-update")
	f.BoolVar(&a.dryRun, "dry-run", false, "Validate and print an API request plan without credentials or network")
	f.BoolVarP(&a.yes, "yes", "y", false, "Confirm API writes and self-update without prompting")
	f.StringVar(&a.expectPlan, "expect-plan", "", "Execute only if the rebuilt request has this SHA-256 plan hash")
	for _, op := range registry.Operations {
		a.addOperation(root, op)
	}
	root.AddCommand(a.authCmd(), a.configCmd(), a.commandsCmd(), a.apiCmd(), a.updateCmd(), a.doctorCmd())
	root.AddCommand(&cobra.Command{Use: "version", Short: "Print version, commit, and build date", Args: cobra.NoArgs, Annotations: effect("read"), RunE: func(cmd *cobra.Command, args []string) error {
		return a.print(map[string]any{"version": build.Version, "commit": build.Commit, "date": build.Date, "catalog_hash": registry.Hash()})
	}})
	root.AddCommand(&cobra.Command{Use: "completion [bash|zsh|fish|powershell]", Short: "Generate shell completion", Args: cobra.ExactArgs(1), Annotations: effect("read"), ValidArgs: []string{"bash", "zsh", "fish", "powershell"}, RunE: func(cmd *cobra.Command, args []string) error {
		switch args[0] {
		case "bash":
			return root.GenBashCompletion(a.out)
		case "zsh":
			return root.GenZshCompletion(a.out)
		case "fish":
			return root.GenFishCompletion(a.out, true)
		case "powershell":
			return root.GenPowerShellCompletionWithDesc(a.out)
		default:
			return fmt.Errorf("unsupported shell %q", args[0])
		}
	}})
	root.CompletionOptions.DisableDefaultCmd = true
	return root
}

func effect(value string) map[string]string { return map[string]string{"effect": value} }

func Run(ctx context.Context, args []string, in io.Reader, out, errOut io.Writer) int {
	return execute(ctx, &app{in: in, out: out, errOut: errOut}, args)
}

func execute(ctx context.Context, a *app, args []string) int {
	root := newRoot(a)
	root.SetArgs(args)
	err := root.ExecuteContext(ctx)
	if err == nil {
		return 0
	}
	a.format = errorFormat(root, args, a.format)
	apiErr := &client.APIError{Kind: "usage", Message: err.Error(), Outcome: "failed"}
	var remote *client.APIError
	if errors.As(err, &remote) {
		apiErr = remote
	}
	bytes, _ := json.Marshal(map[string]any{"contract_version": "agent-cli/v1", "error": apiErr})
	message := string(bytes)
	for _, key := range []string{a.apiKey, a.resolvedKey, os.Getenv("FRESHDESK_API_KEY")} {
		if key != "" {
			encoded, _ := json.Marshal(key)
			message = strings.ReplaceAll(message, string(encoded[1:len(encoded)-1]), "[REDACTED]")
			message = strings.ReplaceAll(message, base64.StdEncoding.EncodeToString([]byte(key+":X")), "[REDACTED]")
		}
	}
	if a.format == "json" || a.format == "yaml" {
		var safe any
		_ = json.Unmarshal([]byte(message), &safe)
		_ = output.Print(a.errOut, a.format, safe)
	} else {
		var safe struct {
			Error client.APIError `json:"error"`
		}
		_ = json.Unmarshal([]byte(message), &safe)
		fmt.Fprintln(a.errOut, "Error:", safe.Error.Message)
	}
	return 1
}

// Recover only the requested error format after parsing fails. The probe knows
// each flag's arity but treats values as strings so a bad --page cannot hide a
// later -o json. It does not execute a command or read configuration.
func errorFormat(root *cobra.Command, args []string, fallback string) string {
	target, rest, err := root.Find(args)
	if err != nil {
		target, rest = root, args
	}
	if target.DisableFlagParsing {
		return fallback
	}
	fs := pflag.NewFlagSet("error-format", pflag.ContinueOnError)
	fs.SetOutput(io.Discard)
	fs.ParseErrorsWhitelist.UnknownFlags = true
	add := func(f *pflag.Flag) {
		if fs.Lookup(f.Name) != nil {
			return
		}
		fs.StringP(f.Name, f.Shorthand, f.DefValue, "")
		fs.Lookup(f.Name).NoOptDefVal = f.NoOptDefVal
	}
	target.Flags().VisitAll(add)
	target.InheritedFlags().VisitAll(add)
	_ = fs.Parse(rest)
	if f := fs.Lookup("output"); f != nil && f.Changed && output.Valid(f.Value.String()) {
		return f.Value.String()
	}
	return fallback
}

func (a *app) print(value any) error { return output.Print(a.out, a.format, value) }

func (a *app) terminal() bool { f, ok := a.in.(*os.File); return ok && term.IsTerminal(int(f.Fd())) }

func (a *app) confirm(action string) error {
	if a.yes {
		return nil
	}
	if a.noInput || !a.terminal() {
		return fmt.Errorf("%s requires --yes when input is non-interactive", action)
	}
	fmt.Fprintf(a.errOut, "%s? [y/N] ", action)
	line, err := bufio.NewReader(a.in).ReadString('\n')
	if err != nil || !strings.EqualFold(strings.TrimSpace(line), "y") && !strings.EqualFold(strings.TrimSpace(line), "yes") {
		return fmt.Errorf("cancelled")
	}
	return nil
}

func (a *app) load() (*config.Config, string, error) {
	path, err := config.Path()
	if err != nil {
		return nil, "", err
	}
	cfg, err := config.Load(path)
	return cfg, path, err
}

func (a *app) store(path string) secrets.Store {
	defaultPath, _ := config.DefaultPath()
	return secrets.Store{Service: secrets.ServiceFor("freshdesk-cli", path, defaultPath), FilePath: filepath.Join(filepath.Dir(path), "secrets.yaml")}
}

func (a *app) profileReadOnly() (bool, error) {
	if a.readOnly {
		return true, nil
	}
	cfg, _, err := a.load()
	if err != nil {
		return false, err
	}
	return cfg.Profiles[cfg.ProfileName(a.profile)].ReadOnly, nil
}

func (a *app) target() (string, error) {
	value := a.domain
	if value == "" {
		value = os.Getenv("FRESHDESK_DOMAIN")
	}
	if value == "" {
		cfg, _, err := a.load()
		if err != nil {
			return "", err
		}
		value = cfg.Profiles[cfg.ProfileName(a.profile)].Domain
	}
	if value == "" {
		return "", fmt.Errorf("set FRESHDESK_DOMAIN, pass --domain, or run freshdesk auth login")
	}
	return client.Domain(value)
}

func (a *app) client() (*client.Client, error) {
	origin, err := a.target()
	if err != nil {
		return nil, err
	}
	key := a.apiKey
	if key == "" {
		key = os.Getenv("FRESHDESK_API_KEY")
	}
	if key == "" {
		cfg, path, err := a.load()
		if err != nil {
			return nil, err
		}
		name := cfg.ProfileName(a.profile)
		p, ok := cfg.Profiles[name]
		if !ok || p.SecretID == "" {
			return nil, fmt.Errorf("no API key: set FRESHDESK_API_KEY or run freshdesk auth login")
		}
		key, err = a.store(path).Load(secrets.Backend(p.TokenStore), p.SecretID)
		if err != nil {
			return nil, err
		}
	}
	if err := validKey(key); err != nil {
		return nil, err
	}
	a.resolvedKey = key
	c := &client.Client{Origin: origin, APIKey: key, Timeout: a.timeout, Transport: a.transport, Version: build.Version}
	if a.verbose {
		c.Log = a.errOut
	}
	return c, nil
}

func validKey(key string) error {
	if len(key) == 0 || len(key) > 4096 || strings.Contains(key, ":") || strings.ContainsFunc(key, func(r rune) bool { return unicode.IsControl(r) || unicode.IsSpace(r) }) {
		return fmt.Errorf("API key must be nonempty and contain no whitespace, colons, or control characters")
	}
	return nil
}
