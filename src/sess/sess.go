package sess

import (
	"crypto/rand"
	"encoding/base64"
	"log"
	"net/http"
	"sync"
	"time"
)

type Sess struct {
	sid string
	v map[interface{}]interface{}
	lifeTime int64
	timeAccessed int64
	m sync.Mutex
	afterFunc func (sess *Sess)
}

type Manager struct {
	pool map[string]*Sess
	m sync.Mutex
	cookieName string
}

func sessId()string{
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return ""
	}
	return base64.URLEncoding.EncodeToString(b)
}

func NewManager(cookie string)*Manager{
	m := &Manager{pool:make(map[string]*Sess),cookieName:cookie}
	return m
}

func (m *Manager)NewSess(w http.ResponseWriter,r *http.Request,lifeTime time.Duration)*Sess{
	s := m.GetSess(r)
	if s != nil {//delete the session already exists
		m.DelSess(s.sid)
	}
	s = &Sess{sid:sessId(),v:make(map[interface{}]interface{})}
	m.pool[s.sid] = s
	c := &http.Cookie{
		Name:m.cookieName,
		HttpOnly:true,
		Value:s.sid,
	}
	if lifeTime == -1 {
		s.lifeTime = int64(3600 * time.Second) //默认一小时
	}else{
		s.lifeTime = int64(lifeTime)
		c.MaxAge = int(lifeTime / time.Second)
	}
	http.SetCookie(w,c)
	m.gc(s.sid)
	return s
}

func (m *Manager)GetSess(r *http.Request)*Sess{
	c,e := r.Cookie(m.cookieName)
	if e != nil {
		return nil
	}
	s,ok := m.pool[c.Value]
	if ok{
		m.gc(s.sid)
		return s
	}
	return nil
}

func (m *Manager)DelSess(sid string){
	s := m.pool[sid]
	if s != nil {
		if s.afterFunc != nil {
			s.afterFunc(s)
		}
		m.m.Lock()
		delete(m.pool,sid)
		m.m.Unlock()
	}
}

func (m *Manager)DestroySess(w http.ResponseWriter,r *http.Request){
	s := m.GetSess(r)
	if s != nil {
		c := &http.Cookie{
			Name:     m.cookieName,
			HttpOnly: true,
			MaxAge:   -1,
		}
		http.SetCookie(w, c)
		m.DelSess(s.sid)
		log.Println("destroySess",s.sid)
	}
}

func (m *Manager)gc(sid string){
	s,ok := m.pool[sid]
	if ok {
		s.update()
		time.AfterFunc(time.Duration(s.lifeTime + 10), func() {
			//log.Printf("%v  %v %v %v address %p\n",s.timeAccessed,s.lifeTime/int64(time.Second),s.timeAccessed+s.lifeTime/int64(time.Second),time.Now().Unix(),s)
			if s.timeAccessed + s.lifeTime/int64(time.Second) <= time.Now().Unix() {
				m.DelSess(sid)
			}
		})
	}
}

func (s *Sess)update(){
	s.m.Lock()
	s.timeAccessed = time.Now().Unix()
	//fmt.Printf("set timeacce %v address %p\n",s.timeAccessed,s)
	s.m.Unlock()
}


func (s *Sess)Set(k,v interface{}){
	s.v[k] = v
}

func (s *Sess)Get(k interface{})interface{}{
	return s.v[k]
}

func (s *Sess)Del(k interface{}){
	delete(s.v,k)
}

func (s *Sess)Id()string{
	return s.sid
}

func (s *Sess)SetAfterFunc(f func (sess * Sess)){
	s.afterFunc = f
}

