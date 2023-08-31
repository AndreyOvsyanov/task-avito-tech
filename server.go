package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/joho/sqltocsv"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// Segment struct
type Segment struct {
	ID   int    `json:"id"`
	Slug string `json:"slug"`
}

// User struct
type User struct {
	ID         int    `json:"id"`
	FIO        string `json:"fio"`
	Created_at string `json:"created_at"`
	Updated_at string `json:"updated_at"`
}

// User's Segments struct
type UserSegments struct {
	UserID   int       `json:"userID"`
	Segments []Segment `json:"segments"`
}

// User's add and remove Segment struct
type UserSegmentsRequest struct {
	UserID         int      `json:"user_id"`
	AddSegments    []string `json:"add_segments"`
	RemoveSegments []string `json:"remove_segments"`
}

var db *sql.DB
var error_db error

// Route's for task
func actionUser(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	user_id, added, removed := params.Get("user_id"), params["add"], params["remove"]
	userAddRemove(user_id, strings.Split(added[0], " "), strings.Split(removed[0], " "))
}

func defaultPage(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Hello Avito!")
	fmt.Fprint(w, "Hello Avito!")
}

func getUserInfo(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	res, err := db.Query("SELECT * FROM user")
	if err != nil {
		log.Fatal(err)
	}

	defer func(result *sql.Rows) { _ = result.Close() }(res)

	var users []User

	for res.Next() {
		var user User
		if err := res.Scan(&user.ID, &user.FIO, &user.Created_at, &user.Updated_at); err != nil {
			log.Fatal(err)
		}

		users = append(users, user)
	}

	usersJson, err := json.MarshalIndent(users, "", "\t")
	if err != nil {
		log.Fatal(err)
	}
	w.Write(usersJson)

}

func createSegment(w http.ResponseWriter, r *http.Request) {
	slug := r.URL.Query().Get("slug")

	var exist bool
	if err := db.QueryRow("SELECT EXISTS(SELECT * FROM segment WHERE slug = ?)", slug).Scan(&exist); err != nil {
		log.Fatal(err)
	}

	if exist {
		fmt.Fprintf(w, "Segment %s is already exist", slug)
	} else {
		if _, err := db.Exec("INSERT INTO segment(slug) VALUES (?)", slug); err != nil {
			fmt.Fprint(w, err.Error())
		}

		fmt.Fprint(w, "Create segment: ", slug)
	}

	w.WriteHeader(http.StatusNoContent)
}

func deleteSegment(w http.ResponseWriter, r *http.Request) {
	slug := r.URL.Query().Get("slug")

	var exist bool
	if err := db.QueryRow("SELECT EXISTS(SELECT * FROM segment WHERE slug = ?)", slug).Scan(&exist); err != nil {
		log.Fatal(err)
	}

	if !exist {
		fmt.Fprintf(w, "Segment %s does not exist", slug)
	} else {
		if _, err := db.Exec("DELETE FROM segment WHERE slug = ?", slug); err != nil {
			fmt.Fprint(w, err.Error())
		}
		fmt.Fprint(w, "Delete segment: ", slug)
	}

	w.WriteHeader(http.StatusNoContent)
}

func getSegmentInfo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	res, err := db.Query("SELECT * FROM segment")
	if err != nil {
		log.Fatal(err)
	}

	defer func(result *sql.Rows) { _ = result.Close() }(res)

	var segments []Segment

	for res.Next() {
		var segment Segment
		if err := res.Scan(&segment.ID, &segment.Slug); err != nil {
			log.Fatal(err)
		}

		segments = append(segments, segment)
	}

	usersJson, err := json.MarshalIndent(segments, "", "\t")
	if err != nil {
		log.Fatal(err)
	}
	w.Write(usersJson)
}

func getUserHistory(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	user_id := params.Get("user_id")
	year := params.Get("year")
	month := params.Get("month")

	// Создаем слайс данных, которые нужно записать в CSV-файл
	query := `SELECT h.user_id, s.slug, h.type_of_operation, h.date_of_operation FROM segment s
   			  JOIN history_operation h ON s.id = h.segment_id 
              WHERE h.user_id = ? AND YEAR(h.date_of_operation) = ? AND MONTH(h.date_of_operation) = ?`

	history, err := db.Query(query, user_id, year, month)

	defer func(result *sql.Rows) { _ = result.Close() }(history)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-type", "text/csv; charset=utf-8")
	w.Header().Set(
		"Content-Disposition", fmt.Sprintf("attachment; filename=\"historyUser(%s).csv\"", user_id))

	sqltocsv.Write(w, history)
}

