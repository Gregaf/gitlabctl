package services

import (
	"context"
	"encoding/json"
	"fmt"
	"gregaf/gitlabctl/internal/config"
	"io"
	"log"
	"log/slog"
	"net/http"
	"slices"

	"golang.org/x/sync/errgroup"
)

type GitlabService struct {
	client        *http.Client
	configuration *config.Config
	logger        *slog.Logger
}

type GitlabProject struct {
	ProjectID     int     `json:"id"`
	Name          string  `json:"name"`
	Description   *string `json:"description"`
	NamespacePath string  `json:"path_with_namespace"`
}

func NewGitlabService(client *http.Client, configuration *config.Config, logger *slog.Logger) *GitlabService {
	return &GitlabService{
		client:        client,
		configuration: configuration,
		logger:        logger,
	}
}

func (gs *GitlabService) BulkHasBranch(gitlabProjects []GitlabProject, targetBranch string) ([]GitlabProject, error) {
	g, ctx := errgroup.WithContext(context.TODO())

	results := make([]bool, len(gitlabProjects))
	for i, project := range gitlabProjects {
		i, project := i, project
		g.Go(func() error {
			result, err := gs.HasBranch(ctx, project, targetBranch)
			if err == nil {
				results[i] = result
			}
			return err
		})
	}

	if err := g.Wait(); err != nil {
		return nil, err
	}

	missingBranchProjects := []GitlabProject{}
	for i, project := range gitlabProjects {
		if results[i] == false {
			missingBranchProjects = append(missingBranchProjects, project)
		}
	}

	return missingBranchProjects, nil
}

func (gs *GitlabService) HasBranch(ctx context.Context, gitlabProject GitlabProject, targetBranch string) (bool, error) {
	url := fmt.Sprintf("%s/projects/%d/repository/branches/%s", gs.configuration.GitlabURL, gitlabProject.ProjectID, targetBranch)
	req, err := gs.getGitlabRequest(ctx, "GET", url, nil)
	if err != nil {
		return false, fmt.Errorf("Error creating branch check request: %w", err)
	}

	res, err := gs.client.Do(req)
	if err != nil {
		return false, fmt.Errorf("Error sending branch check request: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode == http.StatusNotFound {
		return false, nil
	} else if res.StatusCode == http.StatusOK {
		gs.logger.Debug("Branch exists", "branch", gitlabProject.Name)
	} else {
		return false, fmt.Errorf("Error unexpected status code '%d' for '%s' Gitlab project", res.StatusCode, gitlabProject.Name)
	}

	return true, nil
}

func (gs *GitlabService) CreateBranch(gitlabProject GitlabProject, newBranch, baseBranch string) error {
	return fmt.Errorf("unimplemented")
}

func (gs *GitlabService) FindGitlabProjects(projects []string) ([]GitlabProject, error) {
	g, ctx := errgroup.WithContext(context.TODO())

	results := make([][]GitlabProject, len(projects))
	for i, project := range projects {
		i, project := i, project
		g.Go(func() error {
			result, err := gs.FindGitlabProject(ctx, project)
			if err == nil {
				results[i] = result
			}
			return err
		})
	}

	if err := g.Wait(); err != nil {
		return nil, err
	}

	aggGitlabProjects := slices.Concat(results...)

	return aggGitlabProjects, nil
}

func (gs *GitlabService) FindGitlabProject(ctx context.Context, project string) ([]GitlabProject, error) {
	url := fmt.Sprintf("%s/projects?search=%s", gs.configuration.GitlabURL, project)
	req, err := gs.getGitlabRequest(ctx, "GET", url, nil)
	if err != nil {
		log.Println("Erraor creating request", err)
		return nil, fmt.Errorf("failed create request: %w", err)
	}

	res, err := gs.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed http request: %w", err)
	}
	defer res.Body.Close()

	body := &[]GitlabProject{}
	err = json.NewDecoder(res.Body).Decode(body)
	if err != nil {
		log.Println("Error decoding response body", err)
		return nil, fmt.Errorf("failed decode json: %w", err)
	}

	return *body, nil
}

func (gs *GitlabService) getGitlabRequest(ctx context.Context, method, url string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("PRIVATE-TOKEN", gs.configuration.AccessToken)

	return req, nil
}
