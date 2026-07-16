package types

import (
	"database/sql/driver"
	"fmt"

	"github.com/oklog/ulid/v2"
)

type ULID struct {
	ulid.ULID
}

func (u ULID) Value() (driver.Value, error) {
	return u.String(), nil
}

func (u *ULID) Scan(src any) error {
	switch x := src.(type) {
	case nil:
		u.ULID = ulid.ULID{}
		return nil

	case string:
		return u.UnmarshalText([]byte(x))

	case []byte:
		// PostgreSQL TEXT via pgx/database/sql arrive souvent ici
		return u.UnmarshalText(x)

	default:
		return fmt.Errorf("cannot scan ULID from %T", src)
	}
}

func MakeULID() ULID {
	return ULID{ulid.Make()}
}
