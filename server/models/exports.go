package models

type ExportFormat string

//generate:enum
const (
	ExportFormatCSV  ExportFormat = "csv"
	ExportFormatTSV  ExportFormat = "tsv"
	ExportFormatJSON ExportFormat = "json"
	ExportFormatDWC  ExportFormat = "dwc"
)
