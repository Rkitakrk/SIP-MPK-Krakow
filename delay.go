package main

import (
	contxt "context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type vehicles struct {
	LastUpdate int       `json:"lastUpdate,omitempty" bson:"lastUpdate,omitempty"`
	Vehicles   []vehicle `json:"vehicles,omitempty" bson:"vehicles,omitempty"`
}

type vehicle struct {
	IsDeleted bool   `json:"isDeleted,omitempty" bson:"isDeleted,omitempty"`
	Path      []path `json:"path,omitempty" bson:"path,omitempty"`
	Color     string `json:"color,omitempty" bson:"color,omitempty"`
	Heading   int    `json:"heading,omitempty" bson:"heading,omitempty"`
	Latitude  int    `json:"latitude,omitempty" bson:"latitude,omitempty"`
	Name      string `json:"name,omitempty" bson:"name,omitempty"`
	TripId    string `json:"tripId,omitempty" bson:"tripId,omitempty"`
	Id        string `json:"id,omitempty" bson:"id,omitempty"`
	Category  string `json:"category,omitempty" bson:"category,omitempty"`
	Longitude int    `json:"longitude,omitempty" bson:"longitude,omitempty"`
	// TripPassages *tripPassages `json:"tripPassages,omitempty" bson:"tripPassages,omitempty"`
	LastStop    actual  `json:"lastStop,omitempty" bson:"lastStop,omitempty"`
	LastPlanned oldData `json:"lastPlanned,omitempty" bson:"lastPlanned,omitempty"`
	Opoznienie  int     `json:"opoznienie,omitempty" bson:"opoznienie,omitempty"`
}

type path struct {
	Y1     int     `json:"y1,omitempty" bson:"y1,omitempty"`
	Length float64 `json:"length,omitempty" bson:"length,omitempty"`
	X1     int     `json:"x1,omitempty" bson:"x1,omitempty"`
	Y2     int     `json:"y2,omitempty" bson:"y2,omitempty"`
	Angle  int     `json:"angle,omitempty" bson:"angle,omitempty"`
	X2     int     `json:"x2,omitempty" bson:"x2,omitempty"`
}

type tripPassages struct {
	Actual        []actual `json:"actual,omitempty" bson:"actual,omitempty"`
	DirectionText string   `json:"directionText,omitempty" bson:"directionText,omitempty"`
	Old           []actual `json:"old,omitempty" bson:"old,omitempty"`
	RouteName     string   `json:"routeName,omitempty" bson:"routeName,omitempty"`
}

type actual struct {
	ActualTime   string `json:"actualTime,omitempty" bson:"actualTime,omitempty"`
	Status       string `json:"status,omitempty" bson:"status,omitempty"`
	Stop         stop   `json:"stop,omitempty" bson:"stop,omitempty"`
	Stop_seq_num string `json:"stop_seq_num,omitempty" bson:"stop_seq_num,omitempty"`
}

type stop struct {
	Id        string `json:"id,omitempty" bson:"id,omitempty"`
	Name      string `json:"name,omitempty" bson:"name,omitempty"`
	ShortName string `json:"shortName,omitempty" bson:"shortName,omitempty"`
}

type stopPassages struct {
	Old []oldData `json:"old,omitempty" bson:"old,omitempty"`
}

type oldData struct {
	ActualRelativeTime int    `json:"actualRelativeTime,omitempty" bson:"actualRelativeTime,omitempty"`
	Direction          string `json:"direction,omitempty" bson:"direction,omitempty"`
	MixedTime          string `json:"mixedTime,omitempty" bson:"mixedTime,omitempty"`
	Passageid          string `json:"passageid,omitempty" bson:"passageid,omitempty"`
	PatternText        string `json:"patternText,omitempty" bson:"patternText,omitempty"`
	PlannedTime        string `json:"plannedTime,omitempty" bson:"plannedTime,omitempty"`
	RouteId            string `json:"routeId,omitempty" bson:"routeId,omitempty"`
	Status             string `json:"status,omitempty" bson:"status,omitempty"`
	TripId             string `json:"tripId,omitempty" bson:"tripId,omitempty"`
	VehicleId          string `json:"vehicleId,omitempty" bson:"vehicleId,omitempty"`
}

