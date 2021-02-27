package REST

import (
	"data"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
)

func postComments(w http.ResponseWriter,r *http.Request,p httprouter.Params) {
	id,e := strconv.Atoi(r.FormValue("id"))
	if e != nil {
		outPut(w,nil,e)
		return
	}
	com := data.Post{Id:id}
	e = com.GetComment()
	outPut(w,com.Comments(),e)
}

func comment(w http.ResponseWriter,r *http.Request,p httprouter.Params) {
	if p.ByName("method") == "limit" {
		limitComments(w,r,p)
		return
	}
	if p.ByName("method") == "all" {
		allComments(w,r,p)
		return
	}
	id,e := strconv.Atoi(p.ByName("method"))
	if e != nil {
		outPut(w,nil,e)
		return
	}
	if id == -1 {
		outPut(w,data.NumOfComments(),nil)
		return
	}
	com := data.Comment{Id:id}
	e = com.GetComment()
	outPut(w,com,e)
}

func limitComments(w http.ResponseWriter,r *http.Request,p httprouter.Params) {
	outPut(w,nil,fmt.Errorf("error:we have not provide this service"))
}

func allComments(w http.ResponseWriter,r *http.Request,p httprouter.Params) {
	outPut(w,data.AllComments(),nil)
}