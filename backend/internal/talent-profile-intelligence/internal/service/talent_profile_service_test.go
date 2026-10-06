package service

import (
	"context"
	"errors"
	"log/slog"
	"testing"
	"time"

	"github.com/google/uuid"

	"talentiq/talent-profile-intelligence/internal/model"
)

// mockTalentProfileRepository is a fake repository used only
// during unit tests.
//
// It does not connect to PostgreSQL.
type mockTalentProfileRepository struct {
	profiles map[uuid.UUID]*model.TalentProfile

	employeeCodes map[string]bool

	existsByEmployeeCodeResult bool

	existsByEmployeeCodeError error

	createError error

	getByIDError error

	listError error

	updateError error

	deleteError error

	deleteCalled bool
}

// GetByIDIncludingDeleted implements [repository.TalentProfileRepository].
// GetByIDIncludingDeleted simulates retrieving a profile
// even after it has been soft-deleted.
func (m *mockTalentProfileRepository) GetByIDIncludingDeleted(
	ctx context.Context,
	id uuid.UUID,
) (*model.TalentProfile, error) {

	if profile, exists := m.profiles[id]; exists {
		return profile, nil
	}

	return nil, ErrTalentProfileNotFound
}

// newMockTalentProfileRepository creates an empty mock repository.
func newMockTalentProfileRepository() *mockTalentProfileRepository {

	return &mockTalentProfileRepository{
		profiles:      make(map[uuid.UUID]*model.TalentProfile),
		employeeCodes: make(map[string]bool),
	}
}

// Create simulates creating a profile.
func (m *mockTalentProfileRepository) Create(
	ctx context.Context,
	profile *model.TalentProfile,
) (*model.TalentProfile, error) {

	if m.createError != nil {
		return nil, m.createError
	}

	if profile.ID == uuid.Nil {
		profile.ID = uuid.New()
	}

	now := time.Now()

	profile.CreatedAt = now
	profile.UpdatedAt = now

	m.profiles[profile.ID] = profile
	m.employeeCodes[profile.EmployeeCode] = true

	return profile, nil
}

// GetByID simulates retrieving a profile.
func (m *mockTalentProfileRepository) GetByID(
	ctx context.Context,
	id uuid.UUID,
) (*model.TalentProfile, error) {

	if m.getByIDError != nil {
		return nil, m.getByIDError
	}

	profile, exists := m.profiles[id]

	if !exists {
		return nil, ErrTalentProfileNotFound
	}

	return profile, nil
}

// List simulates retrieving paginated profiles.
func (m *mockTalentProfileRepository) List(
	ctx context.Context,
	limit int,
	offset int,
) ([]*model.TalentProfile, int, error) {

	if m.listError != nil {
		return nil, 0, m.listError
	}

	profiles := make(
		[]*model.TalentProfile,
		0,
		len(m.profiles),
	)

	for _, profile := range m.profiles {
		profiles = append(profiles, profile)
	}

	total := len(profiles)

	if offset >= total {
		return []*model.TalentProfile{}, total, nil
	}

	end := offset + limit

	if end > total {
		end = total
	}

	return profiles[offset:end], total, nil
}

// Update simulates updating a profile.
func (m *mockTalentProfileRepository) Update(
	ctx context.Context,
	profile *model.TalentProfile,
) (*model.TalentProfile, error) {

	if m.updateError != nil {
		return nil, m.updateError
	}

	if _, exists := m.profiles[profile.ID]; !exists {
		return nil, ErrTalentProfileNotFound
	}

	profile.UpdatedAt = time.Now()

	m.profiles[profile.ID] = profile

	return profile, nil
}

// Delete simulates soft deleting a profile.
func (m *mockTalentProfileRepository) Delete(
	ctx context.Context,
	id uuid.UUID,
) error {

	if m.deleteError != nil {
		return m.deleteError
	}

	m.deleteCalled = true

	profile, exists := m.profiles[id]

	if !exists {
		return ErrTalentProfileNotFound
	}

	// Simulate the soft-delete behavior.
	profile.ProfileStatus = "INACTIVE"

	return nil
}

