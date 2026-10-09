-- A talent's employment history (company, job title, dates).
-- This is different from "projects": a job can contain many projects.
-- No foreign keys on purpose (same decision as the other talent tables).

CREATE TABLE IF NOT EXISTS professional_experience (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    talent_id    UUID         NOT NULL,           -- which person
    company_name VARCHAR(200) NOT NULL,
    job_title    VARCHAR(150) NOT NULL,
    start_date   DATE         NOT NULL,
    end_date     DATE,                            -- empty = still working there
    description  TEXT         NOT NULL DEFAULT '',
    created_at   TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at   TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT chk_professional_experience_dates
        CHECK (end_date IS NULL OR end_date >= start_date)   -- can't end before it started
);

-- The same job (same company, title and start date) can't be added twice for one person.
CREATE UNIQUE INDEX IF NOT EXISTS idx_professional_experience_unique
    ON professional_experience (talent_id, company_name, job_title, start_date);

CREATE INDEX IF NOT EXISTS idx_professional_experience_talent_id
    ON professional_experience (talent_id);