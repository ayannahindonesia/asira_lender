package cron

import (
	"log"
)

// AutoLoanDisburseConfirm confirms loan disburse status
func AutoLoanDisburseConfirm() func() {
	return func() {
		err := DB.Table("loans").
			Where("disburse_date IS NOT NULL").
			Where("disburse_date != ?", "0001-01-01 00:00:00+00").
			Where("NOW() BETWEEN disburse_date AND disburse_date + make_interval(days => 2)").
			Update("disburse_status", "confirmed").Error

		log.Printf("AutoLoanDisburseConfirm cron executed. error : %v", err)
	}
}
