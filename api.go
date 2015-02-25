package main

import (
	"./mydb"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net"
	"net/http"
	"net/http/fcgi"
	"strconv"
)

// API Version
const version string = "1"

const (
	limit_def  string = "20"  // query default limit Number
	limit_max  string = "100" // query maximum limit Number
	offset_def string = "0"   // query default offset Number
)

type Response_Container struct {
	Meta   Response_Meta `json:"meta"`
	Result interface{}   `json:"result"`
}
type Response_Meta struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}
type Response_Archives struct {
	Id                  string         `json:"id"`
	Title               string         `json:"title"`
	Link                string         `json:"link"`
	Img_tag             string         `json:"img_tag"`
	Description         string         `json:"description"`
	Youtube_video_id    sql.NullString `json:"youtube_video_id"`
	Publish_date        string         `json:"publish_date"`
	Insert_date         string         `json:"insert_date "`
	Insert_year         int            `json:"insert_year"`
	Insert_month        int            `json:"insert_month"`
	Insert_day          int            `json:"insert_day"`
	Impression_num      int            `json:"impression_num"`
	Favorite_num        int            `json:"favorite_num"`
	Comment_num         int            `json:"comment_num"`
	Comment_posted_at   sql.NullString `json:"comment_posted_at"`
	Comment_user_id     sql.NullString `json:"comment_user_id"`
	Evaluate_point      int            `json:"evaluate_point"`
	Total_point         int            `json:"total_point"`
	Channel_name        string         `json:"channel_name"`
	Channel_key         string         `json:"channel_key"`
	Channel_category_id int            `json:"channel_category_id"`
	Channel_description string         `json:"channel_description"`
	Channel_media_id    int            `json:"channel_media_id"`
	Category_name       string         `json:"category_name"`
	Media_name          string         `json:"media_name"`
}

type Response_Channels struct {
	Id                string `json:"id"`
	Name              string `json:"name"`
	Url               string `json:"url"`
	Description       string `json:"description"`
	Key               string `json:"key"`
	Archive_number    int    `json:"archive_number"`
	Impression_number int    `json:"impression_number"`
	Category_id       string `json:"channel_id"`
	Category_name     string `json:"category_name"`
}

// Routing
func main() {
	l, err := net.Listen("tcp", ":4000")
	if err != nil {
		return
	}
	r := mux.NewRouter()
	http.Handle("/", r)
	r.HandleFunc("/v"+version+"/archives.json", ArchivesHandler).Methods("GET")
	r.HandleFunc("/v"+version+"/archives/{id}.json", ArchivesIdHandler)
	r.HandleFunc("/v"+version+"/categories/{id}.json", CategoriesNameHandler).Methods("GET")
	r.HandleFunc("/v"+version+"/channels.json", ChannelsHandler)
	r.HandleFunc("/v"+version+"/channels/{id}.json", ChannelsNameHandler).Methods("GET")
	r.NotFoundHandler = http.HandlerFunc(NotFoundHandler)
	fcgi.Serve(l, nil)
}

func ArchivesHandler(w http.ResponseWriter, r *http.Request) {

	var archives []Response_Archives

	// GetParams
	count := r.URL.Query().Get("count")
	cursor := r.URL.Query().Get("cursor")

	// ValidateParams
	if (!CheckDigit(count) && count != "") || (!CheckDigit(cursor) && cursor != "") {
		ResponseJson(w, "Not Found", 404, nil)
		return
	}
	if count == "" {
		count = limit_def
	}
	if cursor == "" {
		cursor = offset_def
	}

	archives = CreateArchivesQuery("limit " + count + " offset " + cursor)

	defer func() {
		ResponseJson(w, "OK", 200, archives)
	}()
}

func ArchivesIdHandler(w http.ResponseWriter, r *http.Request) {

	var archives []Response_Archives

	vars := mux.Vars(r)
	a := map[string]string{}
	for f, v := range vars {
		a[f] = v
	}
	// ValidateParams
	if !CheckDigit(a["id"]) && a["id"] != "" {
		ResponseJson(w, "Not Found", 404, nil)
		return
	}

	archives = CreateArchivesQuery("where article_id = " + a["id"])

	defer func() {
		ResponseJson(w, "OK", 200, archives)
	}()
}

