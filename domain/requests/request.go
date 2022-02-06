package requests

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Request struct {
	ID      primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	SrcAddr string             `bson:"src_addr" json:"src_addr"`
	SrcPort int                `bson:"src_port" json:"src_port"`
	DstAddr string             `bson:"dst_addr" json:"dst_addr"`
	DstPort int                `bson:"dst_port" json:"dst_port"`

	Method        string    `bson:"method" json:"method"`
	Ts            time.Time `bson:"ts" json:"ts"`
	StatusCode    int       `bson:"status_code" json:"status_code"`
	ContentLength int       `bson:"content_length" json:"content_length"`
	Url           string    `bson:"url" json:"url"`
	UserAgent     string    `bson:"user_agent" json:"user_agent"`
	Body          []byte    `bson:"body" json:"body"`
	Response      []byte    `bson:"response" json:"response"`
}
