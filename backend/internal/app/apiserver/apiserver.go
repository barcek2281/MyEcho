package apiserver

import (
	"net/http"

	"github.com/barcek2281/MyEcho/internal/app/mail"
	storage "github.com/barcek2281/MyEcho/internal/app/store"
	"github.com/gorilla/sessions"
	"github.com/sirupsen/logrus"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
)

func Start(config *Config) error {
	store := storage.New(config.DataBaseURL)
	if err := store.Open(); err != nil { // Ping db
		return err
	}
	session := sessions.NewCookieStore([]byte(config.CookieKey))
	session.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   3600,  // Время жизни сессии в секундах
		HttpOnly: true,  // Куки недоступны через JavaScript
		Secure:   false, // Для HTTP, установи false. Для HTTPS, установи true.
		SameSite: http.SameSiteLaxMode,
	}

	logger := logrus.New()
	logger.SetFormatter(&prefixed.TextFormatter{
		DisableColors:   true,
		TimestampFormat: "2006-01-02 15:04:05",
		FullTimestamp:   true,
		ForceFormatting: true,
	})
	level, err := logrus.ParseLevel(config.LogLevel)
	if err != nil {
		return err
	}
	// f, err := os.OpenFile(config.LogFilePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o666)
	// if err != nil {
	// 	return err
	// } else {
	// 	logger.Out = f
	// }
	

	logger.SetLevel(level)

	sender := mail.NewSender(config.EmailTo, config.EmailToPassword)

	s := newServer(store, session, logger, sender)
	// return http.ListenAndServe("192.168.42.101"+config.BinAddr, s)
	return http.ListenAndServe(config.BinAddr, s)
}
