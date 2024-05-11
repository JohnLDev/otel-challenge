package middlewares

import (
	"context"
	"net/http"

	"github.com/google/uuid"
)

type contextValues string

const (
	CtxRequestId contextValues = "requestId"
)

func RequestId(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), CtxRequestId, uuid.NewString())))
	})
}

func GetRequestId(ctx context.Context) string {
	return ctx.Value(CtxRequestId).(string)
}
