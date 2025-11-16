
CREATE TABLE IF NOT EXISTS users (
    user_id TEXT PRIMARY KEY,
    username TEXT NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT TRUE
);


CREATE TABLE IF NOT EXISTS teams (
    team_name TEXT PRIMARY KEY
);


CREATE TABLE IF NOT EXISTS team_members (
    team_name TEXT REFERENCES teams(team_name) ON DELETE CASCADE,
    user_id   TEXT REFERENCES users(user_id) ON DELETE CASCADE,
    PRIMARY KEY (team_name, user_id)
);


CREATE TABLE IF NOT EXISTS pull_requests (
    pull_request_id   TEXT PRIMARY KEY,
    pull_request_name TEXT NOT NULL,
    author_id         TEXT NOT NULL REFERENCES users(user_id),
    status            TEXT NOT NULL,
    merged_at         TIMESTAMPTZ
);


CREATE UNIQUE INDEX IF NOT EXISTS ux_pr_author_name ON pull_requests (author_id, pull_request_name);


CREATE TABLE IF NOT EXISTS pr_reviewers (
    pull_request_id TEXT REFERENCES pull_requests(pull_request_id) ON DELETE CASCADE,
    reviewer_id     TEXT REFERENCES users(user_id),
    PRIMARY KEY (pull_request_id, reviewer_id)
);


CREATE INDEX IF NOT EXISTS idx_team_members_user ON team_members (user_id);

CREATE INDEX IF NOT EXISTS idx_pr_reviewers_reviewer ON pr_reviewers (reviewer_id);