func getUserSegments(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	user_id := params.Get("user_id")

	q := "SELECT u.user_id, s.id, s.slug FROM user_segments u JOIN segment s ON u.segment_id = s.id WHERE u.user_id = ?"
	rows, err := db.Query(q, user_id)
	if err != nil {
		fmt.Fprint(w, "Такого пользователя нет")
	}

	defer func(result *sql.Rows) { _ = result.Close() }(rows)

	var userSegments UserSegments
	if userSegments.UserID, err = strconv.Atoi(user_id); err != nil {
		log.Fatal(err)
	}

	for rows.Next() {
		var segment Segment
		if err = rows.Scan(&user_id, &segment.ID, &segment.Slug); err != nil {
			log.Fatal(err)
		}

		userSegments.Segments = append(userSegments.Segments, segment)
	}

	var userSegmentsJson []byte
	if userSegmentsJson, err = json.MarshalIndent(userSegments, "", "\t"); err != nil {
		log.Fatal(err)
	}

	w.Write(userSegmentsJson)

	if len(userSegments.Segments) == 0 {
		fmt.Printf("Пользователь %s не cостоит ни в одном сегменте\n", user_id)
	} else {
		fmt.Printf("Пользователь %s cостоит в %d сегментах\n", user_id, len(userSegments.Segments))
		for _, segment := range userSegments.Segments {
			fmt.Println(segment.Slug)
		}
	}
}

// Help's Functions
func updateHistoryUsers(user_id string, segment_id string, type_operation string) {
	if _, err := db.Exec("INSERT INTO history_operation(user_id, segment_id, type_of_operation) VALUES (?, ?, ?)",
		user_id, segment_id, type_operation); err != nil {
		log.Fatal(err)
	}
}

func getIDBySlugSegment(segment string) string {
	var segment_id string
	if err := db.QueryRow("SELECT id FROM segment WHERE slug = ?", segment).Scan(&segment_id); err != nil {
		log.Fatal(err)
	}

	return segment_id
}

func existSegmentAUser(segment string, user_id string) (string, bool) {
	segment_id := getIDBySlugSegment(segment)

	var exist bool
	if err := db.QueryRow("SELECT EXISTS(SELECT * FROM user_segments WHERE user_id = ? AND segment_id = ?)",
		user_id, segment_id).Scan(&exist); err != nil {
		log.Fatal(err)
	}

	return segment_id, exist
}

func userAddRemove(user_id string, added_segments []string, remove_segments []string) {
	for _, added_segment := range added_segments {
		segment_id, exist := existSegmentAUser(added_segment, user_id)
		if !exist {
			if _, err := db.Exec(
				"INSERT INTO user_segments(user_id, segment_id) VALUES (?, ?)",
				user_id, segment_id); err != nil {
				log.Fatal(err)
			}
			fmt.Printf("Cегмент %s добавлен пользователю %s\n", added_segment, user_id)
			updateHistoryUsers(user_id, segment_id, "adding a segment")
		} else {
			fmt.Printf("Сегмент %s уже присутствует у пользователя %s\n", added_segment, user_id)
		}
	}

	for _, remove_segment := range remove_segments {
		segment_id, exist := existSegmentAUser(remove_segment, user_id)
		if exist {
			if _, err := db.Exec(
				"DELETE FROM user_segments WHERE user_id = ? AND segment_id = ?",
				user_id, segment_id); err != nil {
				log.Fatal(err)
			}

			fmt.Printf("Cегмент %s удалён у пользователя %s\n", remove_segment, user_id)
			updateHistoryUsers(user_id, segment_id, "deleting a segment")
		} else {
			fmt.Printf("Cегмента %s нет у пользователя %s\n", remove_segment, user_id)
		}
	}

	if _, err := db.Exec("UPDATE user SET updated_at = CURRENT_TIMESTAMP WHERE id = ?", user_id); err != nil {
		log.Fatal(err)
	}
}

func main() {
	db, error_db = sql.Open("mysql", "root:G12e70891@tcp(localhost:3307)/my_database")
	if error_db != nil {
		log.Fatal(error_db)
	}

	db.SetMaxOpenConns(1000)
	db.SetMaxIdleConns(100)
	db.SetConnMaxIdleTime(time.Minute * time.Duration(3))
	db.SetConnMaxLifetime(time.Hour * time.Duration(1))

	defer func(result *sql.DB) { _ = result.Close() }(db)

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	router := mux.NewRouter()
	router.HandleFunc("/", defaultPage)
	router.HandleFunc("/users", getUserInfo)
	router.HandleFunc("/segments", getSegmentInfo)
	router.HandleFunc("/segment/create", createSegment)
	router.HandleFunc("/segment/delete", deleteSegment)
	router.HandleFunc("/segments/user", getUserSegments)
	router.HandleFunc("/useraddremove", actionUser)
	router.HandleFunc("/history/user", getUserHistory)

	log.Fatal(http.ListenAndServe(":8080", router))

}
