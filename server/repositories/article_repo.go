package repositories

import (
	"log"

	"github.com/lib/pq"

	"github.com/jmoiron/sqlx"
	"github.com/neoreads/backend/server/models"
)

// ArticleRepo Article related data repository
type ArticleRepo struct {
	db      *sqlx.DB
	rootDir string
}

// NewArticleRepo creator for ArticleRepo
func NewArticleRepo(db *sqlx.DB) *ArticleRepo {
	return &ArticleRepo{db: db}
}

func (r *ArticleRepo) AddArticle(a *models.Article) bool {
	var pid string
	// TODO: get PID from jwt credentials
	err := r.db.Get(&pid, "SELECT pid from users where username = $1", a.PID)
	if err != nil {
		log.Printf("error adding article:%v, with err: %v\n", a, err)
		return false
	}
	a.PID = pid
	_, err = r.db.NamedExec("INSERT INTO articles (id, title, content)"+
		" VALUES (:id, :title, :content)", a)
	if err != nil {
		log.Printf("error adding article:%v, with err: %v\n", a, err)
		return false
	}

	_, err = r.db.Exec("INSERT INTO article_people (aid, pid) VALUES ($1, $2)", a.ID, a.PID)
	if err != nil {
		log.Printf("error adding article:%v, with err: %v\n", a, err)
		return false
	}
	return true
}

func (r *ArticleRepo) ModifyArticle(a *models.Article) bool {
	// TODO: add support for modTime
	_, err := r.db.NamedExec("UPDATE articles set title = :title, content = :content, modtime = now() where id = :id", a)
	if err != nil {
		log.Printf("error modifying article:%v, with err: %v\n", a, err)
		return false
	}
	return true
}

func (r *ArticleRepo) GetArticle(artid string) models.Article {
	var article models.Article
	err := r.db.Get(&article, "SELECT a.title, a.content, a.id, p.pid from articles a, article_people p where a.id = p.aid and a.id = $1", artid)
	if err != nil {
		log.Printf("error listing articles from db:%v, with err:%v\n", artid, err)
	}
	return article
}

func (r *ArticleRepo) ListArticles(username string) []models.Article {
	articles := []models.Article{}
	err := r.db.Select(&articles, "SELECT a.*, p.pid from articles a, article_people p where a.id = p.aid and p.pid = (SELECT pid from users where username = $1) order by a.modtime desc", username)
	if err != nil {
		log.Printf("error listing articles from db:%v, with err:%v\n", username, err)
	}
	return articles
}

func (r *ArticleRepo) ListArticlesInCollection(username string, colid string) []models.Article {
	var artids []string
	err := r.db.Select(&artids, "SELECT artid from collections_articles where colid = $1", colid)
	if err != nil {
		log.Printf("error listing articles for collection from db:%v, with err:%v\n", colid, err)
	}

	// TODO: check if using pq.Array and ANY(?) is the best option here
	articles := []models.Article{}
	err = r.db.Select(&articles, "SELECT a.id, a.title, a.modtime, p.pid from articles a, article_people p where a.id = p.aid and p.pid = (SELECT pid from users where username = $1) and a.id = ANY($2) order by a.modtime desc", username, pq.Array(artids))
	if err != nil {
		log.Printf("error listing articles from db:%v, with err:%v\n", username, err)
	}
	return articles
}

func (r *ArticleRepo) RemoveArticle(artid string) bool {
	tx := r.db.MustBegin()
	_, err := tx.Exec("DELETE FROM articles where id = $1", artid)
	if err != nil {
		log.Printf("error removing aritcle from db:%v, with err:%v\n", artid, err)
		return false
	}
	_, err = tx.Exec("DELETE FROM article_people where aid = $1", artid)
	if err != nil {
		log.Printf("error removing artlce_people relation from db:%v, with err:%v\n", artid, err)
		return false
	}
	tx.Commit()
	return true
}
