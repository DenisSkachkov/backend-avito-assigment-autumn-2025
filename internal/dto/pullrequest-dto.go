package dto

type CreatePRDTO struct {
    PullRequestID   string `json:"pull_request_id"`
    PullRequestName string `json:"pull_request_name"`
    AuthorID        string `json:"author_id"`
}

type ReassignReviewerDTO struct {
    PullRequestID string `json:"pull_request_id"`
    OldReviewerID     string `json:"old_reviewer_id"`
}

type MergePRDTO struct {
	PullRequestID   string `json:"pull_request_id"`
}