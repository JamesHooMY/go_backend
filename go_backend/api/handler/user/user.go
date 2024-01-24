package user

import (
	"net/http"

	userSrv "go_backend/api/service/user"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
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
			Username string `json:"username" validate:"required"`
			Password string `json:"password" validate:"required"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"msg": "invalid request body",
			})
			return
		}

		validator := validator.New()
		if err := validator.Struct(req); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"msg": "body validation failed",
			})
			return
		}

		token, err := h.UserService.Login(c.Request.Context(), req.Username, req.Password)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"msg": "login failed",
			})
			return
		}

		c.AbortWithStatusJSON(http.StatusOK, gin.H{
			"token": token,
			"msg":   "login success",
		})
	}
}

func (h *UserHandler) Info() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("token")
		if token == "" {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"msg": "token is required",
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
