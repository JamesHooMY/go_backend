package api

import (
	"fmt"

	userHdl "go_backend/api/handler/user"
	userSrv "go_backend/api/service/user"
	_ "go_backend/docs"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

func InitRouter(router *gin.Engine, db *gorm.DB) *gin.Engine {
	// docs.SwaggerInfo.BasePath = fmt.Sprintf("/api/%s", viper.GetString("server.apiVersion"))
	// router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	userHandler := userHdl.NewUserHandler(userSrv.NewUserService(db))

	v1 := router.Group(fmt.Sprintf("/api/%s", viper.GetString("server.apiVersion")))
	user := v1.Group("/user")
	user.GET("/info", userHandler.Info())
	user.POST("/login", userHandler.Login())

	return router
}
