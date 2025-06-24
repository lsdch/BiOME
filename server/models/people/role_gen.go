// This file is auto-generated *DO NOT EDIT*

package people


import (
	"reflect"

	"github.com/danielgtaylor/huma/v2"
  "github.com/brianvoe/gofakeit/v7"
)



var UserRoleValues = []UserRole{
	Visitor,
	Contributor,
	Maintainer,
	Admin,
}
var UserRoleHierarchy = map[UserRole]int{
  "Visitor": 0,
  "Contributor": 1,
  "Maintainer": 2,
  "Admin": 3,
}

func (u UserRole) IsGreater(v UserRole) bool {
  return UserRoleHierarchy[u] > UserRoleHierarchy[v]
}

func (u UserRole) IsGreaterEqual(v UserRole) bool {
  return UserRoleHierarchy[u] >= UserRoleHierarchy[v]
}

// Register enum in OpenAPI specification
func (u UserRole) Schema(r huma.Registry) *huma.Schema {
  if r.Map()["UserRole"] == nil {
    schemaRef := r.Schema(reflect.TypeOf(""), true, "UserRole")
    schemaRef.Title = "UserRole"
    for _, v := range UserRoleValues {
      schemaRef.Enum = append(schemaRef.Enum, string(v))
    }
    r.Map()["UserRole"] = schemaRef
  }

	return &huma.Schema{Ref: "#/components/schemas/UserRole"}
}

func (m *UserRole) Fake(f *gofakeit.Faker) (any, error) {
	return string(UserRoleValues[f.IntN(len(UserRoleValues) - 1)]), nil
}

// Gel Marshalling
func (m UserRole) MarshalEdgeDBStr() ([]byte, error) {
	return []byte(m), nil
}

func (m *UserRole) UnmarshalEdgeDBStr(data []byte) error {
	*m = UserRole(string(data))
	return nil
}