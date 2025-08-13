package services

import (
	"context"
	"encoding/json"
	"fmt"
	"gregaf/gitlabctl/internal/config"
	"log"
	"log/slog"
	"net/http"
	"slices"

	"golang.org/x/sync/errgroup"
)

type GitlabService struct {
	configuration *config.Config
	logger        *slog.Logger
}

type GitlabProject struct {
	ProjectID     int     `json:"id"`
	Name          string  `json:"name"`
	Description   *string `json:"description"`
	NamespacePath string  `json:"path_with_namespace"`
}

func NewGitlabService(configuration *config.Config, logger *slog.Logger) *GitlabService {
	return &GitlabService{
		configuration: configuration,
		logger:        logger,
	}
}

func (gs *GitlabService) setGitlabHeaders(req *http.Request) {
	req.Header.Set("PRIVATE-TOKEN", gs.configuration.AccessToken)
}

func (gs *GitlabService) BulkHasBranch(gitlabProjects []GitlabProject, targetBranch string) ([]GitlabProject, error) {
	missingBranchProjects := []GitlabProject{}
	for _, project := range gitlabProjects {
		url := fmt.Sprintf("%s/projects/%d/repository/branches/%s", gs.configuration.GitlabURL, project.ProjectID, targetBranch)
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			log.Println("Error creating request", err)
			return nil, fmt.Errorf("Error creating branch check request: %w", err)
		}
		gs.setGitlabHeaders(req)

		client := &http.Client{}
		res, err := client.Do(req)
		if err != nil {
			return nil, fmt.Errorf("Error sending branch check request: %w", err)
		}
		defer res.Body.Close()

		if res.StatusCode == http.StatusNotFound {
			missingBranchProjects = append(missingBranchProjects, project)
		} else if res.StatusCode == http.StatusOK {
			gs.logger.Debug("Branch exists", "branch", project.Name)
		} else {
			return nil, fmt.Errorf("Error unexpected status code '%d' for '%s' Gitlab project", res.StatusCode, project.Name)
		}
	}

	return missingBranchProjects, nil
}

// func (gs *GitlabService) HasBranch(gitlabProjects GitlabProject, targetBranch string) error {

// 	url := fmt.Sprintf("%s/projects/%d/repository/branches/%s", gs.configuration.GitlabURL, project.ProjectID, targetBranch)
// 	req, err := http.NewRequest("GET", url, nil)
// 	if err != nil {
// 		log.Println("Error creating request", err)
// 		return nil, fmt.Errorf("Error creating branch check request: %w", err)
// 	}
// 	gs.setGitlabHeaders(req)

// 	client := &http.Client{}
// 	res, err := client.Do(req)
// 	if err != nil {
// 		return nil, fmt.Errorf("Error sending branch check request: %w", err)
// 	}
// 	defer res.Body.Close()

// 	if res.StatusCode == http.StatusNotFound {
// 		missingBranchProjects = append(missingBranchProjects, project)
// 	} else if res.StatusCode == http.StatusOK {
// 		gs.logger.Debug("Branch exists", "branch", project.Name)
// 	} else {
// 		return nil, fmt.Errorf("Error unexpected status code '%d' for '%s' Gitlab project", res.StatusCode, project.Name)
// 	}

// 	return missingBranchProjects, nil
// }

func (gs *GitlabService) CreateBranch(gitlabProject GitlabProject, newBranch, baseBranch string) error {
	return fmt.Errorf("unimplemented")
}

func (gs *GitlabService) FindGitlabProjects(projects []string) ([]GitlabProject, error) {
	// TODO: Use the context provided to handle error scenario properly.
	g, _ := errgroup.WithContext(context.TODO())

	results := make([][]GitlabProject, len(projects))
	for i, project := range projects {
		i, project := i, project
		g.Go(func() error {
			threadID := fmt.Sprintf("%p", &i)
			gs.logger.Debug("Started worker", "project", project, "workerID", threadID)
			result, err := gs.FindGitlabProject(project)
			if err == nil {
				results[i] = result
			}
			gs.logger.Debug("Completed worker", "project", project, "workerID", threadID, "found", len(result))
			return err
		})
	}
	if err := g.Wait(); err != nil {
		return nil, err
	}

	aggGitlabProjects := slices.Concat(results...)

	return aggGitlabProjects, nil
}

func (gs *GitlabService) FindGitlabProject(project string) ([]GitlabProject, error) {
	url := fmt.Sprintf("%s/projects?search=%s", gs.configuration.GitlabURL, project)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Println("Erraor creating request", err)
		return nil, fmt.Errorf("TODO")
	}

	req.Header.Set("PRIVATE-TOKEN", gs.configuration.AccessToken)

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		log.Println("Error sending request", err)
		return nil, fmt.Errorf("TODO")
	}
	defer res.Body.Close()

	body := &[]GitlabProject{}
	err = json.NewDecoder(res.Body).Decode(body)
	if err != nil {
		log.Println("Error decoding response body", err)
		return nil, fmt.Errorf("Failed JSON Decode %w", err)
	}

	return *body, nil
}
