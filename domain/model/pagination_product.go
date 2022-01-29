package model

type Pagination struct {
	Items       interface{} `json:"items"`
	Count       int         `json:"count"`
	Pages       int         `json:"pages"`
	ItemPerPage int         `json:"item_per_page"`
	Next        bool        `json:"next"`
	Previous    bool        `json:"previous"`
	CurrentPage int         `json:"current_page"`
	PagesList   []int       `json:"pages_list"`
	ShowPerRow  bool        `json:"show_per_row"`
	LastPage    int         `json:"show_per_row"`
	FistPage    int `json:"fist_page"`
}

func (p *Pagination) GetPages() []int {

	var pages []int
	if p.CurrentPage < 10 {
		for i := 2; i <= 10; i++ {
			pages = append(pages, i)
		}

	} else {

		for i :=p.CurrentPage - 8  ; i <= (p.CurrentPage + 1); i++ {
			pages = append(pages, i)
		}
	}

	return pages
}
