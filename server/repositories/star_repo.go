package repositories

import (
	"log"

	"github.com/jmoiron/sqlx"
)

// StarRepo Starorate related data repository
type StarRepo struct {
	db *sqlx.DB
}

// NewStarRepo creator for StarRepo
func NewStarRepo(db *sqlx.DB) *StarRepo {
	return &StarRepo{db: db}
}

// 每日定式更新Stars_quota表
func (r *StarRepo) ReplenishStarsQuota() {
	_, err := r.db.Exec("update Stars_quota set remaining = remaining + 3")
	if err != nil {
		log.Printf("error replenishing Stars quota, with err: %v\n", err)
		return
	}
}

func (r *StarRepo) DoStar() {

}
