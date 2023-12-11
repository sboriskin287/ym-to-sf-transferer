package sf

type SearchRes struct {
	Tracks Tracks
}

type Tracks struct {
	Items []Track
}

type Track struct {
	Id string
}

type LikeTracks struct {
	Ids []string `json:"ids"`
}
