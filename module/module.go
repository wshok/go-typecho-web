package module

import (
	// "com/log"
	// "goblog/libs/drivers"
	// "goblog/libs/memcache"

	"goblog/helper"
	// "goblog/helper"

	"fmt"
	"flag"
	// "time"
	"html/template"
	// "strconv"
	// "database/sql"

	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
	"github.com/cihub/seelog"
)

type LimitRows struct {
	Start int
	Size  int
}

type Article struct {
	Cid          []byte
	Title        []byte
	Slug         []byte
	Created      string
	Modified     []byte
	Text         template.HTML
	Order        []byte
	AuthorId     []byte   `gorm:"column:authorId"`
	Template     []byte
	Ctype        []byte
	Status       []byte
	Password     []byte
	CommentsNum  []byte
	AllowComment []byte
	AllowPing    []byte
	AllowFeed    []byte
	Parent       []byte
	CategoryName []byte
	CategorySlug []byte
	Author       []byte
}


type Content struct {
	Cid          []byte
	Title        []byte
	Slug         []byte
	Created      string
	Modified     []byte
	Text         []byte
	Order        []byte
	AuthorId     []byte `gorm:"column:authorId"`
	Template     []byte
	Type         []byte
	Status       []byte
	Password     []byte
	CommentsNum  []byte
	AllowComment []byte
	AllowPing    []byte
	AllowFeed    []byte
	Parent       []byte
	CategoryName []byte
	CategorySlug []byte
	Author       []byte
}

type Relationship struct {
	Cid  []byte
	Mid  []byte
}

type Meta struct {
	Mid   []byte  `gorm:"column:mid"`
	Name  []byte
	Slug  []byte
	Type  []byte
	Description  []byte
	Count  []byte
	Order  []byte `gorm:"column:order"`
	Parent []byte
}

type User struct {
	Uid   []byte
	Name   []byte
	Password   []byte
	Mail   []byte
	Url   []byte
	Screenname   []byte
	Created   []byte
	Activated   []byte
	Logged   []byte
	Group   []byte
	Authcode   []byte
}

type Filter struct {
	Page     string
	Category string
	Year     string
	Month    string
}

var (
	// db     drivers.DBdriver
	DB     *gorm.DB
	// loger  *logger.Log
)

const (
	EXPIRE_DAY    int = 86400
	EXPIRE_HOUR   int = 3600
	EXPIRE_MINUTE int = 60
)


func init() {
	// loger = logger.New("logs/err.log")

	// db.Connect("config/database.ini")

	DB, _ = opendb()
}


//#############
//########################## start
// ###############


func opendb() (*gorm.DB, error) {

	configFilePath := flag.String("C", "conf/conf.yaml", "config file path")
	// logConfigPath := flag.String("L", "conf/seelog.xml", "log config file path")
	flag.Parse()

	if err := helper.LoadConfiguration(*configFilePath); err != nil {
		seelog.Critical("err parsing config log file", err)
		return nil, nil
	}

	gorm.DefaultTableNameHandler = func (db *gorm.DB, defaultTableName string) string  {
	    return "typecho_" + defaultTableName;
	}

	db, err := gorm.Open("sqlite3", helper.GetConfiguration().DSN)
	if err != nil {
		return nil, err
	}

	db.SingularTable(false)

	//db.LogMode(true)
	// db.AutoMigrate(&Page{}, &Post{}, &Tag{}, &PostTag{}, &User{}, &Comment{}, &Subscriber{}, &Link{}, &SmmsFile{})
	// db.Model(&PostTag{}).AddUniqueIndex("uk_post_tag", "post_id", "tag_id")

  	// defer db.Close()

	return db, err
}

