package REST

import (
	"data"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
)

func limitUsers(w http.ResponseWriter,r *http.Request,p httprouter.Params){
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
	users,e := data.GetUsers(l,right)
	if e != nil {
		outPut(w,nil,e)
		return
	}
	outPut(w,users,nil)

}

func allUsers(w http.ResponseWriter,r *http.Request,p httprouter.Params){
	users := data.AllUsers()
	outPut(w,users,nil)
}

func user(w http.ResponseWriter,r *http.Request,p httprouter.Params)  {
	if p.ByName("method") == "communities" {
		userCommunities(w,r,p)
		return
	}
	if p.ByName("method") == "posts"{
		userPosts(w,r,p)
		return
	}
	if p.ByName("method") == "limit" {
		limitUsers(w,r,p)
		return
	}
	if p.ByName("method") == "all" {
		allUsers(w,r,p)
		return
	}
	id,e := strconv.Atoi(p.ByName("method"))
	if e !=nil {
		outPut(w,nil,fmt.Errorf("no such user id"))
		return
	}
	if id == -1 {
		n := data.NumOfUsers()
		outPut(w,n,nil)
		return
	}
	user := data.User{Id:id}
	e = user.GetUser()
	if e != nil {
		outPut(w,nil,fmt.Errorf("use not found"))
		return
	}
	outPut(w,user,nil)
}