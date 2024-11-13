package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

// The routes() method returns a servemux containing our application routes.
func (app *application) routes() http.Handler {
	router := httprouter.New()
	router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.notFound(w)
	})
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	router.Handler(http.MethodGet, "/static/*filepath", http.StripPrefix("/static", fileServer))

	dynamic := alice.New(app.sessionManager.LoadAndSave)

	router.Handler(http.MethodGet, "/", dynamic.ThenFunc(app.home))
	router.Handler(http.MethodGet, "/vozaci", dynamic.ThenFunc(app.drivers))
	router.Handler(http.MethodGet, "/vozaci/info/:id", dynamic.ThenFunc(app.driver_info))
	router.Handler(http.MethodGet, "/vozaci/edit/:id", dynamic.ThenFunc(app.edit_driver))
	router.Handler(http.MethodPost, "/vozaci/edit/:id", dynamic.ThenFunc(app.edit_driver_post))
	router.Handler(http.MethodGet, "/vozaci/dodaj_vozaca", dynamic.ThenFunc(app.add_driver))
	router.Handler(http.MethodPost, "/vozaci/dodaj_vozaca", dynamic.ThenFunc(app.add_driver_post))
	router.Handler(http.MethodGet, "/kamioni", dynamic.ThenFunc(app.trucks))
	router.Handler(http.MethodGet, "/kamioni/info/:id", dynamic.ThenFunc(app.truck_info))
	router.Handler(http.MethodGet, "/kamioni/edit/:id", dynamic.ThenFunc(app.edit_truck))
	router.Handler(http.MethodPost, "/kamioni/edit/:id", dynamic.ThenFunc(app.edit_truck_post))
	router.Handler(http.MethodGet, "/kamioni/dodaj_kamion", dynamic.ThenFunc(app.add_truck))
	router.Handler(http.MethodPost, "/kamioni/dodaj_kamion", dynamic.ThenFunc(app.add_truck_post))
	router.Handler(http.MethodGet, "/poslovi", dynamic.ThenFunc(app.jobs))
	router.Handler(http.MethodGet, "/poslovi/info/:id", dynamic.ThenFunc(app.job_info))
	router.Handler(http.MethodGet, "/poslovi/edit/:id", dynamic.ThenFunc(app.edit_job))
	router.Handler(http.MethodPost, "/poslovi/edit/:id", dynamic.ThenFunc(app.edit_job_post))
	router.Handler(http.MethodGet, "/poslovi/dodaj_posao", dynamic.ThenFunc(app.add_job))
	router.Handler(http.MethodPost, "/poslovi/dodaj_posao", dynamic.ThenFunc(app.add_job_post))

	standard := alice.New(app.recoverPanic, app.logRequest, secureHeaders)

	return standard.Then(router)
}
