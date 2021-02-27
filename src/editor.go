package main

import (
	"crypto/rand"
	"data"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"io"
	"log"
	"net/http"
	"os"
	"sess"
	"strconv"
)

func imgHandle(w http.ResponseWriter,r *http.Request,p httprouter.Params){
	s,id := sGet(r,"userId")
	if id == nil {
		log.Println(r.RemoteAddr,"does not login")
		return //未登录行为，不予理睬
	}
	dir := "userdata/"+strconv.Itoa(id.(int))
	os.MkdirAll(dir,os.ModePerm)
	name := randomName()
	f,e := os.Create(dir+"/"+name)
	if e != nil {
		log.Println("imgCreateErr:",e)
	}
	r.ParseMultipartForm(4096)
	file,_,e := r.FormFile("file")
	if e != nil{
		log.Println("imgHandleErr:",e)
	}
	io.Copy(f,file)
	fmt.Fprintf(w,`{"link":"/static/%s/%s"}`,dir,name)
	log.Printf(`%s uploadimg: {"link":"/static/%s/%s"}`,r.RemoteAddr,dir,name)
	m := s.Get("imgs").(map[string]bool)
	m[dir+"/"+name] = true
}

func imgDelete(w http.ResponseWriter,r *http.Request,p httprouter.Params) {
	s := session.GetSess(r)
	if s == nil {
		return
	}
	src := struct{
		Src string `json:"src"`
	}{}
	dec := json.NewDecoder(r.Body)
	err := dec.Decode(&src)
	if err != nil {
		log.Println("imgDelDecErr:  ",err)
	}
	//fmt.Println("deleteimg: ",src.Src[8:])
	os.Remove(src.Src[8:])
	imgs := s.Get("imgs").(map[string]bool)
	delete(imgs,src.Src[8:])

}

func randomName()string{
	b := make([]byte, 64)
	if _, err := rand.Read(b); err != nil {
		return ""
	}
	return base64.URLEncoding.EncodeToString(b)
}

func startEditor(s *sess.Sess){
	s.SetAfterFunc(func (s *sess.Sess){
		if imgs := s.Get("imgs");imgs != nil {
			img := imgs.(map[string]bool)
			for path,_ := range img {
				log.Println("deletePath",path)
				os.Remove(path)
			}
		}
	})
	if s.Get("imgs") == nil {
		s.Set("imgs", make(map[string]bool))
	}
}

func endEditor(s *sess.Sess,post *data.Post){
	paths := s.Get("imgs").(map[string]bool)
	var imgs []string
	for path,_ := range paths {
		imgs = append(imgs,path)
	}
	Imgs := data.Imgs{post.Id,imgs}
	Imgs.Put()
	s.Set("imgs",nil)
}
