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
	GetUserList() gin.HandlerFunc
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
// @Router /api/v1/user/login [post]
// @Summary User login
// @Description User login
// @Accept json
// @Produce json
// @Param loginReq body loginReq true "login request"
// @Success 200 {object} LoginResp "success"
// @Failure 400 {object} hdl.ErrorResponse "bad request"
// @Failure 500 {object} hdl.ErrorResponse "internal server error"
func (h *UserHandler) Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req *loginReq
		if err := c.ShouldBindJSON(&req); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, &hdl.ErrorResponse{
				Code: hdl.ErrRequestInvalid,
				Msg:  util.ParseValidateError(err).Error(),
			})
			return
		}

		loginResp, err := h.UserService.Login(c.Request.Context(), req.Email, req.Password)
		if err != nil {
			if errors.Is(err, userRepo.ErrUserNotFound) {
				c.AbortWithStatusJSON(http.StatusBadRequest, &hdl.ErrorResponse{
					Code: hdl.ErrNotFound,
					Msg:  userRepo.ErrUserNotFound.Error(),
				})
				return
			}
			c.AbortWithStatusJSON(http.StatusInternalServerError, &hdl.ErrorResponse{
				Code: hdl.ErrInternalServer,
				Msg:  hdl.ErrInternalServerMsg,
			})
			return
		}

		c.AbortWithStatusJSON(http.StatusOK, &hdl.Response{
			Data: loginResp,
		})
	}
}

type loginReq struct {
	Email    string `form:"email" binding:"required,email,max=50"`
	Password string `form:"password" binding:"required,min=8,max=20,alphanum"`
}

// @Tags User
// @Router /api/v1/user/register [post]
// @Summary User register
// @Description User register
// @Accept json
// @Produce json
// @Param registerReq body registerReq true "register request"
// @Success 204
// @Failure 400 {object} hdl.ErrorResponse "bad request"
// @Failure 403 {object} hdl.ErrorResponse "forbidden"
// @Failure 500 {object} hdl.ErrorResponse "internal server error"
func (h *UserHandler) Register() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req *registerReq
		if err := c.ShouldBindJSON(&req); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, &hdl.ErrorResponse{
				Code: hdl.ErrRequestInvalid,
				Msg:  util.ParseValidateError(err).Error(),
			})
			return
		}

		err := h.UserService.Register(c.Request.Context(), req.Email, req.Password)
		if err != nil {
			if errors.Is(err, userRepo.ErrUserExisted) {
				c.AbortWithStatusJSON(http.StatusForbidden, &hdl.ErrorResponse{
					Code: hdl.ErrForbidden,
					Msg:  userRepo.ErrUserExisted.Error(),
				})
				return
			}
			c.AbortWithStatusJSON(http.StatusInternalServerError, &hdl.ErrorResponse{
				Code: hdl.ErrInternalServer,
				Msg:  hdl.ErrInternalServerMsg,
			})
			return
		}

		c.AbortWithStatus(http.StatusNoContent)
	}
}

type registerReq struct {
	Email    string `form:"email" binding:"required,email,max=50"`
	Password string `form:"password" binding:"required,min=8,max=20,alphanum"`
}

// @Tags User
// @Router /api/v1/user/{id} [get]
// @Summary Get user by id
// @Description Get user by id
// @Accept json
// @Produce json
// @Param id path int true "id"
// @Success 200 {object} UserResp "success"
// @Failure 400 {object} hdl.ErrorResponse "bad request"
// @Failure 500 {object} hdl.ErrorResponse "internal server error"
func (h *UserHandler) GetUserByID() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if id == "" {
			c.AbortWithStatusJSON(http.StatusBadRequest, &hdl.ErrorResponse{
				Code: hdl.ErrRequestInvalid,
				Msg:  "id is required",
			})
			return
		}

		userID, err := strconv.ParseUint(id, 10, 0)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, &hdl.ErrorResponse{
				Code: hdl.ErrRequestInvalid,
				Msg:  "id must be an integer",
			})
			return
		}

		user, err := h.UserService.GetUserByID(c.Request.Context(), uint(userID))
		if err != nil {
			if errors.Is(err, userRepo.ErrUserNotFound) {
				c.AbortWithStatusJSON(http.StatusBadRequest, &hdl.ErrorResponse{
					Code: hdl.ErrNotFound,
					Msg:  userRepo.ErrUserNotFound.Error(),
				})
				return
			}
			c.AbortWithStatusJSON(http.StatusInternalServerError, &hdl.ErrorResponse{
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

// @Tags User
// @Router /api/v1/user/userList [post]
// @Summary Get user list
// @Description Get user list
// @Accept json
// @Produce json
// @Param getUserListReq body getUserListReq true "get user list request"
// @Success 200 {object} UserListResp "success"
// @Failure 400 {object} hdl.ErrorResponse "bad request"
// @Failure 500 {object} hdl.ErrorResponse "internal server error"
func (h *UserHandler) GetUserList() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req *getUserListReq
		if err := c.ShouldBindJSON(&req); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, &hdl.ErrorResponse{
				Code: hdl.ErrRequestInvalid,
				Msg:  util.ParseValidateError(err).Error(),
			})
			return
		}

		userListResp, err := h.UserService.GetUserList(c.Request.Context(), req.Page, req.Limit)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, &hdl.ErrorResponse{
				Code: hdl.ErrInternalServer,
				Msg:  hdl.ErrInternalServerMsg,
			})
			return
		}

		c.AbortWithStatusJSON(http.StatusOK, &hdl.Response{
			Data: userListResp,
		})
	}
}

type getUserListReq struct {
	Page  int `form:"page" binding:"required,min=1"`
	Limit int `form:"limit" binding:"required,min=1"`
}
