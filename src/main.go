package main

import (
	"REST"
	_ "github.com/go-sql-driver/mysql"
	"github.com/julienschmidt/httprouter"
	"html/template"
	"net/http"
	"sess"
)

var session *sess.Manager

func init(){
	session = sess.NewManager("miluSess")
}


func sGet(r *http.Request,key string)(*sess.Sess,interface{}){
	s := session.GetSess(r)
	if s != nil {
		return s,s.Get(key)
	}
	return nil,nil
}

func pageNotF(w http.ResponseWriter,r *http.Request){
	msg := &Message{}
	msg.Get(w,r)
	w.WriteHeader(404)
	t, _ := template.ParseFiles("./html/404.html")
	t.Execute(w,msg)
}

func deleteHandle(w http.ResponseWriter,r *http.Request,p httprouter.Params){
	switch p.ByName("selector") {
	case "comment":
		delComment(w,r,p);
	case "post":
		delPost(w,r,p)
	case "community":
		delComm(w,r,p)
	default:
		pageNotF(w,r)
	}
}

func updateHandle(w http.ResponseWriter,r *http.Request,p httprouter.Params){
	switch p.ByName("selector") {
	case "user":
		userUpdate(w,r,p)
	case "post":
		postUpdate(w,r,p)
	case "community":
		commUpdate(w,r,p)
	default:
		w.WriteHeader(403)
	}
}

func updatePage(w http.ResponseWriter,r *http.Request,p httprouter.Params){
	switch p.ByName("selector") {
	case "post":
		postUpdatePage(w,r,p)
	case "community":
		commUpdatePage(w,r,p)
	default:
		w.WriteHeader(403)
	}
}

func main(){
	r := httprouter.New()
	r.NotFound = http.HandlerFunc(pageNotF)
	r.GET("/api/*param",REST.RestHandler)
	r.GET("/",index)
	r.GET("/login",login)
	r.GET("/signup",signUp)
	r.GET("/logout",logout)
	r.GET("/newcommunity", newComm)
	r.GET("/newpost",newPost)
	r.GET("/post/:id",showPost)
	r.GET("/community/:id",showComm)
	r.GET("/home",home)
	r.GET("/home/:path",homeHandle)
	r.GET("/update/:selector/:id",updatePage)

	r.POST("/login",loginHandle)
	r.POST("/signup",signupHandle)
	r.POST("/newcommunity",commHandle)
	r.POST("/newpost",postHandle)
	r.POST("/imgUpload",imgHandle)
	r.POST("/imgDelete",imgDelete)
	r.POST("/newcomment",comHandle)
	r.POST("/delete/:selector/:id",deleteHandle)
	r.POST("/update/:selector",updateHandle)

	r.ServeFiles("/static/*filepath",http.Dir("./"))
	http.ListenAndServe(":8080",r)
}

func index(w http.ResponseWriter,r *http.Request,p httprouter.Params){
	s := session.GetSess(r)
	var nav string
	if s == nil {
		nav = "./html/public.nav.html"
	}else {
		nav = "./html/private.nav.html"
	}
	t,_ := template.ParseFiles("index.html",nav,"./html/footer.html")
	t.ExecuteTemplate(w,"layout",nil)
}









