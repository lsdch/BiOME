package people

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/edgedb/edgedb-go"
)

type UserRole string // @name UserRole

const (
	Visitor     UserRole = "Visitor"
	Contributor UserRole = "Contributor"
	Maintainer  UserRole = "Maintainer"
	Admin       UserRole = "Admin"
)

func (m UserRole) MarshalEdgeDBStr() ([]byte, error) {
	return []byte(m), nil
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

func (u *User) SetRole(db edgedb.Executor, role UserRole) error {
	if err := db.Execute(context.Background(),
		`update people::User set {
			role := <people::UserRole>$0
		}`,
		string(role),
	); err != nil {
		return fmt.Errorf("Failed to set user role: %v", err)
	}
	u.Role = role
	return nil
}
