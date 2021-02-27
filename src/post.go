package main

import (
	"data"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
)

func postUpdate(w http.ResponseWriter,r *http.Request,p httprouter.Params) {
	s,id := sGet(r,"userId")
	if id == nil {
		w.WriteHeader(403)
		return
	}
	pId,e := strconv.Atoi(r.FormValue("id"))
	if e != nil {
		w.WriteHeader(400)
		fmt.Fprintln(w,"not allow")
		return
	}
	post := &data.Post{Id: pId}
	post.GetPost()
	if post.UserId != id.(int) {
		w.WriteHeader(403)
		return
	}

	post.Title = r.FormValue("title")
	post.Topic = r.FormValue("topic")
	post.Content = r.FormValue("content")
	//log.Println("the post is",post)
	e = post.Update()
	if e != nil{
		msg := &Message{Msg: "出错，修改失败:"+e.Error()}
		msg.Send(w)
		http.Redirect(w,r,"/update/post/"+r.FormValue("id"),303)
		return
	}

	msg := &Message{Msg: "修改成功",Alert: "alert-success"}
	msg.Send(w)
	http.Redirect(w,r,"/update/post/"+r.FormValue("id"),303)

	//paths := s.Get("imgs").(map[string]bool)
	//var imgs []string
	//for path,_ := range paths {
	//	imgs = append(imgs,path)
	//}
	//Imgs := data.Imgs{post.Id,imgs}
	//Imgs.Put()
	//s.Set("imgs",nil)
	endEditor(s,post)
}

func postUpdatePage(w http.ResponseWriter,r *http.Request,p httprouter.Params) {
	s,id := sGet(r,"userId")
	if id == nil {
		w.WriteHeader(403)
		return
	}
	//s.SetAfterFunc(func (s *sess.Sess){
	//	if imgs := s.Get("imgs");imgs != nil {
	//		img := imgs.(map[string]bool)
	//		for path,_ := range img {
	//			log.Println("deletePath",path)
	//			os.Remove(path)
	//		}
	//	}
	//})
	//if s.Get("imgs") == nil {
	//	s.Set("imgs", make(map[string]bool))
	//}
	startEditor(s)
	pId,e := strconv.Atoi(p.ByName("id"))
	if e != nil {
		w.WriteHeader(400)
		fmt.Fprintln(w,"not allow")
		return
	}
	post := &data.Post{Id: pId}
	post.GetPost()
	if post.UserId != id.(int) {
		w.WriteHeader(403)
		return
	}
	msg := &Message{}
	msg.Get(w,r)
	data := struct {
		Msg *Message
		Post *data.Post
		Content template.HTML
	}{msg,post,template.HTML(post.Content)}
	t,_ := template.ParseFiles("html/private.nav.html","html/footer.html","html/postUpdatePage.html")
	t.ExecuteTemplate(w,"layout",data)
}

func delPost(w http.ResponseWriter,r *http.Request,p httprouter.Params){
	_,id := sGet(r,"userId")
	if id == nil {
		w.WriteHeader(403)
		return
	}
	pId,e := strconv.Atoi(p.ByName("id"))
	if e != nil {
		w.WriteHeader(400)
		fmt.Fprintln(w,"not allow")
		return
	}
	post := &data.Post{Id: pId}
	post.GetPost()
	if post.UserId != id.(int) {
		w.WriteHeader(403)
		return
	}
	e = delP(post)
	if e != nil {
		w.WriteHeader(500)
	}
}

func delP(p *data.Post)error{
	p.GetImgs()
	for _,path := range p.Imgs(){
		log.Println("remove imgs:",path)
		e := os.Remove(path)
		if e != nil {
			log.Println("remove imgs error",e)
		}
	}
	return p.Delete()
}

