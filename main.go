package main

import (
	"fmt"
)

const fbToken = "EAAGm0PX4ZCpsBAP875Qq1sMSh3nLx9kGZBcMzEGd6U4tvrzyrILbuhA1O1qyqQG0bKPcihhLKxeHALZBdX1wPlO8P1OAl3su8Bf04PKUcfoEg5Dz5tNseTOKxqWOF6JXUMvBkyg6CjxwkdPiMf2EtTeBO8OL5hDcrzY4zZC4ZAxqgMHMAtYbMeZCXyxoHGFXU0ySmU2xlnNOup5V5sBDkuDIweLypCcZCjbZAjvP5QZBijGi0TTcvnB3qkdpFoIUPedjkSMUVoTMmfdWG6oWMlgt2"

var fbId int64 = 113137066128069

func main() {
	client := NewTinderClient(fbToken, fbId)
	if err := client.Connect(); err != nil {
		panic(err)
	}
	fmt.Println("Got client: ", client)
	recs, err := client.GetRecs()
	if err != nil {
		panic(err)
	}
	model := NewRecsModel(&client)
	model.SetRecs(recs)
	Run(model)
}
