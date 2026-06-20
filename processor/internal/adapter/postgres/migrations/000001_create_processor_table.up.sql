CREATE TABLE repositories (
    id BIGSERIAL PRIMARY KEY,

    owner TEXT NOT NULL,
    repo TEXT NOT NULL,
    description TEXT NOT NULL,
    stars INTEGER NOT NULL,
    forks INTEGER NOT NULL,
    created_at TIMESTAMPTZ NOT NULL,

    added_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),

    CONSTRAINT repositories_owner_repo_unique UNIQUE (owner, repo)
);
