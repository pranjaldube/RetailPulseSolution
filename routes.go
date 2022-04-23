package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

func (h handler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := map[string]string{"status": "success", "message": "server is healthy"}
	json.NewEncoder(w).Encode(response)
}

func (h handler) JobSubmit(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var images_payload ImagesPayload

	err := decoder.Decode(&images_payload)
	checkErr(err)

	if images_payload.Count != len(images_payload.Visits) {
		w.WriteHeader(http.StatusBadRequest)
		response := map[string]string{}
		json.NewEncoder(w).Encode(response)
		return
	}

	var job Job
	job.Status = "ongoing"

	err = h.DB.Create(&job).Error
	checkErr(err)

	// goroutine
	go h.ProcessJob(images_payload, job)

	var response = map[string]int{"job_id": int(job.Model.ID)}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (h handler) JobStatus(w http.ResponseWriter, r *http.Request) {
	jobid := r.URL.Query().Get("jobid")

	if jobid == "" {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		response := map[string]string{"status": "error", "message": "you are missing jobid query parameter in the job status api"}
		json.NewEncoder(w).Encode(response)
		return
	}

	jobid_int, err := strconv.Atoi(jobid)
	checkErr(err)

	var job Job

	if result := h.DB.First(&job, jobid_int); result.Error != nil {
		fmt.Println(result.Error)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		response := map[string]string{}
		json.NewEncoder(w).Encode(response)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if job.Status == "completed" || job.Status == "ongoing" {
		response := map[string]string{"status": job.Status, "job_id": jobid}
		json.NewEncoder(w).Encode(response)
		return
	}

	if job.Status == "failed" {
		images := make([]Image, 0)

		joberrors := make([]map[string]string, 0)

		err := h.DB.Where("job_id = ? AND success = false", int(job.Model.ID)).Find(&images).Error
		checkErr(err)

		for _, image := range images {
			joberrors = append(joberrors, map[string]string{"store_id": image.StoreId, "error": image.ErrorMessage})
		}

		response := map[string]interface{}{"status": job.Status, "job_id": jobid, "error": joberrors}
		json.NewEncoder(w).Encode(response)
		return
	}
}
