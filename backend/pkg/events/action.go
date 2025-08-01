package events

import (
	"context"
	"encoding/json"
	"fmt"
	"gorm.io/gorm"
	"log"
	"time"
)

func LogEvent(ctx context.Context, db *gorm.DB, modelName string, modelID string, obj interface{}, eventAction EventAction) {
	go func() {
		var (
			data []byte
			err  error
		)

		userUID, _ := ctx.Value("userUID").(string)
		if userUID == "" {
			log.Println(fmt.Errorf("user UID is empty"))
			return
		}

		if eventAction.Action != EventDelete || obj != nil {
			data, err = json.Marshal(obj)
			if err != nil {
				fmt.Printf("events marshal error: %v\n", err)
				return
			}
		}

		action, _ := json.Marshal(eventAction)

		eventLog := Event{
			UserUID:   userUID,
			Model:     modelName,
			ModelID:   modelID,
			Action:    string(action),
			Changes:   string(data),
			CreatedAt: time.Now(),
		}

		if err = db.Create(&eventLog).Error; err != nil {
			fmt.Printf("events insert error: %v\n", err)
		}
	}()
}
