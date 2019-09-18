package repositories

import (
	"encoding/json"
	"log"

	"github.com/jmoiron/sqlx"
	"github.com/neoreads/backend/server/models"
)

// NewsRepo
type NewsRepo struct {
	db *sqlx.DB
}

// NewsRepo creator
func NewNewsRepo(db *sqlx.DB) *NewsRepo {
	return &NewsRepo{db: db}
}

func (r *NewsRepo) AddNews(news *models.News) bool {
	tx, err := r.db.Beginx()
	if err != nil {
		log.Printf("error adding news item, can't create transaction: %v\n ", err)
		return false
	}
	_, err = tx.NamedExec("INSERT INTO news (id, kind, link, source, title, summary, content) VALUES (:id, :kind, :link, :source, :title, :summary, :content)", news)
	if err != nil {
		log.Printf("error adding news item, insert error: %v\n ", err)
		return false
	}
	// add tags if neccessary
	tags := news.Tags
	for i := range tags {
		tag := tags[i]
		_, err = tx.NamedExec("INSERT INTO tags (id, kind, tag) VALUES (:id, :kind, :tag) ON CONFLICT (id) DO UPDATE SET kind = :kind, tag = :tag", tag)
		if err != nil {
			log.Printf("Error upserting tag : %#v, with error:%v\n", tag, err)
			return false
		}
		if tag.ID == "" {
			err = tx.Get(&tag.ID, "SELECT ID from tags where kind = $1 and tag = $2", tag.Kind, tag.Tag)
			if err != nil {
				log.Printf("Can not get inserted tag id for %#v\n", tag)
				return false
			}
		}
		_, err = tx.Exec("INSERT INTO news_tags (newsid, tagid) VALUES ($1, $2) ON CONFLICT ON CONSTRAINT news_tags_pkey DO NOTHING", news.ID, tag.ID)
		if err != nil {
			log.Printf("Error upserting news_tags : %#v, %#v, with error:%v\n", news, tag, err)
			return false
		}
	}
	tx.Commit()
	return true
}

func (r *NewsRepo) ListNews() []models.News {
	var news []models.News
	sql := `select n.*, array_to_json(array_remove(array_agg(row_to_json(t.*)::text), null)::json[]) as tagsjson from news n 
	left join news_tags nt on n.id = nt.newsid
	left join tags t on t.id = nt.tagid group by n.id order by n.modtime desc;`
	err := r.db.Select(&news, sql)
	if err != nil {
		log.Printf("Error listing news: %v\n", err)
	}
	//  convert tagsjson to tags
	for i := range news {
		item := &news[i]
		log.Printf("checking tagsjson:%v\n", item.TagsJSON)
		if item.TagsJSON != "" {
			err = json.Unmarshal([]byte(item.TagsJSON), &item.Tags)
			if err != nil {
				log.Printf("Error unmarshaling tagsjson:%v, with error: %v\n", item.TagsJSON, err)
			}
			item.TagsJSON = ""
		}
	}
	return news
}
