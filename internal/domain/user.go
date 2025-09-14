package domain

import (
    "time"
    "go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
    ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
    Name      string             `json:"name" bson:"name"`
    Email     string             `json:"email" bson:"email"`
    Password  string             `json:"-" bson:"password"` // ไม่ส่งออก password เวลา JSON
    CreatedAt time.Time          `json:"createdAt" bson:"createdAt"`
}
