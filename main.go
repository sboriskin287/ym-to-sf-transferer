package main

import (
	"fmt"
	"log"
	"ym-to-spotify-transferer/sf"
	"ym-to-spotify-transferer/ym"
)

func main() {
	ymTracks, err := ym.GetMyTracks()
	if err != nil {
		log.Fatal(err)
	}
	if err = sf.Like(ymTracks); err != nil {
		fmt.Println(err)
	}
}