// ExistsByEmployeeCode simulates checking a duplicate employee code.
func (m *mockTalentProfileRepository) ExistsByEmployeeCode(
	ctx context.Context,
	employeeCode string,
) (bool, error) {

	if m.existsByEmployeeCodeError != nil {
		return false, m.existsByEmployeeCodeError
	}

	return m.existsByEmployeeCodeResult ||
		m.employeeCodes[employeeCode], nil
}

// newTestService creates the service with the mock repository.
func newTestService(
	repository *mockTalentProfileRepository,
) TalentProfileService {

	logger := slog.Default()

	return NewTalentProfileService(
		repository,
		nil,
		logger,
	)
}

// validTalentProfile returns a valid profile that can be
// reused across multiple tests.
func validTalentProfile() *model.TalentProfile {

	return &model.TalentProfile{
		EmployeeCode:         "EMP001",
		FirstName:            "Ashutosh",
		LastName:             "Kumar",
		Email:                "ashutosh@example.com",
		Designation:          "Software Engineer",
		Department:           "Engineering",
		Location:             "Bangalore",
		Summary:              "Backend developer",
		TotalExperienceYears: 2.5,
	}
}

// stringPointer converts a string into a *string.
func stringPointer(value string) *string {
	return &value
}

// TestCreateTalentProfile tests successful profile creation.
func TestCreateTalentProfile(t *testing.T) {

	repository := newMockTalentProfileRepository()

	service := newTestService(repository)

	profile := validTalentProfile()

	result, err := service.Create(
		context.Background(),
		profile,
	)

	if err != nil {
		t.Fatalf(
			"expected no error, got %v",
			err,
		)
	}

	if result == nil {
		t.Fatal("expected profile, got nil")
	}

	if result.ID == uuid.Nil {
		t.Error("expected generated profile ID")
	}

	if result.ProfileStatus != "ACTIVE" {
		t.Errorf(
			"expected ACTIVE status, got %s",
			result.ProfileStatus,
		)
	}
}

// TestCreateTalentProfileWithoutEmployeeCode tests validation
// when employee code is missing.
func TestCreateTalentProfileWithoutEmployeeCode(
	t *testing.T,
) {

	repository := newMockTalentProfileRepository()

	service := newTestService(repository)

	profile := validTalentProfile()

	profile.EmployeeCode = ""

	_, err := service.Create(
		context.Background(),
		profile,
	)

	if err == nil {
		t.Fatal("expected validation error")
	}

	if !errors.Is(
		err,
		ErrInvalidTalentProfile,
	) {
		t.Errorf(
			"expected ErrInvalidTalentProfile, got %v",
			err,
		)
	}
}

// TestCreateDuplicateEmployeeCode verifies that the service
// prevents duplicate employee codes.
func TestCreateDuplicateEmployeeCode(
	t *testing.T,
) {

	repository := newMockTalentProfileRepository()

	repository.existsByEmployeeCodeResult = true

	service := newTestService(repository)

	profile := validTalentProfile()

	_, err := service.Create(
		context.Background(),
		profile,
	)

	if err == nil {
		t.Fatal("expected duplicate employee code error")
	}

	if !errors.Is(
		err,
		ErrEmployeeCodeExists,
	) {
		t.Errorf(
			"expected ErrEmployeeCodeExists, got %v",
			err,
		)
	}
}

// TestGetTalentProfile tests retrieving an existing profile.
func TestGetTalentProfile(t *testing.T) {

	repository := newMockTalentProfileRepository()

	profile := validTalentProfile()

	profile.ID = uuid.New()

	repository.profiles[profile.ID] = profile

	service := newTestService(repository)

	result, err := service.GetByID(
		context.Background(),
		profile.ID,
	)

	if err != nil {
		t.Fatalf(
			"expected no error, got %v",
			err,
		)
	}

	if result.ID != profile.ID {
		t.Errorf(
			"expected ID %s, got %s",
			profile.ID,
			result.ID,
		)
	}
}

// TestGetTalentProfileWithInvalidID verifies validation
// of an empty UUID.
func TestGetTalentProfileWithInvalidID(
	t *testing.T,
) {

	repository := newMockTalentProfileRepository()

	service := newTestService(repository)

	_, err := service.GetByID(
		context.Background(),
		uuid.Nil,
	)

	if err == nil {
		t.Fatal("expected invalid profile error")
	}

	if !errors.Is(
		err,
		ErrInvalidTalentProfile,
	) {
		t.Errorf(
			"expected ErrInvalidTalentProfile, got %v",
			err,
		)
	}
}

