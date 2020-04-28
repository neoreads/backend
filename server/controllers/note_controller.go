package controllers

import (
	"log"
	"net/http"

	"github.com/neoreads/backend/util"

	"github.com/neoreads/backend/server/models"

	"github.com/gin-gonic/gin"
	"github.com/neoreads/backend/server/repositories"
)

// NoteController serves note info
type NoteController struct {
	Repo  *repositories.NoteRepo
	IDGen *util.N64Generator
}

func NewNoteController(r *repositories.NoteRepo) *NoteController {
	return &NoteController{
		Repo:  r,
		IDGen: util.NewN64Generator(8),
	}
}

func (ctrl *NoteController) Star(c *gin.Context) {
	var note models.Note
	if err := c.ShouldBindJSON(&note); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	delta := note.Value
	if delta > 1 {
		delta = 1
	}
	if delta < -1 {
		delta = -1
	}
	note.Value = delta
	// 新建操作
	if note.ID == "" {
		// gen note id
		noteid := ctrl.IDGen.Next()
		note.ID = noteid
		user, _ := c.Get("jwtuser")
		note.PID = user.(*models.User).Pid
		log.Printf("Note to Add: %v\n", note)
		if succ := ctrl.Repo.AddNote(&note); succ {
			c.JSON(http.StatusOK, note)
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "id": noteid})
		}
	} else {
		// 获取当前星级
		stars := ctrl.Repo.GetStars(note.ID)
		note.Value = stars + note.Value
		if note.Value > 0 {
			ctrl.Repo.UpdateStars(note.ID, delta)
		} else {
			ctrl.Repo.RemoveNote(note.ID)
		}
		c.JSON(http.StatusOK, note)
	}

}

func (ctrl *NoteController) AddNote(c *gin.Context) {
	var note models.Note
	if err := c.ShouldBindJSON(&note); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// gen note id
	noteid := ctrl.IDGen.Next()
	note.ID = noteid
	user, _ := c.Get("jwtuser")
	note.PID = user.(*models.User).Pid
	log.Printf("Note to Add: %v\n", note)
	if succ := ctrl.Repo.AddNote(&note); succ {
		c.JSON(http.StatusOK, gin.H{"status": "ok", "id": noteid})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "id": noteid})
	}
}

func (ctrl *NoteController) ModifyNote(c *gin.Context) {
	var note models.Note
	if err := c.ShouldBindJSON(&note); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if ok := ctrl.Repo.ModifyNote(&note); ok {
		c.JSON(http.StatusOK, "")
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to modify note"})
	}
}

func (ctrl *NoteController) RemoveNote(c *gin.Context) {
	id := c.Param("noteid")
	ctrl.Repo.RemoveNote(id)
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func (ctrl *NoteController) ListMyNotesForBook(c *gin.Context) {
	user, _ := c.Get("jwtuser")
	pid := user.(*models.User).Pid
	colid := c.Query("bookid")
	artid := c.Query("chapid")
	notes := ctrl.Repo.ListNotesForPid(pid, colid, artid)
	c.JSON(http.StatusOK, notes)
}

func (ctrl *NoteController) ListMyStarsForArticle(c *gin.Context) {
	user, _ := c.Get("jwtuser")
	pid := user.(*models.User).Pid
	var artids []string
	err := c.ShouldBindJSON(&artids)
	if err != nil {
		c.JSON(http.StatusBadRequest, artids)
	}
	notes := ctrl.Repo.ListStarsForArticles(pid, artids)
	c.JSON(http.StatusOK, notes)
}

func (ctrl *NoteController) ListAllNotes(c *gin.Context) {
	colid := c.Query("colid")
	artid := c.Query("artid")
	notes := ctrl.Repo.ListNotes(colid, artid)
	c.JSON(http.StatusOK, notes)
}

func (ctrl *NoteController) ListNoteCards(c *gin.Context) {
	colid := c.Query("colid")
	artid := c.Query("artid")
	cards := ctrl.Repo.ListNoteCards(colid, artid)
	c.JSON(http.StatusOK, cards)
}
