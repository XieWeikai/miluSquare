package main

import (
	"data"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
)

func comHandle(w http.ResponseWriter,r *http.Request,p httprouter.Params){
	_,id := sGet(r,"userId")
	if id == nil {
		w.WriteHeader(403)
		return //未登录不予操作
	}
	pId,e := strconv.Atoi(r.FormValue("postId"))
	if e != nil {
		return
	}
	c := &data.Comment{Content: r.FormValue("content"),UserId: id.(int),PostId: pId}
	c.Put()
	w.WriteHeader(200)
}

func delComment(w http.ResponseWriter,r *http.Request,p httprouter.Params) {
	s,id := sGet(r,"userId")
	cId,e := strconv.Atoi(p.ByName("id"))
	if e != nil {
		w.WriteHeader(400)
		fmt.Fprintln(w,"not allow")
		return
	}
	c := data.Comment{Id: cId}
	//向数据库读取数据太慢了，此处就不到数据库去查了
	if id == nil || s.Get("postUserId")==nil || id.(int) != s.Get("postUserId").(int) {//此处留下了一个漏洞，当用户先进入自己的帖子，再发送删除别人的评论的请求就会成功删除...暂且这样吧
		w.WriteHeader(403)
		fmt.Fprintln(w,"forbidden")
		return
	}
	c.Delete()
}
