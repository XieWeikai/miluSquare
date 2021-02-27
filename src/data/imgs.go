package data

type Imgs struct {
	PostId int
	Imgs []string
}

func (i *Imgs)Put(){
	for _,path := range i.Imgs {
		db.Exec("insert into imgs (path,post_id) values(?,?) ",path,i.PostId)
	}
}
