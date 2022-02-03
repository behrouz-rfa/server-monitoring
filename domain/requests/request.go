package requests

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Request struct {
	ID      primitive.ObjectID `bson:"_id,omitempty"`
	SrcAddr string             `bson:"src_addr"`
	SrcPort int                `bson:"src_port"`
	DstAddr string             `bson:"dst_addr"`
	DstPort int                `bson:"dst_port"`

	Method        string    `bson:"method"`
	Ts            time.Time `bson:"ts"`
	StatusCode    int       `bson:"status_code"`
	ContentLength int       `bson:"content_length"`
	Url           string    `bson:"url"`
	UserAgent     string    `bson:"user_agent"`
	Body          []byte    `bson:"body"`
	Response      []byte    `bson:"response"`
}
