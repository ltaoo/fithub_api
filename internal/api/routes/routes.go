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
			handler := handlers.NewCoachHandler(db, logger)
			authorized.POST("/student/create", handler.CreateStudent)
			authorized.POST("/student/list", handler.FetchMyStudentList)
			authorized.POST("/student/profile", handler.FetchMyStudentProfile)
		}

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
			api.POST("/workout_plan/profile", handler.FetchWorkoutPlanProfile)
			authorized.POST("/workout_plan/update", handler.UpdateWorkoutPlan)
			authorized.POST("/workout_plan/delete", handler.DeleteWorkoutPlan)
			authorized.POST("/workout_plan/create", handler.CreateWorkoutPlan)
			authorized.POST("/workout_plan/mine", handler.FetchMyWorkoutPlanList)
			authorized.POST("/workout_plan_collection/create", handler.CreateWorkoutPlanCollection)
			authorized.POST("/workout_plan_collection/profile", handler.FetchWorkoutPlanCollectionProfile)
			authorized.POST("/workout_plan_set/list", handler.FetchWorkoutPlanSetList)
			authorized.POST("/workout_plan_set/create", handler.CreateWorkoutPlanSet)
			authorized.POST("/workout_plan_set/update", handler.UpdateWorkoutPlanSet)
		}
		{
			handler := handlers.NewWorkoutDayHandler(db, logger)
			authorized.POST("/workout_day/list", handler.FetchWorkoutDayList)
			authorized.POST("/workout_day/create", handler.CreateWorkoutDay)
			authorized.POST("/workout_day/profile", handler.FetchWorkoutDayProfile)
			authorized.POST("/workout_day/fetch_started", handler.FetchStartedWorkoutDay)
			authorized.POST("/workout_day/start", handler.StartWorkoutDay)
			authorized.POST("/workout_day/give_up", handler.GiveUpWorkoutDay)
			authorized.POST("/workout_day/finish", handler.FinishWorkoutDay)
			authorized.POST("/workout_day/update_steps", handler.UpdateWorkoutDayStepProgress)
			authorized.POST("/workout_day/update_details", handler.UpdateWorkoutDayPlanDetails)
			authorized.POST("/workout_day/delete", handler.DeleteWorkoutDay)
		}
		{
			handler := handlers.NewWorkoutActionHistoryHandler(db, logger)
			authorized.POST("/workout_action_history/list", handler.FetchWorkoutActionHistoryList)
		}
		{
			handler := handlers.NewWorkoutActionHandler(db, logger)
			api.POST("/workout_action/list", handler.FetchWorkoutActionList)
			api.POST("/workout_action/list_by_ids", handler.GetWorkoutActionListByIds)
			api.POST("/workout_action/profile", handler.GetWorkoutAction)
			api.POST("/workout_action/list/by_muscle", handler.GetActionsByMuscle)
			api.POST("/workout_action/list/by_level", handler.FetchWorkoutActionsByLevel)
			api.POST("/workout_action/list/related", handler.FetchRelatedWorkoutActions)

			{
				authorized.POST("/workout_action/create", handler.CreateWorkoutAction)
				authorized.POST("/workout_action/update", handler.UpdateWorkoutActionProfile)
				authorized.POST("/workout_action/delete", handler.DeleteWorkoutAction)
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
		{
			handler := handlers.NewSubscriptionHandler(db, logger)
			{
				authorized.POST("/subscription_plan/list", handler.FetchSubscriptionPlanList)
				authorized.POST("/subscription_plan/create", handler.CreateSubscriptionPlan)
				authorized.POST("/subscription_order/calc", handler.CalcSubscriptionOrderAmount)
			}
		}
		{
			handler := handlers.NewQuizHandler(db, logger)
			{
				authorized.POST("/quiz/list", handler.FetchQuizList)
				authorized.POST("/quiz/create", handler.CreateQuiz)
				authorized.POST("/paper/list", handler.FetchPaperList)
				authorized.POST("/paper/profile", handler.FetchPaperProfile)
				authorized.POST("/paper/create", handler.CreatePaper)
				authorized.POST("/paper/update", handler.UpdatePaper)
				authorized.POST("/exam/running", handler.FetchRunningExam)
				authorized.POST("/exam/start", handler.StartExamWithPaper)
				authorized.POST("/exam/profile", handler.FetchExamProfile)
				authorized.POST("/exam/answer", handler.UpdateQuizAnswer)
				authorized.POST("/exam/complete", handler.CompleteExam)
				authorized.POST("/exam/give_up", handler.GiveUpExam)
				authorized.POST("/exam/result", handler.FetchExamResult)
			}
		}
	}

	return r
}
