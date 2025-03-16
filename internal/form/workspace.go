package form

import (
	"fmt"
	"path/filepath"

	"github.com/phe-lab/ws/internal/log"
	"github.com/phe-lab/ws/internal/utils"

	"github.com/charmbracelet/huh"
)

func ChooseWorkspace(userKeyword, workspaceDir string) (string, error) {
	selectedFile := fmt.Sprintf("%s.code-workspace", userKeyword)
	workspaces, err := utils.FindWorkspaceFiles(workspaceDir, userKeyword)
	if err != nil {
		return "", err
	}

	if len(workspaces) == 0 {
		log.Logger.Info().Str("path", workspaceDir).Msg("No workspace files found")
		return "", nil
	}

	// If there is only one workspace after filtering, select it
	if userKeyword != "" && len(workspaces) == 1 {
		return workspaces[0], nil
	}

	err = huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Select a workspace to open:").
				Filtering(true).
				Options(convertToOptions(workspaces, workspaceDir)...).
				Value(&selectedFile),
		),
	).WithTheme(t).Run()

	if err != nil {
		return "", err
	}

	return selectedFile, nil
}

func convertToOptions(workspaces []string, path string) []huh.Option[string] {
	options := make([]huh.Option[string], len(workspaces))

	for i, workspace := range workspaces {
		relativePath, err := filepath.Rel(path, workspace)
		if err != nil {
			relativePath = workspace
		}
		options[i] = huh.NewOption(utils.ShortenPath(relativePath), workspace)
	}

	return options
}
