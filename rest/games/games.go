package game

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"net/url"
	"strconv"

	"github.com/gorilla/mux"
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
		var offsetInt int64
		var limitInt int64 = 5
		if limitOk {
			value, parseErr := strconv.ParseInt(limit[0], 10, 64)
			if parseErr != nil {
				log.Printf("limit is not defined correctly")
				return
			}
			limitInt = value
		} else {
			log.Printf("limit is not defined correctly")
			return
		}
		if offsetOk {
			value, parseErr := strconv.ParseInt(offset[0], 10, 64)
			if parseErr != nil {
				log.Printf("offset is not defined correctly")
				return
			}
			offsetInt = value
		} else {
			log.Printf("offset is not defined correctly")
			return
		}
		games, _ := searchGames(sgh.Db, limitInt, offsetInt)
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
		createGame(cgh.Db, *game)
	} else {
		log.Println(err)
	}
}

//CreateDeleteGameHandler function
func CreateDeleteGameHandler(db *sql.DB) http.Handler {
	deleteGameHandler := &DeleteGameHandler{Db: db}
	return deleteGameHandler
}

//DeleteGameHandler struct
type DeleteGameHandler struct {
	Db *sql.DB
}

func (dgh *DeleteGameHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, idOk := vars["id"]
	if !idOk {
		log.Printf("id is not defined")
		return
	}
	idint, parseErr := strconv.ParseInt(id, 10, 64)
	if parseErr != nil {
		log.Printf("id is not parsed")
	}
	deleteGame(dgh.Db, idint)
}

//CreateUpdateGameHandler function
func CreateUpdateGameHandler(db *sql.DB) http.Handler {
	updateGameHandler := &UpdateGameHandler{Db: db}
	return updateGameHandler
}

//UpdateGameHandler struct
type UpdateGameHandler struct {
	Db *sql.DB
}

func (ugh *UpdateGameHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, idOk := vars["id"]
	if !idOk {
		log.Printf("id is not defined")
		return
	}
	idint, parseErr := strconv.ParseInt(id, 10, 64)
	if parseErr != nil {
		log.Printf("id is not parsed")
	}
	game := &Game{}
	if err := json.NewDecoder(r.Body).Decode(game); err == nil {
		updateGame(ugh.Db, idint, *game)
	} else {
		log.Println(err)
	}
}

//CreateGetGameHandler function
func CreateGetGameHandler(db *sql.DB) http.Handler {
	getGameHandler := &GetGameHandler{Db: db}
	return getGameHandler
}

//GetGameHandler struct
type GetGameHandler struct {
	Db *sql.DB
}

func (ggh *GetGameHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, idOk := vars["id"]
	if !idOk {
		log.Printf("id is not defined")
		return
	}
	idint, parseErr := strconv.ParseInt(id, 10, 64)
	if parseErr != nil {
		log.Printf("id is not parsed")
	}
	game, _ := getGame(ggh.Db, idint)
	json.NewEncoder(w).Encode(game)
}

/*db functions*/

func deleteGame(db *sql.DB, id int64) error {
	tx, txErr := db.Begin()
	if txErr != nil {
		return txErr
	}
	var rollback = newFalse()
	defer func(tx *sql.Tx, rollback *bool) {
		if rollback != nil && *rollback == true {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				log.Println("Transaction cannot be rollbacked", rollbackErr)
			}
		}
	}(tx, rollback)
	_, execTxErr := tx.Exec("DELETE FROM GAMES WHERE ID = $1", id)
	if execTxErr != nil {
		rollback = newTrue()
		return execTxErr
	}
	if txCommitErr := tx.Commit(); txCommitErr != nil {
		rollback = newTrue()
		return txCommitErr
	}
	return nil
}

func updateGame(db *sql.DB, id int64, game Game) error {
	tx, txErr := db.Begin()
	if txErr != nil {
		return txErr
	}
	var rollback = newFalse()
	defer func(tx *sql.Tx, rollback *bool) {
		if rollback != nil && *rollback == true {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				log.Println("Transaction cannot be rollbacked", rollbackErr)
			}
		}
	}(tx, rollback)
	data, marshallErr := json.Marshal(game)
	if marshallErr != nil {
		rollback = newTrue()
		return marshallErr
	}
	_, execTxErr := tx.Exec("UPDATE GAMES SET BODY = $1 WHERE ID = $2", data, id)
	if execTxErr != nil {
		rollback = newTrue()
		return execTxErr
	}
	if txCommitErr := tx.Commit(); txCommitErr != nil {
		rollback = newTrue()
		return txCommitErr
	}
	return nil
}

func createGame(db *sql.DB, game Game) (int64, error) {
	tx, txErr := db.Begin()
	if txErr != nil {
		return 0, txErr
	}
	var rollback = newFalse()
	defer func(tx *sql.Tx, rollback *bool) {
		if rollback != nil && *rollback == true {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				log.Println("Transaction cannot be rollbacked", rollbackErr)
			}
		}
	}(tx, rollback)
	data, marshallErr := json.Marshal(game)
	if marshallErr != nil {
		rollback = newTrue()
		return 0, marshallErr
	}
	result, execTxErr := tx.Exec("INSERT INTO GAMES (BODY) VALUES($1)", data)
	id, idErr := result.LastInsertId()
	if idErr != nil {
		rollback = newTrue()
		return 0, idErr
	}
	if execTxErr != nil {
		rollback = newTrue()
		return 0, execTxErr
	}
	if txCommitErr := tx.Commit(); txCommitErr != nil {
		rollback = newTrue()
		return 0, txCommitErr
	}
	return id, nil
}

func getGame(db *sql.DB, id int64) (*Game, error) {
	rows, rowsErr := db.Query("SELECT BODY FROM GAMES WHERE ID = $1", id)
	if rowsErr != nil {
		return nil, rowsErr
	}
	defer rows.Close()
	var game = &Game{}
	if rows.Next() {
		var body *[]byte
		rows.Scan(&body)
		if err := json.Unmarshal(*body, game); err == nil {
			game.ID = id
		} else {
			return nil, err
		}
	} else {
		return nil, errors.New("Entity with id = " + strconv.FormatInt(id, 10) + " not found")
	}
	return game, nil
}

func searchGames(db *sql.DB, limit int64, offset int64) (*[]Game, error) {
	sqlQuery := "SELECT ID, BODY FROM GAMES WHERE 1 = 1 LIMIT $1 OFFSET $2"
	rows, rowsErr := db.Query(sqlQuery, limit, offset)
	if rowsErr != nil {
		return nil, rowsErr
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
			return nil, err
		}
	}
	return &games, nil
}

func newTrue() *bool {
	b := true
	return &b
}

func newFalse() *bool {
	b := false
	return &b
}

//Game struct represent rest object for game entity
type Game struct {
	ID               int64    `json:"id"`
	Desc             string   `json:"desc"`
	FirstQuestion    Question `json:"firstQuestion"`
	DoubleQuestion   Question `json:"doubleQuestion"`
	InversedQuestion Question `json:"inversedQuestion"`
}

//Question struct represent rest object for question entity
type Question struct {
	Text    string   `json:"text"`
	Answers []Answer `json:"answers"`
}

//Answer struct represent rest object for answer entity
type Answer struct {
	Text  string `json:"text"`
	Score int64  `json:"score,string,omitempty"`
}
