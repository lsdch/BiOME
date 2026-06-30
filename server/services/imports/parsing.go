package imports

import (
	"fmt"
)

type CSVParseError struct {
	RowNumber int32
	Err       error
}

func (e *CSVParseError) Error() string {
	return fmt.Sprintf("error parsing CSV at row %d: %v", e.RowNumber, e.Err)
}
