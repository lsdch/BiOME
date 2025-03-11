package people

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/danielgtaylor/huma/v2"
	"github.com/geldata/gel-go/geltypes"
)

type UserRole string // @name UserRole

// *CAUTION* the enum order DOES MATTER as it establishes roles hierarchy
// in the generated code. See [./roles_gen.go].

//generate:order-enum
const (
	Visitor     UserRole = "Visitor"
	Contributor UserRole = "Contributor"
	Maintainer  UserRole = "Maintainer"
	Admin       UserRole = "Admin"
)

func (u *User) SetRole(db geltypes.Executor, role UserRole) error {
	if err := db.Execute(context.Background(),
		`#edgeql
			update people::User set { role := <people::UserRole>$0 }
		`, string(role),
	); err != nil {
		return fmt.Errorf("Failed to set user role: %w", err)
	}
	u.Role = role
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

func (m *OptionalUserRole) Schema(r huma.Registry) *huma.Schema {
	return r.Schema(reflect.TypeFor[UserRole](), true, "")
}
