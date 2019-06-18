package main

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"
	"time"

	game "github.com/HDIOES/hundredToOneBackend/rest/games"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/tkanos/gonfig"

	migrate "github.com/rubenv/sql-migrate"
)

type Configuration struct {
	DatabaseUrl        string `json:"databaseUrl"`
	MaxOpenConnections int    `json:"maxOpenConnections"`
	MaxIdleConnections int    `json:"maxIdleConnections"`
	ConnectionTimeout  int    `json:"connectionTimeout"`
	Port               int    `json:"port"`
}

func main() {
	configuration := Configuration{}
	gonfigErr := gonfig.GetConf("dbconfig.json", &configuration)
	if gonfigErr != nil {
		panic(gonfigErr)
	}

	db, err := sql.Open("postgres", configuration.DatabaseUrl)
	if err != nil {
		panic(err)
	}
	//Parody of circuit breaker
	if pingErr := db.Ping(); pingErr == nil {
		log.Printf("Database available!!!")
	} else {
		log.Print(pingErr)
		panic(pingErr)
	}

	db.SetMaxIdleConns(configuration.MaxIdleConnections)
	db.SetMaxOpenConns(configuration.MaxOpenConnections)
	timeout := strconv.Itoa(configuration.ConnectionTimeout) + "s"
	timeoutDuration, durationErr := time.ParseDuration(timeout)
	if durationErr != nil {
		log.Println("Error parsing of timeout parameter")
		panic(durationErr)
	} else {
		db.SetConnMaxLifetime(timeoutDuration)
	}

	log.Println("Configuration has been loaded")

	migrations := &migrate.FileMigrationSource{
		Dir: "migrations",
	}

	if n, err := migrate.Exec(db, "postgres", migrations, migrate.Up); err == nil {
		log.Printf("Applied %d migrations!\n", n)
	} else {
		panic(err)
	}

	router := mux.NewRouter()

	router.Handle("/games", game.CreateSearchGamesHandler(db)).
		Methods("GET")
	router.Handle("/game/{id}", game.CreateGetGameHandler(db)).
		Methods("GET")
	router.Handle("/game", game.CreateCreateGameHandler(db)).
		Methods("POST")
	router.Handle("/game/{id}", game.CreateDeleteGameHandler(db)).
		Methods("DELETE")
	router.Handle("/game/{id}", game.CreateUpdateGameHandler(db)).
		Methods("PUT")

	http.Handle("/", router)

	listenandserveErr := http.ListenAndServe(":"+strconv.Itoa(configuration.Port), corsHandler(router))
	if listenandserveErr != nil {
		panic(err)
	}

}

func corsHandler(h http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "OPTIONS" {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, HEAD, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		} else {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			h.ServeHTTP(w, r)
		}
	}
}
