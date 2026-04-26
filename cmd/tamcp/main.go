// Command tamcp is the tamcp CLI: it runs the MCP server, manages config files,
// and (on Windows) interacts with the Service Control Manager.
package main

import (
	"context"
	"fmt"
	"os"

	"github.com/rangertaha/tamcp/internal"
	"github.com/rangertaha/tamcp/internal/agent"
	"github.com/rangertaha/tamcp/internal/config"
	"github.com/rangertaha/tamcp/internal/winservice"
	cli "github.com/urfave/cli/v3"
)

var (
	options []config.Option
	cfg     *config.Config
)

func main() {
	app := &cli.Command{
		Name:    internal.NAME,
		Usage:   "MCP server for technical-analysis indicators",
		Version: internal.Version,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "config",
				Aliases: []string{"c"},
				Value:   "",
				Usage:   "local config path merged after global and user config",
			},
			&cli.StringFlag{
				Name:    "level",
				Aliases: []string{"l"},
				Value:   "error",
				Usage:   "log level",
			},
			&cli.StringFlag{
				Name:   "log",
				Value:  "stderr",
				Usage:  "log output (stderr, <filename>) — stdout is reserved for MCP",
				Hidden: true,
			},
			&cli.BoolFlag{
				Name:  "debug",
				Value: false,
				Usage: "enable debug mode",
			},
		},
		Before: func(ctx context.Context, cmd *cli.Command) (context.Context, error) {
			options = append(options, config.WithGlobalConfigFile())
			options = append(options, config.WithUserConfigFile())
			options = append(options, config.WithConfigFile(cmd.String("config")))
			options = append(options, config.WithLogLevel(cmd.String("level")))
			options = append(options, config.WithLogFile(cmd.String("log")))
			options = append(options, config.WithDebug(cmd.Bool("debug")))

			c, err := config.New(options...)
			if err != nil {
				return nil, err
			}
			cfg = c
			return ctx, nil
		},
		Commands: []*cli.Command{
			{
				Name:  "server",
				Usage: "run the tamcp MCP server over stdio",
				Action: func(ctx context.Context, cmd *cli.Command) error {
					a, err := agent.New(cfg)
					if err != nil {
						return err
					}
					return a.Run()
				},
			},
			{
				Name:  "init",
				Usage: "create global / user / systemd config files",
				Flags: []cli.Flag{
					&cli.BoolFlag{Name: "global", Aliases: []string{"g"}, Usage: "create global config in /etc/tamcp"},
					&cli.BoolFlag{Name: "user", Aliases: []string{"u"}, Usage: "create user config in ~/.config/tamcp"},
					&cli.BoolFlag{Name: "service", Aliases: []string{"s"}, Usage: "create systemd unit and enable it"},
					&cli.BoolFlag{Name: "all", Aliases: []string{"a"}, Usage: "create all config files"},
					&cli.BoolFlag{Name: "clean", Usage: "remove files created by init"},
				},
				Action: runInit,
			},
			{
				Name:    "version",
				Aliases: []string{"v"},
				Usage:   "print version information",
				Action: func(ctx context.Context, cmd *cli.Command) error {
					fmt.Printf("%s %s (commit %s, built %s, %s)\n",
						internal.NAME, internal.Version, internal.Commit, internal.BuildDate, internal.GoVersion)
					return nil
				},
			},
			{
				Name:  "service",
				Usage: "Windows Service Control Manager integration",
				Commands: []*cli.Command{
					{
						Name:  "install",
						Usage: "install tamcp as a Windows service (auto-start)",
						Action: func(_ context.Context, _ *cli.Command) error {
							if err := winservice.Install("server"); err != nil {
								return err
							}
							fmt.Printf("installed Windows service %q\n", winservice.Name)
							return nil
						},
					},
					{
						Name:  "uninstall",
						Usage: "stop (if running) and remove the tamcp Windows service",
						Action: func(_ context.Context, _ *cli.Command) error {
							if err := winservice.Uninstall(); err != nil {
								return err
							}
							fmt.Printf("uninstalled Windows service %q\n", winservice.Name)
							return nil
						},
					},
					{
						Name:  "start",
						Usage: "start the installed tamcp Windows service",
						Action: func(_ context.Context, _ *cli.Command) error {
							if err := winservice.StartService(); err != nil {
								return err
							}
							fmt.Printf("started Windows service %q\n", winservice.Name)
							return nil
						},
					},
					{
						Name:  "stop",
						Usage: "stop the running tamcp Windows service",
						Action: func(_ context.Context, _ *cli.Command) error {
							if err := winservice.StopService(); err != nil {
								return err
							}
							fmt.Printf("stopped Windows service %q\n", winservice.Name)
							return nil
						},
					},
				},
			},
		},
	}

	if err := app.Run(context.Background(), os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "%s: %v\n", internal.NAME, err)
		os.Exit(1)
	}
}

func runInit(_ context.Context, cmd *cli.Command) error {
	all := cmd.Bool("all")
	if !cmd.Bool("global") && !cmd.Bool("user") && !cmd.Bool("service") && !cmd.Bool("all") && !cmd.Bool("clean") && cmd.Args().First() == "" {
		all = true
	}

	if cmd.Bool("clean") {
		if !cmd.Bool("global") && !cmd.Bool("user") && !cmd.Bool("service") && !cmd.Bool("all") && cmd.Args().First() == "" {
			all = true
		}
		if cmd.Bool("global") || all {
			if err := config.RemoveConfigFile(config.GlobalConfigDir, config.DefaultConfigFileName); err != nil {
				fmt.Println(err.Error())
			}
		}
		if cmd.Bool("user") || all {
			if err := config.RemoveConfigFile(config.UserConfigDir, config.DefaultConfigFileName); err != nil {
				fmt.Println(err.Error())
			}
		}
		if cmd.Bool("service") || all {
			if err := config.RemoveSystemdServiceFile(); err != nil {
				fmt.Println(err.Error())
			}
		}
		if cmd.Args().First() != "" {
			if err := config.RemoveConfigPath(cmd.Args().First()); err != nil {
				fmt.Println(err.Error())
			}
		}
		return nil
	}

	if cmd.Bool("global") || all {
		if err := config.CreateConfigFile(config.GlobalConfigDir, config.DefaultConfigFileName); err != nil {
			fmt.Println(err.Error())
		}
	}
	if cmd.Bool("user") || all {
		if err := config.CreateConfigFile(config.UserConfigDir, config.DefaultConfigFileName); err != nil {
			fmt.Println(err.Error())
		}
	}
	if cmd.Bool("service") || all {
		if err := config.CreateSystemdServiceFile(); err != nil {
			fmt.Println(err.Error())
		}
	}
	if cmd.Args().First() != "" {
		if err := config.CreateConfigPath(cmd.Args().First()); err != nil {
			fmt.Println(err.Error())
		}
	}
	return nil
}
