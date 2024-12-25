package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"hamster/internal/models"
	"hamster/internal/validator"

	"github.com/julienschmidt/httprouter"
)

type jobsCreateForm struct {
	Job         models.Job
	DriverIdMap map[int]string
	TruckIdMap  map[int]string
	validator.Validator
}

type userSignupForm struct {
	Name                string `form:"name"`
	Email               string `form:"email"`
	Password            string `form:"password"`
	validator.Validator `form:"-"`
}

type userLoginForm struct {
	Email               string `form:"email"`
	Password            string `form:"password"`
	validator.Validator `form:"-"`
}

func (app *application) jobs(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	var jobs []*models.Job
	var drivers []*models.Driver
	var trucks []*models.Truck

	// Dodavanj acive poslova TODO: 3 puta get-ujem ima sigurno bolje resenje
	// TODO : jobs i completed jobs su isti hamdleri samo su parametri drugaciji moze da se poboljsa
	params := map[string]any{
		"status": "active",
	}
	err := app.data.Get("jobs", &jobs, params)
	if err != nil {
		app.serverError(w, err)
		return
	}

	// Dodavanje poslova sa pending
	var jobs_pending []*models.Job
	params = map[string]any{
		"status": "pending",
	}
	err = app.data.Get("jobs", &jobs_pending, params)
	if err != nil {
		app.serverError(w, err)
		return
	}
	jobs = append(jobs, jobs_pending...)

	// Dodavanje poslova sa issue
	var jobs_issue []*models.Job
	params = map[string]any{
		"status": "issue",
	}
	err = app.data.Get("jobs", &jobs_issue, params)
	if err != nil {
		app.serverError(w, err)
		return
	}
	jobs = append(jobs, jobs_issue...)

	params = map[string]any{}
	err = app.data.Get("drivers", &drivers, params)
	if err != nil {
		app.serverError(w, err)
		return
	}
	err = app.data.Get("trucks", &trucks, params)
	if err != nil {
		app.serverError(w, err)
		return
	}

	driverMap := make(map[int]string)
	for _, driver := range drivers {
		driverMap[driver.ID] = driver.Name
	}

	truckMap := make(map[int]string)
	for _, truck := range trucks {
		truckMap[truck.ID] = truck.License_plate
	}

	var jobDisplays []*JobDisplay

	for _, job := range jobs {
		jobDisplays = append(jobDisplays, &JobDisplay{
			Job:               job,
			DriverName:        driverMap[job.Driver_id],
			TruckLicensePlate: truckMap[job.Truck_id],
			FormattedDate:     job.StartDate.Format("2006-01-02 15:04"),
			FormattedArrival:  job.EndDate.Format("2006-01-02 15:04"),
		})
	}

	data.JobDisplays = jobDisplays
	app.render(w, http.StatusOK, "jobs.html", data)
}

func (app *application) drivers(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)

	var drivers []*models.Driver

	params := map[string]any{}
	err := app.data.Get("drivers", &drivers, params)
	if err != nil {
		app.serverError(w, err)
		return
	}

	data.Drivers = drivers
	app.render(w, http.StatusOK, "drivers.html", data)

}

func (app *application) trucks(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	params := map[string]any{}
	var trucks []*models.Truck

	err := app.data.Get("trucks", &trucks, params)
	if err != nil {
		app.serverError(w, err)
		return
	}
	data.Trucks = trucks
	app.render(w, http.StatusOK, "trucks.html", data)
}

func (app *application) completed_jobs(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)

	var jobs []*models.Job
	var drivers []*models.Driver
	var trucks []*models.Truck

	params := map[string]any{
		"status": "completed",
	}

	err := app.data.Get("jobs", &jobs, params)
	if err != nil {
		app.serverError(w, err)
		return
	}

	params = map[string]any{}
	err = app.data.Get("drivers", &drivers, params)
	if err != nil {
		app.serverError(w, err)
		return
	}
	err = app.data.Get("trucks", &trucks, params)
	if err != nil {
		app.serverError(w, err)
		return
	}

	driverMap := make(map[int]string)
	for _, driver := range drivers {
		driverMap[driver.ID] = driver.Name
	}

	truckMap := make(map[int]string)
	for _, truck := range trucks {
		truckMap[truck.ID] = truck.License_plate
	}

	var jobDisplays []*JobDisplay

	for _, job := range jobs {
		jobDisplays = append(jobDisplays, &JobDisplay{
			Job:               job,
			DriverName:        driverMap[job.Driver_id],
			TruckLicensePlate: truckMap[job.Truck_id],
			FormattedDate:     job.StartDate.Format("2006-01-02 15:04"),
			FormattedArrival:  job.EndDate.Format("2006-01-02 15:04"),
		})
	}

	flash := app.sessionManager.PopString(r.Context(), "flash")

	data.JobDisplays = jobDisplays
	data.Flash = flash
	app.render(w, http.StatusOK, "jobs.html", data)
}

