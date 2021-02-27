package data

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"os"
)

type User struct {
	Id int `json:"id"`
	Name string `json:"name"`
	Email string `json:"email"`
	password string `json:"password"`
	Date string `json:"date"`

	posts Posts
	cumms Communities
	comments Comments
}

const inf = 0x7FFFFFFF

type Users []*User

var db *sql.DB

func check(e error){
	if e!=nil {
		panic(e)
	}
}

func init()  {
	var err error
	config := struct {
		User string `json:"user"`
		Password string `json:"password"`
		Host string `json:"host"`
		Port string `json:"port"`
		Db string `json:"db"`
	}{}
	file,err := os.Open("data/config.json")
	check(err)

	dec := json.NewDecoder(file)
	err = dec.Decode(&config)
	check(err)

	dbURL := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8", config.User, config.Password, config.Host, config.Port, config.Db)
	db,err = sql.Open("mysql",dbURL)
	check(err)
}

func getUsersFromRows(rows *sql.Rows)(users Users){
	for rows.Next() {
		u := &User{}
		rows.Scan(&u.Id,&u.Name,&u.Email,&u.password,&u.Date)
		users = append(users,u)
	}
	return
}

func (u *User)Password()string{
	return u.password
}

func (u *User)SetPassword(p string){
	u.password = p
}

//to get user by name or user_id
//
//example:
//user := User{Id:1}
//user.GetUser()
func (u *User)GetUser()(err error){
	if u.Id != 0 {
		err = u.getById()
	}else if u.Name != ""{
		err = u.getByName()
	}else{
		err = fmt.Errorf("can not get user")
	}
	return
}

func (u *User)GetUserByEmail()(err error){
	row := db.QueryRow("select * from users where email = ?",u.Email)
	if err = row.Err();err != nil{
		return
	}
	err = row.Scan(&u.Id,&u.Name,&u.Email,&u.password,&u.Date)
	return
}

//return Users that contains all User has the same name
func GetUsersByName(name string)(users Users,err error){
	rows,err := db.Query("select user_id,email,password,date from users where user_name = ?",name)
	if err != nil {
		return
	}
	for rows.Next() {
		u := &User{}
		rows.Scan(&u.Id,&u.Email,&u.password,&u.Date)
		u.Name = name
		users = append(users,u)
	}
	//for _,v := range users{
	//	fmt.Println(v)
	//}
	return
}

func (u *User)getById()(err error){
	row := db.QueryRow("select user_name,email,password,date from users where user_id = ?",u.Id)
	if err = row.Err();err != nil{
		return
	}
	err = row.Scan(&u.Name,&u.Email,&u.password,&u.Date)
	return
}

func (u *User)getByName()(err error){
	users,err := getUsersByName(u.Name)
	if err != nil {
		return
	}
	*u = *users[0];
	return
}

func getUsersByName(name string)(users Users,err error){
	rows,err := db.Query("select user_id,email,password,date from users where user_name = ?",name)
	if err != nil {
		return
	}
	for rows.Next() {
		u := &User{}
		rows.Scan(&u.Id,&u.Email,&u.password,&u.Date)
		u.Name = name
		users = append(users,u)
	}
	//for _,v := range users{
	//	fmt.Println(v)
	//}
	return
}

func NumOfUsers()(n int){
	row := db.QueryRow("select count(*) from users")
	if e:=row.Err();e != nil {
		//fmt.Println(e)
		n = -1
		return
	}
	row.Scan(&n)
	return
}

func AllUsers()(users Users){
	rows,err := db.Query("select *from users order by date desc ")
	if err != nil {
		return
	}
	for rows.Next() {
		u := &User{}
		rows.Scan(&u.Id,&u.Name,&u.Email,&u.password,&u.Date)
		users = append(users,u)
	}
	return
}

//return the communities if GetComm was called and got communities
func (u *User)Communities() Communities{
	return u.cumms
}

//get communities according to the user_id
func (u *User)GetComm()(e error){
	rows,e := db.Query("select * from communities where user_id = ? order by date desc ",u.Id)
	if e != nil {
		return
	}
	for rows.Next() {
		c := &Community{}
		rows.Scan(&c.Id,&c.Name,&c.Desc,&c.Date,&c.BelongTo,&c.UserId)
		u.cumms = append(u.cumms,c)
	}
	return
}

//return the posts if GetPosts was called and got posts
func (u *User)Posts() Posts{
	return u.posts
}

//get communities according to the user_id
func (u *User)GetPosts()(e error){
	rows,e := db.Query("select * from posts where user_id = ? order by date desc ",u.Id)
	if e != nil {
		return
	}
	for rows.Next() {
		c := &Post{}
		rows.Scan(&c.Id,&c.Title,&c.Topic,&c.Content,&c.Date,&c.CommunityId,&c.UserId)
		u.posts = append(u.posts,c)
	}
	return
}

func (u *User)PutUser()(e error){
	st := `insert into users (user_name,email,password) values(?,?,?)`
	res,e := db.Exec(st,u.Name,u.Email,u.password)
	if e != nil {
		return
	}
	id,e := res.LastInsertId()
	if e != nil {
		return
	}
	u.Id = int(id)
	e = u.GetUser()
	return
}

//update the user by user_id
//
//update user_name,email,password
//
//so if you do not want to lose data,
//please use Update()after you use GetUser()
func (u *User)Update()(e error){
	st := `update users set
			user_name = ?,email = ?,password = ?
			where user_id = ?`;

	_,e = db.Exec(st,u.Name,u.Email,u.password,u.Id)

	return
}

func (u *User)Delete()(e error){
	_,e = db.Exec("delete from users where user_id = ?",u.Id)
	return
}

func GetUsers(l,r int)(u Users,e error){
	stat,e := db.Prepare("select *from users order by date desc limit ? offset ?")
	if e != nil {
		return
	}

	rows,e := stat.Query(r-l+1,l)
	if e != nil {
		return
	}
	u = getUsersFromRows(rows)
	return
}

func (u Users)ForEach(handle func(u *User)){
	for _,v := range u {
		handle(v)
	}
}