package models

type ExportFormat string

//generate:enum
const (
	ExportFormatCSV  ExportFormat = "csv"
	ExportFormatJSON ExportFormat = "json"
	ExportFormatDWC  ExportFormat = "dwc"
)
