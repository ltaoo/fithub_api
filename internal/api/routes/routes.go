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
	{
		// 用户处理器
		userHandler := handlers.NewUserHandler(db, logger)

		// 公开路由
		api.POST("/user/list", userHandler.GetUsers)
		api.POST("/user/profile", userHandler.GetUser)

		// 需要认证的路由
		authorized := api.Group("/")
		authorized.Use(middlewares.AuthMiddleware(logger))
		{
			authorized.POST("/user/create", userHandler.CreateUser)
			authorized.POST("/user/update", userHandler.UpdateUser)
			authorized.POST("/user/delete", userHandler.DeleteUser)
		}

		// Action routes
		actionRouter := api.Group("/workout_action")
		{
			actionHandler := handlers.NewWorkoutActionHandler(db, logger)
			actionRouter.POST("/list", actionHandler.GetActions)
			actionRouter.POST("/profile", actionHandler.GetAction)
			actionRouter.POST("/list/by_muscle", actionHandler.GetActionsByMuscle)
			actionRouter.POST("/list/by_level", actionHandler.GetActionsByLevel)
			actionRouter.POST("/list/related", actionHandler.GetRelatedActions)

			// Protected action routes (require authentication)
			actionRouter.Use(middlewares.AuthMiddleware(logger))
			{
				actionRouter.POST("/create", actionHandler.CreateAction)
				actionRouter.POST("/update", actionHandler.UpdateAction)
				actionRouter.POST("/delete", actionHandler.DeleteAction)
			}
		}
	}

	// Coach routes
	coachRouter := api.Group("/coach")
	{
		coachHandler := handlers.NewCoachHandler(db, logger)
		coachRouter.POST("/register", coachHandler.RegisterCoach)
		coachRouter.POST("/login", coachHandler.LoginCoach)
		coachRouter.POST("/send-verification-code", coachHandler.SendVerificationCode)

		// Protected coach routes (require authentication)
		coachRouter.Use(middlewares.AuthMiddleware(logger))
		{
			coachRouter.POST("/profile", coachHandler.GetCoachProfile)
			coachRouter.POST("/update", coachHandler.UpdateCoachProfile)
		}
	}

	// Muscle routes
	muscleRouter := api.Group("/muscle")
	{
		muscleHandler := handlers.NewMuscleHandler(db, logger)
		muscleRouter.POST("/list", muscleHandler.GetMuscles)
		muscleRouter.POST("/profile", muscleHandler.GetMuscle)

		// Protected muscle routes (require authentication)
		muscleRouter.Use(middlewares.AuthMiddleware(logger))
		{
			muscleRouter.POST("/create", muscleHandler.CreateMuscle)
			muscleRouter.POST("/update", muscleHandler.UpdateMuscle)
			muscleRouter.POST("/delete", muscleHandler.DeleteMuscle)
		}
	}

	// Equipment routes
	equipmentRouter := api.Group("/equipment")
	{
		equipmentHandler := handlers.NewEquipmentHandler(db, logger)
		equipmentRouter.POST("/list", equipmentHandler.GetEquipments)
		equipmentRouter.POST("/profile", equipmentHandler.GetEquipment)

		// Protected equipment routes (require authentication)
		equipmentRouter.Use(middlewares.AuthMiddleware(logger))
		{
			equipmentRouter.POST("/create", equipmentHandler.CreateEquipment)
		}
	}

	// Workout plan routes
	workoutPlanRouter := api.Group("/workout_plan")
	{
		workoutPlanHandler := handlers.NewWorkoutPlanHandler(db, logger)
		workoutPlanRouter.POST("/create", workoutPlanHandler.CreateWorkoutPlan)
		workoutPlanRouter.POST("/list", workoutPlanHandler.ListWorkoutPlans)
		workoutPlanRouter.POST("/profile", workoutPlanHandler.GetWorkoutPlan)
		workoutPlanRouter.POST("/update", workoutPlanHandler.UpdateWorkoutPlan)
		workoutPlanRouter.POST("/delete", workoutPlanHandler.DeleteWorkoutPlan)
	}

	return r
}