type delay struct {
	Name          string    `json:"name,omitempty" bson:"name,omitempty"`
	TripId        string    `json:"tripId,omitempty" bson:"tripId,omitempty"`
	Id            string    `json:"id,omitempty" bson:"id,omitempty"`
	Category      string    `json:"category,omitempty" bson:"category,omitempty"`
	StopName      string    `json:"stopName,omitempty" bson:"stopName,omitempty"`
	StopShortName string    `json:"stopShortName,omitempty" bson:"stopShortName,omitempty"`
	ActualTime    time.Time `json:"actualTime,omitempty" bson:"actualTime,omitempty"`
	NumberVehicle int       `json:"numberVehicle,omitempty" bson:"numberVehicle,omitempty"`
	Delay         int       `json:"delay,omitempty" bson:"delay,omitempty"`
	TimeInsert    time.Time `json:"timeInsert,omitempty" bson:"timeInsert,omitempty"`
}
type routeDelay struct {
	NumberVehicle int       `json:"numberVehicle,omitempty" bson:"numberVehicle,omitempty"`
	Delay         int       `json:"delay,omitempty" bson:"delay,omitempty"`
	CountVehicle  int       `json:"countVehicle,omitempty" bson:"countVehicle,omitempty"`
	Id            string    `json:"id,omitempty" bson:"id,omitempty"`
	TimeInsert    time.Time `json:"timeInsert,omitempty" bson:"timeInsert,omitempty"`
}

type lineDelay struct {
	NumberVehicle int       `json:"numberVehicle,omitempty" bson:"numberVehicle,omitempty"`
	Delay         int       `json:"delay,omitempty" bson:"delay,omitempty"`
	CountVehicle  int       `json:"countVehicle,omitempty" bson:"countVehicle,omitempty"`
	TimeInsert    time.Time `json:"timeInsert,omitempty" bson:"timeInsert,omitempty"`
}