func CategoriesNameHandler(w http.ResponseWriter, r *http.Request) {

	var archives []Response_Archives

	vars := mux.Vars(r)
	// GetParams
	count := r.URL.Query().Get("count")
	cursor := r.URL.Query().Get("cursor")

	if count == "" {
		count = limit_def
	}
	if cursor == "" {
		cursor = offset_def
	}
	a := map[string]string{}
	for f, v := range vars {
		a[f] = v
	}
	// ValidateParams
	if (!CheckDigit(count) && count != "") || (!CheckDigit(cursor) && cursor != "") || (!CheckDigit(a["id"]) && a["id"] != "") {
		ResponseJson(w, "Not Found", 404, nil)
		return
	}

	archives = CreateArchivesQuery("where channel_category_id = " + a["id"] + " limit " + count + " offset " + cursor)

	defer func() {
		ResponseJson(w, "OK", 200, archives)
	}()
}

func ChannelsHandler(w http.ResponseWriter, r *http.Request) {

	var channels []Response_Channels

	channels = CreateChannelsQuery("")

	defer func() {
		ResponseJson(w, "OK", 200, channels)
	}()
}

func ChannelsNameHandler(w http.ResponseWriter, r *http.Request) {

	var archives []Response_Archives

	vars := mux.Vars(r)
	// GetParams
	count := r.URL.Query().Get("count")
	cursor := r.URL.Query().Get("cursor")

	if count == "" {
		count = limit_def
	}
	if cursor == "" {
		cursor = offset_def
	}
	a := map[string]string{}
	for f, v := range vars {
		a[f] = v
	}
	// ValidateParams
	if (!CheckDigit(count) && count != "") || (!CheckDigit(cursor) && cursor != "") || (!CheckDigit(a["id"]) && a["id"] != "") {
		ResponseJson(w, "Not Found", 404, nil)
		return
	}

	archives = CreateArchivesQuery("where channel_id = " + a["id"] + " limit " + count + " offset " + cursor)

	defer func() {
		ResponseJson(w, "OK", 200, archives)
	}()
}

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	// http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
	ResponseJson(w, "Not Found", 404, nil)
}

func CreateArchivesQuery(query string) []Response_Archives {

	var archives []Response_Archives

	db := mydb.MyDB{}
	db.Connect()
	defer db.Close()

	rows := db.Query("select article_id, article_title, article_link, article_img_tag, article_description, article_youtube_video_id, article_publish_date, article_insert_date, article_insert_year, article_insert_month, article_insert_day, article_impression_num, favorite_article_num, comment_num, comment_posted_at, comment_user_id, evaluate_point, total_point, channel_name, channel_key, channel_category_id, channel_description, channel_media_id, category_name, media_name from V_information " + query)
	defer rows.Close()

	for rows.Next() {
		var k Response_Archives
		err := rows.Scan(&k.Id, &k.Title, &k.Link, &k.Img_tag, &k.Description, &k.Youtube_video_id, &k.Publish_date, &k.Insert_date, &k.Insert_year, &k.Insert_month, &k.Insert_day, &k.Impression_num, &k.Favorite_num, &k.Comment_num, &k.Comment_posted_at, &k.Comment_user_id, &k.Evaluate_point, &k.Total_point, &k.Channel_name, &k.Channel_key, &k.Channel_category_id, &k.Channel_description, &k.Channel_media_id, &k.Category_name, &k.Media_name)
		if err != nil {
			WriteErrorLogFile(err)
		}
		archives = append(archives, k)
	}
	return archives
}

func CreateChannelsQuery(query string) []Response_Channels {

	var channels []Response_Channels

	db := mydb.MyDB{}
	db.Connect()
	defer db.Close()

	rows := db.Query("select channel_id, channel_name, channel_url, channel_description, channel_key, article_num, article_impression_num, channel_category_id, category_name from V_channel " + query)
	defer rows.Close()

	for rows.Next() {
		var k Response_Channels
		err := rows.Scan(&k.Id, &k.Name, &k.Url, &k.Description, &k.Key, &k.Archive_number, &k.Impression_number, &k.Category_id, &k.Category_name)
		if err != nil {
			WriteErrorLogFile(err)
		}
		channels = append(channels, k)
	}
	return channels
}

func ResponseJson(w http.ResponseWriter, message string, code int, archives interface{}) {

	var (
		container Response_Container
		meta      Response_Meta
	)
	meta = Response_Meta{Message: message, Code: code}
	container.Meta = meta
	container.Result = archives
	outjson, err := json.Marshal(container)
	if err != nil {
		WriteErrorLogFile(err)
	}
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, string(outjson))
}

func WriteErrorLogFile(l error) {
	fmt.Println(l)
}

func CheckDigit(str string) bool {
	flag := false
	if _, err := strconv.Atoi(str); err == nil {
		flag = true
	}
	return flag
}