func ArticleCount(f *Filter) int {
	var (
		val      int
	)

	// DB.Table("contents").Where("status=? AND type=?", "publish", "post").Count(&val)
	// DB.Table("typecho_contents").Select("COUNT(1) AS count").Where("status=? AND type=?", "publish", "post").Row().Scan(&val)
	
	if (f.Page == "index") {
		DB.Model(&Content{}).Where("status=? AND type=?", "publish", "post").Count(&val)

	} else if (f.Page == "category") {

		DB.Model(&Content{}).Where("status=? AND type=? ", "publish", "post").Count(&val)

		DB.Model(&Content{}).Joins("left join typecho_relationships on typecho_relationships.cid = typecho_contents.cid") .Joins("left join typecho_metas on typecho_metas.mid = typecho_relationships.mid") .Where("typecho_contents.status=? AND typecho_contents.type=? AND typecho_metas.type = ?", "publish", "post", "category") .Where("typecho_metas.slug = ?", f.Category) .Count(&val)
	} else if (f.Page == "archive") {

		DB.Model(&Content{}).Where("status=? AND type=?", "publish", "post") .Where("strftime('%Y/%m',datetime(created, 'unixepoch')) = ?", fmt.Sprintf("%s/%s", f.Year, f.Month)).Count(&val)
	}

	return val
}


func ArticleList(l *LimitRows, f *Filter) []Content {
	var (
		val      []Content
	)

	if (f.Page == "index") {
		DB.Model(&Content{}) .Select("typecho_contents.*, typecho_metas.name AS category_name, typecho_metas.slug AS category_slug, typecho_users.screenName as author") .Joins("left join typecho_relationships on typecho_relationships.cid = typecho_contents.cid") .Joins("left join typecho_metas on typecho_metas.mid = typecho_relationships.mid") .Joins("left join typecho_users on typecho_contents.authorId = typecho_users.uid") .Where("typecho_contents.status=? AND typecho_contents.type=? AND typecho_metas.type = ?", "publish", "post", "category") .Order("typecho_contents.cid desc").Offset(l.Start) .Limit(l.Size).Scan(&val)
	} else if (f.Page == "category") {

		DB.Model(&Content{}) .Select("typecho_contents.*, typecho_metas.name AS category_name, typecho_metas.slug AS category_slug, typecho_users.screenName as author") .Joins("left join typecho_relationships on typecho_relationships.cid = typecho_contents.cid") .Joins("left join typecho_metas on typecho_metas.mid = typecho_relationships.mid") .Joins("left join typecho_users on typecho_contents.authorId = typecho_users.uid") .Where("typecho_contents.status=? AND typecho_contents.type=? AND typecho_metas.type = ?", "publish", "post", "category") .Where("typecho_metas.slug = ?", f.Category) .Order("typecho_contents.cid desc").Offset(l.Start) .Limit(l.Size).Scan(&val)
	} else if (f.Page == "archive") {

		DB.Model(&Content{}) .Select("typecho_contents.*, typecho_metas.name AS category_name, typecho_metas.slug AS category_slug, typecho_users.screenName as author") .Joins("left join typecho_relationships on typecho_relationships.cid = typecho_contents.cid") .Joins("left join typecho_metas on typecho_metas.mid = typecho_relationships.mid") .Joins("left join typecho_users on typecho_contents.authorId = typecho_users.uid") .Where("typecho_contents.status=? AND typecho_contents.type=? AND typecho_metas.type = ?", "publish", "post", "category") .Where("strftime('%Y/%m',datetime(typecho_contents.created, 'unixepoch')) = ?", fmt.Sprintf("%s/%s", f.Year, f.Month)) .Order("typecho_contents.cid desc").Offset(l.Start).Limit(l.Size).Scan(&val)
	}

	return val
}


func LeastPosted() []Content {
	var (
		val      []Content
	)

	DB.Model(&Content{}).Select("cid,title,slug").Where("status=? AND type=?", "publish", "post").Order("created desc").Offset(0).Limit(8).Scan(&val)

	return val
}


func Category() []Meta {
	var (
		val      []Meta
	)

	DB.Table("typecho_metas").Where("type=?", "category").Order("order").Scan(&val)

	return val
}

type Result struct {
    Yearmonth string
    Count  int
}

func Archive() []Result{

	var val []Result

	DB.Table("typecho_contents").Select("strftime('%Y/%m',datetime(created, 'unixepoch')) AS yearmonth, COUNT(1) AS count") .Where("status=? AND type=?", "publish", "post") .Group("yearmonth") .Order("created desc") .Scan(&val)
		
	return val
}

