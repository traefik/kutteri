package locker

import (
	"io/ioutil"
	"log"
	"os"
	"time"
)

const ghDateLayout = "2006-01-02T15:04:05-07:00"

// TimeLocker Manage the time lock file.
type TimeLocker struct {
	FilePath string
	HourBack time.Duration
}

// GetLastTime Get the last time.
func (l TimeLocker) GetLastTime() (string, error) {
	if _, err := os.Stat(l.FilePath); err != nil {
		if os.IsNotExist(err) {
			log.Printf("No existing file %s, created one.", l.FilePath)
			return l.SaveLastTime()
		}
		return "", err
	}

	data, err := ioutil.ReadFile(l.FilePath)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

// SaveLastTime Save the current time.
func (l TimeLocker) SaveLastTime() (string, error) {
	srcDate := time.Now().Add(-l.HourBack * time.Hour)
	date := srcDate.Format(ghDateLayout)

	err := ioutil.WriteFile(l.FilePath, []byte(date), 0644)
	if err != nil {
		return "", err
	}

	return date, nil
}
