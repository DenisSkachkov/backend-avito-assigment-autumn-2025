package service

import "errors"

var ErrPullRequestExists = errors.New("PR_EXISTS")

var ErrTeamExists = errors.New("TEAM_EXISTS")

var ErrNotFound = errors.New("NOT_FOUND")

var ErrPullRequestMerged = errors.New("PR_MERGED")

