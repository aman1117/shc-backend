package main

import (
	"fmt"

	"github.com/aj-2000/shc-backend/services"
)

// can't understood how this function is working
func runCronJobs(as *services.AppService) {
	as.CronService.AddFunc("@midnight", func() {
		print("Deactivating all expired subscriptions")
		err := as.SubscriptionService.DeactivateAllExpiredSubscriptions()
		if err != nil {
			fmt.Println(err)
		}
		print("Resetting subscription limits of all active subscriptions")
		err = as.SubscriptionService.ResetSubcriptionLimitsOfAllActiveFreeSubscriptions()
		if err != nil {
			fmt.Println(err)
		}

		//what is non-uploaded files?
		print("Deleting all non-uploaded files")
		err = as.FileService.DeleteAllNonUploadedFiles()
		if err != nil {
			fmt.Println(err)
		}
	})
	as.CronService.Start()
}
