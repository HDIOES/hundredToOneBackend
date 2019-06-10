package game

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
)

//CreateSearchGamesHandler function
func CreateSearchGamesHandler(db *sql.DB) http.Handler {
	searchGamesHandler := &SearchGamesHandler{Db: db}
	return searchGamesHandler
}

//SearchGamesHandler struct
type SearchGamesHandler struct {
	Db *sql.DB
}

func (sgh *SearchGamesHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	/*if vars, parseErr := url.ParseQuery(r.URL.RawQuery); parseErr == nil {
		//limit, limitOk := vars["limit"]
		//offset, offsetOk := vars["offset"]
		//logic of searching games
	} else {
		log.Println(parseErr)
	}*/
}

//CreateCreateGameHandler function
func CreateCreateGameHandler(db *sql.DB) http.Handler {
	createGameHandler := &CreateGameHandler{Db: db}
	return createGameHandler
}

//CreateGameHandler struct
type CreateGameHandler struct {
	Db *sql.DB
}

func (cgh *CreateGameHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if rc, err := r.GetBody(); err == nil {
		var game *Game
		if err := json.NewDecoder(rc).Decode(game); err == nil {
			//database logic of saving game
		} else {
			log.Println(err)
		}
	} else {
		log.Println(err)
	}
}

//Game struct represent rest object for game entity
type Game struct {
	ID               int64
	Desc             string
	FirstQuestion    Question
	DoubleQuestion   Question
	InversedQuestion Question
}

//Question struct represent rest object for question entity
type Question struct {
	Text    string
	answers []Answer
}

//Answer struct represent rest object for answer entity
type Answer struct {
	Text  string
	Score int32
}
