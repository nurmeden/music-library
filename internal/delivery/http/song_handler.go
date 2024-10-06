package http

import (
	"github.com/gin-gonic/gin"
	"github.com/nurmeden/music-library/internal/entity"
	"github.com/nurmeden/music-library/internal/logger"
	"github.com/nurmeden/music-library/internal/usecase"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"net/http"
	"strconv"
)

type SongHandler struct {
	UseCase *usecase.SongUseCase
	Logger  logger.Logger
}

func NewSongHandler(r *gin.Engine, uc *usecase.SongUseCase, logger logger.Logger) {
	handler := &SongHandler{UseCase: uc, Logger: logger}

	// Swagger annotations
	// @Summary Fetch all songs
	// @Description Get a list of songs with optional filters.
	// @Tags songs
	// @Produce json
	// @Param group query string false "Group name"
	// @Param song query string false "Song name"
	// @Param limit query int false "Limit of songs"
	// @Param offset query int false "Offset for pagination"
	// @Success 200 {array} entity.Song
	// @Failure 500 {object} gin.H
	// @Router /songs [get]
	r.GET("/songs", handler.FetchAllSongs)

	// @Summary Add a new song
	// @Description Add a new song to the library.
	// @Tags songs
	// @Accept json
	// @Produce json
	// @Param song body entity.Song true "Song data"
	// @Success 201 {object} gin.H
	// @Failure 400 {object} gin.H
	// @Failure 500 {object} gin.H
	// @Router /songs [post]
	r.POST("/songs", handler.AddSong)

	// @Summary Update an existing song
	// @Description Update a song by its ID.
	// @Tags songs
	// @Accept json
	// @Produce json
	// @Param id path int true "Song ID"
	// @Param song body entity.Song true "Updated song data"
	// @Success 200 {object} gin.H
	// @Failure 400 {object} gin.H
	// @Failure 500 {object} gin.H
	// @Router /songs/{id} [put]
	r.PUT("/songs/:id", handler.UpdateSong)

	// @Summary Delete a song
	// @Description Mark a song as deleted by its ID.
	// @Tags songs
	// @Produce json
	// @Param id path int true "Song ID"
	// @Success 200 {object} gin.H
	// @Failure 400 {object} gin.H
	// @Failure 500 {object} gin.H
	// @Router /songs/{id} [delete]
	r.DELETE("/songs/:id", handler.DeleteSong)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}

func (h *SongHandler) FetchAllSongs(c *gin.Context) {
	h.Logger.Info("Fetching all songs")

	filters := make(map[string]interface{})
	group := c.Query("group")
	if group != "" {
		filters["group_name"] = group
		h.Logger.Debugf("Filter applied: group_name = %s", group)
	}
	song := c.Query("song")
	if song != "" {
		filters["song_name"] = song
		h.Logger.Debugf("Filter applied: song_name = %s", song)
	}

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset, err := strconv.Atoi(c.DefaultQuery("offset", "0"))
	if err != nil {
		h.Logger.Errorf("Error parsing limit or offset: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	h.Logger.Infof("Fetching songs with limit %d and offset %d", limit, offset)
	songs, err := h.UseCase.FetchAll(filters, limit, offset)
	if err != nil {
		h.Logger.Errorf("Error fetching songs: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	h.Logger.Infof("Successfully fetched %d songs", len(songs))
	c.JSON(http.StatusOK, songs)
}

func (h *SongHandler) AddSong(c *gin.Context) {
	h.Logger.Info("Adding a new song")

	var song entity.Song
	if err := c.ShouldBindJSON(&song); err != nil {
		h.Logger.Errorf("Error binding song JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	h.Logger.Debugf("Song to add: %+v", song)
	err := h.UseCase.AddNewSong(&song)
	if err != nil {
		h.Logger.Errorf("Error adding song: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	h.Logger.Info("Song added successfully")
	c.JSON(http.StatusCreated, gin.H{"status": "song added"})
}

func (h *SongHandler) UpdateSong(c *gin.Context) {
	h.Logger.Info("Updating a song")

	var song entity.Song
	if err := c.ShouldBindJSON(&song); err != nil {
		h.Logger.Errorf("Error binding song JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		h.Logger.Errorf("Invalid song ID: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	song.ID = id

	h.Logger.Debugf("Song to update: %+v", song)
	err = h.UseCase.UpdateSong(&song)
	if err != nil {
		h.Logger.Errorf("Error updating song with ID %d: %v", song.ID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	h.Logger.Infof("Song with ID %d updated successfully", song.ID)
	c.JSON(http.StatusOK, gin.H{"status": "song updated"})
}

func (h *SongHandler) DeleteSong(c *gin.Context) {
	h.Logger.Info("Deleting a song")

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		h.Logger.Errorf("Invalid song ID: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	h.Logger.Debugf("Song ID to delete: %d", id)
	err = h.UseCase.DeleteSong(id)
	if err != nil {
		h.Logger.Errorf("Error deleting song with ID %d: %v", id, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	h.Logger.Infof("Song with ID %d deleted successfully", id)
	c.JSON(http.StatusOK, gin.H{"status": "song deleted"})
}
