package routes

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"myapi/config"
	"myapi/internal/models"
	"myapi/pkg/logger"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"
)

func test_catalog_router(t *testing.T) (*gin.Engine, *gorm.DB, *config.Config) {
	t.Helper()
	gin.SetMode(gin.TestMode)
	database, err := gorm.Open(sqlite.Open(filepath.Join(t.TempDir(), "test.db")), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	if err := database.AutoMigrate(&models.Coach{}, &models.WorkoutPlan{}, &models.WorkoutAction{}, &models.WorkoutSchedule{}, &models.CoachContentWithWorkoutPlan{}); err != nil {
		t.Fatal(err)
	}
	sql_database, _ := database.DB()
	t.Cleanup(func() { sql_database.Close() })
	cfg := &config.Config{TokenSecretKey: "test-secret"}
	log := logger.NewLogger("error")
	t.Cleanup(func() { log.Sync() })
	return SetupRouter(database, log, cfg), database, cfg
}
func catalog_request(t *testing.T, router *gin.Engine, url, body, token string) map[string]interface{} {
	t.Helper()
	request := httptest.NewRequest("POST", url, strings.NewReader(body))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", token)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)
	var response map[string]interface{}
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatalf("status=%d body=%s", recorder.Code, recorder.Body.String())
	}
	return response
}
func TestGuestCatalogAndProtectedTraining(t *testing.T) {
	router, database, cfg := test_catalog_router(t)
	for _, plan := range []models.WorkoutPlan{{Id: 1, Status: 1, Title: "public"}, {Id: 2, Status: 2, Title: "private", OwnerId: 7}, {Id: 3, Status: 2, Title: "unowned-private"}, {Id: 4, Status: 1, D: 1, Title: "deleted"}} {
		if err := database.Create(&plan).Error; err != nil {
			t.Fatal(err)
		}
	}
	for _, action := range []models.WorkoutAction{{Id: 1, Status: 1, ZhName: "public-action"}, {Id: 2, Status: 2, ZhName: "private-action"}} {
		if err := database.Create(&action).Error; err != nil {
			t.Fatal(err)
		}
	}
	for _, endpoint := range []string{"/api/workout_plan/list", "/api/workout_plan/profile", "/api/workout_action/list", "/api/workout_action/profile", "/api/workout_schedule/list"} {
		response := catalog_request(t, router, endpoint, `{"id":1,"page":1,"page_size":20}`, "")
		if response["code"] != float64(200) {
			t.Fatalf("%s: %v", endpoint, response)
		}
	}
	for _, plan_id := range []string{"2", "3", "4"} {
		response := catalog_request(t, router, "/api/workout_plan/profile", `{"id":`+plan_id+`}`, "")
		if response["code"] == float64(200) {
			t.Fatalf("guest accessed private/deleted plan %s", plan_id)
		}
	}
	response := catalog_request(t, router, "/api/workout_plan/list", `{"page":1,"page_size":20}`, "")
	records := response["data"].(map[string]interface{})["list"].([]interface{})
	if len(records) != 1 {
		t.Fatalf("guest saw %d plans", len(records))
	}
	response = catalog_request(t, router, "/api/workout_action/profile", `{"id":2}`, "")
	if response["code"] == float64(200) {
		t.Fatal("guest saw private action")
	}
	for _, endpoint := range []string{"/api/workout_day/create", "/api/workout_day/start", "/api/workout_schedule/apply", "/api/workout_plan/create", "/api/student/list", "/api/workout_action_history/list_of_workout_action"} {
		response := catalog_request(t, router, endpoint, `{}`, "")
		if response["code"] != float64(401) {
			t.Fatalf("%s: %v", endpoint, response)
		}
	}
	token, _, err := models.GenerateJWT(7, cfg.TokenSecretKey)
	if err != nil {
		t.Fatal(err)
	}
	response = catalog_request(t, router, "/api/workout_plan/profile", `{"id":2}`, "Bearer "+token)
	if response["code"] != float64(200) {
		t.Fatal("owner lost access to private plan", response)
	}
	response = catalog_request(t, router, "/api/workout_plan/profile", `{"id":2}`, "Bearer invalid")
	if response["code"] == float64(200) {
		t.Fatal("invalid identity gained private access")
	}
	response = catalog_request(t, router, "/api/workout_plan/content/list", `{"workout_plan_id":2,"page":1,"page_size":20}`, "")
	if response["code"] == float64(200) {
		t.Fatal("guest saw content of a private plan")
	}
}
