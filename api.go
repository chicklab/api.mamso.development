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
)

/* APIのバージョン定義 */
const version string = "1"

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

/* ルーティング */
func main() {
	l, err := net.Listen("tcp", ":4000")
	// l, err := net.Listen("unix", "/var/run/go-fcgi.sock")
	if err != nil {
		return
	}
	r := mux.NewRouter()
	http.Handle("/", r)
	r.HandleFunc("/v"+version+"/archives.json", Archives)
	r.HandleFunc("/v"+version+"/archives/{id}.json", Archives_Id)
	// r.HandleFunc("/articles/{category}/", ArticlesCategoryHandler)
	// r.HandleFunc("/articles/{category}/{id:[0-9]+}", ArticleHandler)
	r.NotFoundHandler = http.HandlerFunc(NotFoundHandler)
	fcgi.Serve(l, nil)

}

func Archives(w http.ResponseWriter, r *http.Request) {

	var container Response_Container
	var meta Response_Meta
	var archives []Response_Archives

	vars := mux.Vars(r)
	a := map[string]string{}
	db := mydb.MyDB{}
	db.Connect()
	defer db.Close()

	for f, v := range vars {
		a[f] = v
	}
	rows := db.Query("select article_id, article_title, article_link, article_img_tag, article_description, article_youtube_video_id, article_publish_date, article_insert_date, article_insert_year, article_insert_month, article_insert_day, article_impression_num, favorite_article_num, comment_num, comment_posted_at, comment_user_id, evaluate_point, total_point, channel_name, channel_key, channel_category_id, channel_description, channel_media_id, category_name, media_name from V_information")
	defer rows.Close()

	for rows.Next() {
		var k Response_Archives
		err := rows.Scan(&k.Id, &k.Title, &k.Link, &k.Img_tag, &k.Description, &k.Youtube_video_id, &k.Publish_date, &k.Insert_date, &k.Insert_year, &k.Insert_month, &k.Insert_day, &k.Impression_num, &k.Favorite_num, &k.Comment_num, &k.Comment_posted_at, &k.Comment_user_id, &k.Evaluate_point, &k.Total_point, &k.Channel_name, &k.Channel_key, &k.Channel_category_id, &k.Channel_description, &k.Channel_media_id, &k.Category_name, &k.Media_name)
		if err != nil {
			WriteErrorLogFile(err)
		}
		archives = append(archives, k)
	}

	// JSON return
	defer func() {
		meta = Response_Meta{Message: "OK", Code: 200}
		container.Meta = meta
		container.Result = archives
		outjson, err := json.Marshal(container)
		if err != nil {
			WriteErrorLogFile(err)
		}
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, string(outjson))
	}()
}

func Archives_Id(w http.ResponseWriter, r *http.Request) {

	var container Response_Container
	var meta Response_Meta
	var archives Response_Archives

	vars := mux.Vars(r)
	a := map[string]string{}
	db := mydb.MyDB{}
	db.Connect()
	defer db.Close()

	for f, v := range vars {
		a[f] = v
	}
	rows := db.Query("select article_id, article_title, article_link, article_img_tag, article_description, article_youtube_video_id, article_publish_date, article_insert_date, article_insert_year, article_insert_month, article_insert_day, article_impression_num, favorite_article_num, comment_num, comment_posted_at, comment_user_id, evaluate_point, total_point, channel_name, channel_key, channel_category_id, channel_description, channel_media_id, category_name, media_name from V_information " + "where article_id = " + a["id"])
	defer rows.Close()

	for rows.Next() {
		var k Response_Archives
		err := rows.Scan(&k.Id, &k.Title, &k.Link, &k.Img_tag, &k.Description, &k.Youtube_video_id, &k.Publish_date, &k.Insert_date, &k.Insert_year, &k.Insert_month, &k.Insert_day, &k.Impression_num, &k.Favorite_num, &k.Comment_num, &k.Comment_posted_at, &k.Comment_user_id, &k.Evaluate_point, &k.Total_point, &k.Channel_name, &k.Channel_key, &k.Channel_category_id, &k.Channel_description, &k.Channel_media_id, &k.Category_name, &k.Media_name)
		if err != nil {
			WriteErrorLogFile(err)
		}
		archives = k
	}

	defer func() {
		meta = Response_Meta{Message: "OK", Code: 200}
		container.Meta = meta
		container.Result = archives
		outjson, err := json.Marshal(container)
		if err != nil {
			WriteErrorLogFile(err)
		}
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, string(outjson))
	}()
}

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {

	// http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
	var container Response_Container
	var meta Response_Meta

	meta = Response_Meta{Message: "Not Found", Code: 404}
	container.Meta = meta
	container.Result = nil
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
