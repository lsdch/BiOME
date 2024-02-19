package user_role

import (
	"encoding/json"
)

type UserRole string // @name UserRole

const (
	Guest         UserRole = "Guest"
	Contributor   UserRole = "Contributor"
	ProjectMember UserRole = "ProjectMember"
	Admin         UserRole = "Admin"
)

func (m *UserRole) MarshalEdgeDBStr() ([]byte, error) {
	return []byte(*m), nil
}

func (m *UserRole) UnmarshalEdgeDBStr(data []byte) error {
	*m = UserRole(string(data))
	return nil
}

type OptionalUserRole struct {
	isSet bool
	role  UserRole
} //@name UserRole

func (m OptionalUserRole) MarshalJSON() ([]byte, error) {
	if m.isSet {
		return json.Marshal(m.role)
	}
	return json.Marshal(nil)
}

func (m *OptionalUserRole) UnmarshalEdgeDBStr(data []byte) error {
	if m.isSet {
		m.role = UserRole(string(data))
	}
	return nil
}

func (m *OptionalUserRole) SetMissing(isMissing bool) {
	m.isSet = !isMissing
}
