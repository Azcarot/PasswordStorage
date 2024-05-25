package router

import (
	"net/http"

	"github.com/Azcarot/PasswordStorage/internal/handlers"
	"github.com/Azcarot/PasswordStorage/internal/middleware"
	"github.com/Azcarot/PasswordStorage/internal/utils"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

var Flag utils.Flags

func MakeRouter(flag utils.Flags) *chi.Mux {

	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()
	middleware.Sugar = *logger.Sugar()
	r := chi.NewRouter()
	// ticker := time.NewTicker(2 * time.Second)
	// quit := make(chan struct{})
	// go func() {
	// 	for {
	// 		select {
	// 		case <-ticker.C:
	// 			handlers.ActualiseOrders(flag)
	// 		case <-quit:
	// 			ticker.Stop()
	// 			return
	// 		}
	// 	}
	// }()
	r.Use(middleware.WithLogging)
	r.Route("/api/user", func(r chi.Router) {
		r.Post("/register", http.HandlerFunc(handlers.Registration))
		r.Post("/login", http.HandlerFunc(handlers.LoginUser))
		r.With(middleware.CheckAuthorization).Post("/card/add", http.HandlerFunc(handlers.AddNewCard))
		r.With(middleware.CheckAuthorization).Post("/card/update", http.HandlerFunc(handlers.UpdateCard))
		r.With(middleware.CheckAuthorization).Post("/card/delete", http.HandlerFunc(handlers.DeleteCard))
		r.With(middleware.CheckAuthorization).Get("/card/search", http.HandlerFunc(handlers.SearchBankCard))
		r.With(middleware.CheckAuthorization).Get("/card/all", http.HandlerFunc(handlers.GetAllBankCards))
		r.With(middleware.CheckAuthorization).Get("/card/get", http.HandlerFunc(handlers.GetBankCard))
		r.With(middleware.CheckAuthorization).Post("/lpw/add", http.HandlerFunc(handlers.AddNewLoginPw))
		r.With(middleware.CheckAuthorization).Post("/lpw/update", http.HandlerFunc(handlers.UpdateLoginPW))
		r.With(middleware.CheckAuthorization).Post("/lpw/delete", http.HandlerFunc(handlers.DeleteLoginPW))
		r.With(middleware.CheckAuthorization).Get("/lpw/search", http.HandlerFunc(handlers.SearchLoginPW))
		r.With(middleware.CheckAuthorization).Get("/lpw/all", http.HandlerFunc(handlers.GetAllLoginPWs))
		r.With(middleware.CheckAuthorization).Get("/lpw/get", http.HandlerFunc(handlers.GetLoginPW))
		r.With(middleware.CheckAuthorization).Post("/text/add", http.HandlerFunc(handlers.AddNewText))
		r.With(middleware.CheckAuthorization).Post("/text/update", http.HandlerFunc(handlers.UpdateText))
		r.With(middleware.CheckAuthorization).Post("/text/delete", http.HandlerFunc(handlers.DeleteText))
		r.With(middleware.CheckAuthorization).Get("/text/search", http.HandlerFunc(handlers.SearchText))
		r.With(middleware.CheckAuthorization).Get("/text/all", http.HandlerFunc(handlers.GetAllTexts))
		r.With(middleware.CheckAuthorization).Get("/text/get", http.HandlerFunc(handlers.GetText))
		r.With(middleware.CheckAuthorization).Post("/file/add", http.HandlerFunc(handlers.AddNewFile))
		r.With(middleware.CheckAuthorization).Post("/file/update", http.HandlerFunc(handlers.UpdateFile))
		r.With(middleware.CheckAuthorization).Post("/file/delete", http.HandlerFunc(handlers.DeleteFile))
		r.With(middleware.CheckAuthorization).Get("/file/search", http.HandlerFunc(handlers.SearchFile))
		r.With(middleware.CheckAuthorization).Get("/file/all", http.HandlerFunc(handlers.GetAllFiles))
		r.With(middleware.CheckAuthorization).Get("/file/get", http.HandlerFunc(handlers.GetFile))
		// r.With(middleware.CheckAuthorization).Get("/withdrawals", http.HandlerFunc(handlers.GetWithdrawals))
	})
	return r
}
