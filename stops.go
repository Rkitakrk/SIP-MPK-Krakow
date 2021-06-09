package main

// import (
// 	contxt "context"
// 	"encoding/json"
// 	"fmt"
// 	"io/ioutil"
// 	"net/http"
// )

// type stop struct {
// 	Category  string `json:"category,omitempty" bson:"category,omitempty"`
// 	Id        string `json:"id,omitempty" bson:"id,omitempty"`
// 	Latitude  int    `json:"latitude,omitempty" bson:"latitude,omitempty"`
// 	Longitude int    `json:"longitude,omitempty" bson:"longitude,omitempty"`
// 	Name      string `json:"name,omitempty" bson:"name,omitempty"`
// 	ShortName string `json:"shortName,omitempty" bson:"shortName,omitempty"`
// }

// type stops struct {
// 	Stops []stop `json:"stops,omitempty" bson:"stops,omitempty"`
// }

// // func (app *application) DownloadStopsEndpoint(w http.ResponseWriter, r *http.Request) {
// func (app *application) DownloadStopsEndpoint(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("content-type", "application/json")
// 	var s stops
// 	url_tram := "http://www.ttss.krakow.pl/internetservice/geoserviceDispatcher/services/stopinfo/stops?left=-648000000&bottom=-324000000&right=648000000&top=324000000"
// 	client := &http.Client{}
// 	req, _ := http.NewRequest("GET", url_tram, nil)
// 	res, err := client.Do(req)
// 	if err != nil {
// 		fmt.Fprintf(w, "Error: %s", err.Error())
// 	}
// 	body, err := ioutil.ReadAll(res.Body)

// 	json.Unmarshal([]byte(body), &s)

// 	for _, v := range s.Stops {
// 		fmt.Println(v.Category)
// 		fmt.Println(v.Id)
// 		fmt.Println(v.Latitude)
// 		fmt.Println(v.Longitude)
// 		fmt.Println(v.Name)
// 		fmt.Println(v.ShortName)
// 		fmt.Println("_________________")
// 	}

// 	// id, err := app.DB.Insert(&s)
// 	_, err = app.DB.Insert(&s)
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// 	// fmt.Println(id)
// 	// insertManyResult, err := app.DB.InsertMany(contxt.TODO(), trainers)
// 	fmt.Fprintf(w, string(body))

// }

// func (db *dbModel) Insert(stops *stops) (interface{}, error) {
// 	stopsv2 := []interface{}{}
// 	for _, v := range stops.Stops {
// 		stopsv2 = append(stopsv2, v)
// 	}

// 	collection := db.DB.Database("ttss").Collection("stops")
// 	// collection.Remove(bson.M{})
// 	collection.Drop(nil)
// 	insertManyResult, err := collection.InsertMany(contxt.TODO(), stopsv2)
// 	if err != nil {
// 		return 0, err
// 	}

// 	return insertManyResult.InsertedIDs, nil
// }
