package util

type DataPoint struct {
	Value float64 `json:"value"`
	Unit  string  `json:"unit"`
}

var results []DataPoint

func GetData(text string) []DataPoint {
}