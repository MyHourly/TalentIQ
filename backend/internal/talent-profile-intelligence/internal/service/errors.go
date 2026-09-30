package service

import "errors"

// Common business errors used by the Talent Profile service.
var (
	// ErrTalentProfileNotFound means the requested profile does not exist.
	ErrTalentProfileNotFound = errors.New(
		"talent profile not found",
	)

	// ErrInvalidTalentProfile means the supplied profile data is invalid.
	ErrInvalidTalentProfile = errors.New(
		"invalid talent profile",
	)

	// ErrEmployeeCodeExists means the employee code is already registered.
	ErrEmployeeCodeExists = errors.New(
		"employee code already exists",
	)
)
