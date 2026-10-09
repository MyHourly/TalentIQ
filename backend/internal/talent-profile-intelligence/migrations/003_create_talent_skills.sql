-- This creates a new table: "which talent has which skill, and how good they are".

CREATE TABLE IF NOT EXISTS talent_skills (
    id                UUID PRIMARY KEY DEFAULT gen_random_uuid(),  -- unique id of this row
    talent_id         UUID     NOT NULL,                           -- which person
    skill_id          UUID     NOT NULL,                           -- which skill (e.g. Go, Kafka)
    proficiency_level SMALLINT NOT NULL CHECK (proficiency_level BETWEEN 1 AND 5), -- 1=Beginner ... 5=Master
    created_at        TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at        TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- A person can have a skill only ONCE. This rule is what lets PUT mean
-- "add it, or update it if it's already there".
CREATE UNIQUE INDEX IF NOT EXISTS idx_talent_skills_talent_skill
    ON talent_skills (talent_id, skill_id);

-- Makes "show all skills of this talent" fast later on.
CREATE INDEX IF NOT EXISTS idx_talent_skills_talent_id
    ON talent_skills (talent_id);