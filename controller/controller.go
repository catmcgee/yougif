package controller

import (
	"bytes"
	"catmcgee/model"
	"catmcgee/repository"
	"catmcgee/service"
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
)

func Search(ctx *gin.Context) {
	searchString := ctx.Query("q")
	if len(searchString) == 0 {
		ctx.JSON(http.StatusBadRequest, model.ApiError{Message: "missing or empty query parameter 'q'"})
		return
	}

	searchResultData, err := service.SearchForVideo(ctx.Request.Context(), searchString)
	if err != nil {
		logrus.Println(err)
		ctx.JSON(http.StatusInternalServerError, model.ApiError{Message: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, searchResultData)
}

func GetImage(ctx *gin.Context) {
	id := ctx.Param("id")

	image, err := repository.SelectByIdFrame(ctx.Request.Context(), id)
	if err != nil {
		logrus.Println(err)
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, model.ApiError{Message: fmt.Sprintf("image with id '%s' not found", id)})
			return
		}
		ctx.JSON(http.StatusInternalServerError, model.ApiError{Message: err.Error()})
		return
	}

	ctx.DataFromReader(http.StatusOK, int64(len(image)), "image/jpeg", bytes.NewReader(image), nil)
}