func (app *application) add_job(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)

	var drivers []*models.Driver
	var trucks []*models.Truck

	params := map[string]any{}

	err := app.data.Get("drivers", &drivers, params)
	if err != nil {
		app.serverError(w, err)
		return
	}
	err = app.data.Get("trucks", &trucks, params)
	if err != nil {
		app.serverError(w, err)
		return
	}

	driverMap := make(map[int]string)
	for _, driver := range drivers {
		driverMap[driver.ID] = driver.Name
	}

	truckMap := make(map[int]string)
	for _, truck := range trucks {
		truckMap[truck.ID] = truck.License_plate
	}
	data.Form = jobsCreateForm{
		Job: models.Job{
			Description: "Unesite opis posla",
		},
		DriverIdMap: driverMap,
		TruckIdMap:  truckMap,
	}
	app.render(w, http.StatusOK, "add_job.html", data)
}

func (app *application) add_job_post(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	description := r.PostForm.Get("description")
	driver_id_str := r.PostForm.Get("driver_id")
	truck_id_str := r.PostForm.Get("truck_id")
	scheduledDateStr := r.PostForm.Get("start_date")
	scheduledArrivalDateStr := r.PostForm.Get("end_date")
	truckStartKm_str := r.PostForm.Get("starting_km")

	form := jobsCreateForm{}

	form.CheckField(validator.NotBlank(description), "description", "Ne sme biti prazno")
	form.CheckField(validator.MaxChars(description, 200), "description", "Ne sme biti vise od 200 karaktera")
	form.CheckField(validator.NotBlank(driver_id_str), "driver_id", "Ne sme biti prazno")
	form.CheckField(validator.NotBlank(truck_id_str), "truck_id", "Ne sme biti prazno")
	form.CheckField(validator.NotBlank(scheduledDateStr), "start_date", "Ne sme biti prazno")
	form.CheckField(validator.NotBlank(scheduledArrivalDateStr), "end_date", "Ne sme biti prazno")
	form.CheckField(validator.NotBlank(truckStartKm_str), "starting_km", "Ne sme biti prazno")

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, "add_job.html", data)
	}

	driver_id, err := strconv.Atoi(driver_id_str)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	truck_id, err := strconv.Atoi(truck_id_str)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// Parse the date in the format expected by the 'datetime-local' HTML input type
	scheduled_date, err := time.Parse("2006-01-02T15:04", scheduledDateStr)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	scheduled_arrival_time, err := time.Parse("2006-01-02T15:04", scheduledArrivalDateStr)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	truckStartKm, err := strconv.ParseFloat(truckStartKm_str, 64)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	job := models.Job{
		Description:  description,
		Driver_id:    driver_id,
		Truck_id:     truck_id,
		StartDate:    scheduled_date,
		EndDate:      scheduled_arrival_time,
		TruckStartKm: truckStartKm,
	}

	id, err := app.data.Insert("jobs", job)
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.sessionManager.Put(r.Context(), "flash", "Posao uspesno dodat!")

	http.Redirect(w, r, fmt.Sprintf("/poslovi/info/%d", id), http.StatusSeeOther)
}

func (app *application) add_driver(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	app.render(w, http.StatusOK, "add_driver.html", data)
}

