package main

import (
    "fmt"
)

const fbToken = "EAAGm0PX4ZCpsBAK3ZClKKRYNfWBZBpoVcGpMbZBlEqZAMFNN3W13K1pkafEpQboslBhvFmsRGNyuS8EGGyDJRhJUVAyOZA7NLDuGFEAX3KTSnlhPnaqhWLuSFEZCKDOncHpimZAg7FIPZCZBZBLAZB8MZAZA6NB9fxD00dEbN2GZC8yyOQJuud5bqfJfDxeGZCVhsNWtqqjO8SAcUvIcz9rvTNHtuaPWWHolYQhXhdAMox6SdscR4SIlZBa9O65WZBNnb2XIxtTIZC0k4AshZAhVtFL9lORZBRRlwUYt8raBNnaAZD"

var fbId int64 = 113137066128069

func main() {
	client := NewTinderClient(fbToken, fbId)
	if err := client.Connect(); err != nil {
		panic(err)
	}
	fmt.Println("Got client: ", client)
	recs, err := client.GetRecs()
	if err != nil { panic(err) }
	model := NewRecsModel()
	model.SetRecs(recs)
	Run(model)
}
