package executionerServices

import (
	"log"
	"time"
)

const daysToDelete = 7

func DeletingOldShit() {
	// i am not interested in having any old logs
	ticker := time.NewTicker(time.Hour * 5)
	for range ticker.C {
		err := logRepo.DeleteOldMessages(daysToDelete)
		if err != nil {
			log.Println(err)
		}
	}
}
