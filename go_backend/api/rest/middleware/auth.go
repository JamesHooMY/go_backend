package middleware

import (
	"context"
	"net/http"
	"strings"

	hdl "go_backend/api/rest/handler"
	userRepo "go_backend/api/rest/repo/mysql/user"
	"go_backend/model"
	"go_backend/util"

	"github.com/gin-gonic/gin"
)

type IAuthRdsRepo interface {
	GetLoginToken(ctx context.Context, userID uint) (token string, err error)
}

type IAuthUserQueryRepo interface {
	GetUserByID(ctx context.Context, id uint) (user *model.User, err error)
}

func Auth(authRdsRepo IAuthRdsRepo, authUserQryRepo IAuthUserQueryRepo) gin.HandlerFunc {
	return func(c *gin.Context) {
		headerAuth := c.GetHeader("Authorization")
		headerSplit := strings.Split(headerAuth, "Bearer ")
		if len(headerSplit) != 2 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, &hdl.ErrorResponse{
				Code: hdl.ErrTokenRequired,
				Msg:  "Token required",
			})
			return
		}

		token := headerSplit[1]

		// check if token is valid
		parsedToken, err := util.ParseJwtToken(token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, &hdl.ErrorResponse{
				Code: hdl.ErrInvalidToken,
				Msg:  "Invalid token",
			})
			return
		}

		// check if token is expired
		rdsToken, err := authRdsRepo.GetLoginToken(c.Request.Context(), parsedToken.ID)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, &hdl.ErrorResponse{
				Code: hdl.ErrInternalServer,
				Msg:  hdl.ErrInternalServerMsg,
			})
			return
		}

		// check if token is authorized
		if rdsToken != token {
			c.AbortWithStatusJSON(http.StatusUnauthorized, &hdl.ErrorResponse{
				Code: hdl.ErrUnauthorizedToken,
				Msg:  "Unauthorized token",
			})
			return
		}

		// check if user is existed
		if _, err = authUserQryRepo.GetUserByID(c.Request.Context(), parsedToken.ID); err != nil {
			if err == userRepo.ErrUserNotFound {
				c.AbortWithStatusJSON(http.StatusNotFound, &hdl.ErrorResponse{
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

		// set user id to context
		c.Set("userID", parsedToken.ID)

		c.Next()
	}
}
