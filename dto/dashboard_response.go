package dto

type DashboardResponse struct {
	TotalTechnologies     int `json:"totalTechnologies"`
	GrowingTechnologies   int `json:"growingTechnologies"`
	DecliningTechnologies int `json:"decliningTechnologies"`
	EmergingSkills        int `json:"emergingSkills"`
}
