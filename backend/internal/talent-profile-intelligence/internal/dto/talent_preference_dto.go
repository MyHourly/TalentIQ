package dto

// TalentPreferenceRequest represents the data required
// to create or update a talent's preferences.
//
// TalentID is not included here because it should come
// from the API path parameter rather than the request body.
type TalentPreferenceRequest struct {
	PreferredRole      string `json:"preferred_role"`
	PreferredLocation  string `json:"preferred_location"`
	AvailabilityStatus string `json:"availability_status"`
}

// TalentPreferenceResponse represents the data returned
// to the API client.
type TalentPreferenceResponse struct {
	ID                 string `json:"id"`
	TalentID           string `json:"talent_id"`
	PreferredRole      string `json:"preferred_role"`
	PreferredLocation  string `json:"preferred_location"`
	AvailabilityStatus string `json:"availability_status"`
	CreatedAt          string `json:"created_at"`
	UpdatedAt          string `json:"updated_at"`
}

// TalentAvailabilityRequest contains the availability
// status that can be updated independently.
type TalentAvailabilityRequest struct {
	AvailabilityStatus string `json:"availability_status"`
}