// TestListTalentProfiles tests pagination handling.
func TestListTalentProfiles(t *testing.T) {

	repository := newMockTalentProfileRepository()

	for i := 0; i < 3; i++ {

		profile := validTalentProfile()

		profile.ID = uuid.New()
		profile.EmployeeCode = "EMP00" +
			string(rune('1'+i))

		repository.profiles[profile.ID] = profile
	}

	service := newTestService(repository)

	profiles, total, err := service.List(
		context.Background(),
		2,
		0,
	)

	if err != nil {
		t.Fatalf(
			"expected no error, got %v",
			err,
		)
	}

	if len(profiles) != 2 {
		t.Errorf(
			"expected 2 profiles, got %d",
			len(profiles),
		)
	}

	if total != 3 {
		t.Errorf(
			"expected total 3, got %d",
			total,
		)
	}
}

// TestListTalentProfilesWithInvalidLimit verifies that
// invalid pagination values receive safe defaults.
func TestListTalentProfilesWithInvalidLimit(
	t *testing.T,
) {

	repository := newMockTalentProfileRepository()

	service := newTestService(repository)

	_, _, err := service.List(
		context.Background(),
		-10,
		-5,
	)

	if err != nil {
		t.Fatalf(
			"expected no error, got %v",
			err,
		)
	}
}

// TestUpdateTalentProfile tests successful updates.
func TestUpdateTalentProfile(t *testing.T) {

	repository := newMockTalentProfileRepository()

	profile := validTalentProfile()

	profile.ID = uuid.New()

	repository.profiles[profile.ID] = profile

	service := newTestService(repository)

	profile.Designation = "Senior Software Engineer"

	result, err := service.Update(
		context.Background(),
		profile,
	)

	if err != nil {
		t.Fatalf(
			"expected no error, got %v",
			err,
		)
	}

	if result.Designation !=
		"Senior Software Engineer" {

		t.Errorf(
			"expected updated designation, got %s",
			result.Designation,
		)
	}
}

// TestUpdateTalentProfileWithInvalidID verifies that
// an update cannot happen without a valid ID.
func TestUpdateTalentProfileWithInvalidID(
	t *testing.T,
) {

	repository := newMockTalentProfileRepository()

	service := newTestService(repository)

	profile := validTalentProfile()

	profile.ID = uuid.Nil

	_, err := service.Update(
		context.Background(),
		profile,
	)

	if err == nil {
		t.Fatal("expected invalid profile error")
	}

	if !errors.Is(
		err,
		ErrInvalidTalentProfile,
	) {
		t.Errorf(
			"expected ErrInvalidTalentProfile, got %v",
			err,
		)
	}
}

// TestDeleteTalentProfile tests profile deactivation.
func TestDeleteTalentProfile(t *testing.T) {

	repository := newMockTalentProfileRepository()

	profile := validTalentProfile()

	profile.ID = uuid.New()

	repository.profiles[profile.ID] = profile

	service := newTestService(repository)

	err := service.Delete(
		context.Background(),
		profile.ID,
	)

	if err != nil {
		t.Fatalf(
			"expected no error, got %v",
			err,
		)
	}

	if !repository.deleteCalled {
		t.Error("expected delete repository method to be called")
	}

	if profile.ProfileStatus != "INACTIVE" {
		t.Errorf(
			"expected INACTIVE status, got %s",
			profile.ProfileStatus,
		)
	}
}

// TestDeleteTalentProfileWithInvalidID verifies that
// deletion requires a valid profile ID.
func TestDeleteTalentProfileWithInvalidID(
	t *testing.T,
) {

	repository := newMockTalentProfileRepository()

	service := newTestService(repository)

	err := service.Delete(
		context.Background(),
		uuid.Nil,
	)

	if err == nil {
		t.Fatal("expected invalid profile error")
	}

	if !errors.Is(
		err,
		ErrInvalidTalentProfile,
	) {
		t.Errorf(
			"expected ErrInvalidTalentProfile, got %v",
			err,
		)
	}
}
