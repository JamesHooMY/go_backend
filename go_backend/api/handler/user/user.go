package user

import (
	"errors"
	"net/http"

	hdl "go_backend/api/handler"
	userQryRepo "go_backend/api/repo/mysql/user"
	userSrv "go_backend/api/service/user"
	"go_backend/util"

	"github.com/gin-gonic/gin"
)

type IUserHandler interface {
	Login() gin.HandlerFunc
}

type UserHandler struct {
	UserService userSrv.IUserService
}

func NewUserHandler(userService userSrv.IUserService) IUserHandler {
	return &UserHandler{
		UserService: userService,
	}
}

// @Tags User
// @Router /user/login [post]
// @Summary User login
// @Description User login
// @Accept json
// @Produce json
// @Param username body string true "username"
// @Success 200 {object} LoginResp "success"
// @Failure 400 {object} string "bad request"
// @Failure 500 {object} string "internal server error"
func (h *UserHandler) Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			Email    string `json:"email" binding:"required,email"`
			Password string `json:"password" binding:"required,min=8,max=20,alphanum"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, &hdl.Response{
				Code: hdl.ErrRequestInvalid,
				Msg:  util.ParseValidateError(err).Error(),
			})
			return
		}

		loginResp, err := h.UserService.Login(c.Request.Context(), req.Email, req.Password)
		if err != nil {
			if errors.Is(err, userQryRepo.ErrUserNotFound) {
				c.AbortWithStatusJSON(http.StatusBadRequest, &hdl.Response{
					Code: hdl.ErrNotFound,
					Msg:  userQryRepo.ErrUserNotFound.Error(),
				})
				return
			}
			c.AbortWithStatusJSON(http.StatusInternalServerError, &hdl.Response{
				Code: hdl.ErrInternalServer,
				Msg:  hdl.ErrInternalServerMsg,
			})
			return
		}

		c.AbortWithStatusJSON(http.StatusOK, &hdl.Response{
			Data: gin.H{
				"username": loginResp.Username,
				"token":    loginResp.Token,
			},
		})
	}
}