func (app *application) DownloadDelayEndPoint(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")

	v := new(vehicles)
	v_out := new(vehicles)
	var delays []delay

	url_tram := "http://www.ttss.krakow.pl/internetservice/geoserviceDispatcher/services/vehicleinfo/vehicles?positionType=CORRECTED"
	url_tripPassages := "http://www.ttss.krakow.pl/internetservice/services/tripInfo/tripPassages?"
	url_stop := "http://www.ttss.krakow.pl/internetservice/services/passageInfo/stopPassages/stop?"

	body, err := pushUrl(url_tram)
	if err != nil {
		fmt.Fprintf(w, "Error: %s", err.Error())
	}

	json.Unmarshal([]byte(body), &v)

	for _, v2 := range v.Vehicles {

		url_trip := url_tripPassages + "mode=departure&tripId=" + v2.TripId

		if v2.TripId != "" {
			bodytrip, err := pushUrl(url_trip)
			if err != nil {
				fmt.Fprintf(w, "Error: %s", err.Error())
			}
			var trip tripPassages
			json.Unmarshal([]byte(bodytrip), &trip)
			old := trip.Old
			var last actual

			if len(old) > 0 {
				last = old[len(old)-1]

				v2.LastStop = last
				url_stop_key := url_stop + "mode=departure&stop=" + v2.LastStop.Stop.ShortName
				bodystop, err := pushUrl(url_stop_key)
				if err != nil {
					fmt.Fprintf(w, "Error: %s", err.Error())
				}
				var stopPassage stopPassages
				json.Unmarshal([]byte(bodystop), &stopPassage)

				for _, planned := range stopPassage.Old {
					if v2.TripId == planned.TripId {

						v2.LastPlanned = planned
						if planned.PlannedTime != "" && v2.LastStop.ActualTime != "" {
							plannedHour := convertToInt(w, planned.PlannedTime[:2])
							plannedMinute := convertToInt(w, planned.PlannedTime[3:])
							actualHour := convertToInt(w, v2.LastStop.ActualTime[:2])
							actualMinute := convertToInt(w, v2.LastStop.ActualTime[3:])
							
							a := plannedHour*60 + plannedMinute - (actualHour*60 + actualMinute)
							if a > 12*60 || a < -12*60 {
								b := 24*60 - (plannedHour*60 + plannedMinute) + (actualHour*60 + actualMinute)
								v2.Opoznienie = b
							} else {
								v2.Opoznienie = a
							}
							// fmt.Println(a)
							var d delay
							t := time.Now()
							d.Name = v2.Name
							d.TripId = v2.TripId
							d.Id = v2.Id
							d.Category = v2.Category
							d.StopName = v2.LastStop.Stop.Name
							d.StopShortName = v2.LastStop.Stop.ShortName
							d.ActualTime = time.Date(t.Year(), t.Month(), t.Day(), actualHour, actualMinute, 0, 0, time.UTC)
							d.NumberVehicle = convertToInt(w, v2.LastPlanned.PatternText)
							d.Delay = v2.Opoznienie
							

							fmt.Println(d)

							delays = append(delays, d)

						}

					}
				}
			}
		}
		v_out.Vehicles = append(v_out.Vehicles, v2)
	}

	vv, _ := json.Marshal(delays)
	_, err = app.DB.Insert(delays)
	err = app.DB.updateData()
	if err != nil {
		fmt.Println(err)
	}

	fmt.Fprintf(w, string(vv))

}

func (app *application) SelectDelayEndPoint(w http.ResponseWriter, r *http.Request) {

	data, err := app.DB.Select()
	if err != nil {
		fmt.Println(err)
	}

	bodyBytes, _ := json.Marshal(data)
	fmt.Fprintf(w, string(bodyBytes))

}

func (app *application) SelectAllDelayEndPoint(w http.ResponseWriter, r *http.Request) {

	data, err := app.DB.SelectLineAndRoute()
	if err != nil {
		fmt.Println(err)
	}


	bodyBytes, _ := json.Marshal(data)
	fmt.Fprintf(w, string(bodyBytes))
}

func (db *dbModel) Select() (interface{}, error) {
	collection := db.DB.Database("ttss").Collection("delays")
	t := time.Now()
	timee := time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute()-30, 0, 0, time.UTC)
	// selectManyResult, err := collection.Find(contxt.TODO(), bson.M{"actualTime": { $gte: timee})
	selectManyResult, err := collection.Find(contxt.TODO(), bson.M{"actualTime": bson.M{"$gte": timee}})
	if err != nil {
		return 0, err
	}
	var data []bson.M
	if err = selectManyResult.All(contxt.TODO(), &data); err != nil {
		return 0, err
	}
	return data, err
}

func queryToMongo(coll *mongo.Collection) (interface{}, error) {
	selectManyResult, err := coll.Find(contxt.TODO(), bson.M{})
	if err != nil {
		return 0, err
	}
	var data []bson.M
	if err = selectManyResult.All(contxt.TODO(), &data); err != nil {
		return 0, err
	}
	return data, err

}

