package people

import (
	"context"
	"time"

	"github.com/geldata/gel-go/geltypes"
	"github.com/lsdch/biome/models"
	"github.com/sirupsen/logrus"
)

type UserShortIdentity struct {
	ID    geltypes.UUID `gel:"id" json:"id" format:"uuid"`
	Login string        `gel:"login" json:"login"`
	Name  string        `json:"name" gel:"name"`
	Alias string        `json:"alias" gel:"alias"`
}

// Metadata attached to an item in the database that track updates of the item.
type Meta struct {
	ID          geltypes.UUID                      `gel:"id" json:"-" swaggerignore:"true"`
	Created     time.Time                          `gel:"created" json:"created" binding:"required"`
	Modified    geltypes.OptionalDateTime          `gel:"modified" json:"modified,omitempty"`
	LastUpdated time.Time                          `gel:"lastUpdated" json:"last_updated"`
	CreatedBy   models.Optional[UserShortIdentity] `json:"created_by,omitempty" gel:"created_by"`
	UpdatedBy   models.Optional[UserShortIdentity] `json:"updated_by,omitempty" gel:"updated_by"`
}

func (m *Meta) Save(db geltypes.Executor) {
	if err := db.QuerySingle(context.Background(),
		`#edgeql
			select (<Meta><uuid>$0) { * }
		`, m, m.ID,
	); err != nil {
		logrus.Errorf("Failed to fetch updated Meta infos: %v", err)
	}
}

type MetaWithUser struct {
	Meta          `gel:"$inline" json:",inline"`
	CreatedByUser OptionalUser `json:"created_by_user,omitempty" gel:"created_by_user"`
	UpdatedByUser OptionalUser `json:"updated_by_user,omitempty" gel:"modified_by_user"`
}
