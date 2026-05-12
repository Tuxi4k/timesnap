package deadline_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/Tuxi4k/timesnap/internal/modules/deadline"
	"github.com/glebarez/sqlite"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func setupHandlerApp(t *testing.T) *fiber.App {
	t.Helper()

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open db: %v", err)
	}

	if err := db.AutoMigrate(&deadline.Deadline{}); err != nil {
		t.Fatalf("failed to migrate db: %v", err)
	}

	repo := deadline.NewRepository(db)
	svc := deadline.NewService(repo)
	h := deadline.NewHandler(svc)

	app := fiber.New()
	h.RegisterRoutes(app.Group("deadlines/"))

	return app
}

func doJSONRequest(t *testing.T, app *fiber.App, method, path string, body any) *http.Response {
	t.Helper()

	var payload []byte
	if body != nil {
		var err error
		payload, err = json.Marshal(body)
		if err != nil {
			t.Fatalf("failed to marshal body: %v", err)
		}
	}

	req := httptest.NewRequest(method, path, bytes.NewReader(payload))
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}

	return resp
}

func TestHandler_CRUD_Integration(t *testing.T) {
	app := setupHandlerApp(t)
	dueDate := time.Now().Add(2 * time.Hour).UTC().Format(time.RFC3339Nano)

	createBody := map[string]any{
		"title":    "Integration task",
		"status":   "active",
		"priority": "high",
		"due_date": dueDate,
	}

	createResp := doJSONRequest(t, app, http.MethodPost, "/deadlines/", createBody)
	assert.Equal(t, http.StatusCreated, createResp.StatusCode)

	var created deadline.Deadline
	assert.NoError(t, json.NewDecoder(createResp.Body).Decode(&created))
	assert.NotZero(t, created.ID)
	assert.Equal(t, "Integration task", created.Title)

	getResp := doJSONRequest(t, app, http.MethodGet, "/deadlines/"+toStr(created.ID), nil)
	assert.Equal(t, http.StatusOK, getResp.StatusCode)

	var got deadline.Deadline
	assert.NoError(t, json.NewDecoder(getResp.Body).Decode(&got))
	assert.Equal(t, created.ID, got.ID)

	updateBody := map[string]any{
		"status": "completed",
	}
	updateResp := doJSONRequest(t, app, http.MethodPatch, "/deadlines/"+toStr(created.ID), updateBody)
	assert.Equal(t, http.StatusOK, updateResp.StatusCode)

	var updated deadline.Deadline
	assert.NoError(t, json.NewDecoder(updateResp.Body).Decode(&updated))
	assert.Equal(t, deadline.StatusCompleted, updated.Status)

	deleteResp := doJSONRequest(t, app, http.MethodDelete, "/deadlines/"+toStr(created.ID), nil)
	assert.Equal(t, http.StatusNoContent, deleteResp.StatusCode)

	notFoundResp := doJSONRequest(t, app, http.MethodGet, "/deadlines/"+toStr(created.ID), nil)
	assert.Equal(t, http.StatusNotFound, notFoundResp.StatusCode)
}

func TestHandler_ValidationAndBadParams_Integration(t *testing.T) {
	app := setupHandlerApp(t)

	badIDResp := doJSONRequest(t, app, http.MethodGet, "/deadlines/not-a-number", nil)
	assert.Equal(t, http.StatusBadRequest, badIDResp.StatusCode)

	invalidCreateBody := map[string]any{
		"title":    "",
		"status":   "active",
		"priority": "medium",
		"due_date": time.Now().Add(2 * time.Hour).UTC().Format(time.RFC3339Nano),
	}
	invalidCreateResp := doJSONRequest(t, app, http.MethodPost, "/deadlines/", invalidCreateBody)
	assert.Equal(t, http.StatusUnprocessableEntity, invalidCreateResp.StatusCode)

	invalidUpdateBody := map[string]any{
		"status": "unknown",
	}
	invalidUpdateResp := doJSONRequest(t, app, http.MethodPatch, "/deadlines/1", invalidUpdateBody)
	assert.Equal(t, http.StatusUnprocessableEntity, invalidUpdateResp.StatusCode)
}

func toStr(v uint) string {
	return strconv.FormatUint(uint64(v), 10)
}
