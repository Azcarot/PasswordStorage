// Package middleware - пакет со всеми миддлвейр (добавление секрета, проверка атвторизации, логгирование запросов)
package middleware

import (
	"context"
	"net/http"

	"github.com/Azcarot/PasswordStorage/internal/storage"
)

// AddParamToContext - ручка для добавления секрета сервера в контект запроса
func AddParamToContext(data [16]byte) func(h http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {

		addata := func(res http.ResponseWriter, req *http.Request) {
			ctx := context.WithValue(req.Context(), storage.EncryptionCtxKey, data)
			req = req.WithContext(ctx)
			next.ServeHTTP(res, req)
		}
		return http.HandlerFunc(addata)
	}
}