func Detail(url string) Content {
	var val Content

    DB.Model(&Content{}).Select("typecho_contents.*, typecho_metas.name AS category_name, typecho_metas.slug AS category_slug, typecho_users.screenName as author") .Joins("left join typecho_relationships on typecho_relationships.cid = typecho_contents.cid") .Joins("left join typecho_metas on typecho_metas.mid = typecho_relationships.mid") .Joins("left join typecho_users on typecho_contents.authorId = typecho_users.uid") .Where("typecho_contents.status=? AND typecho_contents.type=? AND typecho_metas.type = ? and typecho_contents.slug = ?", "publish", "post", "category", url) .Scan(&val)

	return val
}

func Page(url string) Content {
	var val Content

    DB.Model(&Content{}).Where("status=? AND type=? and slug = ?", "publish", "page", url) .Scan(&val)

	return val
}


/////
/////////////////////////////////////////////////////  end
/////



// func Options() map[string][]byte {
// 	var (
// 		err      error
// 		val             = make(map[string][]byte)
// 	)


// 	db.Orm.SetTable(db.TableName("options"))
// 	db.Orm.Select("`name`, `user`, `value`")
// 	rows := db.QueryRows()
// 	for _, row := range rows {
// 		val[string(row["name"])] = row["value"]
// 	}

// 	return val
// }



// func About() Article {
// 	var (
// 		err      error
// 		val      Article
// 	)


// 	db.Orm.SetTable(db.TableName("contents"))
// 	db.Orm.Where("slug=? AND type=? AND status=?", "about", "page", "publish")
// 	db.Orm.Limit(1)

// 	row := db.QueryOneRow()

// 	var created int64
// 	created, _ = strconv.ParseInt(string(row["created"]), 10, 64)

// 	val = Article{
// 		Cid:          row["cid"],
// 		Title:        row["title"],
// 		Slug:         row["slug"],
// 		Created:      time.Unix(created, 0).Format("January 02,2006"),
// 		Modified:     row["modified"],
// 		Text:         template.HTML(string(row["text"])),
// 		Order:        row["order"],
// 		AuthorId:     row["authorId"],
// 		Template:     row["template"],
// 		Ctype:        row["type"],
// 		Status:       row["status"],
// 		Password:     row["password"],
// 		CommentsNum:  row["commentsNum"],
// 		AllowComment: row["allowComment"],
// 		AllowPing:    row["allowPing"],
// 		AllowFeed:    row["allowFeed"],
// 		Parent:       row["parent"],
// 	}

// 	return val
// }

// func ArtCount(year, month string) int {
// 	var (
// 		err      error

// 		val      int
// 	)


// 	db.Orm.SetTable(db.TableName("contents"))
// 	db.Orm.Select("COUNT(1) AS total")
// 	db.Orm.Where("status=? AND type=? AND FROM_UNIXTIME(created, '%Y-%m')=?", "publish", "post",
// 		fmt.Sprintf("%s-%s", year, month))

// 	query := db.QueryOneRow()
// 	val, _ = strconv.Atoi(string(query["total"]))


// 	return val
// }


// func ArtList(year, month string, l *LimitRows) []Article {
// 	var (
// 		err      error
// 		val      []Article
// 	)


// 	db.Orm.SetTable(db.TableName("contents") + " AS c")
// 	db.Orm.Select("c.*, m.name AS cat_name, m.slug AS cat_slug, u.screenName")
// 	db.Orm.Join("LEFT", db.TableName("relationships")+" AS r", "c.cid = r.cid")
// 	db.Orm.Join("LEFT", db.TableName("metas")+" AS m", "m.mid = r.mid")
// 	db.Orm.Join("LEFT", db.TableName("users")+" AS u", "c.authorId = u.uid")
// 	db.Orm.Where("c.type =? AND c.status =? AND m.type=? AND FROM_UNIXTIME(c.created, '%Y-%m')=?",
// 		"post", "publish", "category", fmt.Sprintf("%s-%s", year, month))
// 	db.Orm.OrderBy("c.cid desc")
// 	db.Orm.Limit(l.Size, l.Start)

