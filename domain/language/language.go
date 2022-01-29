package language

type Language struct {
	Id           int    `json:"id"`
	Language     string `json:"language"`
	Code         string `json:"code"`
	LanguageCode string `json:"language_code"`
	IsRtl        int8   `json:"is_rtl"`
	IconName     string `json:"icon_name"`
}
