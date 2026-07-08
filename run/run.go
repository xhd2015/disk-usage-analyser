package run

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"disk-usage-analyser/analyse"
	"disk-usage-analyser/server"
	"disk-usage-analyser/tmpfiles"
	"disk-usage-analyser/usagescan"

	"github.com/xhd2015/kool/pkgs/web"
	"github.com/xhd2015/less-gen/flags"
)

const help = `
Usage: disk-usage-analyser <subcommand>

Subcommands:
  analyse [DIR]     Analyse directory tree size and link metrics
  scan [PATH]       List immediate children with recursive directory sizes
  tmp-files scan    Scan temporary file candidates
`

type Options struct {
	Stdout      io.Writer
	Stderr      io.Writer
	HomeDir     string
	StartServer func(context.Context, ServerOptions) error
}

type ServerOptions struct {
	Port      int
	Dev       bool
	Component string
}

func Run(args []string) error {
	return RunWithOptions(context.Background(), args, Options{})
}

func RunWithOptions(ctx context.Context, args []string, opts Options) error {
	stdout := opts.Stdout
	if stdout == nil {
		stdout = os.Stdout
	}
	stderr := opts.Stderr
	if stderr == nil {
		stderr = os.Stderr
	}

	if len(args) > 0 && args[0] == "scan" {
		exitCode, err := usagescan.RunCLI(args[1:], usagescan.CLIOptions{
			Stdout: stdout,
			Stderr: stderr,
		})
		if err != nil {
			return err
		}
		if exitCode != 0 {
			return fmt.Errorf("scan exited with code %d", exitCode)
		}
		return nil
	}

	if len(args) > 0 && args[0] == "analyse" {
		exitCode, err := analyse.RunCLI(args[1:], analyse.CLIOptions{
			Stdout: stdout,
			Stderr: stderr,
		})
		if err != nil {
			return err
		}
		if exitCode != 0 {
			return fmt.Errorf("analyse exited with code %d", exitCode)
		}
		return nil
	}

	if len(args) > 0 && args[0] == "tmp-files" {
		_, exitCode, err := tmpfiles.RunCLI(ctx, args[1:], tmpfiles.CLIOptions{
			Stdout:  stdout,
			Stderr:  stderr,
			HomeDir: opts.HomeDir,
		})
		if err != nil {
			return err
		}
		if exitCode != 0 {
			return fmt.Errorf("tmp-files exited with code %d", exitCode)
		}
		return nil
	}

	var devFlag bool
	var component string
	args, err := flags.
		Bool("--dev", &devFlag).
		String("--component", &component).
		Help("-h,--help", help).
		Parse(args)
	if err != nil {
		return err
	}

	if len(args) > 0 {
		absPath, err := filepath.Abs(args[0])
		if err != nil {
			return fmt.Errorf("invalid path %s: %v", args[0], err)
		}
		server.InitialDir = absPath
		args = args[1:]
	}

	if len(args) > 0 {
		return fmt.Errorf("unrecognized extra args: %s", strings.Join(args, " "))
	}

	if component == "list" {
		fmt.Fprintln(stdout, "Available components: App")
		return nil
	}

	// next port
	port, err := web.FindAvailablePort(8080, 100)
	if err != nil {
		return err
	}

	if opts.StartServer != nil {
		return opts.StartServer(ctx, ServerOptions{
			Port:      port,
			Dev:       devFlag,
			Component: component,
		})
	}

	if component != "" {
		var html string
		if !devFlag {
			html, err = server.FormatTemplateHtml(server.FormatOptions{
				Component: component,
			})
			if err != nil {
				return err
			}
		}
		return server.ServeComponent(port, server.ServeOptions{
			Dev: devFlag,
			Static: server.StaticOptions{
				IndexHtml: html,
			},
			OpenBrowserUrl: func(port int, url string) string {
				if devFlag {
					return fmt.Sprintf("%s/?component=%s", url, component)
				}
				return url
			},
		})
	}

	return server.Serve(port, devFlag)
}
