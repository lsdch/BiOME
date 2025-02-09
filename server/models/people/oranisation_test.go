package people_test

import (
	"fmt"
	"testing"

	"github.com/lsdch/biome/db"
	"github.com/lsdch/biome/models/people"
	"github.com/lsdch/biome/tests"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func SetupOrganisation(t *testing.T) people.Organisation {
	var input = tests.FakeData[people.OrganisationInput](t)
	inst, err := input.Save(db.Client())
	require.NoError(t, err)
	return inst
}

func TestCreateOrganisation(t *testing.T) {
	var input = tests.FakeData[people.OrganisationInput](t)
	fmt.Printf("%+v\n", input)
	inst, err := input.Save(db.Client())
	require.NoError(t, err)
	assert.Equal(t, input.Name, inst.Name)
}

func TestFindOrganisation(t *testing.T) {
	inst := SetupOrganisation(t)
	res, err := people.FindOrganisation(db.Client(), inst.ID)
	require.NoError(t, err)
	assert.Equal(t, inst, res)
}

func TestListOrganisation(t *testing.T) {
	_ = SetupOrganisation(t)
	insts, err := people.ListOrganisations(db.Client())
	require.NoError(t, err)
	assert.NotEmpty(t, insts)
}

func TestDeleteOrganisation(t *testing.T) {
	inst := SetupOrganisation(t)
	deleted, err := inst.Delete(db.Client())
	require.NoError(t, err)
	assert.Equal(t, inst, deleted)
}

func TestUpdateOrganisation(t *testing.T) {
	inst := SetupOrganisation(t)
	input := tests.FakeData[people.OrganisationUpdate](t)
	input.Code.IsSet = true
	input.Description.Null = true
	input.Description.IsSet = true
	updated, err := input.Save(db.Client(), inst.Code)
	require.NoError(t, err)
	assert.Equal(t, input.Code.Value, updated.Code)
}
