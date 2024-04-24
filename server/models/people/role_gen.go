// This file is auto-generated *DO NOT EDIT*

package people


import (
	"reflect"

  "fmt"
	"github.com/danielgtaylor/huma/v2"
  "github.com/go-faker/faker/v4"
	"math/rand"
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
  return UserRoleHierarchy[u] > UserRoleHierarchy[u]
}

func (u UserRole) IsGreaterEqual(v UserRole) bool {
  return UserRoleHierarchy[u] >= UserRoleHierarchy[u]
}

// Register enum in OpenAPI specification
func (u UserRole) Schema(r huma.Registry) *huma.Schema {
	schemaRef := r.Schema(reflect.TypeOf(""), true, "UserRole")
  schemaRef.Title = "UserRole"
  for _, v := range UserRoleValues {
	  schemaRef.Enum = append(schemaRef.Enum, string(v))
  }
  r.Map()["UserRole"] = schemaRef


  schema := r.Schema(reflect.TypeOf(""), true, "UserRole")
  schema.Ref = "#/components/schemas/UserRole"
	return schema
}

func init () {
  // Faker
  faker.AddProvider("UserRole",
    func(v reflect.Value) (interface{}, error) {
      idx := rand.Intn(len(UserRoleValues))
      fmt.Printf("Called provided for UserRole: %s\n", UserRole(UserRoleValues[idx]))
      return string(UserRoleValues[idx]), nil
    })
}

// EdgeDB Marshalling
func (m UserRole) MarshalEdgeDBStr() ([]byte, error) {
	return []byte(m), nil
}

func (m *UserRole) UnmarshalEdgeDBStr(data []byte) error {
	*m = UserRole(string(data))
	return nil
}