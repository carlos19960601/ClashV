package route

import (
	"net/http"

	"github.com/carlos19960601/ClashV/adapter/inbound"
	"github.com/carlos19960601/ClashV/log"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/go-chi/render"
)

var (
	serverAddr = ""
)

func Start(addr string, isDebug bool) {
	if serverAddr != "" {
		return
	}

	serverAddr = addr

	l, err := inbound.Listen("tcp", addr)
	if err != nil {
		log.Errorln("External Controller 监听失败：%s", err)
		return
	}

	serverAddr := l.Addr().String()
	log.Infoln("RESTful API listening at %s", serverAddr)

	if err = http.Serve(l, router(isDebug, true)); err != nil {
		log.Errorln("External Controller server 失败: %s", err)
	}
}

func router(isDebug bool, withAuth bool) *chi.Mux {
	r := chi.NewRouter()
	corsM := cors.New(cors.Options{})

	r.Use(corsM.Handler)
	r.Group(func(r chi.Router) {
		r.Get("/", hello)
		r.Get("/logs", getLogs)
		r.Get("traffic", traffic)
	})

	return r
}

func hello(w http.ResponseWriter, r *http.Request) {
	render.JSON(w, r, render.M{"hello": "ClashX"})
}

func getLogs(w http.ResponseWriter, r *http.Request) {}

func traffic(w http.ResponseWriter, r *http.Request) {}
