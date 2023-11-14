package users

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
