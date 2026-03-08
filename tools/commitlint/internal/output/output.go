// Copyright St3ffn 2025, 2026
// SPDX-License-Identifier: MPL-2.0

package output

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/st3ffn/terraform-provider-aws-appstream/tools/commitlint/internal/model"
)

// Format defines report serialization format.
type Format string

const (
	// FormatText renders a human-readable report.
	FormatText Format = "text"
	// FormatJSON renders JSON suitable for machines.
	FormatJSON Format = "json"
	// FormatGitHub renders GitHub workflow annotations.
	FormatGitHub Format = "github"
	// FormatMarkdown renders markdown suitable for PR comments.
	FormatMarkdown Format = "markdown"
	// FormatRDJSONL renders reviewdog rdjsonl diagnostics.
	FormatRDJSONL Format = "rdjsonl"
)

// ColorMode controls ANSI color usage in text output.
type ColorMode string

const (
	// ColorAuto enables colors only when stdout is a terminal.
	ColorAuto ColorMode = "auto"
	// ColorAlways forces colors.
	ColorAlways ColorMode = "always"
	// ColorNever disables colors.
	ColorNever ColorMode = "never"
)

// Options configures report rendering.
type Options struct {
	Format    Format
	ColorMode ColorMode
	Writer    io.Writer
}

// Write renders report according to options.
func Write(report model.Report, opts Options) error {
	if opts.Writer == nil {
		opts.Writer = os.Stdout
	}

	switch opts.Format {
	case FormatText:
		return writeText(report, opts.Writer, opts.ColorMode)
	case FormatJSON:
		enc := json.NewEncoder(opts.Writer)
		enc.SetIndent("", "  ")
		enc.SetEscapeHTML(false)
		return enc.Encode(report)
	case FormatGitHub:
		return writeGitHub(report, opts.Writer)
	case FormatMarkdown:
		return writeMarkdown(report, opts.Writer)
	case FormatRDJSONL:
		return writeRDJSONL(report, opts.Writer)
	default:
		return fmt.Errorf("unsupported output format %q", opts.Format)
	}
}

func writeText(report model.Report, w io.Writer, colorMode ColorMode) error {
	green := textColorFn(color.New(color.FgGreen), colorMode)
	red := textColorFn(color.New(color.FgRed), colorMode)
	yellow := textColorFn(color.New(color.FgYellow), colorMode)

	if report.InvalidCommits == 0 {
		_, err := fmt.Fprintf(
			w,
			"%s checked=%d skipped=%d invalid=%d\n",
			green("PASS: all checked commits follow Conventional Commits"),
			report.CheckedCommits,
			report.SkippedCommits,
			report.InvalidCommits,
		)
		return err
	}

	if _, err := fmt.Fprintf(
		w,
		"%s checked=%d skipped=%d invalid=%d\n",
		red("FAIL: Conventional Commit violations found"),
		report.CheckedCommits,
		report.SkippedCommits,
		report.InvalidCommits,
	); err != nil {
		return err
	}

	for _, result := range report.Results {
		if result.Valid {
			continue
		}
		if _, err := fmt.Fprintf(
			w,
			"- %s - %s\n",
			yellow(result.ShortHash),
			result.Subject,
		); err != nil {
			return err
		}
		for _, violation := range result.Violations {
			if _, err := fmt.Fprintf(w, "  * [%s] %s\n", violation.Code, violation.Message); err != nil {
				return err
			}
		}
	}

	return nil
}

func writeGitHub(report model.Report, w io.Writer) error {
	for _, result := range report.Results {
		if result.Valid {
			continue
		}
		for _, violation := range result.Violations {
			msg := fmt.Sprintf("%s - %s [%s] %s", result.ShortHash, result.Subject, violation.Code, violation.Message)
			if _, err := fmt.Fprintf(w, "::error title=Conventional Commit violation::%s\n", escapeGitHubAnnotation(msg)); err != nil {
				return err
			}
		}
	}
	return nil
}

func writeMarkdown(report model.Report, w io.Writer) error {
	if report.InvalidCommits == 0 {
		_, err := fmt.Fprintln(w, "✅ All checked commits follow Conventional Commits.")
		return err
	}

	if _, err := fmt.Fprintln(w, "## Conventional Commit Violations"); err != nil {
		return err
	}
	if _, err := fmt.Fprintln(w, ""); err != nil {
		return err
	}
	if _, err := fmt.Fprintln(w, "| Commit | Subject | Violations |"); err != nil {
		return err
	}
	if _, err := fmt.Fprintln(w, "|---|---|---|"); err != nil {
		return err
	}

	for _, result := range report.Results {
		if result.Valid {
			continue
		}
		parts := make([]string, 0, len(result.Violations))
		for _, violation := range result.Violations {
			parts = append(parts, fmt.Sprintf("`%s`: %s", violation.Code, violation.Message))
		}
		if _, err := fmt.Fprintf(w, "| `%s` | %s | %s |\n", result.ShortHash, escapePipe(result.Subject), escapePipe(strings.Join(parts, "<br>"))); err != nil {
			return err
		}
	}

	return nil
}

type rdjsonlDiagnostic struct {
	Message        string        `json:"message"`
	Severity       string        `json:"severity,omitempty"`
	Source         *rdjsonSource `json:"source,omitempty"`
	Code           *rdjsonCode   `json:"code,omitempty"`
	OriginalOutput string        `json:"original_output,omitempty"`
}

type rdjsonSource struct {
	Name string `json:"name"`
	URL  string `json:"url,omitempty"`
}

type rdjsonCode struct {
	Value string `json:"value"`
}

func writeRDJSONL(report model.Report, w io.Writer) error {
	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(false)

	for _, result := range report.Results {
		if result.Valid {
			continue
		}

		for _, violation := range result.Violations {
			diag := rdjsonlDiagnostic{
				Message:  fmt.Sprintf("%s - %s: %s", result.ShortHash, result.Subject, violation.Message),
				Severity: "ERROR",
				Source: &rdjsonSource{
					Name: "commitlint",
					URL:  "https://www.conventionalcommits.org/en/v1.0.0/#specification",
				},
				Code:           &rdjsonCode{Value: violation.Code},
				OriginalOutput: result.Message,
			}

			if err := enc.Encode(diag); err != nil {
				return err
			}
		}
	}

	return nil
}

func textColorFn(c *color.Color, mode ColorMode) func(string) string {
	switch mode {
	case ColorAlways:
		c.EnableColor()
	case ColorNever:
		c.DisableColor()
	case ColorAuto:
		// Keep default auto-detection from fatih/color.
	default:
		c.DisableColor()
	}

	return func(text string) string {
		return c.Sprint(text)
	}
}

func escapeGitHubAnnotation(value string) string {
	r := strings.NewReplacer(
		"%", "%25",
		"\r", "%0D",
		"\n", "%0A",
	)
	return r.Replace(value)
}

func escapePipe(value string) string {
	return strings.ReplaceAll(value, "|", "\\|")
}
