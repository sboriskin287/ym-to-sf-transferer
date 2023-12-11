package ym

type TrackIdRes struct {
	Result struct {
		Library struct {
			Tracks []TrackWithAlbum
		}
	}
}

type TrackRes struct {
	Result []Track
}

type TrackWithAlbum struct {
	Id      string
	AlbumId string
}

type Track struct {
	Title   string
	Artists []Artist
	Albums  []Album
}

type Artist struct {
	Name string
}

type Album struct {
	Title string
}
