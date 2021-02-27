package REST

import (
	"data"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
)

func userCommunities(w http.ResponseWriter,r *http.Request,p httprouter.Params){
	id,e := strconv.Atoi(r.FormValue("id"))
	if e != nil {
		outPut(w,nil,fmt.Errorf("invalid id"))
		return
	}
	user := data.User{Id:id}
	user.GetComm()
	outPut(w,user.Communities(),nil)
}

func limitComms(w http.ResponseWriter,r *http.Request,p httprouter.Params){
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
	comms,e := data.GetCommunities(l,right)
	outPut(w,comms,e)
}

func community(w http.ResponseWriter,r *http.Request,p httprouter.Params){
	if p.ByName("method") == "posts" {
		communityPosts(w,r,p)
		return
	}
	if p.ByName("method") == "limit" {
		limitComms(w,r,p)
		return
	}
	if p.ByName("method") == "all" {
		allComms(w,r,p)
		return
	}
	id,e := strconv.Atoi(p.ByName("method"))
	if e != nil {
		outPut(w,nil,e)
		return
	}
	if id == -1 {
		outPut(w,data.NumOfComms(),nil)
		return
	}
	com := data.Community{Id:id}
	e = com.GetCommunity()
	outPut(w,com,e)
}

func allComms(w http.ResponseWriter,r *http.Request,p httprouter.Params){
	outPut(w,data.AllComms(),nil)
}
