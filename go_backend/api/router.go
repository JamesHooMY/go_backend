package api

import (
	"fmt"

	userHdl "go_backend/api/handler/user"
	userQryRepo "go_backend/api/repo/mysql/user"
	userSrv "go_backend/api/service/user"
	_ "go_backend/docs"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"
)

func InitRouter(router *gin.Engine, db *gorm.DB, rd *redis.ClusterClient) *gin.Engine {
	// middleware
	// * if need cors then uncomment this line
	// router.Use(middleware.Cors())

	// swagger
	// docs.SwaggerInfo.BasePath = fmt.Sprintf("/api/%s", viper.GetString("server.apiVersion"))
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	userHandler := userHdl.NewUserHandler(userSrv.NewUserService(userQryRepo.NewUserQueryRepo(db)))

	v1 := router.Group(fmt.Sprintf("/api/%s", viper.GetString("server.apiVersion")))
	user := v1.Group("/user")
	user.POST("/login", userHandler.Login())

	return router
}
