// Copyright St3ffn 2025, 2026
// SPDX-License-Identifier: MPL-2.0

package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/st3ffn/terraform-provider-aws-appstream/tools/commitlint/internal/output"
	"github.com/st3ffn/terraform-provider-aws-appstream/tools/commitlint/internal/run"
)

const (
	envPrefix = "COMMITLINT"
)

var errViolationsFound = errors.New("conventional commit violations found")

type usageError struct {
	cause error
}

func (e usageError) Error() string {
	return e.cause.Error()
}

func (e usageError) Unwrap() error {
	return e.cause
}

type failLevel string

const (
	failLevelError failLevel = "error"
	failLevelNone  failLevel = "none"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	cmd := newRootCommand()
	cmd.SilenceUsage = true
	cmd.SilenceErrors = true

	if err := cmd.ExecuteContext(ctx); err != nil {
		if errors.Is(err, errViolationsFound) {
			os.Exit(1)
		}

		_, _ = fmt.Fprintf(os.Stderr, "commitlint error: %v\n", err)
		var uerr usageError
		if errors.As(err, &uerr) {
			_, _ = fmt.Fprintln(os.Stderr)
			_ = cmd.Usage()
		}
		os.Exit(2)
	}
}

func newRootCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "commitlint",
		Short: "Validate commit messages against Conventional Commits",
		Long: "Validate commit messages against Conventional Commits.\n\n" +
			"Default behavior (no selection flags): lint commits in current-branch..base,\n" +
			"where base defaults to 'main'. This is equivalent to a PR-style branch diff.",
		RunE: runCommand,
	}
	cmd.SetFlagErrorFunc(func(_ *cobra.Command, err error) error {
		return usageError{cause: err}
	})

	flags := cmd.Flags()
	flags.String("repo", ".", "Repository path")
	flags.Bool("all", false, "Lint all commits in repository history")
	flags.String("branch", "", "Lint all commits reachable from branch tip")
	flags.String("branch-diff", "", "Lint commits in base..branch-diff")
	flags.String("base", "main", "Base branch/ref for default mode or --branch-diff")
	flags.String("from", "", "From commit SHA/revision (inclusive; must be used with --to)")
	flags.String("to", "", "To commit SHA/revision (inclusive; must be used with --from)")
	flags.String("format", string(output.FormatText), "Output format: text|json|github|markdown|rdjsonl")
	flags.String("fail-level", string(failLevelError), "Fail behavior on lint findings: error|none")
	flags.String("color", string(output.ColorAuto), "Color mode for text format: auto|always|never")
	flags.Bool("no-color", false, "Disable colored text output")
	flags.Bool("include-merge-commits", false, "Validate merge commit messages instead of skipping them")

	cmd.AddCommand(newCompletionCommand(cmd))

	return cmd
}

func newCompletionCommand(root *cobra.Command) *cobra.Command {
	cmd := &cobra.Command{
		Use:       "completion [bash|zsh|fish|powershell]",
		Short:     "Generate shell completion scripts",
		Args:      cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
		ValidArgs: []string{"bash", "zsh", "fish", "powershell"},
		RunE: func(cmd *cobra.Command, args []string) error {
			switch args[0] {
			case "bash":
				return root.GenBashCompletionV2(cmd.OutOrStdout(), true)
			case "zsh":
				return root.GenZshCompletion(cmd.OutOrStdout())
			case "fish":
				return root.GenFishCompletion(cmd.OutOrStdout(), true)
			case "powershell":
				return root.GenPowerShellCompletionWithDesc(cmd.OutOrStdout())
			default:
				return fmt.Errorf("unsupported shell %q", args[0])
			}
		},
	}

	return cmd
}

func runCommand(cmd *cobra.Command, _ []string) error {
	cfg, err := buildConfig(cmd)
	if err != nil {
		return err
	}

	format := output.Format(cfg.GetString("format"))
	if !isValidOutputFormat(format) {
		return usageError{cause: fmt.Errorf("unsupported output format %q", format)}
	}
	level := failLevel(cfg.GetString("fail-level"))
	if !isValidFailLevel(level) {
		return usageError{cause: fmt.Errorf("unsupported fail level %q", level)}
	}

	colorMode := output.ColorMode(cfg.GetString("color"))
	if cfg.GetBool("no-color") {
		colorMode = output.ColorNever
	}
	if !isValidColorMode(colorMode) {
		return usageError{cause: fmt.Errorf("unsupported color mode %q", colorMode)}
	}

	ctx := cmd.Context()

	report, err := run.Lint(ctx, run.Options{
		RepoPath:            cfg.GetString("repo"),
		All:                 cfg.GetBool("all"),
		Branch:              cfg.GetString("branch"),
		BranchDiff:          cfg.GetString("branch-diff"),
		Base:                cfg.GetString("base"),
		From:                cfg.GetString("from"),
		To:                  cfg.GetString("to"),
		IncludeMergeCommits: cfg.GetBool("include-merge-commits"),
	})
	if err != nil {
		return err
	}

	if err := output.Write(report, output.Options{
		Format:    format,
		ColorMode: colorMode,
		Writer:    os.Stdout,
	}); err != nil {
		return fmt.Errorf("render output: %w", err)
	}

	if report.InvalidCommits > 0 && level == failLevelError {
		return errViolationsFound
	}

	return nil
}

func buildConfig(cmd *cobra.Command) (*viper.Viper, error) {
	cfg := viper.New()
	cfg.SetEnvPrefix(envPrefix)
	cfg.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	cfg.AutomaticEnv()

	if err := cfg.BindPFlags(cmd.Flags()); err != nil {
		return nil, fmt.Errorf("bind flags: %w", err)
	}
	if err := cfg.BindPFlags(cmd.InheritedFlags()); err != nil {
		return nil, fmt.Errorf("bind inherited flags: %w", err)
	}

	return cfg, nil
}

func isValidOutputFormat(format output.Format) bool {
	switch format {
	case output.FormatText, output.FormatJSON, output.FormatGitHub, output.FormatMarkdown, output.FormatRDJSONL:
		return true
	default:
		return false
	}
}

func isValidColorMode(mode output.ColorMode) bool {
	switch mode {
	case output.ColorAuto, output.ColorAlways, output.ColorNever:
		return true
	default:
		return false
	}
}

func isValidFailLevel(level failLevel) bool {
	switch level {
	case failLevelError, failLevelNone:
		return true
	default:
		return false
	}
}
