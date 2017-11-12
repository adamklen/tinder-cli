package main

import (
    "fmt"
)

const fbToken = "EAAGm0PX4ZCpsBAJRBWg5I7GNI3garFATRLssFch0gmhDthCzB0QzIp4spgYG39YbakZAQsHUPQrqRoHrVYrIOTrt9GvdTybl2NLlhLBjtjqjQp84mndubWSKHVtqbZAHNa9LTkz22nSMAz9e78GJZCZCqEVwawdMBZAyZCcD3O6S5r6hxZAQSia7lf9JZBbr0JtfoLQZCFsBvLfmVcAk20aDvo5Om2vinivcCqEOam7Nneaq1cMDq3fg5WaYT2FkYKZCGkkKwgrMb2nRiNOnochyEa6WGz8jDQJuv8ZD"

var fbId int64 = 113137066128069

func main() {
	client := NewTinderClient(fbToken, fbId)
	if err := client.Connect(); err != nil {
		panic(err)
	}
	fmt.Println("Got client: ", client)
	recs, err := client.GetRecs()
	if err != nil { panic(err) }
	model := NewRecsModel(&client)
	model.SetRecs(recs)
	Run(model)
}
