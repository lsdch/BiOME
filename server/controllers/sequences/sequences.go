package sequences

import (
	"darco/proto/router"
)

func RegisterRoutes(r router.Router) {
	RegisterGeneRoutes(r)
	RegisterSeqDBRoutes(r)
}
