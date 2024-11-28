package people_test

import (
	"darco/proto/db"
	"darco/proto/models/people"
	"darco/proto/tests"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func SetupInstitution(t *testing.T) people.Institution {
	var input = tests.FakeData[people.InstitutionInput](t)
	inst, err := input.Create(db.Client())
	require.NoError(t, err)
	return inst
}

func TestCreateInstitution(t *testing.T) {
	var input = tests.FakeData[people.InstitutionInput](t)
	fmt.Printf("%+v\n", input)
	inst, err := input.Create(db.Client())
	require.NoError(t, err)
	assert.Equal(t, input.Name, inst.Name)
}

func TestFindInstitution(t *testing.T) {
	inst := SetupInstitution(t)
	res, err := people.FindInstitution(db.Client(), inst.ID)
	require.NoError(t, err)
	assert.Equal(t, inst, res)
}

func TestListInstitution(t *testing.T) {
	_ = SetupInstitution(t)
	insts, err := people.ListInstitutions(db.Client())
	require.NoError(t, err)
	assert.NotEmpty(t, insts)
}

func TestDeleteInstitution(t *testing.T) {
	inst := SetupInstitution(t)
	deleted, err := inst.Delete(db.Client())
	require.NoError(t, err)
	assert.Equal(t, inst, deleted)
}

func TestUpdateInstitution(t *testing.T) {
	inst := SetupInstitution(t)
	input := tests.FakeData[people.InstitutionUpdate](t)
	input.Code.IsSet = true
	input.Description.Null = true
	input.Description.IsSet = true
	updated, err := input.Save(db.Client(), inst.Code)
	require.NoError(t, err)
	assert.Equal(t, input.Code.Value, updated.Code)
}
