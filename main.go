package main

import (
    "fmt"
)

const fbToken = "EAAGm0PX4ZCpsBAI5jZBpxDVp2FOID2c5u4gjZCOQ0CvF2aUlW0Eg4mXvlX4Vyx2k5J1XuqZCrFoWloNZBfuMQikStMEPU38wY7jKWpK8D8cgZBwE6Hf6b8qWe2DKb0vOKez7pZAQZBYZC3tRl3mXNJIQyvTiG4VncZCXSUJdwQuUSGRaarZB8rjZBpuKoBT7cSwlHwE9PR0rc59rEbiL6DpAGVFL0KKIrA0PIYVys9vZCfyCO1HfKNqZCgFWZA5gD5xPYLKAwhlWP1cQG9zZAPZB1RNMXmHA2sGJLIFWcvM4ZD"

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