func (db *dbModel) SelectLineAndRoute() (interface{}, error) {
	collLineDelayOne := db.DB.Database("ttss").Collection("lineDelayOne")
	collLineDelayFive := db.DB.Database("ttss").Collection("lineDelayFive")
	collLineDelayFifteen := db.DB.Database("ttss").Collection("lineDelayFifteen")
	collLineDelayOneHour := db.DB.Database("ttss").Collection("lineDelayOneHour")
	collLineDelaySixHour := db.DB.Database("ttss").Collection("lineDelaySixHour")

	collRouteDelayOne := db.DB.Database("ttss").Collection("routeDelayOne")
	collRouteDelayFive := db.DB.Database("ttss").Collection("routeDelayFive")
	collRouteDelayFifteen := db.DB.Database("ttss").Collection("routeDelayFifteen")
	collRouteDelayOneHour := db.DB.Database("ttss").Collection("routeDelayOneHour")
	collRouteDelaySixHour := db.DB.Database("ttss").Collection("routeDelaySixHour")

	dataLineDelayOne, err := queryToMongo(collLineDelayOne)
	dataLineDelayFive, err := queryToMongo(collLineDelayFive)
	dataLineDelayFifteen, err := queryToMongo(collLineDelayFifteen)
	dataLineDelayOneHour, err := queryToMongo(collLineDelayOneHour)
	dataLineDelaySixHour, err := queryToMongo(collLineDelaySixHour)

	dataRouteDelayOne, err := queryToMongo(collRouteDelayOne)
	dataRouteDelayFive, err := queryToMongo(collRouteDelayFive)
	dataRouteDelayFifteen, err := queryToMongo(collRouteDelayFifteen)
	dataRouteDelayOneHour, err := queryToMongo(collRouteDelayOneHour)
	dataRouteDelaySixHour, err := queryToMongo(collRouteDelaySixHour)
	if err != nil {
		return 0, err
	}
	dataAll := map[string]interface{}{
		"LineDelayOne":     dataLineDelayOne,
		"LineDelayFive":    dataLineDelayFive,
		"LineDelayFifteen": dataLineDelayFifteen,
		"LineDelayOneHour": dataLineDelayOneHour,
		"LineDelaySixHour": dataLineDelaySixHour,

		"RouteDelayOne":     dataRouteDelayOne,
		"RouteDelayFive":    dataRouteDelayFive,
		"RouteDelayFifteen": dataRouteDelayFifteen,
		"RouteDelayOneHour": dataRouteDelayOneHour,
		"RouteDelaySixHour": dataRouteDelaySixHour,
	}
	/*
	var plik = daneZpliku
	selectWithForLine(plik.LineDelayOne)
	...
	selectWithForLine(plik.LineDelaySixHour)

	selectWithForRoute(plik.RouteDelayOne)
	...
	selectWithForRoute(plik.RouteDelaySixHour)


	func selectWithForLine(data intefrace{}){
		for _, v := range data{
			fmt.Println(v.numberVehicle)
			fmt.Println(v.delay)
		}
	}
	func selectWithForRoute(data intefrace{}){
		for _, v := range data{
			fmt.Println(v.numberVehicle)
			fmt.Println(v.id)
			fmt.Println(v.delay)
		}
	}
	
	
	
	*/
	
	return dataAll, nil
}

//Główna funkcja insertowa danych do bazdy danych
func (db *dbModel) Insert(delays []delay) (interface{}, error) {
	stopsv2 := []interface{}{}
	listRoute := []interface{}{}
	listTrams := []interface{}{}
	var route routeDelay
	Trams := removeDuplicateValues(delays)
	// fmt.Println(Trams)
	t := time.Now()
	now := time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), 0, 0, time.UTC)

	for _, v := range delays {
		route.Id = v.Id
		route.Delay = v.Delay
		route.NumberVehicle = v.NumberVehicle
		route.TimeInsert = now

		v.TimeInsert = now
		stopsv2 = append(stopsv2, v)
		listRoute = append(listRoute, route)
	}
	//Pętla zlicza wszystkie opóźnienia na liniach i je podsuwuje
	for _, vv := range Trams {
		for _, v := range delays {
			if v.NumberVehicle == vv.NumberVehicle {
				vv.CountVehicle = vv.CountVehicle + 1
				vv.Delay = vv.Delay + v.Delay
			}
		}
		vv.Delay = vv.Delay / vv.CountVehicle
		vv.TimeInsert = now
		listTrams = append(listTrams, vv)
	}

	// fmt.Println(listTrams...)

	// collection := db.DB.Database("ttss").Collection("delays")
	//Nawiązanie połączenia z tabelami w bazie ttss
	collectionRoute := db.DB.Database("ttss").Collection("routeDelay")
	collectionLine := db.DB.Database("ttss").Collection("lineDelay")
	//Insertowanie dancyh do serwera
	err := insertDelay(collectionRoute, listRoute)
	err = insertDelay(collectionLine, listTrams)

	//Sprawdzanie błędów po insercie
	if err != nil {
		return nil, err
	}

	// insertManyResult, err := collection.InsertMany(contxt.TODO(), stopsv2)
	// if err != nil {
	// 	return 0, err
	// }

	return "", nil
}

