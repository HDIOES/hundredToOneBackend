package game

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"strconv"
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
	if vars, parseErr := url.ParseQuery(r.URL.RawQuery); parseErr == nil {
		limit, limitOk := vars["limit"]
		offset, offsetOk := vars["offset"]
		sqlQuery := "SELECT ID, BODY FROM GAMES WHERE 1 = 1"
		countOfParameters := 1
		args := make([]interface{}, 0)
		if limitOk {
			sqlQuery += " LIMIT $" + strconv.Itoa(countOfParameters)
			countOfParameters++
			args = append(args, limit)
		} else {
			sqlQuery += " LIMIT 5"
		}
		if offsetOk {
			sqlQuery += " OFFSET $" + strconv.Itoa(countOfParameters)
			countOfParameters++
			args = append(args, offset)
		}
		rows, rowsErr := sgh.Db.Query(sqlQuery, args...)
		if rowsErr != nil {
			log.Println(rowsErr)
		}
		defer rows.Close()
		games := []Game{}
		for rows.Next() {
			var ID *int64
			var body *[]byte
			rows.Scan(&ID, &body)
			var game = &Game{}
			if err := json.Unmarshal(*body, game); err == nil {
				game.ID = *ID
				games = append(games, *game)
			} else {
				log.Println(err)
			}
		}
		json.NewEncoder(w).Encode(games)
	} else {
		log.Println(parseErr)
	}
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
	game := &Game{}
	if err := json.NewDecoder(r.Body).Decode(game); err == nil {
		//database logic of saving game
		tx, txErr := cgh.Db.Begin()
		if txErr != nil {
			log.Println("Transaction start failed: ", txErr)
			err = txErr
			return
		}
		defer func(tx *sql.Tx) {
			if r := recover(); r != nil {
				tx.Rollback()
			}
		}(tx)
		if data, err := json.Marshal(game); err == nil {
			_, execTxErr := tx.Exec("INSERT INTO GAMES (BODY) VALUES($1)", data)
			if execTxErr != nil {
				log.Println("Query cannot be executed: ", execTxErr)
				err = execTxErr
				panic(execTxErr)
			}
			if txCommitErr := tx.Commit(); txCommitErr != nil {
				log.Println("Transaction cannot be commited: ", txCommitErr)
				err = txCommitErr
				panic(txCommitErr)
			}
		} else {
			log.Println(err)
		}
	} else {
		log.Println(err)
	}
}

//Game struct represent rest object for game entity
type Game struct {
	ID               int64     `json:"id"`
	Desc             string    `json:"desc"`
	FirstQuestion    *Question `json:"firstQuestion"`
	DoubleQuestion   *Question `json:"doubleQuestion"`
	InversedQuestion *Question `json:"inversedQuestion"`
}

//Question struct represent rest object for question entity
type Question struct {
	Text    string    `json:"text"`
	answers []*Answer `json:"answers"`
}

//Answer struct represent rest object for answer entity
type Answer struct {
	Text  string `json:"text"`
	Score int32  `json:"score"`
}
