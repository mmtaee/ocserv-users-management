package audit_log

import (
	"context"
	"encoding/json"
	"fmt"
	"gorm.io/gorm"
	"log"
	"time"
)

func AuditLogHandler(ctx context.Context, db *gorm.DB, modelName string, modelID string, obj interface{}, eventAction AuditLogAction) {
	userUID, _ := ctx.Value("userUID").(string)
	username, _ := ctx.Value("username").(string)
	if userUID == "" {
		log.Println(fmt.Errorf("user UID is empty"))
		return
	}

	go func() {
		c, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()

		var (
			data []byte
			err  error
		)

		if eventAction.Action != EventDelete || obj != nil {
			data, err = json.Marshal(obj)
			if err != nil {
				fmt.Printf("events marshal error: %v\n", err)
				return
			}
		}

		action, _ := json.Marshal(eventAction)

		eventLog := AuditLog{
			UserUID:   userUID,
			Username:  username,
			Model:     modelName,
			ModelID:   modelID,
			Action:    string(action),
			Changes:   string(data),
			CreatedAt: time.Now(),
		}

		if err = db.WithContext(c).Create(&eventLog).Error; err != nil {
			fmt.Printf("events insert error: %v\n", err)
		}
	}()
}
