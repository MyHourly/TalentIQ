-- TPI-013.1
-- Creates the talent_preferences table.
--
-- This table stores the career/work preferences of a talent.
-- The main profile remains in talent_profiles, while preference
-- information is maintained separately.

CREATE TABLE IF NOT EXISTS talent_preferences (
    -- Unique ID for the preference record.
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    -- ID of the talent whose preferences are stored.
    -- We intentionally keep this as UUID without a database
    -- foreign key because TalentIQ follows service-level ownership
    -- between microservices.
    talent_id UUID NOT NULL,

    -- The role the talent prefers.
    preferred_role VARCHAR(150),

    -- The location where the talent prefers to work.
    preferred_location VARCHAR(150),

    -- Current availability of the talent.
    -- Example values: AVAILABLE, BUSY, NOTICE_PERIOD.
    availability_status VARCHAR(50),

    -- Time when the preference record was created.
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,

    -- Time when the preference record was last updated.
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- A talent normally has one current preference record.
-- This prevents duplicate preference records for the same talent.
CREATE UNIQUE INDEX IF NOT EXISTS idx_talent_preferences_talent_id
ON talent_preferences (talent_id);