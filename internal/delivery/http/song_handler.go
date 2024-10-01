package http

import (
	"github.com/gin-gonic/gin"
	"github.com/nurmeden/music-library/internal/entity"
	"github.com/nurmeden/music-library/internal/usecase"
	"net/http"
	"strconv"
)

type SongHandler struct {
	UseCase *usecase.SongUseCase
}

func NewSongHandler(r *gin.Engine, uc *usecase.SongUseCase) {
	handler := &SongHandler{UseCase: uc}
	r.GET("/songs", handler.FetchAllSongs)
	r.POST("/songs", handler.AddSong)
	r.PUT("/songs/:id", handler.UpdateSong)
	r.DELETE("/songs/:id", handler.DeleteSong)
}

func (h *SongHandler) FetchAllSongs(c *gin.Context) {
	filters := make(map[string]interface{})
	group := c.Query("group")
	if group != "" {
		filters["group_name"] = group
	}
	song := c.Query("song")
	if song != "" {
		filters["song_name"] = song
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	songs, err := h.UseCase.FetchAll(filters, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, songs)
}

func (h *SongHandler) AddSong(c *gin.Context) {
	var song entity.Song
	if err := c.ShouldBindJSON(&song); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.UseCase.AddNewSong(&song)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"status": "song added"})
}

func (h *SongHandler) UpdateSong(c *gin.Context) {
	var song entity.Song
	if err := c.ShouldBindJSON(&song); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, _ := strconv.Atoi(c.Param("id"))
	song.ID = id

	err := h.UseCase.UpdateSong(&song)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "song updated"})
}

func (h *SongHandler) DeleteSong(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	err := h.UseCase.DeleteSong(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "song deleted"})
}
