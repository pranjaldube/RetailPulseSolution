package main

import "gorm.io/gorm"

// json responses
type StorePayload struct {
	StoreId   string   `json:"store_id"`
	ImageUrls []string `json:"image_url"`
	VisitTime string   `json:"visit_time"`
}

type ImagesPayload struct {
	Count  int            `json:"count"`
	Visits []StorePayload `json:"visits"`
}

type JobError struct {
	StoreId string `json:"store_id"`
	Error   string `json:"error"`
}

type JobStatusResponse struct {
	JobId  string     `json:"job_id"`
	Status string     `json:"status"`
	Error  []JobError `json:"error"`
}

// db models
type Job struct {
	gorm.Model
	Status string `json:"status"`
}

type Image struct {
	gorm.Model
	JobId        int    `json:"job_id"`
	Job          Job    `json:"job" gorm:"foreignKey:JobId"`
	StoreId      string `json:"store_id"`
	Url          string `json:"url"`
	Perimeter    int    `json:"perimeter"`
	Success      bool   `json:"success"`
	ErrorMessage string `json:"error"`
}
