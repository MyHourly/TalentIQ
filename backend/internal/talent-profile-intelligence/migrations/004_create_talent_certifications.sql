-- Creates the table for "which talent earned which certification, and when".
-- No foreign keys on purpose (same decision as talent_skills / talent_preferences).

CREATE TABLE IF NOT EXISTS talent_certifications (
    id               UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    talent_id        UUID        NOT NULL,        -- which person
    certification_id UUID        NOT NULL,        -- which certification
    issued_at        DATE        NOT NULL,        -- date earned
    expires_at       DATE,                        -- empty = never expires
    credential_url   TEXT,                        -- link to the certificate, optional
    status           VARCHAR(30) NOT NULL DEFAULT 'ACTIVE',   -- ACTIVE or EXPIRED
    created_at       TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT chk_talent_certifications_dates
        CHECK (expires_at IS NULL OR expires_at >= issued_at)   -- can't expire before it was issued
);

-- The same certificate (same issue date) can't be added twice for one person.
CREATE UNIQUE INDEX IF NOT EXISTS idx_talent_certifications_unique
    ON talent_certifications (talent_id, certification_id, issued_at);

CREATE INDEX IF NOT EXISTS idx_talent_certifications_talent_id
    ON talent_certifications (talent_id);