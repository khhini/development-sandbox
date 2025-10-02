package domain

import (
	"time"
)

type HailReport struct {
	Timestamp   time.Time `bigquery:"timestamp" json:"timestamp"`
	Time        string    `bigquery:"time" json:"time"`
	Size        int64     `bigquery:"size" json:"size"`
	Location    string    `bigquery:"location" json:"location"`
	County      string    `bigquery:"county" json:"county"`
	State       string    `bigquery:"state" json:"state"`
	Latitude    float32   `bigquery:"latitude" json:"latitude"`
	Longitude   float32   `bigquery:"longitude" json:"longitude"`
	Comments    string    `bigquery:"comments" json:"comments"`
	ReportPoint string    `bigquery:"report_point" json:"report_point"`
}
