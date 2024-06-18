// This file is safe to edit. Once it exists it will not be overwritten

package handler

import (
	"crypto/tls"
	"net/http"
	"os"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
	"github.com/sirupsen/logrus"

	"telegram-bot/internal/handler/operations"
	custMid "telegram-bot/internal/middleware"
	"telegram-bot/internal/worker"
	"telegram-bot/pkg/logger"
	log "telegram-bot/pkg/logger"
	"telegram-bot/pkg/smap"

	tgbotapi "github.com/Syfaro/telegram-bot-api"
	formatters "github.com/fabienm/go-logrus-formatters"

	graylog "github.com/gemnasium/logrus-graylog-hook/v3"
)

const (
	MSGS_BUFFER  = 10
	SERVICE_NAME = "telegram-bot"
)

//go:generate swagger generate server --target ../../../telegram-bot --name TelegramBot --spec ../../api/api.yml --server-package internal/handler

func configureFlags(api *operations.TelegramBotAPI) {
	// api.CommandLineOptionsGroups = []swag.CommandLineOptionsGroup{ ... }
}

func configureAPI(api *operations.TelegramBotAPI) http.Handler {
	// configure the api here
	api.ServeError = errors.ServeError

	fmter := formatters.NewGelf(SERVICE_NAME)
	hooks := []logrus.Hook{graylog.NewGraylogHook(os.Getenv("GRAYLOG_HOST"), map[string]interface{}{})}
	logger := log.NewLogger(fmter, hooks)

	api.Logger = logger.Printf

	api.JSONConsumer = runtime.JSONConsumer()

	api.JSONProducer = runtime.JSONProducer()

	msgs := make(chan string, MSGS_BUFFER)
	users := smap.NewSTMap(100)
	botApi, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_BOT_KEY"))
	if err != nil {
		logger.WithFields(
			logger.TraceWrap(err),
		).Panic(err)
	}
	bot := worker.NewTelBot(logger, users, botApi)
	go bot.InitTelBotConsumer()
	go bot.InitTelBotProducer(msgs)

	api.PostSendHandler = operations.PostSendHandlerFunc(func(params operations.PostSendParams) middleware.Responder {
		err := SendMsg(msgs, params.Msg)
		if err != nil {
			logger.WithFields(
				logger.TraceWrap(err),
			).Error(err)
			return operations.NewPostSendInternalServerError()
		}
		return operations.NewPostSendOK()
	})

	api.PreServerShutdown = func() {}

	api.ServerShutdown = func() {}

	return setupGlobalMiddleware(logger, api.Serve(setupMiddlewares))
}

// The TLS configuration before HTTPS server starts.
func configureTLS(tlsConfig *tls.Config) {
	// Make all necessary changes to the TLS configuration here.
}

// As soon as server is initialized but not run yet, this function will be called.
// If you need to modify a config, store server instance to stop it individually later, this is the place.
// This function can be called multiple times, depending on the number of serving schemes.
// scheme value will be set accordingly: "http", "https" or "unix"
func configureServer(s *http.Server, scheme, addr string) {
}

// The middleware configuration is for the handler executors. These do not apply to the swagger.json document.
// The middleware executes after routing but before authentication, binding and validation
func setupMiddlewares(handler http.Handler) http.Handler {
	return handler
}

// The middleware configuration happens before anything, this middleware also applies to serving the swagger.json document.
// So this is a good place to plug in a panic handling middleware, logging and metrics
func setupGlobalMiddleware(logger logger.Logger, handler http.Handler) http.Handler {
	md := &custMid.Middleware{
		L: logger,
	}
	return md.LogRequests(md.RecoverMiddleware(handler))
}
