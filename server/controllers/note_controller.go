package controllers

import (
	"log"
	"net/http"

	"github.com/neoreads-backend/go/util"

	"github.com/neoreads-backend/go/server/models"

	"github.com/gin-gonic/gin"
	"github.com/neoreads-backend/go/server/repositories"
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

func (ctrl *NoteController) AddNote(c *gin.Context) {
	var note models.Note
	if err := c.ShouldBindJSON(&note); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// gen note id
	noteid := ctrl.IDGen.Next()
	note.ID = noteid
	// TODO: currently credentials include only users.username, but notes are linked with people.id(or users.pid)
	user, _ := c.Get("id")
	note.PID = user.(*models.Credential).Username
	log.Printf("Note to Add: %v\n", note)
	ctrl.Repo.AddNote(&note)
	c.JSON(http.StatusOK, gin.H{"status": "ok", "id": noteid})
}

func (ctrl *NoteController) RemoveNote(c *gin.Context) {
	id := c.Param("noteid")
	ctrl.Repo.RemoveNote(id)
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func (ctrl *NoteController) ListNotes(c *gin.Context) {
	user, _ := c.Get("id")
	username := user.(*models.Credential).Username
	bookid := c.Query("bookid")
	chapid := c.Query("chapid")
	notes := ctrl.Repo.ListNotes(username, bookid, chapid)
	c.JSON(http.StatusOK, notes)
}
