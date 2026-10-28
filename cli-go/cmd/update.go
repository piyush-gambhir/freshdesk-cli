package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/piyush-gambhir/freshdesk-cli/cli-go/internal/build"
	"github.com/piyush-gambhir/freshdesk-cli/cli-go/internal/update"
	"github.com/spf13/cobra"
)

func (a *app) updateCmd() *cobra.Command {
	var check bool
	c := &cobra.Command{Use: "update", Short: "Check or install a checksummed GitHub release", Long: "Downloads the release for this platform and verifies its SHA-256 checksum before replacing the binary.\nSource installs in a Go bin directory print the source update command instead.\nThere is no background release check. --check never writes a cache or replaces files.", Args: cobra.NoArgs, Annotations: effect("conditional"), RunE: func(cmd *cobra.Command, args []string) error {
		if !check {
			if a.dryRun {
				return fmt.Errorf("use update --check; API request previews do not apply to self-update")
			}
			ro, err := a.profileReadOnly()
			if err != nil {
				return err
			}
			if ro {
				return fmt.Errorf("self-update is blocked by read-only mode")
			}
		}
		source := update.GitHub(a.transport)
		release, err := source.Latest(cmd.Context())
		if err != nil {
			return err
		}
		result := map[string]any{"current_version": build.Version, "latest_version": release.Version(), "release_url": release.URL, "updated": false, "update_available": update.Newer(release.Version(), build.Version) || !update.IsRelease(build.Version)}
		if check || result["update_available"] == false {
			return a.print(result)
		}
		exe, err := os.Executable()
		if err != nil {
			return err
		}
		exe, err = filepath.EvalSymlinks(exe)
		if err != nil {
			return err
		}
		home, _ := os.UserHomeDir()
		if update.Method(exe, os.Getenv, home) == update.MethodGo {
			result["update_command"] = update.SourceUpdate
			return a.print(result)
		}
		if err := a.confirm("Update freshdesk to " + release.Tag); err != nil {
			return err
		}
		if err := update.CheckWritable(exe, runtime.GOOS); err != nil {
			return err
		}
		bin, err := source.Fetch(cmd.Context(), release, runtime.GOOS, runtime.GOARCH)
		if err != nil {
			return err
		}
		if err := update.Replace(exe, bin, runtime.GOOS); err != nil {
			return err
		}
		result["updated"] = true
		return a.print(result)
	}}
	c.Flags().BoolVar(&check, "check", false, "Only check the latest release; no file writes")
	return c
}

func (a *app) doctorCmd() *cobra.Command {
	var online bool
	c := &cobra.Command{Use: "doctor", Short: "Check local configuration, optionally verify API access", Args: cobra.NoArgs, Annotations: effect("read"), RunE: func(cmd *cobra.Command, args []string) error {
		status := a.statusCmd()
		status.SetContext(cmd.Context())
		if online {
			_ = status.Flags().Set("verify", "true")
		}
		return status.RunE(status, nil)
	}}
	c.Flags().BoolVar(&online, "online", false, "Verify credentials with GET /agents/me")
	return c
}
