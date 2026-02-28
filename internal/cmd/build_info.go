package cmd

import (
	"encoding/json"
	"fmt"
	"runtime"

	"github.com/spf13/cobra"
	"github.com/steveyegge/gastown/internal/version"
)

var buildInfoJSON bool

// BuildInfo holds structured version and build metadata.
type BuildInfo struct {
	Version   string `json:"version"`
	Build     string `json:"build"`
	Commit    string `json:"commit,omitempty"`
	Branch    string `json:"branch,omitempty"`
	BuildTime string `json:"build_time,omitempty"`
	GoVersion string `json:"go_version"`
	OS        string `json:"os"`
	Arch      string `json:"arch"`
}

var buildInfoCmd = &cobra.Command{
	Use:     "build-info",
	GroupID: GroupDiag,
	Short:   "Print structured version and build information",
	Long: `Print version, commit SHA, build time, Go version, and OS/arch.

Useful for debugging and CI pipelines. Use --json for machine-readable output.`,
	RunE: runBuildInfo,
}

func init() {
	buildInfoCmd.Flags().BoolVar(&buildInfoJSON, "json", false, "Output as JSON")
	rootCmd.AddCommand(buildInfoCmd)
}

func runBuildInfo(cmd *cobra.Command, args []string) error {
	info := BuildInfo{
		Version:   Version,
		Build:     Build,
		Commit:    resolveCommitHash(),
		Branch:    resolveBranch(),
		BuildTime: BuildTime,
		GoVersion: runtime.Version(),
		OS:        runtime.GOOS,
		Arch:      runtime.GOARCH,
	}

	if buildInfoJSON {
		enc := json.NewEncoder(cmd.OutOrStdout())
		enc.SetIndent("", "  ")
		return enc.Encode(info)
	}

	fmt.Fprintf(cmd.OutOrStdout(), "Version:    %s\n", info.Version)
	fmt.Fprintf(cmd.OutOrStdout(), "Build:      %s\n", info.Build)
	if info.Commit != "" {
		fmt.Fprintf(cmd.OutOrStdout(), "Commit:     %s\n", version.ShortCommit(info.Commit))
	}
	if info.Branch != "" {
		fmt.Fprintf(cmd.OutOrStdout(), "Branch:     %s\n", info.Branch)
	}
	if info.BuildTime != "" {
		fmt.Fprintf(cmd.OutOrStdout(), "Build Time: %s\n", info.BuildTime)
	}
	fmt.Fprintf(cmd.OutOrStdout(), "Go Version: %s\n", info.GoVersion)
	fmt.Fprintf(cmd.OutOrStdout(), "OS/Arch:    %s/%s\n", info.OS, info.Arch)
	return nil
}
