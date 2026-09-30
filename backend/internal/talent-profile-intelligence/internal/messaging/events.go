package messaging

import (
	"time"

	"github.com/google/uuid"
)

// TalentProfileEvent represents a change related to a talent profile.
//
// The event contains metadata required by other services
// to understand what changed and which talent was affected.
type TalentProfileEvent struct {
	EventID      uuid.UUID `json:"event_id"`
	EventType    string    `json:"event_type"`
	TalentID     uuid.UUID `json:"talent_id"`
	EmployeeCode string    `json:"employee_code"`
	OccurredAt   time.Time `json:"occurred_at"`
}
