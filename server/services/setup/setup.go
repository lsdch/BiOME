package setup

import (
	"context"
	"darco/proto/models/people"
	"darco/proto/models/settings"

	"github.com/edgedb/edgedb-go"
)

func Init(db *edgedb.Client, superadmin people.SuperAdminInput) error {
	return db.Tx(context.Background(), func(context context.Context, tx *edgedb.Tx) error {
		return InitTx(tx, superadmin)
	})
}

func InitTx(tx *edgedb.Tx, superadmin people.SuperAdminInput) error {
	admin, err := superadmin.Save(tx)
	if err != nil {
		return err
	}
	return settings.Setup(tx, admin.ID)
}