func insertDelay(coll *mongo.Collection, tab []interface{}) error {
	t := time.Now()
	//Usuwanie wszystkich rekordów starszych o 6 godzin
	timee := time.Date(t.Year(), t.Month(), t.Day(), t.Hour()-6, t.Minute(), 0, 0, time.UTC)
	_, err := coll.DeleteMany(contxt.TODO(), bson.M{"timeInsert": bson.M{"$lt": timee}})
	//Wrzucanie nowych rekordów do bazy danych
	_, err = coll.InsertMany(contxt.TODO(), tab)

	if err != nil {
		return err
	}
	return nil
}
func (db *dbModel) updateData() error {
	var lineDelayAll []lineDelay
	var lineDelayOne []lineDelay
	var lineDelayFive []lineDelay
	var lineDelayFifteen []lineDelay
	var lineDelayOneHour []lineDelay
	var lineDelaySixHour []lineDelay
	var lineDelayOneRemove []interface{}
	var lineDelayFiveRemove []interface{}
	var lineDelayFifteenRemove []interface{}
	var lineDelayOneHourRemove []interface{}
	var lineDelaySixHourRemove []interface{}

	var routeDelayAll []routeDelay
	var routeDelayOne []routeDelay
	var routeDelayFive []routeDelay
	var routeDelayFifteen []routeDelay
	var routeDelayOneHour []routeDelay
	var routeDelaySixHour []routeDelay
	var routeDelayOneRemove []interface{}
	var routeDelayFiveRemove []interface{}
	var routeDelayFifteenRemove []interface{}
	var routeDelayOneHourRemove []interface{}
	var routeDelaySixHourRemove []interface{}
	//Tworzenie aktualnej daty
	t := time.Now()
	timee := time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), 0, 0, time.UTC)
	//Nawiązanie połącznia z kolekcjami z mango
	collLineDelay := db.DB.Database("ttss").Collection("lineDelay")
	collRouteDelay := db.DB.Database("ttss").Collection("routeDelay")
	selectLine, err := collLineDelay.Find(contxt.TODO(), bson.M{})
	selectRoute, err := collRouteDelay.Find(contxt.TODO(), bson.M{})
	//Wrzycenie pobranych dancyh do struktur
	if err = selectLine.All(contxt.TODO(), &lineDelayAll); err != nil {
		return err
	}
	if err = selectRoute.All(contxt.TODO(), &routeDelayAll); err != nil {
		return err
	}
	//Podział danych na strefy czasowe 
	for _, v := range lineDelayAll {
		a := timee.Sub(v.TimeInsert)

		if a <= time.Minute*1 {
			lineDelayOne = append(lineDelayOne, v)
		}
		if a <= time.Minute*5 {
			lineDelayFive = append(lineDelayFive, v)
		}
		if a <= time.Minute*15 {
			lineDelayFifteen = append(lineDelayFifteen, v)
		}
		if a <= time.Minute*60 {
			lineDelayOneHour = append(lineDelayOneHour, v)
		}
		if a <= time.Minute*60*6 {
			lineDelaySixHour = append(lineDelaySixHour, v)
		}
	}

	for _, v := range routeDelayAll {
		a := timee.Sub(v.TimeInsert)

		if a <= time.Minute*1 {
			routeDelayOne = append(routeDelayOne, v)
		}
		if a <= time.Minute*5 {
			routeDelayFive = append(routeDelayFive, v)
		}
		if a <= time.Minute*15 {
			routeDelayFifteen = append(routeDelayFifteen, v)
		}
		if a <= time.Minute*60 {
			routeDelayOneHour = append(routeDelayOneHour, v)
		}
		if a <= time.Minute*60*6 {
			routeDelaySixHour = append(routeDelaySixHour, v)
		}
	}
	//Usuwanie duklikatów, czyli usunięcie powtórzeń tramwajów
	lineDelayOneRemove = removeDuplicateValuesLineDelay(lineDelayOne)
	lineDelayFiveRemove = removeDuplicateValuesLineDelay(lineDelayFive)
	lineDelayFifteenRemove = removeDuplicateValuesLineDelay(lineDelayFifteen)
	lineDelayOneHourRemove = removeDuplicateValuesLineDelay(lineDelayOneHour)
	lineDelaySixHourRemove = removeDuplicateValuesLineDelay(lineDelaySixHour)

	routeDelayOneRemove = removeDuplicateValuesRouteDelay(routeDelayOne)
	routeDelayFiveRemove = removeDuplicateValuesRouteDelay(routeDelayFive)
	routeDelayFifteenRemove = removeDuplicateValuesRouteDelay(routeDelayFifteen)
	routeDelayOneHourRemove = removeDuplicateValuesRouteDelay(routeDelayOneHour)
	routeDelaySixHourRemove = removeDuplicateValuesRouteDelay(routeDelaySixHour)

	collLineDelayOne := db.DB.Database("ttss").Collection("lineDelayOne")
	collLineDelayFive := db.DB.Database("ttss").Collection("lineDelayFive")
	collLineDelayFifteen := db.DB.Database("ttss").Collection("lineDelayFifteen")
	collLineDelayOneHour := db.DB.Database("ttss").Collection("lineDelayOneHour")
	collLineDelaySixHour := db.DB.Database("ttss").Collection("lineDelaySixHour")

	collRouteDelayOne := db.DB.Database("ttss").Collection("routeDelayOne")
	collRouteDelayFive := db.DB.Database("ttss").Collection("routeDelayFive")
	collRouteDelayFifteen := db.DB.Database("ttss").Collection("routeDelayFifteen")
	collRouteDelayOneHour := db.DB.Database("ttss").Collection("routeDelayOneHour")
	collRouteDelaySixHour := db.DB.Database("ttss").Collection("routeDelaySixHour")

	//Insertowanie aktualnych danych do bazdy danych
	insertDelayUpdate(collLineDelayOne, lineDelayOneRemove)
	insertDelayUpdate(collLineDelayFive, lineDelayFiveRemove)
	insertDelayUpdate(collLineDelayFifteen, lineDelayFifteenRemove)
	insertDelayUpdate(collLineDelayOneHour, lineDelayOneHourRemove)
	insertDelayUpdate(collLineDelaySixHour, lineDelaySixHourRemove)

	insertDelayUpdate(collRouteDelayOne, routeDelayOneRemove)
	insertDelayUpdate(collRouteDelayFive, routeDelayFiveRemove)
	insertDelayUpdate(collRouteDelayFifteen, routeDelayFifteenRemove)
	insertDelayUpdate(collRouteDelayOneHour, routeDelayOneHourRemove)
	insertDelayUpdate(collRouteDelaySixHour, routeDelaySixHourRemove)

	// fmt.Println("******************")
	// fmt.Println(lineDelayOneRemove)
	// fmt.Println("******************")
	// fmt.Println(lineDelayFiveRemove)
	// fmt.Println("******************")
	// fmt.Println(lineDelayFifteenRemove)
	// fmt.Println("******************")
	// fmt.Println(lineDelayOneHourRemove)
	// fmt.Println("******************")
	// fmt.Println(lineDelaySixHourRemove)

	return nil
}