func (app *application) add_driver_post(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	name := r.PostForm.Get("name")
	licenseNumber := r.PostForm.Get("license_number")
	phoneNumber := r.PostForm.Get("phone_number")

	driver := models.Driver{
		Name:           name,
		License_number: licenseNumber,
		Phone_number:   phoneNumber,
		Status:         models.Available,
	}

	// Insert into the database
	id, err := app.data.Insert("drivers", driver)
	if err != nil {
		app.serverError(w, err)
		return
	}

	// Redirect to the driver info page
	http.Redirect(w, r, fmt.Sprintf("/vozaci/info/%d", id), http.StatusSeeOther)
}

func (app *application) add_truck(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	app.render(w, http.StatusOK, "add_truck.html", data)
}

func (app *application) add_truck_post(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	model := r.PostForm.Get("model")
	licensePlate := r.PostForm.Get("license_plate")

	truck := models.Truck{
		Model:         model,
		License_plate: licensePlate,
		Status:        models.Available,
	}

	id, err := app.data.Insert("trucks", truck)
	if err != nil {
		app.serverError(w, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/kamioni/info/%d", id), http.StatusSeeOther)
}

func (app *application) driver_info(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)

	url_params := httprouter.ParamsFromContext(r.Context())
	id, err := strconv.Atoi(url_params.ByName("id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}
	// url_params := httprouter.ParamsFromContext(r.Context())
	// name := url_params.ByName("name")

	params := map[string]any{
		"id": id,
	}
	var drivers []*models.Driver
	err = app.data.Get("drivers", &drivers, params)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}

	data.Drivers = drivers
	app.render(w, http.StatusOK, "view.html", data)
}

func (app *application) truck_info(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)

	url_params := httprouter.ParamsFromContext(r.Context())
	id, err := strconv.Atoi(url_params.ByName("id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	params := map[string]any{
		"id": id,
	}
	var trucks []*models.Truck
	err = app.data.Get("trucks", &trucks, params)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}

	data.Trucks = trucks
	app.render(w, http.StatusOK, "view.html", data)
}

// TODO : Mapa driver truck iD kako bi na info pageu pisalo
func (app *application) job_info(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)

	url_params := httprouter.ParamsFromContext(r.Context())
	id, err := strconv.Atoi(url_params.ByName("id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	params := map[string]any{
		"id": id,
	}
	var jobs []*models.Job
	err = app.data.Get("jobs", &jobs, params)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}

	data.Jobs = jobs
	app.render(w, http.StatusOK, "view.html", data)
}

func (app *application) edit_job(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)

	url_params := httprouter.ParamsFromContext(r.Context())
	id, err := strconv.Atoi(url_params.ByName("id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	params := map[string]any{
		"id": id,
	}
	var jobs []*models.Job
	var drivers []*models.Driver
	var trucks []*models.Truck

	err = app.data.Get("jobs", &jobs, params)
	if err != nil {
		app.serverError(w, err)
		return
	}
	job := jobs[0] // Get vraca niz ali selektujem po jedom ID-u

	params = map[string]any{}
	err = app.data.Get("drivers", &drivers, params)
	if err != nil {
		app.serverError(w, err)
		return
	}
	err = app.data.Get("trucks", &trucks, params)
	if err != nil {
		app.serverError(w, err)
		return
	}

	// Create maps for quick lookup
	driverMap := make(map[int]string)
	for _, driver := range drivers {
		driverMap[driver.ID] = driver.Name
	}

	truckMap := make(map[int]string)
	for _, truck := range trucks {
		truckMap[truck.ID] = truck.License_plate
	}

	var jobDisplays []*JobDisplay
	driverName := string(driverMap[job.Driver_id])
	truckLicensePlate := string(truckMap[job.Truck_id])

	jobDisplays = append(jobDisplays, &JobDisplay{
		Job:               job,
		DriverName:        driverName,
		TruckLicensePlate: truckLicensePlate,
		FormattedDate:     job.StartDate.Format("2006-01-02 15:04"),
		FormattedArrival:  job.EndDate.Format("2006-01-02 15:04"),
		IsEdit:            true,
		TruckIdMap:        truckMap,
		DriverIdMap:       driverMap,
	})

	data.JobDisplays = jobDisplays
	app.render(w, http.StatusOK, "edit_job.html", data)
}

