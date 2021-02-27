package REST

import (
	"data"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
)

func userPosts(w http.ResponseWriter,r *http.Request,p httprouter.Params){
	id,e := strconv.Atoi(r.FormValue("id"))
	if e != nil {
		outPut(w,nil,fmt.Errorf("invalid id"))
		return
	}
	user := data.User{Id:id}
	user.GetPosts()
	outPut(w,user.Posts(),nil)
}

func communityPosts(w http.ResponseWriter,r *http.Request,p httprouter.Params) {
	id,e := strconv.Atoi(r.FormValue("id"))
	if e != nil {
		outPut(w,nil,e)
		return
	}
	comm := data.Community{Id:id}
	e = comm.GetPost()
	outPut(w,comm.Posts(),e)
}

func post(w http.ResponseWriter,r *http.Request,p httprouter.Params){
	if p.ByName("method") == "comments" {
		postComments(w,r,p)
		return
	}
	if p.ByName("method") == "limit" {
		limitPosts(w,r,p)
		return
	}
	if p.ByName("method") == "all" {
		allPosts(w,r,p)
		return
	}
	id,e := strconv.Atoi(p.ByName("method"))
	if e != nil {
		outPut(w,nil,e)
		return
	}
	if id == -1 {
		outPut(w,data.NumOfPosts(),nil)
		return
	}
	com := data.Post{Id:id}
	e = com.GetPost()
	outPut(w,com,e)
}

func limitPosts(w http.ResponseWriter,r *http.Request,p httprouter.Params){
	l,e := strconv.Atoi(r.FormValue("l"))
	if e != nil {
		outPut(w,nil,fmt.Errorf("invalid request"))
		return
	}
	right,e := strconv.Atoi(r.FormValue("r"))
	if e != nil {
		outPut(w,nil,fmt.Errorf("invalid request"))
		return
	}
	posts,e := data.GetPosts(l,right)
	outPut(w,posts,e)
}

func allPosts(w http.ResponseWriter,r *http.Request,p httprouter.Params){
	outPut(w,data.AllPost(),nil)
}
