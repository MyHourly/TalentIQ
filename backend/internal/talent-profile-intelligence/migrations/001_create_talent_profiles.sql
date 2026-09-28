CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE IF NOT EXISTS talent_profiles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    employee_code VARCHAR(50) NOT NULL UNIQUE,

    first_name VARCHAR(100) NOT NULL,

    last_name VARCHAR(100),

    email VARCHAR(255) NOT NULL UNIQUE,

    phone VARCHAR(30),

    designation VARCHAR(150),

    department VARCHAR(150),

    location VARCHAR(150),

    summary TEXT,

    total_experience_years NUMERIC(4,1) DEFAULT 0,

    profile_status VARCHAR(30) NOT NULL DEFAULT 'ACTIVE',

    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    deleted_at TIMESTAMP NULL
);

CREATE INDEX IF NOT EXISTS idx_talent_profiles_employee_code
    ON talent_profiles(employee_code);

CREATE INDEX IF NOT EXISTS idx_talent_profiles_email
    ON talent_profiles(email);

CREATE INDEX IF NOT EXISTS idx_talent_profiles_department
    ON talent_profiles(department);

CREATE INDEX IF NOT EXISTS idx_talent_profiles_status
    ON talent_profiles(profile_status);