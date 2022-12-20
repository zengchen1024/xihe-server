package repository

type WuKongPictureListOption struct {
	CountPerPage int
	PageNum      int
}

type WuKongPictures struct {
	Pictures []string `json:"pictures"`
	Total    int      `json:"total"`
}

type WuKong interface {
	ListSamples(string, []int) ([]string, error)
	ListPictures(string, *WuKongPictureListOption) (WuKongPictures, error)
}
