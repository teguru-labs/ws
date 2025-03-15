package ws

import (
	"fmt"
	"os"

	"github.com/phe-lab/ws/internal/config"
	"github.com/phe-lab/ws/internal/exception"
	"github.com/phe-lab/ws/internal/form"
	"github.com/phe-lab/ws/internal/log"
	"github.com/phe-lab/ws/internal/utils"

	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"
)

var debug bool

var rootCmd = &cobra.Command{
	Use:   "ws [filename]",
	Short: "A simple CLI tool to quickly open VSCode Workspace",
	Long: `
██╗    ██╗ ██████╗ ██████╗ ██╗  ██╗███████╗██████╗  █████╗  ██████╗███████╗     ██████╗██╗     ██╗
██║    ██║██╔═══██╗██╔══██╗██║ ██╔╝██╔════╝██╔══██╗██╔══██╗██╔════╝██╔════╝    ██╔════╝██║     ██║
██║ █╗ ██║██║   ██║██████╔╝█████╔╝ ███████╗██████╔╝███████║██║     █████╗      ██║     ██║     ██║
██║███╗██║██║   ██║██╔══██╗██╔═██╗ ╚════██║██╔═══╝ ██╔══██║██║     ██╔══╝      ██║     ██║     ██║
╚███╔███╔╝╚██████╔╝██║  ██║██║  ██╗███████║██║     ██║  ██║╚██████╗███████╗    ╚██████╗███████╗██║
 ╚══╝╚══╝  ╚═════╝ ╚═╝  ╚═╝╚═╝  ╚═╝╚══════╝╚═╝     ╚═╝  ╚═╝ ╚═════╝╚══════╝     ╚═════╝╚══════╝╚═╝

A simple CLI tool to quickly open VSCode Workspace.

Environment Variables:
  - VSCODE_WS_PATH: The directory path containing the "*.code-workspace" files. Default: ~/code-workspaces
  - VSCODE_WS_DEBUG: Enable the debug mode. Example: VSCODE_WS_DEBUG=true`,
	Example: `  # List the workspaces
  ws

  # Open the workspace with the filename "simple-scrollspy.code-workspace"
  ws simple-scrollspy

  # Set the logging level to "debug"
  ws --debug`,
	Args: cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		cfg := config.LoadConfig(debug)
		logger := log.InitLogger(debug)
		workspaceDir := cfg.GetWorkspacePath()

		userKeyword := ""
		if len(args) > 0 && args[0] != "" {
			userKeyword = args[0]
		}

		if err := ensureWorkspaceDir(workspaceDir); err != nil {
			logger.Error().Msg(err.Error())
			if err != exception.ErrUserAborted {
				os.Exit(1)
			}
			return
		}

		logger.Debug().Msg("Application is running in debug mode")
		logger.Debug().Str("envVar", "VSCODE_WS_PATH").Msg("Using the default workspace directory as the environment variable was not found")
		logger.Debug().Str("workspacePath", workspaceDir).Msg("")

		workspace, err := form.ChooseWorkspace(userKeyword, workspaceDir)
		if err != nil && err != huh.ErrUserAborted {
			logger.Error().Msg(err.Error())
			os.Exit(1)
		}

		if workspace != "" {
			log.Logger.Debug().Str("file", workspace).Msg("Selected workspace")
			utils.OpenWorkspace(workspace)
		}
	},
}

func init() {
	rootCmd.PersistentFlags().BoolVar(&debug, "debug", false, "enable debug mode")
}

func SetVersionInfo(version, commit, date string) {
	rootCmd.Version = fmt.Sprintf("%s (Built on %s from Git SHA %s)", version, date, commit)
}

func ensureWorkspaceDir(path string) error {
	log.Logger.Debug().Msg("Validate workspace path")
	if err := utils.ValidateWorkspacePath(path); err != nil {
		if err == exception.ErrNotExist {
			log.Logger.Warn().
				Str("default", "~/code-workspaces").
				Msg("Declare environment variable `VSCODE_WS_PATH=~/your/path` to change the workspace directory")

			confirm, err := form.ConfirmCreateDirectory(path)
			if err != nil || !confirm {
				if err == huh.ErrUserAborted || !confirm {
					return exception.ErrUserAborted
				}
				return err
			}

			if err = os.MkdirAll(path, 0755); err != nil {
				return err
			}
			return nil
		}
		return err
	}

	return nil
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Oops. An error while executing ws '%s'\n", err)
		os.Exit(1)
	}
}
