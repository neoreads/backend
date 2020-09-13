package repositories

import (
	"log"
	"strings"

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
	_, err := r.db.NamedExec("INSERT INTO articles (id, kind, title, content)"+
		" VALUES (:id, :kind, :title, :content)", a)
	if err != nil {
		log.Printf("error adding article:%v, with err: %v\n", a, err)
		return false
	}

	pids := a.PID
	if strings.Contains(pids, ",") {
		pidarr := strings.Split(pids, ",")
		for i := range pidarr {
			pid := pidarr[i]
			_, err = r.db.Exec("INSERT INTO article_people (aid, pid) VALUES ($1, $2)", a.ID, pid)
			if err != nil {
				log.Printf("error adding article:%v, with err: %v\n", a, err)
				return false
			}
		}

	} else {
		_, err = r.db.Exec("INSERT INTO article_people (aid, pid) VALUES ($1, $2)", a.ID, a.PID)
		if err != nil {
			log.Printf("error adding article:%v, with err: %v\n", a, err)
			return false
		}
	}
	return true
}

func (r *ArticleRepo) ModifyArticle(a *models.Article) bool {
	// TODO: add support for modTime
	_, err := r.db.NamedExec("UPDATE articles set kind = :kind, title = :title, content = :content, modtime = now() where id = :id", a)
	if err != nil {
		log.Printf("error modifying article:%v, with err: %v\n", a, err)
		return false
	}
	return true
}

func (r *ArticleRepo) GetArticle(artid string) models.Article {
	var article models.Article
	err := r.db.Get(&article, "SELECT a.kind, a.title, a.content, a.id, ap.pid, p.fullname as author from articles a, article_people ap, people p where a.id = ap.aid and ap.pid = p.id and a.id = $1", artid)
	if err != nil {
		log.Printf("error listing articles from db:%v, with err:%v\n", artid, err)
	}
	return article
}

func (r *ArticleRepo) SearchArticles(kind models.ArticleKind, pid string) []models.Article {
	articles := []models.Article{}
	if pid != "" {
		err := r.db.Select(&articles, "SELECT a.*, ap.pid, p.fullname as author from articles a, article_people ap, people p where a.id = ap.aid and ap.pid = p.id and a.kind = $1 and ap.pid = $2 order by a.modtime desc", kind, pid)
		if err != nil {
			log.Printf("error searching articles from db:%v:%v, with err:%v\n", kind, pid, err)
		}
	} else {
		err := r.db.Select(&articles, "SELECT a.*, ap.pid, p.fullname as author from articles a, article_people ap, people p where a.id = ap.aid and ap.pid = p.id and a.kind = $1 order by a.modtime desc", kind)
		if err != nil {
			log.Printf("error searching articles from db:%v, with err:%v\n", kind, err)
		}
	}
	return articles
}

// TODO: 把SearchArticles 和 SearchFavArticles 合并起来
// TODO: 找一个QueryBuilder替代直接拼凑SQL语句的方案，或者自己写一个查询组装器
func (r *ArticleRepo) SearchFavArticles(kind models.ArticleKind, mypid string, pid string) []models.Article {
	// note.ptype 是positionType, ptype = 3 表示文章级别的笔记；
	// note.ntype 是noteType.  ntype = 0 表示是这个笔记的类型是“收藏”。
	articles := []models.Article{}
	if pid != "" {
		err := r.db.Select(&articles, "SELECT a.*, ap.pid, p.fullname as author from articles a, notes n, article_people ap, people p where n.artid = a.id and n.ptype = 3 and n.ntype = 0 and a.id = ap.aid and ap.pid = p.id and a.kind = $1 and ap.pid = $2 order by a.modtime desc", kind, pid)
		if err != nil {
			log.Printf("error searching articles from db:%v:%v, with err:%v\n", kind, pid, err)
		}
	} else {
		err := r.db.Select(&articles, "SELECT a.*, ap.pid, p.fullname as author from articles a, notes n, article_people ap, people p where n.artid = a.id and n.ptype = 3 and n.ntype = 0 and a.id = ap.aid and ap.pid = p.id and a.kind = $1 order by a.modtime desc", kind)
		if err != nil {
			log.Printf("error searching articles from db:%v, with err:%v\n", kind, err)
		}
	}
	return articles
}

func (r *ArticleRepo) ListArticles(pid string) []models.Article {
	articles := []models.Article{}
	err := r.db.Select(&articles, "SELECT a.*, p.pid from articles a, article_people p where a.id = p.aid and p.pid = $1 order by a.modtime desc", pid)
	if err != nil {
		log.Printf("error listing articles from db:%v, with err:%v\n", pid, err)
	}
	return articles
}

func (r *ArticleRepo) ListArticlesInCollection(pid string, colid string) []models.Article {
	var artids []string
	err := r.db.Select(&artids, "SELECT artid from collections_articles where colid = $1", colid)
	if err != nil {
		log.Printf("error listing articles for collection from db:%v, with err:%v\n", colid, err)
	}

	// TODO: check if using pq.Array and ANY(?) is the best option here
	articles := []models.Article{}
	err = r.db.Select(&articles, "SELECT a.id, a.title, a.modtime, p.pid from articles a, article_people p where a.id = p.aid and p.pid = $1 and a.id = ANY($2) order by a.modtime desc", pid, pq.Array(artids))
	if err != nil {
		log.Printf("error listing articles from db:%v, with err:%v\n", pid, err)
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
