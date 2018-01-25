package frontend

import (
	"encoding/json"
	"net/http"

	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"

	"github.com/hwgo/pher/httperr"
	"github.com/hwgo/pher/log"
	"github.com/hwgo/pher/metrics"
	"github.com/hwgo/pher/tracing"

	"github.com/hwgo/customer"
)

// Server implements jaeger-demo-frontend service
type Server struct {
	hostPort string
	tracer   opentracing.Tracer
	logger   log.Factory
	bestETA  *bestETA
}

func NewServer(hostPort string) *Server {
	logger := log.Service(ServiceName)
	tracer := tracing.Init(ServiceName, metrics.Namespace(ServiceName, nil), logger)

	return &Server{
		hostPort: hostPort,
		tracer:   tracer,
		logger:   logger,
		bestETA:  newBestETA(tracer, logger),
	}
}

func (s *Server) Run() error {
	mux := s.createServeMux()
	s.logger.Bg().Info("Starting", zap.String("address", "http://"+s.hostPort))
	return http.ListenAndServe(s.hostPort, mux)
}

func (s *Server) createServeMux() http.Handler {
	mux := tracing.NewServeMux(s.tracer)
	mux.Handle("/", http.HandlerFunc(s.home))
	mux.Handle("/dispatch", http.HandlerFunc(s.dispatch))
	return mux
}

func (s *Server) home(w http.ResponseWriter, r *http.Request) {
	s.logger.For(r.Context()).Info("HTTP", zap.String("method", r.Method), zap.Stringer("url", r.URL))
	http.ServeFile(w, r, "/Users/luotao/.go/src/github.com/hwgo/frontend/web_assets/index.html")
}

func (s *Server) dispatch(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	s.logger.For(ctx).Info("HTTP request received", zap.String("method", r.Method), zap.Stringer("url", r.URL))
	if err := r.ParseForm(); httperr.HandleError(w, err, http.StatusBadRequest) {
		s.logger.For(ctx).Error("bad request", zap.Error(err))
		return
	}

	customerID := r.Form.Get("customer")
	if customerID == "" {
		http.Error(w, "Missing required 'customer' parameter", http.StatusBadRequest)
		return
	}

	customerClient := customer.NewClient(s.tracer, s.logger)
	defer customerClient.Close()

	customerClient.LoggerFactory().For(ctx).Info("xxoo @ frontend")
	user := customerClient.Get(ctx)

	s.logger.For(ctx).Info("Load Customer From gRPC", zap.String("name", user.Name))

	// // TODO distinguish between user errors (such as invalid customer ID) and server failures
	response, err := s.bestETA.Get(ctx, customerID)
	if httperr.HandleError(w, err, http.StatusInternalServerError) {
		s.logger.For(ctx).Error("request failed", zap.Error(err))
		return
	}

	data, err := json.Marshal(response)
	if httperr.HandleError(w, err, http.StatusInternalServerError) {
		s.logger.For(ctx).Error("cannot marshal response", zap.Error(err))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}
