package rest

import (
	"fmt"

	userHdl "go_backend/app/api/rest/v1/handler/user"
	"go_backend/app/api/rest/v1/middleware"
	userRepo "go_backend/app/repo/mysql/user"
	userRdsRepo "go_backend/app/repo/redis/user"
	userSrv "go_backend/app/service/user"
	_ "go_backend/docs"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"
)

func InitRouter(router *gin.Engine, db *gorm.DB, rds *redis.Client) *gin.Engine {
	// middleware
	// * if need cors then uncomment this line
	// router.Use(middleware.Cors())

	// swagger
	// docs.SwaggerInfo.BasePath = fmt.Sprintf("/api/%s", viper.GetString("server.apiVersion"))
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	userHandler := userHdl.NewUserHandler(userSrv.NewUserService(
		userRepo.NewUserQueryRepo(db),
		userRepo.NewUserCommandRepo(db),
		userRdsRepo.NewUserRedisRepo(rds),
	))

	v1 := router.Group(fmt.Sprintf("/api/%s", viper.GetString("server.apiVersion")))
	v1.POST("/login", userHandler.Login())
	v1.POST("/register", userHandler.Register())

	user := v1.Group("/user").Use(middleware.Auth(
		userRdsRepo.NewUserRedisRepo(rds),
		userRepo.NewUserQueryRepo(db),
	))
	user.GET("/:id", userHandler.GetUserByID())
	user.POST("/userList", userHandler.GetUserList())
	user.PUT("/:id", userHandler.UpdateUserByID())
	user.DELETE("/:id", userHandler.DeleteUserByID())

	return router
}
