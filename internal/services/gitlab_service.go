package services

import (
	"encoding/json"
	"fmt"
	"gregaf/gitlabctl/internal/config"
	"log"
	"log/slog"
	"net/http"
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

func (gs *GitlabService) HasBranch(gitlabProjects []GitlabProject, targetBranch string) ([]GitlabProject, error) {
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

func (gs *GitlabService) FindGitlabProjects(projects []string) ([]GitlabProject, error) {
	ch := make(chan []GitlabProject, len(projects))

	for _, project := range projects {
		go func() {
			url := fmt.Sprintf("%s/projects?search=%s", gs.configuration.GitlabURL, project)
			req, err := http.NewRequest("GET", url, nil)
			if err != nil {
				log.Println("Erraor creating request", err)
				return
			}

			req.Header.Set("PRIVATE-TOKEN", gs.configuration.AccessToken)

			client := &http.Client{}
			res, err := client.Do(req)
			if err != nil {
				log.Println("Error sending request", err)
				return
			}
			defer res.Body.Close()

			var body []GitlabProject
			err = json.NewDecoder(res.Body).Decode(&body)
			if err != nil {
				log.Println("Error decoding response body", err)
			}

			ch <- body
		}()
	}

	gitlabProjects := []GitlabProject{}
	for range projects {
		res := <-ch
		gitlabProjects = append(gitlabProjects, res...)
	}

	return gitlabProjects, nil
}