func (app *application) edit_job_post(w http.ResponseWriter, r *http.Request) {
	url_params := httprouter.ParamsFromContext(r.Context())
	id, err := strconv.Atoi(url_params.ByName("id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	err = r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	description := r.PostForm.Get("description")
	driver_id_str := r.PostForm.Get("driver_id")
	truck_id_str := r.PostForm.Get("truck_id")
	scheduledDateStr := r.PostForm.Get("start_date")
	scheduledArrivalDateStr := r.PostForm.Get("end_date")
	status := r.PostForm.Get("status")
	truckStartKm_str := r.PostForm.Get("starting_km")

	form := jobsCreateForm{}

	form.CheckField(validator.NotBlank(description), "description", "Ne sme biti prazno")
	form.CheckField(validator.MaxChars(description, 200), "description", "Ne sme biti vise od 200 karaktera")
	form.CheckField(validator.NotBlank(driver_id_str), "driver_id", "Ne sme biti prazno")
	form.CheckField(validator.NotBlank(truck_id_str), "truck_id", "Ne sme biti prazno")
	form.CheckField(validator.NotBlank(scheduledDateStr), "start_date", "Ne sme biti prazno")
	form.CheckField(validator.NotBlank(scheduledArrivalDateStr), "end_date", "Ne sme biti prazno")
	form.CheckField(validator.NotBlank(status), "status", "Ne sme biti prazno")
	form.CheckField(validator.NotBlank(truckStartKm_str), "starting_km", "Ne sme biti prazno")

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, "edit_job.html", data)
	}

	driver_id, err := strconv.Atoi(driver_id_str)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	truck_id, err := strconv.Atoi(truck_id_str)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// Parse the date in the format expected by the 'datetime-local' HTML input type
	scheduled_date, err := time.Parse("2006-01-02T15:04", scheduledDateStr)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	scheduled_arrival_time, err := time.Parse("2006-01-02T15:04", scheduledArrivalDateStr)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	truckStartKm, err := strconv.ParseFloat(truckStartKm_str, 64)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	job := models.Job{
		Description:  description,
		Driver_id:    driver_id,
		Truck_id:     truck_id,
		StartDate:    scheduled_date,
		EndDate:      scheduled_arrival_time,
		Status:       models.JobStatus(status),
		TruckStartKm: truckStartKm,
	}

	err = app.data.Update("jobs", id, job)
	if err != nil {
		app.serverError(w, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/poslovi/info/%d", id), http.StatusSeeOther)
}

func (app *application) edit_truck(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)

	url_params := httprouter.ParamsFromContext(r.Context())
	id, err := strconv.Atoi(url_params.ByName("id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	params := map[string]any{
		"id": id,
	}
	var trucks []*models.Truck

	err = app.data.Get("trucks", &trucks, params)
	if err != nil {
		app.serverError(w, err)
		return
	}

	if len(trucks) == 0 {
		app.notFound(w)
		return
	}
	truck := trucks[0]

	data.Truck = truck
	app.render(w, http.StatusOK, "edit_truck.html", data)
}

func (app *application) edit_truck_post(w http.ResponseWriter, r *http.Request) {
	url_params := httprouter.ParamsFromContext(r.Context())
	id, err := strconv.Atoi(url_params.ByName("id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	err = r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	model := r.PostForm.Get("model")
	license_plate := r.PostForm.Get("license_plate")
	status := models.Status(r.PostForm.Get("status"))

	truck := models.Truck{
		Model:         model,
		License_plate: license_plate,
		Status:        status,
	}

	err = app.data.Update("trucks", id, truck)
	if err != nil {
		app.serverError(w, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/kamioni/info/%d", id), http.StatusSeeOther)
}

func (app *application) edit_driver(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)

	url_params := httprouter.ParamsFromContext(r.Context())
	id, err := strconv.Atoi(url_params.ByName("id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	params := map[string]any{
		"id": id,
	}
	var drivers []*models.Driver

	err = app.data.Get("drivers", &drivers, params)
	if err != nil {
		app.serverError(w, err)
		return
	}

	if len(drivers) == 0 {
		app.notFound(w)
		return
	}
	driver := drivers[0]

	data.Driver = driver
	app.render(w, http.StatusOK, "edit_driver.html", data)
}

func (app *application) edit_driver_post(w http.ResponseWriter, r *http.Request) {
	url_params := httprouter.ParamsFromContext(r.Context())
	id, err := strconv.Atoi(url_params.ByName("id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	err = r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	name := r.PostForm.Get("name")
	license_number := r.PostForm.Get("license_number")
	phone_number := r.PostForm.Get("phone_number")
	status := models.Status(r.PostForm.Get("status"))

	driver := models.Driver{
		Name:           name,
		License_number: license_number,
		Phone_number:   phone_number,
		Status:         status,
	}

	err = app.data.Update("drivers", id, driver)
	if err != nil {
		app.serverError(w, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/vozaci/info/%d", id), http.StatusSeeOther)
}

func (app *application) add_user(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = userSignupForm{}
	app.render(w, http.StatusOK, "signup.html", data)
}

func (app *application) add_user_post(w http.ResponseWriter, r *http.Request) {
	var form userSignupForm

	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form.CheckField(validator.NotBlank(form.Name), "name", "This field cannot be blank")
	form.CheckField(validator.NotBlank(form.Email), "email", "This field cannot be blank")
	form.CheckField(validator.Matches(form.Email, validator.EmailRX), "email", "This field must be a valid email address")
	form.CheckField(validator.NotBlank(form.Password), "password", "This field cannot be blank")
	form.CheckField(validator.MinChars(form.Password, 8), "password", "This field must be at least 8 characters long")

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, "signup.html", data)
		return
	}

	// Try to create a new user record in the database. If the email already
	// exists then add an error message to the form and re-display it.
	err = app.users.Insert(form.Name, form.Email, form.Password)
	if err != nil {
		if errors.Is(err, models.ErrDuplicateEmail) {
			form.AddFieldError("email", "Email address is already in use")
			data := app.newTemplateData(r)
			data.Form = form
			app.render(w, http.StatusUnprocessableEntity, "signup.html", data)
		} else {
			app.serverError(w, err)
		}
		return
	}
	// Otherwise add a confirmation flash message to the session confirming that
	// their signup worked.
	app.sessionManager.Put(r.Context(), "flash", "Your signup was successful. Please log in.")
	// And redirect the user to the login page.
	http.Redirect(w, r, "/korisnik/prijava", http.StatusSeeOther)
}

func (app *application) user_login(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = userLoginForm{}
	app.render(w, http.StatusOK, "login.html", data)
}
func (app *application) user_login_post(w http.ResponseWriter, r *http.Request) {
	// Decode the form data into the userLoginForm struct.
	var form userLoginForm
	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	// Do some validation checks on the form. We check that both email and
	// password are provided, and also check the format of the email address as
	// a UX-nicety (in case the user makes a typo).
	form.CheckField(validator.NotBlank(form.Email), "email", "This field cannot be blank")
	form.CheckField(validator.Matches(form.Email, validator.EmailRX), "email", "This field must be a valid email address")
	form.CheckField(validator.NotBlank(form.Password), "password", "This field cannot be blank")
	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, "login.html", data)
		return
	}
	// Check whether the credentials are valid. If they're not, add a generic
	// non-field error message and re-display the login page.
	id, err := app.users.Authenticate(form.Email, form.Password)
	if err != nil {
		if errors.Is(err, models.ErrInvalidCredentials) {
			form.AddNonFieldError("Email or password is incorrect")
			data := app.newTemplateData(r)
			data.Form = form
			app.render(w, http.StatusUnprocessableEntity, "login.html", data)
		} else {
			app.serverError(w, err)
		}
		return
	}
	// Use the RenewToken() method on the current session to change the session
	// ID. It's good practice to generate a new session ID when the
	// authentication state or privilege levels changes for the user (e.g. login
	// and logout operations).
	err = app.sessionManager.RenewToken(r.Context())
	if err != nil {
		app.serverError(w, err)
		return
	}
	// Add the ID of the current user to the session, so that they are now
	// 'logged in'.
	app.sessionManager.Put(r.Context(), "authenticatedUserID", id)
	// Redirect the user to the create snippet page.
	http.Redirect(w, r, "/poslovi", http.StatusSeeOther)
}
func (app *application) user_logout_post(w http.ResponseWriter, r *http.Request) {
	// Use the RenewToken() method on the current session to change the session
	// ID again.
	err := app.sessionManager.RenewToken(r.Context())
	if err != nil {
		app.serverError(w, err)
		return
	}
	// Remove the authenticatedUserID from the session data so that the user is
	// 'logged out'.
	app.sessionManager.Remove(r.Context(), "authenticatedUserID")
	// Add a flash message to the session to confirm to the user that they've been
	// logged out.
	app.sessionManager.Put(r.Context(), "flash", "You've been logged out successfully!")
	// Redirect the user to the application home page.
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *application) finish_job(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)

	url_params := httprouter.ParamsFromContext(r.Context())
	id, err := strconv.Atoi(url_params.ByName("id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	params := map[string]any{
		"id": id,
	}
	var jobs []*models.Job
	var drivers []*models.Driver
	var trucks []*models.Truck

	err = app.data.Get("jobs", &jobs, params)
	if err != nil {
		app.serverError(w, err)
		return
	}
	job := jobs[0] // Get vraca niz ali selektujem po jedom ID-u

	params = map[string]any{}
	err = app.data.Get("drivers", &drivers, params)
	if err != nil {
		app.serverError(w, err)
		return
	}
	err = app.data.Get("trucks", &trucks, params)
	if err != nil {
		app.serverError(w, err)
		return
	}

	// Create maps for quick lookup
	driverMap := make(map[int]string)
	for _, driver := range drivers {
		driverMap[driver.ID] = driver.Name
	}

	truckMap := make(map[int]string)
	for _, truck := range trucks {
		truckMap[truck.ID] = truck.License_plate
	}

	var jobDisplays []*JobDisplay
	driverName := string(driverMap[job.Driver_id])
	truckLicensePlate := string(truckMap[job.Truck_id])

	jobDisplays = append(jobDisplays, &JobDisplay{
		Job:               job,
		DriverName:        driverName,
		TruckLicensePlate: truckLicensePlate,
		FormattedDate:     job.StartDate.Format("2006-01-02 15:04"),
		FormattedArrival:  job.EndDate.Format("2006-01-02 15:04"),
		IsEdit:            true,
		TruckIdMap:        truckMap,
		DriverIdMap:       driverMap,
	})

	data.JobDisplays = jobDisplays
	app.render(w, http.StatusOK, "finish_job.html", data)
}
func (app *application) finish_job_post(w http.ResponseWriter, r *http.Request) {
	url_params := httprouter.ParamsFromContext(r.Context())
	id, err := strconv.Atoi(url_params.ByName("id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	err = r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	end_date_str := r.PostForm.Get("end_date")
	ending_km_str := r.PostForm.Get("ending_km")
	fuel_spent_str := r.PostForm.Get("fuel_spent")
	expenses_str := r.PostForm.Get("expenses")
	revenue_str := r.PostForm.Get("revenue")

	form := jobsCreateForm{}

	form.CheckField(validator.NotBlank(end_date_str), "end_date", "Ne sme biti prazno")
	form.CheckField(validator.NotBlank(ending_km_str), "ending_km", "Ne sme biti prazno")
	form.CheckField(validator.NotBlank(fuel_spent_str), "fuel_spent", "Ne sme biti prazno")
	form.CheckField(validator.NotBlank(expenses_str), "expenses", "Ne sme biti prazno")
	form.CheckField(validator.NotBlank(revenue_str), "revenue", "Ne sme biti prazno")

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, "finish_job.html", data)
	}

	// Parse the date in the format expected by the 'datetime-local' HTML input type
	end_date, err := time.Parse("2006-01-02T15:04", end_date_str)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	ending_km, err := strconv.ParseFloat(ending_km_str, 64)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	fuel_spent, err := strconv.ParseFloat(fuel_spent_str, 64)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	expenses, err := strconv.ParseFloat(expenses_str, 64)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	revenue, err := strconv.ParseFloat(revenue_str, 64)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	job := models.Job{
		EndDate:    end_date,
		Status:     models.JobStatus("completed"),
		TruckEndKm: ending_km,
		Fuel_spent: fuel_spent,
		Expenses:   expenses,
		Revenue:    revenue,
	}

	err = app.data.Update("jobs", id, job)
	if err != nil {
		app.serverError(w, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/poslovi/info/%d", id), http.StatusSeeOther)
}
