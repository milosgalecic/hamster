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

	router.HandlerFunc(http.MethodGet, "/", app.home)
	router.HandlerFunc(http.MethodGet, "/vozaci", app.drivers)
	router.HandlerFunc(http.MethodGet, "/vozaci/info/:id", app.driver_info)
	router.HandlerFunc(http.MethodGet, "/vozaci/edit/:id", app.edit_driver)
	router.HandlerFunc(http.MethodPost, "/vozaci/edit/:id", app.edit_driver_post)
	router.HandlerFunc(http.MethodGet, "/vozaci/dodaj_vozaca", app.add_driver)
	router.HandlerFunc(http.MethodPost, "/vozaci/dodaj_vozaca", app.add_driver_post)
	router.HandlerFunc(http.MethodGet, "/kamioni", app.trucks)
	router.HandlerFunc(http.MethodGet, "/kamioni/info/:id", app.truck_info)
	router.HandlerFunc(http.MethodGet, "/kamioni/edit/:id", app.edit_truck)
	router.HandlerFunc(http.MethodPost, "/kamioni/edit/:id", app.edit_truck_post)
	router.HandlerFunc(http.MethodGet, "/kamioni/dodaj_kamion", app.add_truck)
	router.HandlerFunc(http.MethodPost, "/kamioni/dodaj_kamion", app.add_truck_post)
	router.HandlerFunc(http.MethodGet, "/poslovi", app.jobs)
	router.HandlerFunc(http.MethodGet, "/poslovi/info/:id", app.job_info)
	router.HandlerFunc(http.MethodGet, "/poslovi/edit/:id", app.edit_job)
	router.HandlerFunc(http.MethodPost, "/poslovi/edit/:id", app.edit_job_post)
	router.HandlerFunc(http.MethodGet, "/poslovi/dodaj_posao", app.add_job)
	router.HandlerFunc(http.MethodPost, "/poslovi/dodaj_posao", app.add_job_post)

	standard := alice.New(app.recoverPanic, app.logRequest, secureHeaders)

	return standard.Then(router)
}
