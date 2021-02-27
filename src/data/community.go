package data

import "database/sql"

type Community struct {
	Id int `json:"id"`
	Name string `json:"name"`
	Desc string `json:"description"`
	Date string `json:"date"`
	BelongTo string `json:"belong_to"`
	UserId int `json:"user_id"`

	posts Posts
}

type Communities []*Community

func getCommsFromRows(rows *sql.Rows)(c Communities){
	for rows.Next() {
		co := &Community{}
		rows.Scan(&co.Id,&co.Name,&co.Desc,&co.Date,&co.BelongTo,&co.UserId)
		c = append(c,co)
	}
	return
}

func (c *Community)GetCommunity()(e error){
	row := db.QueryRow("select * from communities where community_id = ?",c.Id)
	if e = row.Err();e != nil {
		return
	}
	e = row.Scan(&c.Id,&c.Name,&c.Desc,&c.Date,&c.BelongTo,&c.UserId)
	return
}

func GetCommsByName(name string)(c Communities,e error){
	rows,e := db.Query("select * from communities where name = ?",name)
	if e != nil {
		return
	}
	for rows.Next() {
		co := &Community{}
		rows.Scan(&co.Id,&co.Name,&co.Desc,&co.Date,&co.BelongTo,&co.UserId)
		c = append(c,co)
	}
	return
}

func NumOfComms()(n int){
	row := db.QueryRow("select count(*) from communities")
	if row.Err() != nil {
		return -1
	}
	row.Scan(&n)
	return
}

func AllComms()(c Communities){
	rows,err := db.Query("select *from communities order by date desc ")
	if err != nil {
		return
	}
	for rows.Next() {
		co := &Community{}
		rows.Scan(&co.Id,&co.Name,&co.Desc,&co.Date,&co.BelongTo,&co.UserId)
		c = append(c,co)
	}
	return
}

func (c *Community)Posts()Posts{
	return c.posts
}

func (c *Community)GetPost()(e error){
	rows,e := db.Query("select *from posts where community_id = ? order by date desc ",c.Id)
	if e != nil {
		return
	}
	c.posts = getPostFromRows(rows)
	return
}

func (c *Community)Put()(e error){
	st := `insert into communities (name,description,belong_to,user_id) values(?,?,?,?)`
	res,e := db.Exec(st,c.Name,c.Desc,c.BelongTo,c.UserId)
	if e != nil {
		return
	}
	id,e := res.LastInsertId()
	if e != nil {
		return
	}
	c.Id = int(id)
	e = c.GetCommunity()
	return
}

func (u *Community)Update()(e error){
	st := `update communities set
			name=?,description = ?,belong_to = ?,user_id = ?
			where community_id = ?`;

	_,e = db.Exec(st,u.Name,u.Desc,u.BelongTo,u.UserId,u.Id)

	return
}

func GetCommunities(l,r int)(u Communities,e error){
	stat,e := db.Prepare("select *from communities order by date desc limit ? offset ?")
	if e != nil {
		return
	}

	rows,e := stat.Query(r-l+1,l)
	if e != nil {
		return
	}
	u = getCommsFromRows(rows)
	return
}

func (u *Community)Delete()(e error){
	_,e = db.Exec("delete from communities where community_id = ?",u.Id)
	return
}

func (c Communities)ForEach(handle func(c *Community)){
	for _,v := range c {
		handle(v)
	}
}
