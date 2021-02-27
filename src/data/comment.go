package data

import "database/sql"

type Comment struct {
	Id int `json:"id"`
	Date string `json:"date"`
	Content string `json:"content"`
	PostId int `json:"post_id"`
	UserId int `json:"user_id"`
}

type Comments []*Comment

func getCommentsFromRows(rows *sql.Rows)(c Comments){
	for rows.Next(){
		co := &Comment{}
		rows.Scan(&co.Id,&co.Date,&co.Content,&co.PostId,&co.UserId)
		c = append(c,co)
	}
	return
}

func (c *Comment)GetComment()(e error){
	row := db.QueryRow("select * from comments where comment_id = ?",c.Id)
	if e = row.Err();e != nil {
		return
	}
	e = row.Scan(&c.Id,&c.Date,&c.Content,&c.PostId,&c.UserId)
	return
}

func NumOfComments()(n int){
	row := db.QueryRow("select count(*) from comments")
	if row.Err() != nil {
		return -1
	}
	row.Scan(&n)
	return
}

func AllComments()(p Comments){
	rows,err := db.Query("select *from comments order by date desc ")
	if err != nil {
		return
	}
	p = getCommentsFromRows(rows)
	return
}

func (c *Comment)Put()(e error){
	st := `insert into comments (content,post_id,user_id) values(?,?,?)`
	res,e := db.Exec(st,c.Content,c.PostId,c.UserId)
	if e != nil {
		return
	}
	id,e := res.LastInsertId()
	if e != nil {
		return
	}
	c.Id = int(id)
	e = c.GetComment()
	return
}

func (p *Comment)Update()(e error){
	st := `update comments set
			content = ?,post_id =? ,user_id = ?
			where comment_id = ?`;

	_,e = db.Exec(st,p.Content,p.PostId,p.UserId)

	return
}

func (p *Comment)Delete()(e error){
	_,e = db.Exec("delete from comments where comment_id = ?",p.Id)
	return
}

func (c Comments)ForEach(handle func(c *Comment)){
	for _,v := range c {
		handle(v)
	}
}