//Insertowanie do bazy dancyh, najpierw usunięcie wszystkich rekordów w kolekcjach, a następnie nadpisanie nowymi rekordami
func insertDelayUpdate(coll *mongo.Collection, d []interface{}) {
	//Usuwanie wszystkich rekordów z kolekcji
	_, err := coll.DeleteMany(contxt.TODO(), bson.M{})
	//Wrzucanie nowych rekordów do bazy danych
	_, err = coll.InsertMany(contxt.TODO(), d)
	//Zwracanie błędu
	if err != nil {

	}

}

//Remove duplicate values Line Delay
//*********************************
func removeDuplicateValuesLineDelay(d []lineDelay) []interface{} {
	var dOut []interface{}
	var dd []lineDelay
	var line lineDelay
	keys := make(map[int]bool)
	list := []int{}

	for _, dd := range d {
		if _, value := keys[dd.NumberVehicle]; !value {
			keys[dd.NumberVehicle] = true
			list = append(list, dd.NumberVehicle)
		}
	}
	sort.Ints(list)
	for _, v := range list {
		line.NumberVehicle = v
		dd = append(dd, line)
	}

	for _, vv := range dd {
		for _, v := range d {
			if v.NumberVehicle == vv.NumberVehicle {
				vv.CountVehicle = vv.CountVehicle + 1
				vv.Delay = vv.Delay + v.Delay
			}
		}
		vv.Delay = vv.Delay / vv.CountVehicle
		t := time.Now()
		now := time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), 0, 0, time.UTC)
		vv.TimeInsert = now
		dOut = append(dOut, vv)
	}

	return dOut
}

