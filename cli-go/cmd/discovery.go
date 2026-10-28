package cmd

import (
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/piyush-gambhir/freshdesk-cli/cli-go/internal/build"
	"github.com/piyush-gambhir/freshdesk-cli/cli-go/internal/registry"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

type flagInfo struct {
	Name        string `json:"name"`
	Shorthand   string `json:"shorthand,omitempty"`
	Type        string `json:"type"`
	Default     string `json:"default"`
	Role        string `json:"role"`
	Description string `json:"description"`
	Repeatable  bool   `json:"repeatable"`
}

type commandInfo struct {
	ID              string              `json:"id"`
	Path            []string            `json:"path"`
	Usage           string              `json:"usage"`
	Summary         string              `json:"summary"`
	Effect          string              `json:"effect"`
	PossibleEffects []string            `json:"possible_effects"`
	Conditional     bool                `json:"conditional"`
	Flags           []flagInfo          `json:"flags"`
	Operation       *registry.Operation `json:"operation,omitempty"`
}

func flagRole(name string) string {
	switch name {
	case "api-key":
		return "credential"
	case "domain":
		return "connection"
	case "profile":
		return "profile"
	case "read-only", "no-input", "yes", "dry-run", "expect-plan":
		return "safety"
	case "input-json", "set", "body-file", "attach", "upload", "multipart":
		return "input"
	case "output":
		return "local-output"
	default:
		return "query"
	}
}

func flagsInfo(fs *pflag.FlagSet) []flagInfo {
	list := []flagInfo{}
	fs.VisitAll(func(f *pflag.Flag) {
		if f.Hidden || f.Name == "help" {
			return
		}
		list = append(list, flagInfo{Name: f.Name, Shorthand: f.Shorthand, Type: f.Value.Type(), Default: f.DefValue, Role: flagRole(f.Name), Description: f.Usage, Repeatable: f.Value.Type() == "stringArray"})
	})
	return list
}

func describe(c *cobra.Command) commandInfo {
	path := strings.Fields(strings.TrimPrefix(c.CommandPath(), "freshdesk "))
	id := "freshdesk." + strings.Join(path, ".")
	info := commandInfo{ID: id, Path: path, Usage: c.UseLine(), Summary: c.Short, Effect: c.Annotations["effect"], Flags: flagsInfo(c.LocalNonPersistentFlags()), PossibleEffects: []string{c.Annotations["effect"]}}
	if op, ok := registry.Find(c.Annotations["operation"]); ok {
		info.ID = op.ID
		info.Operation = &op
	}
	if info.Effect == "conditional" {
		info.Conditional = true
		info.PossibleEffects = []string{"read", "write"}
		if c.Name() == "update" && c.Parent().Name() == "freshdesk" {
			info.PossibleEffects = []string{"read", "local-write"}
		}
	}
	return info
}

func catalog(root *cobra.Command) []commandInfo {
	list := []commandInfo{}
	var walk func(*cobra.Command)
	walk = func(c *cobra.Command) {
		if c.Runnable() && c.Annotations["effect"] != "" {
			list = append(list, describe(c))
		}
		for _, sub := range c.Commands() {
			if !sub.Hidden {
				walk(sub)
			}
		}
	}
	walk(root)
	sort.Slice(list, func(i, j int) bool { return list[i].ID < list[j].ID })
	return list
}

func (a *app) commandsCmd() *cobra.Command {
	c := &cobra.Command{Use: "commands", Short: "Discover commands, flags, and effects offline", Args: cobra.NoArgs, Annotations: effect("read"), RunE: func(cmd *cobra.Command, args []string) error {
		return a.print(map[string]any{"contract_version": "agent-cli/v1", "catalog_hash": registry.Hash(), "cli": map[string]string{"id": "freshdesk", "binary": "freshdesk", "version": build.Version, "commit": build.Commit},
			"provider":     map[string]any{"env_prefixes": []string{"FRESHDESK_"}, "auth_methods": []string{"api_key"}, "required_env": map[string]any{"all_of": []string{"FRESHDESK_DOMAIN", "FRESHDESK_API_KEY"}}, "config_files": []string{"~/.config/freshdesk-cli/config.yaml", "~/.config/freshdesk-cli/secrets.yaml"}, "verify_argv": []string{"auth", "status", "--verify", "-o", "json"}, "login_argv": []string{"auth", "login"}, "controls": map[string]any{"read_only": "--read-only", "read_only_env": "FRESHDESK_READ_ONLY", "no_input": "--no-input", "no_input_env": "FRESHDESK_NO_INPUT", "json_argv": []string{"-o", "json"}, "credential_precedence": []string{"--api-key", "FRESHDESK_API_KEY", "selected profile"}, "background_requests": false}},
			"global_flags": flagsInfo(cmd.Root().PersistentFlags()), "commands": catalog(cmd.Root())})
	}}
	c.AddCommand(&cobra.Command{Use: "describe ID", Short: "Describe a stable command ID offline", Args: cobra.ExactArgs(1), Annotations: effect("read"), RunE: func(cmd *cobra.Command, args []string) error {
		for _, info := range catalog(cmd.Root()) {
			if info.ID == args[0] || strings.Join(info.Path, " ") == args[0] {
				return a.print(info)
			}
		}
		return fmt.Errorf("unknown command ID")
	}})
	c.AddCommand(&cobra.Command{Use: "resolve -- ARGUMENTS", Short: "Resolve arguments with the real parser, without executing", DisableFlagParsing: true, Annotations: effect("read"), RunE: func(cmd *cobra.Command, args []string) error {
		for i, arg := range args {
			if arg == "--" {
				args = args[i+1:]
				break
			}
		}
		a.format = "json"
		probe := NewRoot(strings.NewReader(""), io.Discard, io.Discard)
		target, rest, err := probe.Find(args)
		if err != nil || target == probe || !target.Runnable() || target.Annotations["effect"] == "" {
			return fmt.Errorf("arguments do not name a runnable command")
		}
		if err := target.ParseFlags(rest); err != nil {
			return err
		}
		positionals := target.Flags().Args()
		if err := target.ValidateArgs(positionals); err != nil {
			return err
		}
		info := describe(target)
		resolved := info.Effect
		if target.CommandPath() == "freshdesk api request" {
			resolved = classify(strings.ToUpper(positionals[0]), positionals[1]).Effect
		}
		if target.CommandPath() == "freshdesk update" {
			check, _ := target.Flags().GetBool("check")
			resolved = "local-write"
			if check {
				resolved = "read"
			}
		}
		flags := map[string]string{}
		target.Flags().Visit(func(f *pflag.Flag) {
			v := f.Value.String()
			if flagRole(f.Name) == "credential" {
				v = "[REDACTED]"
			}
			flags[f.Name] = v
		})
		return a.print(map[string]any{"contract_version": "agent-cli/v1", "catalog_hash": registry.Hash(), "command_id": info.ID, "command": target.CommandPath(), "canonical_path": info.Path, "effect": resolved, "effects": []string{resolved}, "conditional": info.Conditional, "flags": flags, "args": positionals, "dry_run_supported": info.Operation != nil || target.CommandPath() == "freshdesk api request"})
	}})
	return c
}
