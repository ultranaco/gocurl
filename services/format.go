package services

import (
	"fmt"
	"time"
)

type Format struct{}

func (service Format) ProgressFromDateTime(start, finish, current time.Time) float64 {
	unit := (float64(finish.UnixNano()) - float64(start.UnixNano())) * 0.01
	progress := (float64(current.UnixNano()) - float64(start.UnixNano())) / unit
	return progress
}

func (service Format) TimeProgressLoad(progressInput float64) string {

	progress := int(progressInput)

	var progressBar string

	for iterator := 0; iterator < 25; iterator++ {
		if iterator < (int(progress) / 4) {
			progressBar += "█"
		} else {
			progressBar += "▁"
		}
	}

	return fmt.Sprintf("[%s] %v%%", progressBar, progress)
}
