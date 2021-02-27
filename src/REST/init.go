package REST

import (
	"data"
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
	"os"
)

type message struct {
	Msg string `json:"msg"`
	Value interface{} `json:"value"`
}

const format  = true
const debug = false

var ro = httprouter.New()
type myserver struct {}

func (_ myserver)ServeHTTP(w http.ResponseWriter,r *http.Request){
	//var invalid = true
	//for _,addr := range config.Address {
	//	if strings.HasPrefix(r.RemoteAddr,addr){
	//		invalid = false
	//		break
	//	}
	//
	//}
	//if invalid {
	//	log.Println("invalid address",r.RemoteAddr," url: ",r.URL," host: ",r.Host)
	//	fmt.Fprintln(w,"invlid address")
	//	return
	//}
	log.Println(r.RemoteAddr," url: ",r.URL," host: ",r.Host)
	ro.ServeHTTP(w,r)
}

var config struct{Address []string `json:"address"`}

//const configDir  = "./src/REST/config.json"

func init(){
	ro.GET("/api/:item/:method",item)

	//configjson,err := os.Open(configDir)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//dec := json.NewDecoder(configjson)
	//dec.Decode(&config)
	//fmt.Println(config)
}

func item(w http.ResponseWriter,r *http.Request,p httprouter.Params){
	switch p.ByName("item") {
	case "user":user(w,r,p)
	case "community":community(w,r,p)
	case "post":post(w,r,p)
	case "comment":comment(w,r,p)
	}
}

func RESTServer(){
	s := myserver{}
	http.ListenAndServe("",s)
}

func RestHandler(w http.ResponseWriter,r *http.Request,p httprouter.Params){
	s := myserver{}
	s.ServeHTTP(w,r)
}

func outPut(w http.ResponseWriter,v interface{},e error){
	enc := json.NewEncoder(w)
	if format {
		enc.SetIndent("","\t")
	}
	msg := message{}
	if e != nil {
		msg.Msg = e.Error()
		enc.Encode(msg)
		return
	}

	msg.Value = v
	if _,ok := v.(data.Post);ok {
		enc.SetEscapeHTML(false)
	}
	if _,ok := v.(data.Posts);ok {
		enc.SetEscapeHTML(false)
	}
	enc.Encode(msg)

	if debug {
		e := json.NewEncoder(os.Stdout)
		e.SetIndent("","\t")
		e.Encode(msg)
		fmt.Println(v)
	}
}
