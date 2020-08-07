package restapi

import (
	"go.uber.org/zap"
	"net/http"
	"payment-system/internal/common"
	"time"
)

func Logger(logger *zap.Logger, inner http.Handler, name string) http.Handler {
	log := logger.With(zap.String(common.Request, name))
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		inner.ServeHTTP(w, r)
		checkedEntry := log.Check(zap.DebugLevel, "")
		if checkedEntry != nil {
			checkedEntry.Write(zap.String("", r.Method), zap.String("", r.RequestURI), zap.Duration("", time.Since(start)))
		}
	})
}