func showPost(w http.ResponseWriter,r *http.Request,p httprouter.Params) {
	id,e := strconv.Atoi(p.ByName("id"))
	if e != nil{
		pageNotF(w,r)
		return
	}
	post := &data.Post{Id: id}
	e = post.GetPost()
	if e != nil {
		pageNotF(w,r)
		return
	}
	post.GetComment()
	comm := &data.Community{Id:post.CommunityId}
	comm.GetCommunity()
	s,uid := sGet(r,"userId")
	data := struct{
		Post *data.Post
		Comm *data.Community
		Delete bool
		Content template.HTML
		Login bool
	}{
		Post: post,
		Comm: comm,
		Content: template.HTML(post.Content),
		Login: false,
	}
	//fmt.Println(comm,"\n-------------\n",post)
	var t *template.Template
	if s == nil {
		t,e = template.ParseFiles("html/showPost.html","html/public.nav.html","html/footer.html","html/comment.html")
	}else {
		data.Login = true
		data.Delete = uid.(int) == post.UserId
		t,e = template.ParseFiles("html/showPost.html","html/private.nav.html","html/footer.html","html/comment.html")
		s.Set("postUserId",post.UserId)//此处留下了一个漏洞，当用户先进入自己的帖子，再发送删除别人的评论的请求就会成功删除...暂且这样吧
	}
	if e != nil {
		log.Println("templateErr:",e)
	}
	e = t.ExecuteTemplate(w,"layout",data)
	if e != nil {
		log.Println("ExecErr:",e)
	}
}

func newPost(w http.ResponseWriter,r *http.Request,p httprouter.Params){
	s := session.GetSess(r)

	if s == nil {
		msg := Message{Msg: "请登录后再发帖子"}
		msg.Send(w)
		http.Redirect(w,r,"/login",303)
		return
	}
	//s.SetAfterFunc(func (s *sess.Sess){
	//	if imgs := s.Get("imgs");imgs != nil {
	//		img := imgs.(map[string]bool)
	//		for path,_ := range img {
	//			log.Println("deletePath",path)
	//			os.Remove(path)
	//		}
	//	}
	//})
	//if s.Get("imgs") == nil {
	//	s.Set("imgs", make(map[string]bool))
	//}
	startEditor(s)
	t,err := template.ParseFiles("html/newPost.html","html/private.nav.html","html/footer.html")
	if err != nil {
		panic(err)
	}
	msg := &Message{}
	msg.Get(w,r)
	commId := r.FormValue("commID")
	if cid,e := strconv.Atoi(commId);commId != "" {
		if e != nil{
			pageNotF(w,r)
			return
		}
		comm := &data.Community{Id: cid}
		e = comm.GetCommunity()
		if e != nil {
			pageNotF(w,r)
			return
		}
	}
	data := struct {
		Msg *Message
		Id string
	}{msg,commId}

	t.ExecuteTemplate(w,"layout",data)
}

func postHandle(w http.ResponseWriter,r *http.Request,p httprouter.Params){
	s,id := sGet(r,"userId")
	if id == nil {
		msg := &Message{Msg: "登录后再做该操作"}
		msg.Send(w)
		http.Redirect(w,r,"/login",303)
		return
	}
	commId,e := strconv.Atoi(r.FormValue("communityId"))
	if e !=nil {
		msg := &Message{Msg: "未知错误:"+e.Error()}
		msg.Send(w)
		http.Redirect(w,r,"/newpost",303)
		return
	}
	post := &data.Post{
		Title: r.FormValue("title"),
		Topic: r.FormValue("topic"),
		CommunityId:commId,
		Content: r.FormValue("content"),
		UserId: id.(int),
	}
	e = post.Put()
	if e != nil {
		msg := &Message{Msg:"发帖错误:"+e.Error()}
		msg.Send(w)
		http.Redirect(w,r,"/newpost",303)
		return
	}
	msg := &Message{Msg: "发帖成功",Alert: "alert-success"}
	msg.Send(w)
	http.Redirect(w,r,"/newpost",303)
	//paths := s.Get("imgs").(map[string]bool)
	//var imgs []string
	//for path,_ := range paths {
	//	imgs = append(imgs,path)
	//}
	//Imgs := data.Imgs{post.Id,imgs}
	//Imgs.Put()
	//s.Set("imgs",nil)
	endEditor(s,post)
}