package data

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

type Post struct {
	Id int `json:"id"`
	Title string `json:"title"`
	Topic string `json:"topic"`
	Content string `json:"content"`
	Date string `json:"date"`
	CommunityId int `json:"community_id"`
	UserId int `json:"user_id"`

	comments Comments
	imgs []string
}

type Posts []*Post

func (p *Post)GetImgs()[]string{
	rows,_ := db.Query("select path from imgs where post_id = ?",p.Id)
	var path string
	for rows.Next() {
		rows.Scan(&path)
		p.imgs = append(p.imgs,path)
	}
	return p.imgs
}

func (p *Post)Imgs()[]string{
	return p.imgs
}

func getPostFromRows(rows *sql.Rows)(p Posts){
	for rows.Next(){
		po := &Post{}
		rows.Scan(&po.Id,&po.Title,&po.Topic,&po.Content,&po.Date,&po.CommunityId,&po.UserId)
		p = append(p,po)
	}
	return
}

func (p *Post)GetPost()(e error){
	row := db.QueryRow("select * from posts where post_id = ?",p.Id)
	if e = row.Err();e != nil {
		return
	}
	e = row.Scan(&p.Id,&p.Title,&p.Topic,&p.Content,&p.Date,&p.CommunityId,&p.UserId)
	return
}

func NumOfPosts()(n int){
	row := db.QueryRow("select count(*) from posts")
	if row.Err() != nil {
		return -1
	}
	row.Scan(&n)
	return
}

func AllPost()(p Posts){
	rows,err := db.Query("select *from posts order by date desc ")
	if err != nil {
		return
	}
	p = getPostFromRows(rows)
	return
}

func (p *Post)Comments()Comments{
	return p.comments
}

func (p *Post)GetComment()(e error){
	rows,e := db.Query("select *from comments where post_id = ? order by date desc ",p.Id)
	if e != nil {
		return
	}
	p.comments = getCommentsFromRows(rows)
	return
}

func (p *Post)Put()(e error){
	st := `insert into posts (title,topic,content,community_id,user_id) values(?,?,?,?,?)`
	res,e := db.Exec(st,p.Title,p.Topic,p.Content,p.CommunityId,p.UserId)
	if e != nil {
		return
	}
	id,e := res.LastInsertId()
	if e != nil {
		return
	}
	p.Id = int(id)
	e = p.GetPost()
	return
}

func (p *Post)Update()(e error){
	st := `update posts set
			title = ?,topic = ?,content = ?,community_id = ?,user_id = ?
			where post_id = ?`;

	_,e = db.Exec(st,p.Title,p.Topic,p.Content,p.CommunityId,p.UserId,p.Id)

	return
}

func GetPosts(l,r int)(p Posts,e error){
	stat,e := db.Prepare("select *from posts order by date desc limit ? offset ?")
	if e != nil {
		return
	}

	rows,e := stat.Query(r-l+1,l)
	if e != nil {
		return
	}
	p = getPostFromRows(rows)
	return
}

func (p *Post)Delete()(e error){
	_,e = db.Exec("delete from posts where post_id = ?",p.Id)
	return
}

func (c Posts)ForEach(handle func(c *Post)){
	for _,v := range c {
		handle(v)
	}
}