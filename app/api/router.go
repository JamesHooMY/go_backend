package rest

import (
	"fmt"

	"go_backend/app/api/graphql/gql"
	"go_backend/app/api/graphql/gql/generated"
	userHdl "go_backend/app/api/rest/v1/handler/user"
	"go_backend/app/api/rest/v1/middleware"
	userRepo "go_backend/app/repo/mysql/user"
	userRdsRepo "go_backend/app/repo/redis/user"
	userSrv "go_backend/app/service/user"
	_ "go_backend/docs"

	// _ "github.com/99designs/gqlgen"
	_ "github.com/99designs/gqlgen/graphql/introspection"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gin-contrib/pprof"

	// "github.com/99designs/gqlgen/handler"
	gqlHdl "go_backend/app/api/graphql"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"
)

func InitRouter(router *gin.Engine, db *gorm.DB, rds *redis.Client) *gin.Engine {
	// pprof
	pprof.Register(router)

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

	// graphql playground
	if viper.GetBool("server.gqlPlayground") {
		router.GET(viper.GetString("server.gqlPlaygroundPath"), func(c *gin.Context) {
			playground.Handler("GraphQL playground", viper.GetString("server.gqlPath")).ServeHTTP(c.Writer, c.Request)
		})
	}

	// graphql
	gqlHandler := gqlHdl.NewGqlHandler(handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{
		Resolvers: &gql.Resolver{
			UserService: userSrv.NewUserService(
				userRepo.NewUserQueryRepo(db),
				userRepo.NewUserCommandRepo(db),
				userRdsRepo.NewUserRedisRepo(rds),
			),
		},
	})))
	router.POST(viper.GetString("server.gqlPath"), gqlHandler.GqlHdl())

	// grpc client
	// grpcClt := grpcClient.NewGrpcClient()
	// grpc := router.Group("/grpc")
	// grpc.POST(viper.GetString("server.grpcPath"), func(c *gin.Context) {
	// 	grpcClt.UserSrvClient.CreateUser(c)
	// })

	return router
}
