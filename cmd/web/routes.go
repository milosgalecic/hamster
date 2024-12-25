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

	router.Handler(http.MethodGet, "/korisnik/registracija", dynamic.ThenFunc(app.add_user))
	router.Handler(http.MethodPost, "/korisnik/registracija", dynamic.ThenFunc(app.add_user_post))
	router.Handler(http.MethodGet, "/korisnik/prijava", dynamic.ThenFunc(app.user_login))
	router.Handler(http.MethodPost, "/korisnik/prijava", dynamic.ThenFunc(app.user_login_post))

	// Protected (authenticated-only) application routes, using a new "protected"
	// middleware chain which includes the requireAuthentication middleware.
	protected := dynamic.Append(app.requireAuthentication)

	router.Handler(http.MethodGet, "/", protected.ThenFunc(app.jobs))
	router.Handler(http.MethodGet, "/vozaci", protected.ThenFunc(app.drivers))
	router.Handler(http.MethodGet, "/vozaci/info/:id", protected.ThenFunc(app.driver_info))
	router.Handler(http.MethodGet, "/vozaci/edit/:id", protected.ThenFunc(app.edit_driver))
	router.Handler(http.MethodPost, "/vozaci/edit/:id", protected.ThenFunc(app.edit_driver_post))
	router.Handler(http.MethodGet, "/vozaci/dodaj_vozaca", protected.ThenFunc(app.add_driver))
	router.Handler(http.MethodPost, "/vozaci/dodaj_vozaca", protected.ThenFunc(app.add_driver_post))
	router.Handler(http.MethodGet, "/kamioni", protected.ThenFunc(app.trucks))
	router.Handler(http.MethodGet, "/kamioni/info/:id", protected.ThenFunc(app.truck_info))
	router.Handler(http.MethodGet, "/kamioni/edit/:id", protected.ThenFunc(app.edit_truck))
	router.Handler(http.MethodPost, "/kamioni/edit/:id", protected.ThenFunc(app.edit_truck_post))
	router.Handler(http.MethodGet, "/kamioni/dodaj_kamion", protected.ThenFunc(app.add_truck))
	router.Handler(http.MethodPost, "/kamioni/dodaj_kamion", protected.ThenFunc(app.add_truck_post))
	router.Handler(http.MethodGet, "/poslovi", protected.ThenFunc(app.completed_jobs))
	router.Handler(http.MethodGet, "/poslovi/zavrsi/:id", protected.ThenFunc(app.finish_job))
	router.Handler(http.MethodPost, "/poslovi/zavrsi/:id", protected.ThenFunc(app.finish_job_post))
	router.Handler(http.MethodGet, "/poslovi/info/:id", protected.ThenFunc(app.job_info))
	router.Handler(http.MethodGet, "/poslovi/edit/:id", protected.ThenFunc(app.edit_job))
	router.Handler(http.MethodPost, "/poslovi/edit/:id", protected.ThenFunc(app.edit_job_post))
	router.Handler(http.MethodGet, "/poslovi/dodaj_posao", protected.ThenFunc(app.add_job))
	router.Handler(http.MethodPost, "/poslovi/dodaj_posao", protected.ThenFunc(app.add_job_post))
	router.Handler(http.MethodPost, "/korisnik/odjava", protected.ThenFunc(app.user_logout_post))

	standard := alice.New(app.recoverPanic, app.logRequest, secureHeaders)

	return standard.Then(router)
}
