-- TPI-003
-- Creates the main Talent Profile table.

-- UUID generation is used for profile IDs.
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE IF NOT EXISTS talent_profiles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    -- Internal employee identifier.
    employee_code VARCHAR(50) NOT NULL UNIQUE,

    -- Basic personal information.
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100),

    -- Contact information.
    email VARCHAR(255) NOT NULL UNIQUE,
    phone VARCHAR(30),

    -- Professional information.
    designation VARCHAR(150),
    department VARCHAR(150),
    location VARCHAR(150),

    -- Short professional summary.
    summary TEXT,

    -- Total professional experience in years.
    total_experience_years NUMERIC(4,1) NOT NULL DEFAULT 0,

    -- Current profile state.
    profile_status VARCHAR(30) NOT NULL DEFAULT 'ACTIVE',

    -- Audit timestamps.
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    -- Soft-delete timestamp.
    deleted_at TIMESTAMP NULL
);

-- Indexes support common lookup operations.
CREATE INDEX IF NOT EXISTS idx_talent_profiles_employee_code
    ON talent_profiles(employee_code);

CREATE INDEX IF NOT EXISTS idx_talent_profiles_email
    ON talent_profiles(email);

CREATE INDEX IF NOT EXISTS idx_talent_profiles_department
    ON talent_profiles(department);

CREATE INDEX IF NOT EXISTS idx_talent_profiles_status
    ON talent_profiles(profile_status);