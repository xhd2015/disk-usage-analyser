package run

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"disk-usage-analyser/analyse"
	"disk-usage-analyser/explain"
	"disk-usage-analyser/server"
	"disk-usage-analyser/skill"
	"disk-usage-analyser/tmpfiles"
	"disk-usage-analyser/usagescan"

	"github.com/xhd2015/kool/pkgs/web"
	lessflags "github.com/xhd2015/less-flags"
)

const help = `
Usage: disk-usage-analyser <subcommand>

Subcommands:
  analyse [DIR]     Analyse directory tree size and link metrics
  scan [PATH]       Walk a directory and emit a size tree (text or --json); use scan --inspect FILE to query offline
  explain [PATH]    Explain reclaim kind, size breakdown, and safe-to-reclaim advice for a path
  tmp-files scan    Scan temporary file candidates
  skill --show [topic]
                    Print embedded analyse-my-disk-space skill or nested topic
  skill --install   Install skill files to agent skill directories
  skill --list      List skill name and topic paths
  install …         Alias of skill --install

Server options:
  --dev             Run the web UI in development mode
  --component NAME  Serve a single component (use --component list to list)

Run disk-usage-analyser <command> --help for command-specific options.
Run disk-usage-analyser skill --help for skill actions.
Run disk-usage-analyser skill --install --help for install flags.
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

	if len(args) > 0 && args[0] == "skill" {
		return skill.Handle(args[1:])
	}
	if len(args) > 0 && args[0] == "install" {
		// top-level install alias → skill --install
		return skill.HandleInstall(args[1:])
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

	if len(args) > 0 && args[0] == "explain" {
		exitCode, err := explain.RunCLI(args[1:], explain.CLIOptions{
			Stdout: stdout,
			Stderr: stderr,
		})
		if err != nil {
			return err
		}
		if exitCode != 0 {
			return fmt.Errorf("explain exited with code %d", exitCode)
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
	args, err := lessflags.
		Bool("--dev", &devFlag).
		String("--component", &component).
		HelpFunc("-h,--help", func() {
			txt := strings.TrimPrefix(help, "\n")
			fmt.Fprint(stdout, txt)
			if !strings.HasSuffix(txt, "\n") {
				fmt.Fprintln(stdout)
			}
		}).
		HelpNoExit().
		Parse(args)
	if err == lessflags.ErrHelp {
		return nil
	}
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
