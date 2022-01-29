package nodes

import "go.mongodb.org/mongo-driver/bson/primitive"

type Node struct {
	ID        primitive.ObjectID `bson:"_id"`
	SrcPort string `json:"src_port"`
	DstPort string `json:"dst_port"`
	SrcIp string `json:"src_ip"`
	DstIp string `json:"dst_ip"`
	SrcMac string `json:"src_mac"`
	DstMac string `json:"dst_mac"`
	Protocol string `json:"protocol"`
	Timestamp int64 `json:"timestamp"`
}
