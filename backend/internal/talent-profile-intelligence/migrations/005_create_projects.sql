-- Two tables, following the HLD: "projects" is the project itself,
-- "project_members" links a talent to it (with the role they played).
-- No foreign keys on purpose (same decision as the other talent tables).

CREATE TABLE IF NOT EXISTS projects (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name        VARCHAR(200) NOT NULL,
    description TEXT         NOT NULL DEFAULT '',
    domain      VARCHAR(100) NOT NULL DEFAULT '',   -- e.g. Banking, E-Commerce
    client_name VARCHAR(150) NOT NULL DEFAULT '',
    start_date  DATE         NOT NULL,
    end_date    DATE,                               -- empty until the project is completed
    status      VARCHAR(30)  NOT NULL DEFAULT 'IN_PROGRESS',
    created_at  TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT chk_projects_status CHECK (status IN ('IN_PROGRESS', 'COMPLETED')),
    CONSTRAINT chk_projects_dates  CHECK (end_date IS NULL OR end_date >= start_date)
);

CREATE TABLE IF NOT EXISTS project_members (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    project_id UUID         NOT NULL,
    talent_id  UUID         NOT NULL,
    role       VARCHAR(100) NOT NULL DEFAULT '',    -- e.g. Backend Developer
    start_date DATE         NOT NULL,
    end_date   DATE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT chk_project_members_dates CHECK (end_date IS NULL OR end_date >= start_date)
);

-- A talent is a member of a project only once.
CREATE UNIQUE INDEX IF NOT EXISTS idx_project_members_project_talent
    ON project_members (project_id, talent_id);

-- "All projects of this talent" will be needed by later features.
CREATE INDEX IF NOT EXISTS idx_project_members_talent_id
    ON project_members (talent_id);