package api

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"github.com/ekristen/pipeliner/pkg/api/hub"
	"github.com/ekristen/pipeliner/pkg/api/routes"
	"github.com/ekristen/pipeliner/pkg/controllers"
	"github.com/ekristen/pipeliner/pkg/models"
	"github.com/ekristen/pipeliner/pkg/store"
)

type apiServer struct {
	ctx          context.Context
	db           *gorm.DB
	log          *logrus.Entry
	port         int
	cookieSecret string

	storage *store.Uploader

	notifyPipelines chan *models.Pipeline
	notifyBuilds    chan *models.Pipeline

	controllerPiplines *controllers.Pipelines
	controllerBuilds   *controllers.Builds

	notifyDatabaseChange chan []byte
}

// NewServer --
func NewServer(
	ctx context.Context,
	db *gorm.DB,
	port int,
	cookieSecret string,
	storage *store.Uploader,
	cP *controllers.Pipelines,
	cB *controllers.Builds,
	nP chan *models.Pipeline,
	nB chan *models.Pipeline,
	notifyDatabaseChange chan []byte,
) *apiServer {
	return &apiServer{
		ctx: ctx,
		db:  db,
		log: logrus.WithField("component", "api"),

		port:         port,
		cookieSecret: cookieSecret,

		storage: storage,

		notifyPipelines: nP,
		notifyBuilds:    nB,

		controllerPiplines: cP,
		controllerBuilds:   cB,

		notifyDatabaseChange: notifyDatabaseChange,
	}
}

func (s *apiServer) Run() error {
	r := mux.NewRouter()

	routes.NewRoutes(s.ctx, s.db, r, s.log, s.storage, s.notifyPipelines, s.notifyBuilds, s.controllerPiplines, s.controllerBuilds)

	store := sessions.NewCookieStore([]byte(s.cookieSecret))
	store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   60 * 60 * 168,
		HttpOnly: true,
	}

	wsHub := hub.New(s.ctx, s.notifyDatabaseChange)

	//log := logrus.WithField("component", "server")
	//rlogger := middleware.NewLogger(log, store)
	//r.Use(rlogger.Middleware)
	r.Path("/v1/ws").HandlerFunc(wsHub.WSHandler)

	http.Handle("/", r)

	srv := &http.Server{
		Addr: fmt.Sprintf(":%d", s.port),
	}

	go wsHub.Run()

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.log.Fatalf("listen: %s\n", err)
		}
	}()
	s.log.Info("Starting API Server")

	<-s.ctx.Done()

	s.log.Info("Shutting down API Server")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		s.log.WithError(err).Error("Unable to shutdown the API server properly")
		return err
	}

	return nil
}
