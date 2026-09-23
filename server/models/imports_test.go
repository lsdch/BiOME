package models

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestImportBatchFileIdentity(t *testing.T) {
	input := ImportBatchWithFileInput{File: File{Name: "source.tsv", Size: 123, ContentType: "text/tab-separated-values"}}
	// The manager initializes the ID before storing the file.
	id := input.ID()
	require.Equal(t, id, input.ID())
	stored := input.WithFileHash("digest")
	params := stored.ToParams(uuid.New())
	require.Equal(t, id, params.ID)
	require.Equal(t, input.FileKey(), (ImportBatch{ID: params.ID}).FileKey())
	require.Equal(t, "digest", params.ImportedFileHash)
	require.Equal(t, input.File.Name, params.ImportedFileName)
	require.Equal(t, input.File.Size, params.ImportedFileSize)
}

func TestMaterializationRequiresEveryResolution(t *testing.T) {
	ready := MaterializationReadyCheck{Taxonomy: true, Methods: true, Fixatives: true, Bibliography: true}
	require.True(t, ready.IsReady())
	for _, field := range []string{"taxonomy", "methods", "fixatives", "bibliography"} {
		t.Run(field, func(t *testing.T) {
			check := ready
			switch field {
			case "taxonomy":
				check.Taxonomy = false
			case "methods":
				check.Methods = false
			case "fixatives":
				check.Fixatives = false
			case "bibliography":
				check.Bibliography = false
			}
			require.False(t, check.IsReady())
			require.Len(t, check.AppError().Errors, 1)
		})
	}
}
