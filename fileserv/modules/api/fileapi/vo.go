package fileapi

type ProductVO struct {
	Id      int64  `json:"id"`
	Name    string `json:"name"`
	Creator string `json:"creator"`
	Created string `json:"created"`
}
