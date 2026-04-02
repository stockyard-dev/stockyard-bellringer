package store
import ("database/sql";"fmt";"os";"path/filepath";"time";_ "modernc.org/sqlite")
type DB struct{db *sql.DB}
type Notification struct {
	ID string `json:"id"`
	Title string `json:"title"`
	Body string `json:"body"`
	URL string `json:"url"`
	Icon string `json:"icon"`
	Channel string `json:"channel"`
	SentCount int `json:"sent_count"`
	ClickCount int `json:"click_count"`
	CreatedAt string `json:"created_at"`
}
func Open(d string)(*DB,error){if err:=os.MkdirAll(d,0755);err!=nil{return nil,err};db,err:=sql.Open("sqlite",filepath.Join(d,"bellringer.db")+"?_journal_mode=WAL&_busy_timeout=5000");if err!=nil{return nil,err}
db.Exec(`CREATE TABLE IF NOT EXISTS notifications(id TEXT PRIMARY KEY,title TEXT NOT NULL,body TEXT DEFAULT '',url TEXT DEFAULT '',icon TEXT DEFAULT '',channel TEXT DEFAULT 'default',sent_count INTEGER DEFAULT 0,click_count INTEGER DEFAULT 0,created_at TEXT DEFAULT(datetime('now')))`)
return &DB{db:db},nil}
func(d *DB)Close()error{return d.db.Close()}
func genID()string{return fmt.Sprintf("%d",time.Now().UnixNano())}
func now()string{return time.Now().UTC().Format(time.RFC3339)}
func(d *DB)Create(e *Notification)error{e.ID=genID();e.CreatedAt=now();_,err:=d.db.Exec(`INSERT INTO notifications(id,title,body,url,icon,channel,sent_count,click_count,created_at)VALUES(?,?,?,?,?,?,?,?,?)`,e.ID,e.Title,e.Body,e.URL,e.Icon,e.Channel,e.SentCount,e.ClickCount,e.CreatedAt);return err}
func(d *DB)Get(id string)*Notification{var e Notification;if d.db.QueryRow(`SELECT id,title,body,url,icon,channel,sent_count,click_count,created_at FROM notifications WHERE id=?`,id).Scan(&e.ID,&e.Title,&e.Body,&e.URL,&e.Icon,&e.Channel,&e.SentCount,&e.ClickCount,&e.CreatedAt)!=nil{return nil};return &e}
func(d *DB)List()[]Notification{rows,_:=d.db.Query(`SELECT id,title,body,url,icon,channel,sent_count,click_count,created_at FROM notifications ORDER BY created_at DESC`);if rows==nil{return nil};defer rows.Close();var o []Notification;for rows.Next(){var e Notification;rows.Scan(&e.ID,&e.Title,&e.Body,&e.URL,&e.Icon,&e.Channel,&e.SentCount,&e.ClickCount,&e.CreatedAt);o=append(o,e)};return o}
func(d *DB)Delete(id string)error{_,err:=d.db.Exec(`DELETE FROM notifications WHERE id=?`,id);return err}
func(d *DB)Count()int{var n int;d.db.QueryRow(`SELECT COUNT(*) FROM notifications`).Scan(&n);return n}
