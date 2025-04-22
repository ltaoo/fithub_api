package routes

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"myapi/config"
	"myapi/internal/api/handlers"
	"myapi/internal/api/middlewares"
	"myapi/pkg/logger"
)

// SetupRouter 配置API路由
func SetupRouter(db *gorm.DB, logger *logger.Logger, cfg *config.Config) *gin.Engine {
	// 设置Gin模式
	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()

	// 使用中间件
	r.Use(gin.Recovery())
	r.Use(gin.Logger())

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})

	// API路由组
	api := r.Group("/api")
	authorized := api.Group("/")
	authorized.Use(middlewares.AuthMiddleware(logger))
	{
		// 用户处理器
		userHandler := handlers.NewUserHandler(db, logger)

		// 公开路由
		api.POST("/user/list", userHandler.GetUsers)
		api.POST("/user/profile", userHandler.GetUser)

		{
			authorized := api.Group("/user")
			authorized.Use(middlewares.AuthMiddleware(logger))
			authorized.POST("/user/create", userHandler.CreateUser)
			authorized.POST("/user/update", userHandler.UpdateUser)
			authorized.POST("/user/delete", userHandler.DeleteUser)
		}
		{

			handler := handlers.NewWorkoutPlanHandler(db, logger)
			api.POST("/workout_plan/list", handler.FetchWorkoutPlanList)
			api.POST("/workout_plan/profile", handler.GetWorkoutPlan)
			authorized.POST("/workout_plan/update", handler.UpdateWorkoutPlan)
			authorized.POST("/workout_plan/delete", handler.DeleteWorkoutPlan)
			authorized.POST("/workout_plan/create", handler.CreateWorkoutPlan)
			authorized.POST("/workout_plan/mine", handler.FetchMyWorkoutPlanList)
		}
		{
			handler := handlers.NewWorkoutDayHandler(db, logger)
			authorized.POST("/workout_day/create", handler.CreateWorkoutDay)
			authorized.POST("/workout_day/profile", handler.FetchWorkoutDayProfile)
			authorized.POST("/workout_day/fetch_started", handler.FetchRunningWorkoutDay)
			authorized.POST("/workout_day/start", handler.StartWorkoutDay)
			authorized.POST("/workout_day/give_up", handler.GiveUpWorkoutDay)
			authorized.POST("/workout_day/finish", handler.FinishWorkoutDay)
			authorized.POST("/workout_day/update_steps", handler.UpdateWorkoutDaySteps)
			authorized.POST("/workout_day/delete", handler.DeleteWorkoutDay)
		}
		{
			handler := handlers.NewWorkoutActionHistoryHandler(db, logger)
			authorized.POST("/workout_action_history/list", handler.FetchWorkoutActionHistoryList)
		}
		{
			handler := handlers.NewWorkoutActionHandler(db, logger)
			api.POST("/workout_action/list", handler.GetWorkoutActionList)
			api.POST("/workout_action/list_by_ids", handler.GetWorkoutActionListByIds)
			api.POST("/workout_action/profile", handler.GetWorkoutAction)
			api.POST("/workout_action/list/by_muscle", handler.GetActionsByMuscle)
			api.POST("/workout_action/list/by_level", handler.GetActionsByLevel)
			api.POST("/workout_action/list/related", handler.GetRelatedActions)

			{
				authorized.POST("/workout_action/create", handler.CreateAction)
				authorized.POST("/workout_action/update", handler.UpdateAction)
				authorized.POST("/workout_action/delete", handler.DeleteAction)
			}
		}
		{
			handler := handlers.NewCoachHandler(db, logger)
			api.POST("/coach/register", handler.RegisterCoach)
			api.POST("/coach/login", handler.LoginCoach)
			api.POST("/coach/send-verification-code", handler.SendVerificationCode)
			{
				authorized.POST("/coach/profile", handler.GetCoachProfile)
				authorized.POST("/coach/update", handler.UpdateCoachProfile)
			}
		}
		{
			handler := handlers.NewMuscleHandler(db, logger)
			api.POST("/muscle/list", handler.FetchMuscleList)
			api.POST("/muscle/profile", handler.GetMuscleProfile)

			{
				authorized.POST("/muscle/create", handler.CreateMuscle)
				authorized.POST("/muscle/update", handler.UpdateMuscle)
				authorized.POST("/muscle/delete", handler.DeleteMuscle)
			}
		}
		{
			handler := handlers.NewEquipmentHandler(db, logger)
			api.POST("/equipment/list", handler.FetchEquipmentList)
			api.POST("/equipment/profile", handler.GetEquipment)
			{
				authorized.POST("/equipment/create", handler.CreateEquipment)
			}
		}
	}

	return r
}
