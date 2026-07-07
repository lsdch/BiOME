package models

type OccurrencesAtPoint struct {
	Coordinates Coordinates               `json:"coordinates"`
	Samplings   []SamplingWithOccurrences `json:"samplings"`
}
