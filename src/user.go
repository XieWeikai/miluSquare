package main

import (
	"crypto/md5"
	"data"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"html/template"
	"io"
	"log"
	"net/http"
	"time"
)

func userUpdate(w http.ResponseWriter,r *http.Request,p httprouter.Params){
	_,id := sGet(r,"userId")
	if id == nil {
		pageNotF(w,r)
		return
	}
	u := &data.User{Id: id.(int)}
	u.GetUser()
	var v string
	if v = r.FormValue("userName");v != "" {
		u.Name = v
	}
	if v = r.FormValue("email");v != ""{
		u.Email = v
	}
	if v = r.FormValue("password");v != ""{
		u.SetPassword(pw(v))
	}
	e := u.Update()
	if e != nil {
		msg := &Message{Msg: "修改出错:"+e.Error()}
		msg.Send(w)
		http.Redirect(w,r,"/home/info",303)
	}else{
		msg := &Message{Msg: "修改成功",Alert: "alert-success"}
		msg.Send(w)
		http.Redirect(w,r,"/home/info",303)
	}
}

func homeHandle(w http.ResponseWriter,r *http.Request,p httprouter.Params) {
	_,id := sGet(r,"userId")
	if id == nil {
		msg := &Message{Msg: "请先登录"}
		msg.Send(w)
		http.Redirect(w,r,"/login",303)
		return
	}
	switch p.ByName("path") {
	case "info":
		info(w,r,id.(int))
		return
	case "post":
		homePost(w,r,id.(int))
	case "community":
		homeComm(w,r,id.(int))
	default:
		pageNotF(w,r)
	}

}


func homeComm(w http.ResponseWriter,r *http.Request,id int) {
	u := &data.User{Id: id}
	u.GetComm()
	t,_ := template.ParseFiles("html/private.nav.html","html/footer.html","html/homeComm.html")
	t.ExecuteTemplate(w,"layout",u.Communities())
}

func homePost(w http.ResponseWriter,r *http.Request,id int) {
	t,_ := template.ParseFiles("html/private.nav.html","html/footer.html","html/homePost.html")
	e := t.ExecuteTemplate(w,"layout",id)
	if e != nil {
		log.Println("ExecErr:",e)
	}
}

func info(w http.ResponseWriter,r *http.Request,id int){
	var t *template.Template
	t,e := t.ParseFiles("html/private.nav.html","html/footer.html")
	u := &data.User{Id: id}
	u.GetUser()
	msg := &Message{}
	msg.Get(w,r)
	data := struct {
		Msg *Message
		User *data.User
	}{msg,u}
	t,e = t.ParseFiles("html/info.html")
	if e != nil {
		log.Println("parseErr:",e)
		return
	}
	e = t.ExecuteTemplate(w,"layout",data)
	if e != nil {
		log.Println("execErr:",e)
		return
	}
}

func home(w http.ResponseWriter,r *http.Request,p httprouter.Params) {
	_,id := sGet(r,"userId")
	if id == nil {
		msg := &Message{Msg: "登录后才能进行相关操作"}
		msg.Send(w)
		http.Redirect(w,r,"/login",303)
		return
	}
	u := &data.User{Id: id.(int)}
	u.GetUser()
	t,_ := template.ParseFiles("html/home.html","html/private.nav.html","html/footer.html")
	t.ExecuteTemplate(w,"layout",u)
}

func login(w http.ResponseWriter,r *http.Request,p httprouter.Params){
	t,_ := template.ParseFiles("./html/login.html","./html/public.nav.html","./html/footer.html")
	msg := &Message{}
	msg.Get(w,r)
	t.ExecuteTemplate(w,"layout",msg)
}

func loginHandle(w http.ResponseWriter,r *http.Request,p httprouter.Params) {
	u := data.User{Email: r.FormValue("email")}
	e := u.GetUserByEmail()
	if e != nil || u.Password() != pw(r.FormValue("password")) {
		msg := Message{Msg: "用户名或密码错误"}
		msg.Send(w)
		http.Redirect(w,r,"/login",303)
		return
	}
	var lifetime time.Duration = -1
	if r.FormValue("remember") == "on" {
		lifetime = 2592000 * time.Second
	}
	log.Println("new session from",r.RemoteAddr," lifetime is",lifetime)
	s := session.NewSess(w,r,lifetime)
	s.Set("userId",u.Id)
	http.Redirect(w,r,"/",303)
}

func logout(w http.ResponseWriter,r *http.Request,p httprouter.Params){
	session.DestroySess(w,r)
	http.Redirect(w,r,"/login",303)
}

func signUp(w http.ResponseWriter,r *http.Request,p httprouter.Params){
	t,_ := template.ParseFiles("./html/signup.html","./html/public.nav.html","./html/footer.html")
	msg := &Message{}
	msg.Get(w,r)
	t.ExecuteTemplate(w,"layout",msg)
}


func signupHandle(w http.ResponseWriter,r *http.Request,p httprouter.Params) {
	u := data.User{
		Name:r.FormValue("username"),
		Email:r.FormValue("email"),
	}

	u.SetPassword(pw(r.FormValue("password")))
	e := u.PutUser()
	if e != nil{
		msg := &Message{Msg: "注册失败，邮箱已使用"}
		msg.Send(w)
		http.Redirect(w,r,"/signup",303)
		return
	}
	http.Redirect(w,r,"/",303)
}

func pw(p string)string{
	h := md5.New()
	io.WriteString(h, "加密的密码")

	//pwmd5等于e10adc3949ba59abbe56e057f20f883e
	pwmd5 :=fmt.Sprintf("%x", h.Sum(nil))

	//指定两个 salt： salt1 = @#$%   salt2 = ^&*()
	salt1 := "#*"
	salt2 := "!*)"

	//salt1+用户名+salt2+MD5拼接
	io.WriteString(h, salt1)
	io.WriteString(h, p)
	io.WriteString(h, salt2)
	io.WriteString(h, pwmd5)
	return fmt.Sprintf("%x", h.Sum(nil))
}
