package handler

import (
	"github.com/gin-gonic/gin"
	"go-web-wire-starter/internal/pkg/request"
	"go-web-wire-starter/internal/pkg/response"
	"go-web-wire-starter/internal/service"
	"go.uber.org/zap"
	"strconv"
)

type MediaHandler struct {
	logger       *zap.Logger
	mediaService *service.MediaService
}

func NewMediaHandler(logger *zap.Logger, mediaService *service.MediaService) *MediaHandler {
	return &MediaHandler{
		logger:       logger,
		mediaService: mediaService,
	}
}

// ImageUpload 上传图片
func (h *MediaHandler) ImageUpload(c *gin.Context) {
	var form request.ImageUpload
	if err := c.ShouldBind(&form); err != nil {
		response.FailByParams(c, request.GetErrorMsg(form, err))
		return
	}

	media, err := h.mediaService.SaveImage(c, &form)
	if err != nil {
		response.FailByBussiness(c, err.Error())
		return
	}

	response.Success(c, media)
}

// 根据媒体ID获取媒体信息
func (h *MediaHandler) GetUrlById(c *gin.Context) {
	idStr := c.Query("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.FailByErr(c, err)
		return
	}

	media, err := h.mediaService.GetUrlById(c, id)
	if err != nil {
		response.FailByBussiness(c, err.Error())
		return
	}

	response.Success(c, media)
}