//Remove duplicate values Route Delay
//*********************************
func removeDuplicateValuesRouteDelay(d []routeDelay) []interface{} {
	var dOut []interface{}
	var dd []routeDelay
	var route routeDelay
	keys := make(map[string]bool)
	list := []string{}

	for _, dd := range d {
		if _, value := keys[dd.Id]; !value {
			keys[dd.Id] = true
			list = append(list, dd.Id)
		}
	}

	for _, v := range list {
		route.Id = v
		dd = append(dd, route)
	}

	for _, vv := range dd {
		for _, v := range d {
			if v.Id == vv.Id {
				vv.CountVehicle = vv.CountVehicle + 1
				vv.Delay = vv.Delay + v.Delay
				vv.NumberVehicle = v.NumberVehicle
			}
		}
		vv.Delay = vv.Delay / vv.CountVehicle
		t := time.Now()
		now := time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), 0, 0, time.UTC)
		vv.TimeInsert = now
		dOut = append(dOut, vv)
	}

	return dOut
}

//Funkacja separuje linie tramwajowe
func removeDuplicateValues(trams []delay) []lineDelay {
	keys := make(map[int]bool)
	var listTrams []lineDelay
	var Tram lineDelay
	list := []int{}

	//funkcja słuzy do wyciągnięcia wszystkich aktualnie lini tramwajowych.
	for _, tram := range trams {
		if _, value := keys[tram.NumberVehicle]; !value {
			keys[tram.NumberVehicle] = true
			list = append(list, tram.NumberVehicle)
		}
	}
	// Sortowanie lini tramwajowych
	sort.Ints(list)
	for _, v := range list {
		Tram.NumberVehicle = v
		listTrams = append(listTrams, Tram)
	}
	return listTrams
}

//Funkcja, która pobiera dane z TTSS
func pushUrl(s string) ([]byte, error) {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", s, nil)
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(res.Body)
	return body, err
}

//Funkacja, która konwertuje format danych
func convertToInt(w http.ResponseWriter, s string) int {
	integer, err := strconv.Atoi(s)
	if err != nil {
		fmt.Fprintf(w, "Error: %s", err.Error())
	}
	return integer
}
