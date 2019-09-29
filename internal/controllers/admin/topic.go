package admin

import (
	"net/http"
	"satellity/internal/durable"
	"satellity/internal/models"
	"satellity/internal/views"
	"time"

	"github.com/dimfeld/httptreemux"
)

type topicImpl struct {
	database *durable.Database
}

func registerAdminTopic(database *durable.Database, router *httptreemux.TreeMux) {
	impl := &topicImpl{database: database}

	router.GET("/admin/topics", impl.index)
}

func (impl *topicImpl) index(w http.ResponseWriter, r *http.Request, params map[string]string) {
	ctx := models.WrapContext(r.Context(), impl.database)
	offset, _ := time.Parse(time.RFC3339Nano, r.URL.Query().Get("offset"))
	if topics, err := models.ReadTopics(ctx, offset); err != nil {
		views.RenderErrorResponse(w, r, err)
	} else {
		views.RenderTopics(w, r, topics)
	}
}
