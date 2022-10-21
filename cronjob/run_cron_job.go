package cronjob

import (
	"fmt"
	"gorm.io/gorm"
	"os"
	"seadeals-backend/service"
)

type RunCronJobHelper interface {
	RunCronJobs()
}

type runCronJobHelper struct {
	db           *gorm.DB
	orderService service.OrderService
}

type RunCronJobConfig struct {
	DB           *gorm.DB
	OrderService service.OrderService
}

func NewCronJob(c *RunCronJobConfig) RunCronJobHelper {
	if os.Getenv("ENV") == "testing" {
		fmt.Println("disable cron")
		return nil
	}

	return &runCronJobHelper{
		db:           c.DB,
		orderService: c.OrderService,
	}
}

func (r *runCronJobHelper) RunCronJobs() {
	if os.Getenv("ENV") == "testing" {
		fmt.Println("disable cron")
		return
	}
	fmt.Println("in handler")
	r.orderService.RunCronJobs()

}
