package user

import (
	"errors"
	"net/http"
	"strconv"

	hdl "go_backend/api/rest/handler"
	userRepo "go_backend/api/rest/repo/mysql/user"
	userSrv "go_backend/api/rest/service/user"
	"go_backend/util"

	"github.com/gin-gonic/gin"
)

type IUserHandler interface {
	Login() gin.HandlerFunc
	Register() gin.HandlerFunc
	GetUserByID() gin.HandlerFunc
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
			Email    string `json:"email" binding:"required,email,max=50"`
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
			if errors.Is(err, userRepo.ErrUserNotFound) {
				c.AbortWithStatusJSON(http.StatusBadRequest, &hdl.Response{
					Code: hdl.ErrNotFound,
					Msg:  userRepo.ErrUserNotFound.Error(),
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

// @Tags User
// @Router /user/register [post]
// @Summary User register
// @Description User register
// @Accept json
// @Produce json
// @Param username body string true "username"
// @Param email body string true "email"
// @Param password body string true "password"
// @Success 200 {object} string "success"
// @Failure 400 {object} string "bad request"
// @Failure 500 {object} string "internal server error"
func (h *UserHandler) Register() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			Email    string `json:"email" binding:"required,email,max=50"`
			Password string `json:"password" binding:"required,min=8,max=20,alphanum"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, &hdl.Response{
				Code: hdl.ErrRequestInvalid,
				Msg:  util.ParseValidateError(err).Error(),
			})
			return
		}

		err := h.UserService.Register(c.Request.Context(), req.Email, req.Password)
		if err != nil {
			if errors.Is(err, userRepo.ErrUserExisted) {
				c.AbortWithStatusJSON(http.StatusBadRequest, &hdl.Response{
					Code: hdl.ErrForbidden,
					Msg:  userRepo.ErrUserExisted.Error(),
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
			Data: "success",
		})
	}
}

// @Tags User
// @Router /user/{id} [get]
// @Summary Get user by id
// @Description Get user by id
// @Accept json
// @Produce json
// @Param id path int true "id"
// @Success 200 {object} UserResp "success"
// @Failure 400 {object} string "bad request"
// @Failure 500 {object} string "internal server error"
func (h *UserHandler) GetUserByID() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if id == "" {
			c.AbortWithStatusJSON(http.StatusBadRequest, &hdl.Response{
				Code: hdl.ErrRequestInvalid,
				Msg:  "id is required",
			})
			return
		}

		userID, err := strconv.ParseUint(id, 10, 0)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, &hdl.Response{
				Code: hdl.ErrRequestInvalid,
				Msg:  "id must be an integer",
			})
			return
		}

		user, err := h.UserService.GetUserByID(c.Request.Context(), uint(userID))
		if err != nil {
			if errors.Is(err, userRepo.ErrUserNotFound) {
				c.AbortWithStatusJSON(http.StatusBadRequest, &hdl.Response{
					Code: hdl.ErrNotFound,
					Msg:  userRepo.ErrUserNotFound.Error(),
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
			Data: user,
		})
	}
}
