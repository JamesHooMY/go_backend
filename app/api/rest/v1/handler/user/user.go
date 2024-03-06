package user

import (
	"errors"
	"net/http"
	"strconv"

	"go_backend/model"
	"go_backend/util"

	v1 "go_backend/app/api/rest/v1"
	userRepo "go_backend/app/repo/mysql/user"
	userSrv "go_backend/app/service/user"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type IUserHandler interface {
	Login() gin.HandlerFunc
	Register() gin.HandlerFunc
	GetUserByID() gin.HandlerFunc
	GetUserList() gin.HandlerFunc
	UpdateUserByID() gin.HandlerFunc
	DeleteUserByID() gin.HandlerFunc
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
// @Router /api/v1/login [post]
// @Summary User login
// @Description User login
// @Accept json
// @Produce json
// @Param loginReq body loginReq true "login request"
// @Success 200 {object} LoginResp "success"
// @Failure 400 {object} v1.ErrorResponse "bad request"
// @Failure 401 {object} v1.ErrorResponse "unauthorized"
// @Failure 404 {object} v1.ErrorResponse "not found"
// @Failure 500 {object} v1.ErrorResponse "internal server error"
func (h *UserHandler) Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req *loginReq
		if err := c.ShouldBindJSON(&req); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, &v1.ErrorResponse{
				Code: v1.ErrRequestInvalid,
				Msg:  util.ParseValidateError(err).Error(),
			})
			return
		}

		loginResp, err := h.UserService.Login(c.Request.Context(), req.Email, req.Password)
		if err != nil {
			if errors.Is(err, userSrv.ErrPasswordIncorrect) {
				c.AbortWithStatusJSON(http.StatusUnauthorized, &v1.ErrorResponse{
					Code: v1.ErrInvalidPassword,
					Msg:  userSrv.ErrPasswordIncorrect.Error(),
				})
				return
			}

			if errors.Is(err, userRepo.ErrUserNotFound) {
				c.AbortWithStatusJSON(http.StatusNotFound, &v1.ErrorResponse{
					Code: v1.ErrNotFound,
					Msg:  userRepo.ErrUserNotFound.Error(),
				})
				return
			}

			c.AbortWithStatusJSON(http.StatusInternalServerError, &v1.ErrorResponse{
				Code: v1.ErrInternalServer,
				Msg:  v1.ErrInternalServerMsg,
			})
			return
		}

		c.JSON(http.StatusOK, &v1.Response{
			Data: loginResp,
		})
	}
}

type loginReq struct {
	Email    string `form:"email" binding:"required,email,max=50"`
	Password string `form:"password" binding:"required,min=8,max=20,alphanum"`
}

// @Tags User
// @Router /api/v1/register [post]
// @Summary User register
// @Description User register
// @Accept json
// @Produce json
// @Param registerReq body registerReq true "register request"
// @Success 204
// @Failure 400 {object} v1.ErrorResponse "bad request"
// @Failure 403 {object} v1.ErrorResponse "forbidden"
// @Failure 500 {object} v1.ErrorResponse "internal server error"
func (h *UserHandler) Register() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req *registerReq
		if err := c.ShouldBindJSON(&req); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, &v1.ErrorResponse{
				Code: v1.ErrRequestInvalid,
				Msg:  util.ParseValidateError(err).Error(),
			})
			return
		}

		err := h.UserService.Register(c.Request.Context(), req.Email, req.Password)
		if err != nil {
			if errors.Is(err, userRepo.ErrUserExisted) {
				c.AbortWithStatusJSON(http.StatusForbidden, &v1.ErrorResponse{
					Code: v1.ErrForbidden,
					Msg:  userRepo.ErrUserExisted.Error(),
				})
				return
			}

			c.AbortWithStatusJSON(http.StatusInternalServerError, &v1.ErrorResponse{
				Code: v1.ErrInternalServer,
				Msg:  v1.ErrInternalServerMsg,
			})
			return
		}

		c.JSON(http.StatusNoContent, nil)
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
// @Failure 400 {object} v1.ErrorResponse "bad request"
// @Failure 404 {object} v1.ErrorResponse "not found"
// @Failure 500 {object} v1.ErrorResponse "internal server error"
func (h *UserHandler) GetUserByID() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, err := strconv.ParseUint(c.Param("id"), 10, 0)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, &v1.ErrorResponse{
				Code: v1.ErrRequestInvalid,
				Msg:  "id must be an integer",
			})
			return
		}

		user, err := h.UserService.GetUserByID(c.Request.Context(), uint(userID))
		if err != nil {
			if errors.Is(err, userRepo.ErrUserNotFound) {
				c.AbortWithStatusJSON(http.StatusNotFound, &v1.ErrorResponse{
					Code: v1.ErrNotFound,
					Msg:  userRepo.ErrUserNotFound.Error(),
				})
				return
			}
			c.AbortWithStatusJSON(http.StatusInternalServerError, &v1.ErrorResponse{
				Code: v1.ErrInternalServer,
				Msg:  v1.ErrInternalServerMsg,
			})
			return
		}

		c.JSON(http.StatusOK, &v1.Response{
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
// @Failure 400 {object} v1.ErrorResponse "bad request"
// @Failure 500 {object} v1.ErrorResponse "internal server error"
func (h *UserHandler) GetUserList() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req *getUserListReq
		if err := c.ShouldBindJSON(&req); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, &v1.ErrorResponse{
				Code: v1.ErrRequestInvalid,
				Msg:  util.ParseValidateError(err).Error(),
			})
			return
		}

		userListResp, err := h.UserService.GetUserList(c.Request.Context(), req.Page, req.Limit)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, &v1.ErrorResponse{
				Code: v1.ErrInternalServer,
				Msg:  v1.ErrInternalServerMsg,
			})
			return
		}

		c.JSON(http.StatusOK, &v1.Response{
			Data: userListResp,
		})
	}
}

