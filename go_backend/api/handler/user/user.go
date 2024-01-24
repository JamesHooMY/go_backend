package user

import (
	"net/http"

	hdl "go_backend/api/handler"
	userSrv "go_backend/api/service/user"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	UserService userSrv.UserService
}

func NewUserHandler(userService userSrv.UserService) *UserHandler {
	return &UserHandler{
		UserService: userService,
	}
}

func (h *UserHandler) Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			Username string `json:"username" binding:"required"`
			Password string `json:"password" binding:"required"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, &hdl.Response{
				Status: http.StatusBadRequest,
				Msg:    hdl.ParseValidateError(err, &req).Error(),
			})
			return
		}

		token, err := h.UserService.Login(c.Request.Context(), req.Username, req.Password)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, &hdl.Response{
				Status: http.StatusInternalServerError,
				Msg:    "login failed",
			})
			return
		}

		c.AbortWithStatusJSON(http.StatusOK, &hdl.Response{
			Status: http.StatusOK,
			Data:   &LoginResp{Token: token},
		})
	}
}

type LoginResp struct {
	Token string `json:"token"`
}

// @Tags User
// @Router /user/info [get]
// @Summary Get user info summary
// @Description Get user info description
// @Accept json
// @Produce json
// @Param token header string true "token"
// @Success 200 {object} string "ok"
// @Failure 400 {object} string "bad request"
// @Failure 500 {object} string "internal server error"
func (h *UserHandler) Info() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("token")
		if token == "" {
			c.AbortWithStatusJSON(http.StatusBadRequest, &hdl.Response{
				Status: http.StatusBadRequest,
				Msg:    "token is empty",
			})
			return
		}

		name, err := h.UserService.Info(c.Request.Context(), token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"msg": "get info failed",
			})
			return
		}

		c.AbortWithStatusJSON(http.StatusOK, gin.H{
			"info": name,
			"msg":  "get info success",
		})
	}
}
