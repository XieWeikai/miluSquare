package main

import (
	"net/http"
	"net/url"
	"strings"
)

type Message struct{
	Alert string
	Exist bool
	Msg string
}

func (m *Message)Send(w http.ResponseWriter){
	if m.Alert == "" {
		m.Alert = "alert-danger"
	}
	c := &http.Cookie{
		Name:"message",
		Value:url.QueryEscape(m.Msg + ";;" + m.Alert),
		Path: "/",
	}
	http.SetCookie(w,c)
	//log.Println("msg set",m.Msg,"!!!!!!!!!!")
}
func (m *Message)Get(w http.ResponseWriter,r *http.Request){
	c,err := r.Cookie("message")
	if err != nil {
		m.Exist = false
		return
	}
	m.Exist = true
	data,_ := url.QueryUnescape(c.Value)
	s := strings.Split(data,";;")
	m.Msg,m.Alert = s[0],s[1]
	c.MaxAge = -1
	//c.Expires = time.Now().Add(-1 * time.Second)
	c.Path = "/"
	http.SetCookie(w,c)
	//log.Println(c.Path,"  msg get set path",m.Msg,"!!!!!!!!!!","\n",c)
}