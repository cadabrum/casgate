package types

import (
	"fmt"
	"strings"
)

type UserRole string

const (
	UserRoleUnknown     UserRole = "unknown"
	UserRoleGlobalAdmin UserRole = "admin"
	UserRolePartner     UserRole = "partner"
	UserRoleDistributor UserRole = "distributor"
)

type SubscriptionStateName string

func (s SubscriptionStateName) String() string {
	return string(s)
}

const (
	SubscriptionNew           SubscriptionStateName = "New"
	SubscriptionPending       SubscriptionStateName = "Pending"
	SubscriptionPreAuthorized SubscriptionStateName = "PreAuthorized"
	SubscriptionUnauthorized  SubscriptionStateName = "Unauthorized"
	SubscriptionAuthorized    SubscriptionStateName = "Authorized"
	SubscriptionStarted       SubscriptionStateName = "Started"
	SubscriptionPreFinished   SubscriptionStateName = "PreFinished"
	SubscriptionFinished      SubscriptionStateName = "Finished"
	SubscriptionCancelled     SubscriptionStateName = "Cancelled"
)

type SubscriptionStateNames []SubscriptionStateName

func (s SubscriptionStateNames) Contains(name SubscriptionStateName) bool {
	for _, value := range s {
		if value == name {
			return true
		}
	}

	return false
}

func (s SubscriptionStateNames) String() string {
	var strs []string
	for _, state := range s {
		strs = append(strs, state.String())
	}
	return strings.Join(strs, ", ")
}

type SubscriptionFieldName string

const (
	SubscriptionFieldNameName        SubscriptionFieldName = "Name"
	SubscriptionFieldNameDisplayName SubscriptionFieldName = "Display Name"
	SubscriptionFieldNameStartDate   SubscriptionFieldName = "Start Date"
	SubscriptionFieldNameEndDate     SubscriptionFieldName = "End Date"
	SubscriptionFieldNameSubUser     SubscriptionFieldName = "Sub user"
	SubscriptionFieldNameSubPlan     SubscriptionFieldName = "Sub plan"
	SubscriptionFieldNameDiscount    SubscriptionFieldName = "Discount"
	SubscriptionFieldNameDescription SubscriptionFieldName = "Description"
	SubscriptionFieldNameComment     SubscriptionFieldName = "Comment"
)

type SubscriptionFieldNames []SubscriptionFieldName

func (s SubscriptionFieldNames) Contains(name SubscriptionFieldName) bool {
	if s == nil {
		return false
	}

	for _, value := range s {
		if value == name {
			return true
		}
	}

	return false
}

type SubscriptionFieldPermissions map[UserRole]SubscriptionFieldNames
type SubscriptionTransitions map[UserRole]SubscriptionStateNames
type SubscriptionState struct {
	FieldPermissions SubscriptionFieldPermissions
	Transitions      SubscriptionTransitions
}

var SubscriptionStateMap = map[SubscriptionStateName]SubscriptionState{
	SubscriptionNew: {
		FieldPermissions: SubscriptionFieldPermissions{
			UserRolePartner: {
				SubscriptionFieldNameName,
				SubscriptionFieldNameDisplayName,
				SubscriptionFieldNameSubUser,
				SubscriptionFieldNameSubPlan,
				SubscriptionFieldNameDiscount,
				SubscriptionFieldNameDescription,
			},
		},
		Transitions: SubscriptionTransitions{
			UserRolePartner: SubscriptionStateNames{SubscriptionPending},
		},
	},
	SubscriptionPending: {
		FieldPermissions: SubscriptionFieldPermissions{
			UserRolePartner: {
				SubscriptionFieldNameDisplayName,
				SubscriptionFieldNameSubPlan,
				SubscriptionFieldNameDiscount,
				SubscriptionFieldNameDescription,
				SubscriptionFieldNameComment,
			},
		},
		Transitions: nil,
	},
	SubscriptionPreAuthorized: {
		FieldPermissions: SubscriptionFieldPermissions{
			UserRolePartner: {
				SubscriptionFieldNameDisplayName,
				SubscriptionFieldNameDescription,
			},
		},
		Transitions: SubscriptionTransitions{
			UserRolePartner: SubscriptionStateNames{SubscriptionAuthorized, SubscriptionCancelled},
		},
	},
	SubscriptionUnauthorized: {
		FieldPermissions: SubscriptionFieldPermissions{
			UserRolePartner: {
				SubscriptionFieldNameDisplayName,
				SubscriptionFieldNameSubPlan,
				SubscriptionFieldNameDiscount,
				SubscriptionFieldNameDescription,
			},
		},
		Transitions: SubscriptionTransitions{
			UserRolePartner: SubscriptionStateNames{SubscriptionPending, SubscriptionCancelled},
		},
	},
	SubscriptionAuthorized: {
		FieldPermissions: SubscriptionFieldPermissions{
			UserRoleDistributor: {
				SubscriptionFieldNameDisplayName,
				SubscriptionFieldNameStartDate,
				SubscriptionFieldNameDescription,
			},
		},
		Transitions: SubscriptionTransitions{
			UserRoleDistributor: SubscriptionStateNames{SubscriptionStarted, SubscriptionCancelled},
		},
	},
	SubscriptionStarted: {
		FieldPermissions: SubscriptionFieldPermissions{
			UserRolePartner: {
				SubscriptionFieldNameDisplayName,
				SubscriptionFieldNameDescription,
			},
		},
		Transitions: SubscriptionTransitions{
			UserRolePartner: SubscriptionStateNames{SubscriptionPreFinished},
		},
	},
	SubscriptionPreFinished: {
		FieldPermissions: SubscriptionFieldPermissions{
			UserRoleDistributor: {
				SubscriptionFieldNameDisplayName,
				SubscriptionFieldNameEndDate,
				SubscriptionFieldNameDescription,
			},
		},
		Transitions: SubscriptionTransitions{
			UserRoleDistributor: SubscriptionStateNames{SubscriptionFinished},
		},
	},
	SubscriptionFinished: {
		FieldPermissions: SubscriptionFieldPermissions{
			UserRolePartner: {
				SubscriptionFieldNameDisplayName,
				SubscriptionFieldNameDescription,
			},
		},
		Transitions: nil,
	},
	SubscriptionCancelled: {
		FieldPermissions: SubscriptionFieldPermissions{
			UserRolePartner: {
				SubscriptionFieldNameDisplayName,
				SubscriptionFieldNameDescription,
			},
		},
		Transitions: nil,
	},
}

func NewStateChangeForbiddenError(statusName SubscriptionStateName) error {
	return fmt.Errorf("You are not allowed to change the state to %s", statusName)
}