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
	authorized.Use(middlewares.AuthMiddleware(logger, cfg))
	{
		// 用户处理器
		// userHandler := handlers.NewUserHandler(db, logger)

		// 公开路由
		// api.POST("/user/list", userHandler.GetUsers)
		// api.POST("/user/profile", userHandler.GetUser)

		{
			handler := handlers.NewCoachHandler(db, logger, cfg)
			api.POST("/auth/web_register", handler.RegisterCoach)
			api.POST("/auth/web_login", handler.LoginCoach)
			api.GET("/ping", handler.FetchVersion)
			// api.POST("/coach/send-verification-code", handler.SendVerificationCode)

			authorized.POST("/auth/profile", handler.FetchCoachProfile)
			authorized.POST("/auth/update_profile", handler.UpdateCoachProfile)
			authorized.POST("/auth/refresh_token", handler.RefreshToken)
			authorized.POST("/auth/create_account", handler.CreateAccount)
			authorized.POST("/refresh_workout_stats", handler.RefreshCoachStats)
			authorized.POST("/refresh_workout_action_stats", handler.RefreshWorkoutActionStats)
			authorized.POST("/student/create", handler.CreateStudent)
			authorized.POST("/student/list", handler.FetchStudentList)
			authorized.POST("/student/profile", handler.FetchStudentProfile)
			authorized.POST("/student/update", handler.UpdateStudentProfile)
			authorized.POST("/student/delete", handler.DeleteStudent)
			authorized.POST("/student/auth_url", handler.BuildStudentAuthURL)
			authorized.POST("/student/url_with_token", handler.DeleteStudent)
			authorized.POST("/student/to_friend", handler.StudentToFriend)
			authorized.POST("/friend/add", handler.AddFriend)
			authorized.POST("/coach/profile", handler.FetchCoachProfileInWechat)
			authorized.POST("/content/create", handler.CreateArticle)
			authorized.POST("/content/update", handler.UpdateArticle)
			authorized.POST("/content/list", handler.FetchArticleList)
			authorized.POST("/content/profile", handler.FetchArticleProfile)
			authorized.POST("/follow", handler.FollowCoach)
			authorized.POST("/my/follower/list", handler.FetchMyFollowerList)
			authorized.POST("/my/following/list", handler.FetchMyFollowingList)

			// 管理后台
			authorized.POST("/coach/list", handler.FetchCoachList)
			authorized.POST("/coach/create", handler.CreateCoach)
			authorized.POST("/coach/content/list", handler.FetchCoachContentList)
			authorized.POST("/coach/content/create", handler.CreateCoachContent)
			authorized.POST("/admin/coach/auth_url", handler.BuildCoachAuthURLInAdmin)
			authorized.POST("/admin/coach/profile", handler.FetchCoachProfileInAdmin)
		}
		{

			handler := handlers.NewWorkoutPlanHandler(db, logger)
			authorized.POST("/workout_plan/profile", handler.FetchWorkoutPlanProfile)
			authorized.POST("/workout_plan/list", handler.FetchWorkoutPlanList)
			authorized.POST("/workout_plan/update", handler.UpdateWorkoutPlan)
			authorized.POST("/workout_plan/delete", handler.DeleteWorkoutPlan)
			authorized.POST("/workout_plan/create", handler.CreateWorkoutPlan)
			authorized.POST("/workout_plan/mine", handler.FetchMyWorkoutPlanList)
			// authorized.POST("/workout_plan/stats", handler.FetchMyWorkoutPlanList)
			authorized.POST("/workout_plan/content/list", handler.FetchContentListOfWorkoutPlan)
			authorized.POST("/workout_plan/content/profile", handler.FetchContentProfileOfWorkoutPlan)
			authorized.POST("/workout_plan/content/create", handler.CreateContentWithWorkoutPlan)
			// 周期计划
			authorized.POST("/workout_schedule/list", handler.FetchWorkoutScheduleList)
			authorized.POST("/workout_schedule/create", handler.CreateWorkoutSchedule)
			authorized.POST("/workout_schedule/update", handler.UpdateWorkoutSchedule)
			authorized.POST("/workout_schedule/profile", handler.FetchWorkoutScheduleProfile)
			authorized.POST("/workout_schedule/apply", handler.ApplyWorkoutSchedule)
			authorized.POST("/workout_schedule/cancel", handler.CancelWorkoutSchedule)
			authorized.POST("/workout_schedule/enabled", handler.FetchAppliedWorkoutScheduleList)
			// 计划合集
			authorized.POST("/workout_plan_set/list", handler.FetchWorkoutPlanSetList)
			authorized.POST("/workout_plan_set/create", handler.CreateWorkoutPlanSet)
			authorized.POST("/workout_plan_set/update", handler.UpdateWorkoutPlanSet)
		}
		{
			handler := handlers.NewWorkoutDayHandler(db, logger)
			authorized.POST("/workout_day/list", handler.FetchWorkoutDayList)
			authorized.POST("/workout_day/create", handler.CreateWorkoutDay)
			authorized.POST("/workout_day/create_free", handler.CreateFreeWorkoutDay)
			authorized.POST("/workout_day/update", handler.UpdateWorkoutDay)
			authorized.POST("/workout_day/profile", handler.FetchWorkoutDayProfile)
			authorized.POST("/workout_day/has_started", handler.CheckHasStartedWorkoutDay)
			authorized.POST("/workout_day/started_list", handler.FetchStartedWorkoutDay)
			authorized.POST("/workout_day/finished_list", handler.FetchFinishedWorkoutDayList)
			authorized.POST("/workout_day/start", handler.StartWorkoutDay)
			authorized.POST("/workout_day/give_up", handler.GiveUpWorkoutDay)
			authorized.POST("/workout_day/finish", handler.FinishWorkoutDay)
			authorized.POST("/workout_day/result", handler.FetchWorkoutDayResult)
			authorized.POST("/workout_day/continue", handler.ContinueWorkoutDay)
			authorized.POST("/workout_day/update_steps", handler.UpdateWorkoutDayStepProgress)
			authorized.POST("/workout_day/update_details", handler.UpdateWorkoutDayPlanDetails)
			authorized.POST("/workout_day/delete", handler.DeleteWorkoutDay)
			authorized.POST("/student/workout_day/list", handler.FetchMyStudentWorkoutDayList)
			authorized.POST("/student/workout_day/profile", handler.FetchStudentWorkoutDayProfile)
			authorized.POST("/admin/workout_day/refresh_250630", handler.RefreshWorkoutDayRecords250630)
		}
		{
			handler := handlers.NewWorkoutActionHistoryHandler(db, logger)
			authorized.POST("/workout_action_history/create", handler.CreateWorkoutHistory)
			authorized.POST("/workout_action_history/list_of_workout_day", handler.FetchWorkoutActionHistoryListOfWorkoutDay)
			authorized.POST("/workout_action_history/list_of_workout_action", handler.FetchWorkoutActionHistoryListOfWorkoutAction)
			authorized.POST("/student/workout_action_history/list", handler.FetchStudentWorkoutActionHistoryListOfWorkoutDay)
		}
		{
			handler := handlers.NewWorkoutActionHandler(db, logger)
			authorized.POST("/workout_action/list", handler.FetchWorkoutActionList)
			authorized.POST("/workout_action/list_by_ids", handler.FetchWorkoutActionListByIds)
			authorized.POST("/workout_action/list/by_muscle", handler.GetActionsByMuscle)
			authorized.POST("/workout_action/list/by_level", handler.FetchWorkoutActionsByLevel)
			authorized.POST("/workout_action/list/cardio", handler.FetchCardioWorkoutActionList)
			authorized.POST("/workout_action/list/related", handler.FetchRelatedWorkoutActions)
			authorized.POST("/workout_action/profile", handler.GetWorkoutAction)
			authorized.POST("/workout_action/update_idx", handler.UpdateWorkoutActionIdx)
			authorized.POST("/workout_action/create", handler.CreateWorkoutAction)
			authorized.POST("/workout_action/update", handler.UpdateWorkoutActionProfile)
			authorized.POST("/workout_action/delete", handler.DeleteWorkoutAction)
			authorized.POST("/workout_action/content/create", handler.CreateContentWithWorkoutAction)
			authorized.POST("/workout_action/content/list", handler.FetchContentListOfWorkoutAction)
		}
		{
			handler := handlers.NewMuscleHandler(db, logger)
			authorized.POST("/muscle/list", handler.FetchMuscleList)
			authorized.POST("/muscle/profile", handler.FetchMuscleProfile)
			authorized.POST("/muscle/create", handler.CreateMuscle)
			authorized.POST("/muscle/update", handler.UpdateMuscle)
			authorized.POST("/muscle/delete", handler.DeleteMuscle)
		}
		{
			handler := handlers.NewEquipmentHandler(db, logger)
			authorized.POST("/equipment/list", handler.FetchEquipmentList)
			authorized.POST("/equipment/profile", handler.FetchEquipment)
			authorized.POST("/equipment/create", handler.CreateEquipment)
			authorized.POST("/equipment/update", handler.UpdateEquipment)
			authorized.POST("/equipment/delete", handler.DeleteEquipment)
		}
		{
			handler := handlers.NewSubscriptionHandler(db, logger)
			authorized.POST("/subscription_plan/list", handler.FetchSubscriptionPlanList)
			authorized.POST("/subscription_plan/create", handler.CreateSubscriptionPlan)
			authorized.POST("/subscription_order/calc", handler.CalcSubscriptionOrderAmount)
			authorized.POST("/subscription/list", handler.FetchSubscriptionList)
		}
		{
			handler := handlers.NewQuizHandler(db, logger)
			authorized.POST("/quiz/list", handler.FetchQuizList)
			authorized.POST("/quiz/create", handler.CreateQuiz)
			authorized.POST("/paper/list", handler.FetchPaperList)
			authorized.POST("/paper/profile", handler.FetchPaperProfile)
			authorized.POST("/paper/create", handler.CreatePaper)
			authorized.POST("/paper/update", handler.UpdatePaper)
			authorized.POST("/exam/running", handler.FetchRunningExam)
			authorized.POST("/exam/list", handler.FetchExamList)
			authorized.POST("/exam/start", handler.StartExamWithPaper)
			authorized.POST("/exam/profile", handler.FetchExamProfile)
			authorized.POST("/exam/answer", handler.UpdateQuizAnswer)
			authorized.POST("/exam/complete", handler.CompleteExam)
			authorized.POST("/exam/give_up", handler.GiveUpExam)
			authorized.POST("/exam/result", handler.FetchExamResult)
		}
		{
			handler := handlers.NewReportHandler(db, logger)
			authorized.POST("/report/create", handler.CreateReport)
			authorized.POST("/report/profile", handler.FetchReportProfile)
			authorized.POST("/report/list", handler.FetchReportList)
			authorized.POST("/report/list_of_mine", handler.FetchMineReportList)
		}
		{
			handler := handlers.NewGiftCardHandler(db, logger)
			authorized.POST("/gift_card/create", handler.CreateGiftCard)
			authorized.POST("/gift_card/create_reward", handler.CreateGiftCardReward)
			authorized.POST("/gift_card/list", handler.FetchGiftCardList)
			authorized.POST("/gift_card/reward_list", handler.FetchGiftCardRewardList)
			authorized.POST("/gift_card/profile", handler.FetchGiftCardProfile)
			authorized.POST("/gift_card/using", handler.UsingGiftCard)
			authorized.POST("/gift_card/send", handler.SendGiftCard)
		}
		{
			handler := handlers.NewMediaResourceHandler(db, logger, cfg)
			api.POST("/media/qiniu_token", handler.BuildQiniuToken)
			authorized.POST("/media/create", handler.CreateMediaResource)
			authorized.POST("/media/list", handler.FetchMediaResourceList)
			authorized.POST("/media/delete", handler.DeleteMediaResource)
		}
	}

	return r
}
