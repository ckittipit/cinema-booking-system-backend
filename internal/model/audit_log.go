// structure ของ model สำหรับเก็บข้อมูล audit log ใน MongoDB โดยมีฟิลด์ต่าง ๆ ที่เกี่ยวข้องกับเหตุการณ์ที่เกิดขึ้น เช่น ประเภทของเหตุการณ์, ผู้ใช้ที่เกี่ยวข้อง, ข้อมูลเพิ่มเติมในรูปแบบของ metadata และเวลาที่เกิดเหตุการณ์นั้น ๆ
package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AuditLog struct {
	ID         primitive.ObjectID  `bson:"_id,omitempty" json:"id"`
	EventType  string              `bson:"event_type" json:"event_type"`
	UserID     *primitive.ObjectID `bson:"user_id,omitempty" json:"user_id,omitempty"`
	ShowtimeID *primitive.ObjectID `bson:"showtime_id,omitempty" json:"showtime_id,omitempty"`
	SeatID     string              `bson:seat_id,omitempty" json:"seat_id,omitempty"`
	BookingID  *primitive.ObjectID `bson:"booking_id,omitempty" json:"booking_id,omitempty"`
	Message    string              `bson:"message" json:"message"`
	Metadata   map[string]any      `bson:"metadata,omitempty" json:"metadata,omitempty"`
	CreatedAt  time.Time           `bson:"created_at" json:"created_at"`
}
