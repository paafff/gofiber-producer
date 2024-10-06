package models

import "time"

type Activity struct {
	UserID    uint      `bson:"user_id"`
	Email     string    `bson:"email"`
	LoginTime time.Time `bson:"login_time"`
}
