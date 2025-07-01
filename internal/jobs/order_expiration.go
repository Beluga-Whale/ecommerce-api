package jobs

import (
	"log"
	"time"

	"github.com/Beluga-Whale/ecommerce-api/internal/models"
	"github.com/Beluga-Whale/ecommerce-api/internal/services"
	"github.com/robfig/cron"
	"gorm.io/gorm"
)

func StartOrderExpirationJob(db *gorm.DB, orderService *services.OrderService) {
	c := cron.New()

	c.AddFunc("@every 15m", func() {
		// NOTE -เวลาปัจจุบัน
		now := time.Now()

		var orders []models.Order
		// NOTE - หา order จาก สถานะเป้น peding และ payment_expire_at
		if err := db.Model(&models.Order{}).
			Where("status = ? AND payment_expire_at <= ?", models.Pending, now).
			Find(&orders).Error; err != nil {
			log.Printf("Error finding expired orders",err)
			return
		}

		for _, order := range orders{
			err := orderService.CancelOrderAndRestoreStock(order.ID)
			if err != nil{
				log.Printf("Failed to cancel order %d: %v\n",order.ID, err)
			}else {
				// NOTE - cancel สำเสร็จและทำการคืน stock เรียบบร้อย
				log.Printf("Cancelled order %d due to timeout\n", order.ID)
			}
		}
	})

	c.Start()
	log.Printf("Order expiration cron jon started")
}

