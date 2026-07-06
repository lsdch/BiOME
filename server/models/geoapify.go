package models

import (
	"fmt"

	"github.com/lsdch/biome/db/biomedb"
)

type GeoapifyCoords struct {
	Lat float32 `json:"lat"`
	Lon float32 `json:"lon"`
}

type GeoapifyUsage biomedb.GeoapifyUsage

func GeoapifyUsageFromDB(u biomedb.GeoapifyUsage) GeoapifyUsage {
	return GeoapifyUsage(u)
}

type GeoapifyStatus struct {
	APIKey        string `json:"-"`
	Available     bool   `json:"available"`
	HasApiKey     bool   `json:"has_api_key"`
	TodayRequests int32  `json:"requests"`
	Limit         int32  `json:"limit"`
}

type GeoapifyRequestDenial error

var (
	ErrNoAPIKey      = GeoapifyRequestDenial(fmt.Errorf("Geoapify API key is not set"))
	ErrLimitExceeded = GeoapifyRequestDenial(fmt.Errorf("Geoapify usage limit exceeded"))
)

func (s *GeoapifyStatus) AllowRequests(n int32) GeoapifyRequestDenial {
	if !s.HasApiKey {
		return ErrNoAPIKey
	}
	if s.TodayRequests+n > s.Limit {
		return ErrLimitExceeded
	}
	return nil
}

type GeoapifyPendingResponse struct {
	ID     string `json:"id"`
	Status string `json:"status"`
	URL    string `json:"url"`
}

type GeoapifyResult struct {
	Formatted    string  `json:"formatted"`
	Municipality string  `json:"municipality"`
	City         string  `json:"city"`
	County       string  `json:"county"`
	State        string  `json:"state"`
	Country      string  `json:"country"`
	CountryCode  string  `json:"country_code"`
	Lat          float64 `json:"lat"`
	Lon          float64 `json:"lon"`
	PostalCode   string  `json:"postcode"`
	Street       string  `json:"street"`
	HouseNumber  string  `json:"housenumber"`
	Suburb       string  `json:"suburb"`
}

type ReverseGeoCodeResponse struct {
	Results []GeoapifyResult       `json:"results"`
	Query   map[string]interface{} `json:"query"`
}
