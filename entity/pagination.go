package entity

type Pagination struct {
	PerPage     int   `json:"per_page"`
	Page        int   `json:"page"`
	MaxPage     int   `json:"max_page"`
	TotalData   int64 `json:"total_data"`
	DataPerPage any   `json:"data_per_page"`
}

func (p *Pagination) GetOffset() int {
	return (p.GetPage() - 1) * p.GetPerPage()
}

func (p *Pagination) GetPerPage() int {
	if p.PerPage <= 0 {
		p.PerPage = 10
	}
	return p.PerPage
}

func (p *Pagination) GetPage() int {
	if p.Page <= 0 {
		p.Page = 1
	}
	return p.Page
}