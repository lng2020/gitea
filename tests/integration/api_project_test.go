// Copyright 2022 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package integration

import (
	"fmt"
	"net/http"
	"net/url"
	"testing"

	auth_model "code.gitea.io/gitea/models/auth"
	repo_model "code.gitea.io/gitea/models/repo"
	"code.gitea.io/gitea/models/unittest"
	user_model "code.gitea.io/gitea/models/user"
	api "code.gitea.io/gitea/modules/structs"
	"code.gitea.io/gitea/tests"
	"github.com/stretchr/testify/assert"
)

func TestAPIListProjects(t *testing.T) {
	defer tests.PrepareTestEnv(t)()

	user := unittest.AssertExistsAndLoadBean(t, &user_model.User{ID: 2})
	org := unittest.AssertExistsAndLoadBean(t, &user_model.User{ID: 17})
	repo := unittest.AssertExistsAndLoadBean(t, &repo_model.Repository{ID: 1})

	session := loginUser(t, user.Name)
	token := getTokenForLoggedInUser(t, session, auth_model.AccessTokenScopeReadUser, auth_model.AccessTokenScopeReadOrganization, auth_model.AccessTokenScopeReadRepository, auth_model.AccessTokenScopeReadIssue)

	testCase := []struct {
		link        string
		expectedLen int
	}{
		{fmt.Sprintf("/api/v1/user/%s/projects", user.Name), 1},
		{fmt.Sprintf("/api/v1/orgs/%s/projects", org.Name), 1},
		{fmt.Sprintf("/api/v1/repos/%s/%s/projects", user.Name, repo.Name), 1},
	}

	for _, test := range testCase {
		link, _ := url.Parse(test.link)
		link.RawQuery = url.Values{"token": {token}}.Encode()

		var apiProjects []*api.Project
		req := NewRequest(t, "GET", link.String())
		resp := session.MakeRequest(t, req, http.StatusOK)
		DecodeJSON(t, resp, &apiProjects)
		assert.Len(t, apiProjects, test.expectedLen)
	}
}

func TestAPICreateProject(t *testing.T) {
	defer tests.PrepareTestEnv(t)()

	user := unittest.AssertExistsAndLoadBean(t, &user_model.User{ID: 2})
	repo := unittest.AssertExistsAndLoadBean(t, &repo_model.Repository{ID: 1})
	org := unittest.AssertExistsAndLoadBean(t, &user_model.User{ID: 17})

	session := loginUser(t, user.Name)
	token := getTokenForLoggedInUser(t, session, auth_model.AccessTokenScopeWriteUser, auth_model.AccessTokenScopeWriteRepository, auth_model.AccessTokenScopeWriteIssue)

	payload := api.NewProjectPayload{
		Title:       "Test Project",
		Description: "Test Project Description",
		BoardType:   uint8(1),
	}

	expectedProject := api.Project{
		Title:       "Test Project",
		Description: "Test Project Description",
		BoardType:   uint8(1),
	}

	testCase := []struct {
		link string
	}{
		{fmt.Sprintf("/api/v1/user/%s/projects", user.Name)},
		{fmt.Sprintf("/api/v1/repos/%s/%s/projects", user.Name, repo.Name)},
		{fmt.Sprintf("/api/v1/orgs/%s/projects", org.Name)},
	}

	for _, test := range testCase {
		link, _ := url.Parse(test.link)
		link.RawQuery = url.Values{"token": {token}}.Encode()

		req := NewRequestWithJSON(t, "POST", link.String(), payload)
		resp := session.MakeRequest(t, req, http.StatusCreated)

		var apiProject api.Project
		DecodeJSON(t, resp, &apiProject)
		assert.Equal(t, expectedProject.Title, apiProject.Title)
		assert.Equal(t, expectedProject.Description, apiProject.Description)
		assert.Equal(t, expectedProject.BoardType, apiProject.BoardType)
	}
}

func TestAPIGetProject(t *testing.T) {
	defer tests.PrepareTestEnv(t)()

	user := unittest.AssertExistsAndLoadBean(t, &user_model.User{ID: 2})
	session := loginUser(t, user.Name)
	token := getTokenForLoggedInUser(t, session, auth_model.AccessTokenScopeReadUser, auth_model.AccessTokenScopeReadRepository, auth_model.AccessTokenScopeReadIssue)

	testCase := []struct {
		projectID       int
		expectedProject api.Project
	}{
		{1, api.Project{ID: 1, Title: "First project"}},
		{2, api.Project{ID: 2, Title: "Second project"}},
	}

	for _, test := range testCase {
		link, _ := url.Parse(fmt.Sprintf("/api/v1/projects/%d", test.projectID))
		link.RawQuery = url.Values{"token": {token}}.Encode()

		req := NewRequest(t, "GET", link.String())
		resp := session.MakeRequest(t, req, http.StatusOK)

		var apiProject api.Project
		DecodeJSON(t, resp, &apiProject)
		assert.Equal(t, test.expectedProject.ID, apiProject.ID)
		assert.Equal(t, test.expectedProject.Title, apiProject.Title)
	}
}

func TestAPIUpdateProject(t *testing.T) {
	defer tests.PrepareTestEnv(t)()

	user := unittest.AssertExistsAndLoadBean(t, &user_model.User{ID: 2})
	session := loginUser(t, user.Name)
	token := getTokenForLoggedInUser(t, session, auth_model.AccessTokenScopeWriteUser, auth_model.AccessTokenScopeWriteRepository, auth_model.AccessTokenScopeWriteIssue)

	payload := api.UpdateProjectPayload{
		Title:       "Edited test Project",
		Description: "Edited test Project Description",
	}

	expectedProject := api.Project{
		Title:       "Edited test Project",
		Description: "Edited test Project Description",
	}

	testCase := []int{1, 2}

	for _, test := range testCase {
		link, _ := url.Parse(fmt.Sprintf("/api/v1/projects/%d", test))
		link.RawQuery = url.Values{"token": {token}}.Encode()

		req := NewRequestWithJSON(t, "PATCH", link.String(), payload)
		resp := session.MakeRequest(t, req, http.StatusOK)

		var apiProject api.Project
		DecodeJSON(t, resp, &apiProject)
		assert.Equal(t, expectedProject.Title, apiProject.Title)
		assert.Equal(t, expectedProject.Description, apiProject.Description)
	}
}

func TestAPIDeleteProject(t *testing.T) {
}
