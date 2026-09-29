package dto

import (
	"time"

	"github.com/google/uuid"

	"talentiq/talent-profile-intelligence/internal/model"
)

// TalentProfileCreateRequest represents the JSON body
// accepted when creating a talent profile.
type TalentProfileCreateRequest struct {
	EmployeeCode         string  `json:"employee_code"`
	FirstName            string  `json:"first_name"`
	LastName             string  `json:"last_name"`
	Email                string  `json:"email"`
	Phone                *string `json:"phone"`
	Designation          string  `json:"designation"`
	Department           string  `json:"department"`
	Location             *string `json:"location"`
	Summary              *string `json:"summary"`
	TotalExperienceYears float64 `json:"total_experience_years"`
}

// TalentProfileUpdateRequest represents the fields that
// can be changed through the update API.
type TalentProfileUpdateRequest struct {
	FirstName            string  `json:"first_name"`
	LastName             string  `json:"last_name"`
	Email                string  `json:"email"`
	Phone                *string `json:"phone"`
	Designation          string  `json:"designation"`
	Department           string  `json:"department"`
	Location             *string `json:"location"`
	Summary              *string `json:"summary"`
	TotalExperienceYears float64 `json:"total_experience_years"`
	ProfileStatus        string  `json:"profile_status"`
}

// TalentProfileResponse represents a talent profile returned
// to API consumers.
type TalentProfileResponse struct {
	ID                   uuid.UUID `json:"id"`
	EmployeeCode         string    `json:"employee_code"`
	FirstName            string    `json:"first_name"`
	LastName             string    `json:"last_name"`
	Email                string    `json:"email"`
	Phone                *string   `json:"phone"`
	Designation          string    `json:"designation"`
	Department           string    `json:"department"`
	Location             *string   `json:"location"`
	Summary              *string   `json:"summary"`
	TotalExperienceYears float64   `json:"total_experience_years"`
	ProfileStatus        string    `json:"profile_status"`
	CreatedAt            time.Time `json:"created_at"`
	UpdatedAt            time.Time `json:"updated_at"`
}

// TalentProfileListResponse represents a paginated
// collection of talent profiles.
type TalentProfileListResponse struct {
	Data  []*TalentProfileResponse `json:"data"`
	Page  int                      `json:"page"`
	Limit int                      `json:"limit"`
	Total int                      `json:"total"`
}

// ToModel converts the create request DTO into the
// internal TalentProfile model.
func (r *TalentProfileCreateRequest) ToModel() *model.TalentProfile {
	return &model.TalentProfile{
		EmployeeCode:         r.EmployeeCode,
		FirstName:            r.FirstName,
		LastName:             r.LastName,
		Email:                r.Email,
		Phone:                *r.Phone,
		Designation:          r.Designation,
		Department:           r.Department,
		Location:             *r.Location,
		Summary:              *r.Summary,
		TotalExperienceYears: r.TotalExperienceYears,
	}
}

// ToModel converts the update request DTO into the
// internal TalentProfile model.
func (r *TalentProfileUpdateRequest) ToModel(
	id uuid.UUID,
	employeeCode string,
) *model.TalentProfile {

	return &model.TalentProfile{
		ID:                   id,
		EmployeeCode:         employeeCode,
		FirstName:            r.FirstName,
		LastName:             r.LastName,
		Email:                r.Email,
		Phone:                *r.Phone,
		Designation:          r.Designation,
		Department:           r.Department,
		Location:             *r.Location,
		Summary:              *r.Summary,
		TotalExperienceYears: r.TotalExperienceYears,
		ProfileStatus:        r.ProfileStatus,
	}
}

// FromModel converts the internal model into an API response DTO.
func FromTalentProfileModel(
	profile *model.TalentProfile,
) *TalentProfileResponse {

	return &TalentProfileResponse{
		ID:                   profile.ID,
		EmployeeCode:         profile.EmployeeCode,
		FirstName:            profile.FirstName,
		LastName:             profile.LastName,
		Email:                profile.Email,
		Phone:                &profile.Phone,
		Designation:          profile.Designation,
		Department:           profile.Department,
		Location:             &profile.Location,
		Summary:              &profile.Summary,
		TotalExperienceYears: profile.TotalExperienceYears,
		ProfileStatus:        profile.ProfileStatus,
		CreatedAt:            profile.CreatedAt,
		UpdatedAt:            profile.UpdatedAt,
	}
}