// 	query := db.QueryRows()
// 	for _, row := range query {
// 		var created int64
// 		created, _ = strconv.ParseInt(string(row["created"]), 10, 64)
// 		art := Article{
// 			Cid:          row["cid"],
// 			Title:        row["title"],
// 			Slug:         row["slug"],
// 			Created:      time.Unix(created, 0).Format("January 02,2006"),
// 			Modified:     row["modified"],
// 			Text:         template.HTML(string(helper.ReadMore(row["text"]))),
// 			Order:        row["order"],
// 			AuthorId:     row["authorId"],
// 			Template:     row["template"],
// 			Ctype:        row["type"],
// 			Status:       row["status"],
// 			Password:     row["password"],
// 			CommentsNum:  row["commentsNum"],
// 			AllowComment: row["allowComment"],
// 			AllowPing:    row["allowPing"],
// 			AllowFeed:    row["allowFeed"],
// 			Parent:       row["parent"],
// 			CategoryName: row["cat_name"],
// 			CategorySlug: row["cat_slug"],
// 			Author:       row["screenName"],
// 		}
// 		val = append(val, art)
// 	}

// 	return val
// }


// func ArticleInfo(article string) Article {
// 	var (
// 		err      error

// 		val      Article
// 	)



// 	db.Orm.SetTable(db.TableName("metas") + " AS m")
// 	db.Orm.Select("c.*, m.name AS cat_name, m.slug AS cat_slug, u.screenName")
// 	db.Orm.Join("LEFT", db.TableName("relationships")+" AS r", "m.mid = r.mid")
// 	db.Orm.Join("LEFT", db.TableName("contents")+" AS c", "c.cid = r.cid AND c.status = 'publish'")
// 	db.Orm.Join("LEFT", db.TableName("users")+" AS u", "c.authorId = u.uid")
// 	db.Orm.Where("c.status=? AND c.slug=?", "publish", article)

// 	row := db.QueryOneRow()

// 	var created int64
// 	created, _ = strconv.ParseInt(string(row["created"]), 10, 64)

// 	val = Article{
// 		Cid:          row["cid"],
// 		Title:        row["title"],
// 		Slug:         row["slug"],
// 		Created:      time.Unix(created, 0).Format("January 02,2006"),
// 		Modified:     row["modified"],
// 		Text:         template.HTML(string(row["text"])),
// 		Order:        row["order"],
// 		AuthorId:     row["authorId"],
// 		Template:     row["template"],
// 		Ctype:        row["type"],
// 		Status:       row["status"],
// 		Password:     row["password"],
// 		CommentsNum:  row["commentsNum"],
// 		AllowComment: row["allowComment"],
// 		AllowPing:    row["allowPing"],
// 		AllowFeed:    row["allowFeed"],
// 		Parent:       row["parent"],
// 		CategoryName: row["cat_name"],
// 		CategorySlug: row["cat_slug"],
// 		Author:       row["screenName"],
// 	}


// 	return val
// }


// func CateArtCount(cate string) int {
// 	var (
// 		err      error
// 		val      int
// 	)


// 	db.Orm.SetTable(db.TableName("metas") + " AS m")
// 	db.Orm.Select("COUNT(1) AS count")
// 	db.Orm.Join("LEFT", db.TableName("relationships")+" AS r", "m.mid = r.mid")
// 	db.Orm.Join("LEFT", db.TableName("contents")+" AS c", "c.cid = r.cid AND c.status = 'publish'")
// 	db.Orm.Where("m.slug = ? AND m.type =? ", cate, "category")
// 	db.Orm.OrderBy("c.cid desc")

// 	query := db.QueryOneRow()
// 	val, _ = strconv.Atoi(string(query["count"]))

// 	return val
// }


// func CateArtList(cate string, l *LimitRows) []Article {
// 	var (
// 		err      error
// 		val      []Article
// 	)


// 	db.Orm.SetTable(db.TableName("contents") + " AS c")
// 	db.Orm.Select("c.*, m.name AS cat_name, m.slug AS cat_slug, u.screenName")
// 	db.Orm.Join("LEFT", db.TableName("relationships")+" AS r", "c.cid = r.cid")
// 	db.Orm.Join("LEFT", db.TableName("metas")+" AS m", "m.mid = r.mid")
// 	db.Orm.Join("LEFT", db.TableName("users")+" AS u", "c.authorId = u.uid")
// 	db.Orm.Where("m.slug = ? AND m.type=? AND m.type =? AND c.status =? ", cate, "category", "category", "publish")
// 	db.Orm.OrderBy("c.cid desc")
// 	db.Orm.Limit(l.Size, l.Start)

