package institution

import (
	"darco/proto/models/institution"
	"net/http"

	_ "darco/proto/models/validations"

	"github.com/edgedb/edgedb-go"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// List institutions
// swagger:route GET /people/institutions/
// @Summary List Institutions
// @Tags People
// @Success 200 {array} institution.Institution
// @Router /people/institutions [get]
func List(ctx *gin.Context) {
	institutions, err := institution.List()
	if err != nil {
		ctx.Error(err)
	} else {
		ctx.JSON(http.StatusOK, institutions)
	}
}

// @Summary Create institution
// @Description Register a new institution that people work in.
// @id CreateInstitution
// @tags People
// @Accept json
// @Produce json
// @Success 202 {object} institution.Institution
// @Failure 400 {object} validations.FieldErrors
// @Router /people/institutions [post]
// @Param data body institution.InstitutionInput true "Institution informations"
func Create(ctx *gin.Context, db *edgedb.Client) {
	var inst institution.InstitutionInput
	if err := ctx.ShouldBindJSON(&inst); err != nil {
		ctx.Error(err)
		return
	}
	created, err := inst.Create(db)
	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
	}
	ctx.JSON(http.StatusAccepted, created)
}

// @Summary Delete institution
// @id DeleteInstitution
// @tags People
// @Accept json
// @Produce json
// @Success 202 "Delete successful"
// @Failure 404 "Institution does not exist"
// @Router /people/institutions/{acronym} [delete]
// @Param acronym path string true "Institution short name"
func Delete(ctx *gin.Context, db *edgedb.Client) {
	acronym := ctx.Param("acronym")
	logrus.Debugf("Deleting Institution : %v", acronym)
	inst, err := institution.Find(db, acronym)
	if err != nil {
		ctx.Error(err)
		return
	}
	if err := inst.Delete(db); err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
	}
}

// @Summary Update institution
// @id UpdateInstitution
// @tags People
// @Accept json
// @Produce json
// @Success 202 {object} institution.Institution
// @Failure 400 {object} validations.FieldErrors
// @Router /people/institutions/ [patch]
// @Param data body institution.Institution true "Institution informations"
func Update(ctx *gin.Context, db *edgedb.Client) {
	var inst institution.Institution
	if err := ctx.ShouldBindJSON(&inst); err != nil {
		ctx.Error(err)
		return
	}
	if err := inst.Update(db); err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
	}
	ctx.JSON(http.StatusAccepted, inst)
}