type getUserListReq struct {
	Page  int `form:"page" binding:"required,min=1"`
	Limit int `form:"limit" binding:"required,min=1"`
}

// @Tags User
// @Router /api/v1/user/{id} [put]
// @Summary Update user by id
// @Description Update user by id
// @Accept json
// @Produce json
// @Param id path int true "id"
// @Param updateUserReq body updateUserReq true "update user request"
// @Success 204
// @Failure 400 {object} v1.ErrorResponse "bad request"
// @Failure 403 {object} v1.ErrorResponse "forbidden"
// @Failure 404 {object} v1.ErrorResponse "not found"
// @Failure 500 {object} v1.ErrorResponse "internal server error"
func (h *UserHandler) UpdateUserByID() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, err := strconv.ParseUint(c.Param("id"), 10, 0)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, &v1.ErrorResponse{
				Code: v1.ErrRequestInvalid,
				Msg:  "id must be an integer",
			})
			return
		}

		// make sure user can only update his/her own info
		if userID != uint64(c.GetInt("userID")) {
			c.AbortWithStatusJSON(http.StatusForbidden, &v1.ErrorResponse{
				Code: v1.ErrForbidden,
				Msg:  "Not allowed to update other user's info",
			})
			return
		}

		var req *updateUserReq
		if err = c.ShouldBindJSON(&req); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, &v1.ErrorResponse{
				Code: v1.ErrRequestInvalid,
				Msg:  util.ParseValidateError(err).Error(),
			})
			return
		}

		err = h.UserService.UpdateUser(c.Request.Context(), &model.User{
			Model: gorm.Model{
				ID: uint(userID),
			},
			Mobile:   req.Mobile,
			Name:     req.Name,
			Age:      req.Age,
			Password: req.Password,
		})
		if err != nil {
			if errors.Is(err, userRepo.ErrUserNotFound) {
				c.AbortWithStatusJSON(http.StatusNotFound, &v1.ErrorResponse{
					Code: v1.ErrNotFound,
					Msg:  userRepo.ErrUserNotFound.Error(),
				})
				return
			}
			c.AbortWithStatusJSON(http.StatusInternalServerError, &v1.ErrorResponse{
				Code: v1.ErrInternalServer,
				Msg:  v1.ErrInternalServerMsg,
			})
			return
		}

		c.JSON(http.StatusNoContent, nil)
	}
}

type updateUserReq struct {
	Email    string `form:"email,omitempty" binding:"omitempty,email,max=50"`
	Mobile   string `form:"mobile,omitempty" binding:"omitempty,max=11"`
	Name     string `form:"name,omitempty" binding:"omitempty,max=20"`
	Age      int    `form:"age,omitempty" binding:"omitempty,min=1,max=150"`
	Password string `form:"password,omitempty" binding:"omitempty,min=8,max=20,alphanum"`
}

// @Tags User
// @Router /api/v1/user/{id} [delete]
// @Summary Delete user by id
// @Description Delete user by id
// @Accept json
// @Produce json
// @Param id path int true "id"
// @Success 204
// @Failure 400 {object} v1.ErrorResponse "bad request"
// @Failure 403 {object} v1.ErrorResponse "forbidden"
// @Failure 404 {object} v1.ErrorResponse "not found"
// @Failure 500 {object} v1.ErrorResponse "internal server error"
func (h *UserHandler) DeleteUserByID() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, err := strconv.ParseUint(c.Param("id"), 10, 0)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, &v1.ErrorResponse{
				Code: v1.ErrRequestInvalid,
				Msg:  "id must be an integer",
			})
			return
		}

		// make sure user can only delete his/her own info
		if userID != uint64(c.GetInt("userID")) {
			c.AbortWithStatusJSON(http.StatusForbidden, &v1.ErrorResponse{
				Code: v1.ErrForbidden,
				Msg:  "Not allowed to delete other user",
			})
			return
		}

		err = h.UserService.DeleteUserByID(c.Request.Context(), uint(userID))
		if err != nil {
			if errors.Is(err, userRepo.ErrUserNotFound) {
				c.AbortWithStatusJSON(http.StatusNotFound, &v1.ErrorResponse{
					Code: v1.ErrNotFound,
					Msg:  userRepo.ErrUserNotFound.Error(),
				})
				return
			}
			c.AbortWithStatusJSON(http.StatusInternalServerError, &v1.ErrorResponse{
				Code: v1.ErrInternalServer,
				Msg:  v1.ErrInternalServerMsg,
			})
			return
		}

		c.JSON(http.StatusNoContent, nil)
	}
}
