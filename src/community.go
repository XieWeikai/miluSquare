package main

import (
	"data"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

func delComm(w http.ResponseWriter,r *http.Request,p httprouter.Params){
	_,id := sGet(r,"userId")
	if id == nil {
		w.WriteHeader(403)
		log.Println("delCommErr:","does not login")
		return
	}
	pId,e := strconv.Atoi(p.ByName("id"))
	if e != nil {
		w.WriteHeader(400)
		fmt.Fprintln(w,"not allow")
		log.Println("delCommErr:","parse id err:"+e.Error())
		return
	}
	comm := &data.Community{Id: pId}
	comm.GetCommunity()
	comm.GetPost()
	if comm.UserId != id.(int) {
		w.WriteHeader(403)
		log.Println("delCommErr","comm.UserId(",comm.UserId, ") != userId(",id,")")
		return
	}
	e = delC(comm)
	if e != nil {
		w.WriteHeader(500)
	}
}

func delC(c *data.Community)(e error){
	c.GetPost()
	c.Posts().ForEach(func(c *data.Post) {
		e = delP(c)
		if e != nil {
			return
		}
	})
	if e != nil {
		return
	}
	e = c.Delete()
	return
}

func commUpdatePage(w http.ResponseWriter,r *http.Request,p httprouter.Params){
	_,id := sGet(r,"userId")
	if id == nil {
		msg := &Message{Msg: "请先登录"}
		msg.Send(w)
		http.Redirect(w,r,"/login",303)
		return
	}

	cId,e := strconv.Atoi(p.ByName("id"))
	if e != nil {
		w.WriteHeader(400)
		return
	}
	c := &data.Community{Id: cId}
	c.GetCommunity()
	if c.UserId != id.(int) {
		w.WriteHeader(403)
		return
	}
	msg := &Message{}
	msg.Get(w,r)
	data := struct {
		Comm *data.Community
		Msg *Message
	}{c,msg}
	t,_ := template.ParseFiles("html/private.nav.html","html/footer.html","html/commUpdatePage.html")
	t.ExecuteTemplate(w,"layout",data)
}

func commUpdate(w http.ResponseWriter,r *http.Request,p httprouter.Params){
	_,id := sGet(r,"userId")
	if id == nil {
		msg := &Message{Msg: "请先登录"}
		msg.Send(w)
		http.Redirect(w,r,"/login",303)
		return
	}

	cId,e := strconv.Atoi(r.FormValue("id"))
	if e != nil {
		w.WriteHeader(400)
		return
	}
	c := &data.Community{Id: cId}
	c.GetCommunity()
	if c.UserId != id.(int) {
		w.WriteHeader(403)
		return
	}
	c.Name = r.FormValue("name")
	c.BelongTo = r.FormValue("belong")
	c.Desc = r.FormValue("desc")

	msg := &Message{}
	e = c.Update()
	if e != nil {
		msg.Msg = "修改出错:"+e.Error()
	}else{
		msg.Msg = "修改成功"
		msg.Alert = "alert-success"
	}
	msg.Send(w)
	http.Redirect(w,r,"/update/community/"+r.FormValue("id"),303)
}

func showComm(w http.ResponseWriter,r *http.Request,p httprouter.Params) {
	s,uid := sGet(r,"userId")
	id,e := strconv.Atoi(p.ByName("id"))
	if e != nil{
		msg := &Message{Msg: "未知post_id,无法查看"}
		msg.Send(w)
		pageNotF(w,r)
		return
	}
	c := &data.Community{Id: id}
	e = c.GetCommunity()
	if e != nil {
		msg := &Message{Msg: "未知错误:"+e.Error()}
		msg.Send(w)
		pageNotF(w,r)
		return
	}
	c.GetPost()
	data := struct {
		Comm *data.Community
		Posts data.Posts
		Uid int
	}{
		Comm: c,
		Posts: c.Posts(),
	}
	var t *template.Template
	if s == nil {
		t,e = template.ParseFiles("html/showComm.html","html/public.nav.html","html/footer.html")
	}else {
		t,e = template.ParseFiles("html/showComm.html","html/private.nav.html","html/footer.html")
		data.Uid = uid.(int)
	}

	e = t.ExecuteTemplate(w,"layout",data)
	if e != nil {
		log.Println("ExecErr:"+e.Error())
	}
}

func newComm(w http.ResponseWriter,r *http.Request,p httprouter.Params) {
	s := session.GetSess(r)
	if s == nil {
		msg := Message{Msg: "请先登录后再创建"}
		msg.Send(w)
		http.Redirect(w,r,"/login",303)
		return
	}
	msg := &Message{}
	msg.Get(w,r)
	t,err := template.ParseFiles("html/newCommunity.html","html/private.nav.html","html/footer.html")
	if err != nil {
		panic(err)
	}
	t.ExecuteTemplate(w,"layout",msg)
}

func commHandle(w http.ResponseWriter,r *http.Request,p httprouter.Params){
	s,id := sGet(r,"userId")
	if s == nil {
		msg := &Message{Msg: "请登录后再执行该操作"}
		msg.Send(w)
		http.Redirect(w,r,"/login",303)
		return
	}
	comm := data.Community{Name: r.FormValue("name"),Desc: r.FormValue("desc"),BelongTo: r.FormValue("belong"),UserId: id.(int)}
	e := comm.Put()
	msg := &Message{}
	if e == nil{
		msg.Msg = "创建成功"
		msg.Alert = "alert-success"
	}else {
		msg.Msg = "创建失败:" + e.Error()
		log.Println(e)
	}
	msg.Send(w)
	http.Redirect(w,r,"/newcommunity",303)
}
