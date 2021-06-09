package main

import (
	contxt "context"
	"errors"
	"fmt"
	"net/http"
	"runtime"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/context"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/mitchellh/mapstructure"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type dbModel struct {
	DB *mongo.Client
}
type application struct {
	DB *dbModel
}

//go get github.com/gorilla/context
//go get github.com/gorilla/handlers

type CustomJWTClain struct {
	Id string `json:"id"`
	jwt.StandardClaims
}

var JWT_SECRET []byte = []byte("thekitadeveloper")

func ValidateJWT(t string) (interface{}, error) {
	token, err := jwt.Parse(t, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected siging method %v", token.Header["alg"])
		}
		return JWT_SECRET, nil
	})
	if err != nil {
		return nil, errors.New(`{"message": "` + err.Error() + `"}`)
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		var tokenData CustomJWTClain
		mapstructure.Decode(claims, &tokenData)
		return tokenData, nil
	} else {
		return nil, errors.New(`{"message": "invalid token"}`)
	}
}

func ValidateMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		authorizationHeader := request.Header.Get("Authorization")
		if authorizationHeader != "" {
			bearerToken := strings.Split(authorizationHeader, " ")
			if len(bearerToken) == 2 {
				decoded, err := ValidateJWT(bearerToken[1])
				if err != nil {
					response.Header().Add("Content-Type", "application/josn")
					response.WriteHeader(500)
					response.Write([]byte(`{"message": "` + err.Error() + `"}`))
					return
				}
				//next stpe here
				context.Set(request, "decoded", decoded)
				next(response, request)
			}
		} else {

			response.Header().Add("content-type", "application/josn")
			response.WriteHeader(500)
			response.Write([]byte(`{"message": " Auth header is required"}`))
			return

		}
	})
}

func (app *application) RootEndpoint(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("Content-type", "application/json")
	response.Write([]byte(`{"message": "Hello world"}`))
}

func main() {
	fmt.Println("Starting application...")
	//Nawiązanie połączenia z bazą danych
	db, err := openDB()
	if err != nil {
		fmt.Println(err)
	}

	defer db.Disconnect(contxt.TODO())
	//Uruchomienie pobierania danych z serwera TTSS na osobnym wątku/
	runtime.GOMAXPROCS(2)
	go downloadData()

	app := &application{
		DB: &dbModel{DB: db},
	}
	// fmt.Println(authors)
	// data, _ := json.Marshal(authors)
	// fmt.Println(string(data))

	router := mux.NewRouter()
	router.HandleFunc("/", app.RootEndpoint).Methods("GET")
	// router.HandleFunc("/a", app.DownloadStopsEndpoint).Methods("GET")
	router.HandleFunc("/getdelay", app.DownloadDelayEndPoint).Methods("GET")
	router.HandleFunc("/selectdelay", app.SelectDelayEndPoint).Methods("GET")

	router.HandleFunc("/selectAllDelay", app.SelectAllDelayEndPoint).Methods("GET")
	// router.HandleFunc("/register", AuthorRegisterEndpoint).Methods("POST")
	// router.HandleFunc("/login", AuthorLoginEndpoint).Methods("POST")

	// router.HandleFunc("/authors", AuthorRetrieveAllEndpoint).Methods("GET")
	// router.HandleFunc("/author/{id}", AuthorRetrieveEndpoint).Methods("GET")
	// router.HandleFunc("/author/{id}", AuthorDeleteEndpoint).Methods("DELETE")
	// router.HandleFunc("/author/{id}", AuthorUpdateEndpoint).Methods("PUT")

	// router.HandleFunc("/article", ArticleRetrieveAllEndpoint).Methods("GET")
	// router.HandleFunc("/article/{id}", ArticleRetrieveEndpoint).Methods("GET")
	// router.HandleFunc("/article/{id}", ValidateMiddleware(ArticleDeleteEndpoint)).Methods("DELETE")
	// router.HandleFunc("/article/{id}", ValidateMiddleware(ArticleUpdateEndpoint)).Methods("PUT")
	// router.HandleFunc("/article", ValidateMiddleware(ArticleCreateEndpoint)).Methods("POST")
	methods := handlers.AllowedMethods(
		[]string{
			"GET",
			"POST",
			"PUT",
			"DELETE",
		},
	)
	headers := handlers.AllowedHeaders(
		[]string{
			"Content-Type",
			"Authorization",
			"X-Requested-With",
		},
	)
	origins := handlers.AllowedOrigins(
		[]string{
			"*",
		},
	)
	http.ListenAndServe(
		":12345",
		handlers.CORS(headers, methods, origins)(router),
	)

}

func downloadData() {
	for {
		client := &http.Client{}
		req, _ := http.NewRequest("GET", "http://localhost:12345/getdelay", nil)
		_, _ = client.Do(req)

		// _, _ = ioutil.ReadAll(res.Body)
		time.Sleep(time.Second * 60)
		fmt.Println("test")
	}
}

func openDB() (*mongo.Client, error) {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(contxt.TODO(), clientOptions)
	if err != nil {
		return nil, err
	}
	fmt.Println("Connected to MongoDB!")

	err = client.Ping(contxt.TODO(), nil)
	if err != nil {
		return nil, err
	}
	return client, nil
}

// type flights struct {
// 	Id      string `json:"id,omitempty"`
// 	Time    string `json:"time,omitempty"`
// 	Origin  string `json:"origin,omitempty"`
// 	Flight  string `json:"flight,omitempty"`
// 	Arrival string `json:"arrival,omitempty"`
// 	Remarks string `json:"remarks,omitempty"`
// }

// var data []flights = []flights{
// 	flights{
// 		Id:      "1",
// 		Time:    "16:00",
// 		Origin:  "Berlin",
// 		Flight:  "FA007",
// 		Arrival: "15:55",
// 		Remarks: "is ok",
// 	},
// 	flights{
// 		Id:      "1",
// 		Time:    "16:00",
// 		Origin:  "Berlin",
// 		Flight:  "FA007",
// 		Arrival: "15:55",
// 		Remarks: "is ok",
// 	},
// }
