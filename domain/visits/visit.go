package visits

type Visit struct {
	Url       string `json:"url"`
	Ip        string `json:"ip"`
	UserId    int64  `json:"user_id"`
	CreatedAt string `json:"created_at"`
}
