// Copyright 2017 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package integration

import (
	"fmt"
	"net/http"
	"net/url"
	"testing"

	repo_model "code.gitea.io/gitea/models/repo"
	"code.gitea.io/gitea/models/unittest"
	user_model "code.gitea.io/gitea/models/user"
	api "code.gitea.io/gitea/modules/structs"
	"code.gitea.io/gitea/tests"
	"github.com/stretchr/testify/assert"
)

func TestAPIListProjectBoards(t *testing.T) {
	defer tests.PrepareTestEnv(t)()

	user := unittest.AssertExistsAndLoadBean(t, &user_model.User{ID: 2})
	project := unittest.AssertExistsAndLoadBean(t, &repo_model.Repository{ID: 1})

	session := loginUser(t, user.Name)
	token := getTokenForLoggedInUser(t, session)

	link, _ := url.Parse(fmt.Sprintf("/api/v1/projects/%d/boards", project.ID))
	link.RawQuery = url.Values{"token": {token}}.Encode()

	var apiBoards []*api.ProjectBoard
	req := NewRequest(t, "GET", link.String())
	resp := session.MakeRequest(t, req, http.StatusOK)
	DecodeJSON(t, resp, &apiBoards)
	assert.Len(t, apiBoards, 3)
}

func TestAPICreateProjectBoard(t *testing.T) {
	defer tests.PrepareTestEnv(t)()

	user := unittest.AssertExistsAndLoadBean(t, &user_model.User{ID: 2})
	project := unittest.AssertExistsAndLoadBean(t, &repo_model.Repository{ID: 1})

	session := loginUser(t, user.Name)
	token := getTokenForLoggedInUser(t, session)

	link, _ := url.Parse(fmt.Sprintf("/api/v1/projects/%d/boards", project.ID))
	link.RawQuery = url.Values{"token": {token}}.Encode()

	payload := &api.NewProjectBoardPayload{
		Title:   "Test Board",
		Color:   "#000000",
		Default: false,
	}

	req := NewRequestWithJSON(t, "POST", link.String(), payload)
	resp := session.MakeRequest(t, req, http.StatusCreated)
	var apiBoard *api.ProjectBoard
	DecodeJSON(t, resp, &apiBoard)
	assert.Equal(t, payload.Title, apiBoard.Title)
	assert.Equal(t, payload.Color, apiBoard.Color)
	assert.Equal(t, payload.Default, apiBoard.Default)
}

func TestAPIGetProjectBoard(t *testing.T) {
	defer tests.PrepareTestEnv(t)()

	user := unittest.AssertExistsAndLoadBean(t, &user_model.User{ID: 2})
	project := unittest.AssertExistsAndLoadBean(t, &repo_model.Repository{ID: 1})

	session := loginUser(t, user.Name)
	token := getTokenForLoggedInUser(t, session)

	testCase := []struct {
		boardID       int64
		expectedBoard *api.ProjectBoard
	}{
		{1, &api.ProjectBoard{Title: "Test Board 1", Color: "#000000", Default: true}},
		{2, &api.ProjectBoard{Title: "Test Board 2", Color: "#000000", Default: false}},
		{3, &api.ProjectBoard{Title: "Test Board 3", Color: "#000000", Default: false}},
	}

	for _, test := range testCase {
		link, _ := url.Parse(fmt.Sprintf("/api/v1/projects/%d/boards/%d", project.ID, test.boardID))
		link.RawQuery = url.Values{"token": {token}}.Encode()

		req := NewRequest(t, "GET", link.String())
		resp := session.MakeRequest(t, req, http.StatusOK)
		var apiBoard *api.ProjectBoard
		DecodeJSON(t, resp, apiBoard)
		assert.Equal(t, test.expectedBoard.Title, apiBoard.Title)
		assert.Equal(t, test.expectedBoard.Color, apiBoard.Color)
		assert.Equal(t, test.expectedBoard.Default, apiBoard.Default)
	}
}

func TestAPIGetProjectBoardReqPermission(t *testing.T) {
}

func TestAPIUpdateProjectBoard(t *testing.T) {
	defer tests.PrepareTestEnv(t)()

	user := unittest.AssertExistsAndLoadBean(t, &user_model.User{ID: 2})
	project := unittest.AssertExistsAndLoadBean(t, &repo_model.Repository{ID: 1})

	session := loginUser(t, user.Name)
	token := getTokenForLoggedInUser(t, session)

	payload := &api.UpdateProjectBoardPayload{
		Title: "Edit Test Board",
		Color: "#000000",
	}

	expected := &api.ProjectBoard{
		Title: payload.Title,
		Color: payload.Color,
	}

	testCase := []int{1, 2, 3}
	for _, boardID := range testCase {
		link, _ := url.Parse(fmt.Sprintf("/api/v1/projects/%d/boards/%d", project.ID, boardID))
		link.RawQuery = url.Values{"token": {token}}.Encode()

		req := NewRequestWithJSON(t, "PATCH", link.String(), payload)
		resp := session.MakeRequest(t, req, http.StatusOK)
		var apiBoard *api.ProjectBoard
		DecodeJSON(t, resp, apiBoard)
		assert.Equal(t, expected.Title, apiBoard.Title)
		assert.Equal(t, expected.Color, apiBoard.Color)
		assert.Equal(t, expected.Default, apiBoard.Default)
	}
}

func TestAPIDeleteProjectBoard(t *testing.T) {
	defer tests.PrepareTestEnv(t)()

	user := unittest.AssertExistsAndLoadBean(t, &user_model.User{ID: 2})
	project := unittest.AssertExistsAndLoadBean(t, &repo_model.Repository{ID: 1})

	session := loginUser(t, user.Name)
	token := getTokenForLoggedInUser(t, session)

	testCase := []int{1, 2, 3}
	for _, boardID := range testCase {
		link, _ := url.Parse(fmt.Sprintf("/api/v1/projects/%d/boards/%d", project.ID, boardID))
		link.RawQuery = url.Values{"token": {token}}.Encode()

		req := NewRequest(t, "DELETE", link.String())
		resp := session.MakeRequest(t, req, http.StatusNoContent)
		DecodeJSON(t, resp, nil)

		// Check if board is deleted
		link, _ = url.Parse(fmt.Sprintf("/api/v1/projects/%d/boards/%d", project.ID, boardID))
		req = NewRequest(t, "GET", link.String())
		resp = session.MakeRequest(t, req, http.StatusNotFound)
	}
}