// 	query := db.QueryRows()
// 	for _, row := range query {
// 		var created int64
// 		created, _ = strconv.ParseInt(string(row["created"]), 10, 64)
// 		art := Article{
// 			Cid:          row["cid"],
// 			Title:        row["title"],
// 			Slug:         row["slug"],
// 			Created:      time.Unix(created, 0).Format("January 02,2006"),
// 			Modified:     row["modified"],
// 			Text:         template.HTML(string(helper.ReadMore(row["text"]))),
// 			Order:        row["order"],
// 			AuthorId:     row["authorId"],
// 			Template:     row["template"],
// 			Ctype:        row["type"],
// 			Status:       row["status"],
// 			Password:     row["password"],
// 			CommentsNum:  row["commentsNum"],
// 			AllowComment: row["allowComment"],
// 			AllowPing:    row["allowPing"],
// 			AllowFeed:    row["allowFeed"],
// 			Parent:       row["parent"],
// 			CategoryName: row["cat_name"],
// 			CategorySlug: row["cat_slug"],
// 			Author:       row["screenName"],
// 		}
// 		val = append(val, art)
// 	}


// 	return val
// }


// func SearchArtCount(keyword string) int {
// 	var (
// 		err      error

// 		val      int
// 	)


// 	db.Orm.SetTable(db.TableName("contents"))
// 	db.Orm.Select("COUNT(1) AS total")
// 	db.Orm.Where("status=? AND type=? AND (title LIKE ? OR text LIKE ? ) ",
// 		"publish", "post", fmt.Sprintf("%%%s%%", keyword), fmt.Sprintf("%%%s%%", keyword))

// 	query := db.QueryOneRow()
// 	val, _ = strconv.Atoi(string(query["total"]))


// 	return val
// }


// func SearchArtList(keyword string, l *LimitRows) []Article {
// 	var (
// 		err      error
// 		update   bool
// 		cacheKey string = fmt.Sprintf("search-keyword:%s-page-%d-%d", keyword, l.Start, l.Start+l.Size)
// 		val      []Article
// 	)


// 	db.Orm.SetTable(db.TableName("contents") + " AS c")
// 	db.Orm.Select("c.*, m.name AS cat_name, m.slug AS cat_slug, u.screenName")
// 	db.Orm.Join("LEFT", db.TableName("relationships")+" AS r", "c.cid = r.cid")
// 	db.Orm.Join("LEFT", db.TableName("metas")+" AS m", "m.mid = r.mid")
// 	db.Orm.Join("LEFT", db.TableName("users")+" AS u", "c.authorId = u.uid")
// 	db.Orm.Where("c.status=? AND c.type=? AND m.type =? AND (c.title LIKE ? OR c.text LIKE ? ) ",
// 		"publish", "post", "category", fmt.Sprintf("%%%s%%", keyword), fmt.Sprintf("%%%s%%", keyword))
// 	db.Orm.OrderBy("c.cid desc")
// 	db.Orm.Limit(l.Size, l.Start)

// 	query := db.QueryRows()
// 	for _, row := range query {
// 		var created int64
// 		created, _ = strconv.ParseInt(string(row["created"]), 10, 64)
// 		art := Article{
// 			Cid:          row["cid"],
// 			Title:        row["title"],
// 			Slug:         row["slug"],
// 			Created:      time.Unix(created, 0).Format("January 02,2006"),
// 			Modified:     row["modified"],
// 			Text:         template.HTML(string(helper.ReadMore(row["text"]))),
// 			Order:        row["order"],
// 			AuthorId:     row["authorId"],
// 			Template:     row["template"],
// 			Ctype:        row["type"],
// 			Status:       row["status"],
// 			Password:     row["password"],
// 			CommentsNum:  row["commentsNum"],
// 			AllowComment: row["allowComment"],
// 			AllowPing:    row["allowPing"],
// 			AllowFeed:    row["allowFeed"],
// 			Parent:       row["parent"],
// 			CategoryName: row["cat_name"],
// 			CategorySlug: row["cat_slug"],
// 			Author:       row["screenName"],
// 		}
// 		val = append(val, art)
// 	}


// 	return val
// }



// func errlog(txt string) {
// 	// loger.Out(txt)
// }
