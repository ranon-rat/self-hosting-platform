package executioner

const daysToDelete = 7

func DeletingOldShit() {
	ticker := time.NewTicker(time.Hour * 5)
	for range ticker.C {
		err := logRepo.DeleteOldMessages(daysToDelete)
		if err != nil {
			log.Println(err)
		}
	}
}