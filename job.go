package main

import (
	"fmt"
	"image"
	_ "image/jpeg"
	"io"
	"math/rand"
	"net/http"
	"os"
	"path"
	"time"
)

func DownloadImage(url string) (int, error) {
	response, err := http.Get(url)
	checkErr(err)

	defer response.Body.Close()

	filepath := path.Base(url)

	file, err := os.Create(filepath)
	checkErr(err)

	defer file.Close()

	_, err = io.Copy(file, response.Body)
	checkErr(err)

	var perimeter int = 0

	if reader, err := os.Open(filepath); err != nil {
		checkErr(err)
		return perimeter, err
	} else {
		defer reader.Close()
		if im, _, err := image.DecodeConfig(reader); err != nil {
			checkErr(err)
			return perimeter, err
		} else {
			perimeter = 2 * (im.Width + im.Height)
			fmt.Println(im.Width, im.Height)
			return perimeter, err
		}
	}
}

func (h handler) ProcessJob(images_payload ImagesPayload, job Job) {
	success := true

	for _, visit := range images_payload.Visits {
		for _, url := range visit.ImageUrls {
			var image Image

			perimeter, err := DownloadImage(url)

			n := int64(1 + rand.Float64()*3)
			time.Sleep(time.Duration(n) * time.Millisecond)

			image.StoreId = visit.StoreId
			image.Url = url
			image.Perimeter = perimeter
			image.Job = job
			image.JobId = int(job.Model.ID)

			if err != nil {
				image.Success = false
				image.ErrorMessage = err.Error()

				success = false
			} else {
				image.Success = true
				image.ErrorMessage = ""
			}

			if err := h.DB.Create(&image).Error; err != nil {
				panic(err)
			}
		}
	}

	if success {
		job.Status = "completed"
	} else {
		job.Status = "failed"
	}
	h.DB.Save(&job)
}
