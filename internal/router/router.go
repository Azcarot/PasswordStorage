package router

import (
	"net/http"

	"github.com/Azcarot/PasswordStorage/internal/cfg"
	"github.com/Azcarot/PasswordStorage/internal/handlers"
	"github.com/Azcarot/PasswordStorage/internal/middleware"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

// Flag - полученные флаги
var Flag cfg.Flags

// MakeRouter - создание chi роутера со всеми ручками и миддварями
func MakeRouter(flag cfg.Flags) *chi.Mux {

	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()
	middleware.Sugar = *logger.Sugar()
	r := chi.NewRouter()
	r.Use(middleware.WithLogging)
	r.Route("/api/user", func(r chi.Router) {
		r.Post("/register", http.HandlerFunc(handlers.Registration))
		r.Post("/login", http.HandlerFunc(handlers.LoginUser))
		r.With(middleware.CheckAuthorization, middleware.AddParamToContext(flag.SecretKey)).Post("/card/add", http.HandlerFunc(handlers.AddNewCard))
		r.With(middleware.CheckAuthorization, middleware.AddParamToContext(flag.SecretKey)).Post("/card/update", http.HandlerFunc(handlers.UpdateCard))
		r.With(middleware.CheckAuthorization).Post("/card/delete", http.HandlerFunc(handlers.DeleteCard))
		r.With(middleware.CheckAuthorization, middleware.AddParamToContext(flag.SecretKey)).Get("/card/search", http.HandlerFunc(handlers.SearchBankCard))
		r.With(middleware.CheckAuthorization, middleware.AddParamToContext(flag.SecretKey)).Get("/card/sync", http.HandlerFunc(handlers.SyncBankData))
		r.With(middleware.CheckAuthorization, middleware.AddParamToContext(flag.SecretKey)).Get("/card/all", http.HandlerFunc(handlers.GetAllBankCards))
		r.With(middleware.CheckAuthorization, middleware.AddParamToContext(flag.SecretKey)).Get("/card/get", http.HandlerFunc(handlers.GetBankCard))
		r.With(middleware.CheckAuthorization, middleware.AddParamToContext(flag.SecretKey)).Post("/lpw/add", http.HandlerFunc(handlers.AddNewLoginPw))
		r.With(middleware.CheckAuthorization, middleware.AddParamToContext(flag.SecretKey)).Post("/lpw/update", http.HandlerFunc(handlers.UpdateLoginPW))
		r.With(middleware.CheckAuthorization).Post("/lpw/delete", http.HandlerFunc(handlers.DeleteLoginPW))
		r.With(middleware.CheckAuthorization, middleware.AddParamToContext(flag.SecretKey)).Get("/lpw/search", http.HandlerFunc(handlers.SearchLoginPW))
		r.With(middleware.CheckAuthorization, middleware.AddParamToContext(flag.SecretKey)).Get("/lpw/sync", http.HandlerFunc(handlers.SyncLPWData))
		r.With(middleware.CheckAuthorization, middleware.AddParamToContext(flag.SecretKey)).Get("/lpw/all", http.HandlerFunc(handlers.GetAllLoginPWs))
		r.With(middleware.CheckAuthorization, middleware.AddParamToContext(flag.SecretKey)).Get("/lpw/get", http.HandlerFunc(handlers.GetLoginPW))
		r.With(middleware.CheckAuthorization, middleware.AddParamToContext(flag.SecretKey)).Post("/text/add", http.HandlerFunc(handlers.AddNewText))
		r.With(middleware.CheckAuthorization, middleware.AddParamToContext(flag.SecretKey)).Post("/text/update", http.HandlerFunc(handlers.UpdateText))
		r.With(middleware.CheckAuthorization).Post("/text/delete", http.HandlerFunc(handlers.DeleteText))
		r.With(middleware.CheckAuthorization, middleware.AddParamToContext(flag.SecretKey)).Get("/text/search", http.HandlerFunc(handlers.SearchText))
		r.With(middleware.CheckAuthorization, middleware.AddParamToContext(flag.SecretKey)).Get("/text/all", http.HandlerFunc(handlers.GetAllTexts))
		r.With(middleware.CheckAuthorization, middleware.AddParamToContext(flag.SecretKey)).Get("/text/sync", http.HandlerFunc(handlers.SyncTextData))
		r.With(middleware.CheckAuthorization, middleware.AddParamToContext(flag.SecretKey)).Get("/text/get", http.HandlerFunc(handlers.GetText))
		r.With(middleware.CheckAuthorization, middleware.AddParamToContext(flag.SecretKey)).Post("/file/add", http.HandlerFunc(handlers.AddNewFile))
		r.With(middleware.CheckAuthorization, middleware.AddParamToContext(flag.SecretKey)).Post("/file/update", http.HandlerFunc(handlers.UpdateFile))
		r.With(middleware.CheckAuthorization).Post("/file/delete", http.HandlerFunc(handlers.DeleteFile))
		r.With(middleware.CheckAuthorization, middleware.AddParamToContext(flag.SecretKey)).Get("/file/search", http.HandlerFunc(handlers.SearchFile))
		r.With(middleware.CheckAuthorization, middleware.AddParamToContext(flag.SecretKey)).Get("/file/sync", http.HandlerFunc(handlers.SyncFileData))
		r.With(middleware.CheckAuthorization, middleware.AddParamToContext(flag.SecretKey)).Get("/file/all", http.HandlerFunc(handlers.GetAllFiles))
		r.With(middleware.CheckAuthorization, middleware.AddParamToContext(flag.SecretKey)).Get("/file/get", http.HandlerFunc(handlers.GetFile))
	})
	return r
}
