package listing

type Preview struct {
	Id        uint   `json:"id"`
	Name      string `json:"name"`
	Loved     bool   `json:"loved"`
	Thumbnail string `json:"thumbnail"`
}
