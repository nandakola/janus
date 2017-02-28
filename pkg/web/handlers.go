package web

import (
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/hellofresh/janus/pkg/errors"
	"github.com/hellofresh/janus/pkg/response"
)

// Home handler is just a nice home page message
func Home() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		response.JSON(w, http.StatusOK, "Welcome to Janus")
	}
}

// NotFound handler is called when no route is matched
func NotFound(w http.ResponseWriter, r *http.Request) {
	notfoundError := errors.ErrRouteNotFound
	response.JSON(w, notfoundError.Code, notfoundError)
}

// RecoveryHandler handler is used when a panic happens
func RecoveryHandler(w http.ResponseWriter, r *http.Request, err interface{}) {
	switch internalErr := err.(type) {
	case *errors.Error:
		log.WithFields(log.Fields{"code": internalErr.Code, "error": internalErr.Error()}).
			Warning("Internal error hadled")
		response.JSON(w, internalErr.Code, internalErr.Error())
	default:
		log.WithField("error", err).Error("Internal server error handled")
		response.JSON(w, http.StatusInternalServerError, err)
	}
}

// Heartbeat normally is used by the load balancers to identify if the application is alive
func Heartbeat() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		response.JSON(w, http.StatusOK, nil)
	}
}
