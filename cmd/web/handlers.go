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

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	var jobs []*models.Job
	var drivers []*models.Driver
	var trucks []*models.Truck

	params := map[string]any{
		"status": "active",
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
			FormattedDate:     job.Scheduled_date.Format("2006-01-02 15:04"),
			FormattedArrival:  job.Scheduled_arrival_time.Format("2006-01-02 15:04"),
		})
	}

	app.render(w, http.StatusOK, "home.html", &templateData{
		JobDisplays: jobDisplays,
	})
}

func (app *application) drivers(w http.ResponseWriter, r *http.Request) {

	var drivers []*models.Driver

	params := map[string]any{}
	err := app.data.Get("drivers", &drivers, params)
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.render(w, http.StatusOK, "drivers.html", &templateData{
		Drivers: drivers,
	})

}

func (app *application) trucks(w http.ResponseWriter, r *http.Request) {
	params := map[string]any{}
	var trucks []*models.Truck

	err := app.data.Get("trucks", &trucks, params)
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.render(w, http.StatusOK, "trucks.html", &templateData{
		Trucks: trucks,
	})

}

func (app *application) jobs(w http.ResponseWriter, r *http.Request) {
	var jobs []*models.Job
	var drivers []*models.Driver
	var trucks []*models.Truck

	params := map[string]any{}
	err := app.data.Get("jobs", &jobs, params)
	if err != nil {
		app.serverError(w, err)
		return
	}
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
			FormattedDate:     job.Scheduled_date.Format("2006-01-02 15:04"),
			FormattedArrival:  job.Scheduled_arrival_time.Format("2006-01-02 15:04"),
		})
	}

	flash := app.sessionManager.PopString(r.Context(), "flash")

	app.render(w, http.StatusOK, "home.html", &templateData{
		JobDisplays: jobDisplays,
		Flash:       flash,
	})
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
	scheduledDateStr := r.PostForm.Get("scheduled_date")
	scheduledArrivalDateStr := r.PostForm.Get("scheduled_arrival_time")
	distance_str := r.PostForm.Get("distance")
	package_size_str := r.PostForm.Get("package_size")
	package_weight_str := r.PostForm.Get("package_weight")
	client_name := r.PostForm.Get("client_name")
	start_location := r.PostForm.Get("start_location")
	destination_location := r.PostForm.Get("destination_location")

	form := jobsCreateForm{}

	form.CheckField(validator.NotBlank(description), "description", "Ne sme biti prazno")
	form.CheckField(validator.MaxChars(description, 200), "description", "Ne sme biti vise od 200 karaktera")
	form.CheckField(validator.NotBlank(driver_id_str), "driver_id", "Ne sme biti prazno")
	form.CheckField(validator.NotBlank(truck_id_str), "truck_id", "Ne sme biti prazno")
	form.CheckField(validator.NotBlank(scheduledDateStr), "scheduled_date", "Ne sme biti prazno")
	form.CheckField(validator.NotBlank(scheduledArrivalDateStr), "scheduled_arrival_time", "Ne sme biti prazno")
	form.CheckField(validator.NotBlank(distance_str), "distance", "Ne sme biti prazno")
	form.CheckField(validator.NotBlank(package_size_str), "package_size", "Ne sme biti prazno")
	form.CheckField(validator.NotBlank(client_name), "client_name", "Ne sme biti prazno")
	form.CheckField(validator.NotBlank(start_location), "start_location", "Ne sme biti prazno")
	form.CheckField(validator.NotBlank(destination_location), "destination_location", "Ne sme biti prazno")
	form.CheckField(validator.NotBlank(package_weight_str), "package_weight", "Ne sme biti prazno")

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

	distance, err := strconv.ParseFloat(distance_str, 64)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	package_size, err := strconv.ParseFloat(package_size_str, 64)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	package_weight, err := strconv.ParseFloat(package_weight_str, 64)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	job := models.Job{
		Description:            description,
		Driver_id:              driver_id,
		Truck_id:               truck_id,
		Scheduled_date:         scheduled_date,
		Scheduled_arrival_time: scheduled_arrival_time,
		Distance:               distance,
		Package_size:           package_size,
		Client_name:            client_name,
		Start_location:         start_location,
		Destination_location:   destination_location,
		Package_weight:         package_weight,
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

	averageConsumption, err := strconv.ParseFloat(r.PostForm.Get("average_consumption"), 64)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	kmTraveled, err := strconv.ParseFloat(r.PostForm.Get("km_traveled"), 64)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	driver := models.Driver{
		Name:                name,
		License_number:      licenseNumber,
		Phone_number:        phoneNumber,
		Status:              models.Available,
		Average_consumption: averageConsumption,
		Km_traveled:         kmTraveled,
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

	kmTraveled, err := strconv.ParseFloat(r.PostForm.Get("km_traveled"), 64)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	averageConsumption, err := strconv.ParseFloat(r.PostForm.Get("average_consumption"), 64)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	truck := models.Truck{
		Model:               model,
		License_plate:       licensePlate,
		Status:              models.Available,
		Km_traveled:         kmTraveled,
		Average_consumption: averageConsumption,
	}

	id, err := app.data.Insert("trucks", truck)
	if err != nil {
		app.serverError(w, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/kamioni/info/%d", id), http.StatusSeeOther)
}

func (app *application) driver_info(w http.ResponseWriter, r *http.Request) {
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

	app.render(w, http.StatusOK, "view.html", &templateData{
		Drivers: drivers,
	})
}

func (app *application) truck_info(w http.ResponseWriter, r *http.Request) {
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

	app.render(w, http.StatusOK, "view.html", &templateData{
		Trucks: trucks,
	})
}

func (app *application) job_info(w http.ResponseWriter, r *http.Request) {
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

	app.render(w, http.StatusOK, "view.html", &templateData{
		Jobs: jobs,
	})
}

func (app *application) edit_job(w http.ResponseWriter, r *http.Request) {
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

	for _, job := range jobs {
		jobDisplays = append(jobDisplays, &JobDisplay{
			Job:               job,
			DriverName:        driverMap[job.Driver_id],
			TruckLicensePlate: truckMap[job.Truck_id],
			FormattedDate:     job.Scheduled_date.Format("2006-01-02 15:04"),
			FormattedArrival:  job.Scheduled_arrival_time.Format("2006-01-02 15:04"),
			IsEdit:            true,
			TruckIdMap:        truckMap,
			DriverIdMap:       driverMap,
		})
	}

	app.render(w, http.StatusOK, "edit_job.html", &templateData{
		JobDisplays: jobDisplays,
	})
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
	driver_id, err := strconv.Atoi(r.PostForm.Get("driver_id"))
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	truck_id, err := strconv.Atoi(r.PostForm.Get("truck_id"))
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	scheduledDateStr := r.PostForm.Get("scheduled_date")
	scheduled_date, err := time.Parse("2006-01-02T15:04", scheduledDateStr)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	scheduledArrivalStr := r.PostForm.Get("scheduled_arrival_time")
	scheduled_arrival_time, err := time.Parse("2006-01-02T15:04", scheduledArrivalStr)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	status := models.JobStatus(r.PostForm.Get("status"))
	distance, err := strconv.ParseFloat(r.PostForm.Get("distance"), 64)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	package_size, err := strconv.ParseFloat(r.PostForm.Get("package_size"), 64)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	client_name := r.PostForm.Get("client_name")
	start_location := r.PostForm.Get("start_location")
	destination_location := r.PostForm.Get("destination_location")

	package_weight, err := strconv.ParseFloat(r.PostForm.Get("package_weight"), 64)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	job := models.Job{
		Description:            description,
		Driver_id:              driver_id,
		Truck_id:               truck_id,
		Scheduled_date:         scheduled_date,
		Scheduled_arrival_time: scheduled_arrival_time,
		Status:                 status,
		Distance:               distance,
		Package_size:           package_size,
		Client_name:            client_name,
		Start_location:         start_location,
		Destination_location:   destination_location,
		Package_weight:         package_weight,
	}

	err = app.data.Update("jobs", id, job)
	if err != nil {
		app.serverError(w, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/poslovi/info/%d", id), http.StatusSeeOther)
}

func (app *application) edit_truck(w http.ResponseWriter, r *http.Request) {
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

	app.render(w, http.StatusOK, "edit_truck.html", &templateData{
		Truck: truck,
	})
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

	km_traveled, err := strconv.ParseFloat(r.PostForm.Get("km_traveled"), 64)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	average_consumption, err := strconv.ParseFloat(r.PostForm.Get("average_consumption"), 64)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	truck := models.Truck{
		Model:               model,
		License_plate:       license_plate,
		Status:              status,
		Km_traveled:         km_traveled,
		Average_consumption: average_consumption,
	}

	err = app.data.Update("trucks", id, truck)
	if err != nil {
		app.serverError(w, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/kamioni/info/%d", id), http.StatusSeeOther)
}

func (app *application) edit_driver(w http.ResponseWriter, r *http.Request) {
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

	app.render(w, http.StatusOK, "edit_driver.html", &templateData{
		Driver: driver,
	})
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

	km_traveled, err := strconv.ParseFloat(r.PostForm.Get("km_traveled"), 64)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	average_consumption, err := strconv.ParseFloat(r.PostForm.Get("average_consumption"), 64)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	driver := models.Driver{
		Name:                name,
		License_number:      license_number,
		Phone_number:        phone_number,
		Status:              status,
		Km_traveled:         km_traveled,
		Average_consumption: average_consumption,
	}

	err = app.data.Update("drivers", id, driver)
	if err != nil {
		app.serverError(w, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/vozaci/info/%d", id), http.StatusSeeOther)
}
