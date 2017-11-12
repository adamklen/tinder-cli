package main

import (
    "fmt"
)

const fbToken = "EAAGm0PX4ZCpsBALr8Jy3uDtueXaZBCQutKM6wDDv3aHk4pUkmUIig91UrN1KjBoetrKI5ZAcxeWSZCuVJrJBXmbw75ZC4yl2dTxlyZCcY72O3dbRfycvQiWAqjwkjWtsmeLlIpqCEyinZB8zkdrp0KVMUamAaaAq4eUZC8uFBGQ3g32O4aI3oVABPgxPjpMLHTzeAWpY8G7QVXtDPBlwXZCeYIY97yFLzbEaQ3aZCNY1DW9RLOj8n0aGA1aWiiiS8ZCUIF5CENO2n0gibr5ajTEslnkDCTpJTEnGpoZD"

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
