package usecases

import (
	"Beside-Mom-BE/modules/entities"
	"Beside-Mom-BE/modules/repositories"
)

type AppUseCase interface {
	CreateAppointment(app *entities.Appointment) (*entities.Appointment, error)
	GetAppByID(id string) (map[string]interface{}, error)
	GetAppByUserID(userID string) ([]map[string]interface{}, error)
	GetAppInProgressByUserID(userID string) ([]map[string]interface{}, error)
	GetAllApp() ([]map[string]interface{}, error)
	UpdateAppByID(id string, app *entities.Appointment) (*entities.Appointment, error)
	DeleteAppByID(id string) error
}

type AppUseCaseImpl struct {
	repo repositories.AppRepository
}

func NewAppUseCase(repo repositories.AppRepository) *AppUseCaseImpl {
	return &AppUseCaseImpl{repo: repo}
}

func (u *AppUseCaseImpl) CreateAppointment(app *entities.Appointment) (*entities.Appointment, error) {
	createdApp, err := u.repo.CreateAppointment(app)
	if err != nil {
		return nil, err
	}

	return createdApp, err
}

func (u *AppUseCaseImpl) GetAppByID(id string) (map[string]interface{}, error) {
	app, err := u.repo.GetAppByID(id)
	if err != nil {
		return nil, err
	}

	appData := map[string]interface{}{
		"id":          app.ID,
		"title":       app.Title,
		"date":        app.Date,
		"start_time":  app.StartTime,
		"building":    app.Building,
		"requirement": app.Requirement,
		"status":      app.Status,
		"user_id":     app.User.ID,
		"name":        app.User.Firstname + " " + app.User.Lastname,
	}

	return appData, nil
}
func (u *AppUseCaseImpl) GetAppInProgressByUserID(userID string) ([]map[string]interface{}, error) {
	apps, err := u.repo.GetAppInProgressByUserID(userID)
	if err != nil {
		return nil, err
	}

	var appsList []map[string]interface{}
	for _, app := range apps {
		appData := map[string]interface{}{
			"id":          app.ID,
			"title":       app.Title,
			"date":        app.Date,
			"start_time":  app.StartTime,
			"building":    app.Building,
			"requirement": app.Requirement,
			"doctor":      app.Doctor,
			"status":      app.Status,
			"user_id":     app.User.ID,
			"name":        app.User.Firstname + " " + app.User.Lastname,
		}

		appsList = append(appsList, appData)
	}

	return appsList, nil
}

func (u *AppUseCaseImpl) GetAppByUserID(userID string) ([]map[string]interface{}, error) {
	apps, err := u.repo.GetAppByUserID(userID)
	if err != nil {
		return nil, err
	}

	var appsList []map[string]interface{}
	for _, app := range apps {
		appData := map[string]interface{}{
			"id":          app.ID,
			"title":       app.Title,
			"date":        app.Date,
			"start_time":  app.StartTime,
			"building":    app.Building,
			"requirement": app.Requirement,
			"doctor":      app.Doctor,
			"status":      app.Status,
			"user_id":     app.User.ID,
			"name":        app.User.Firstname + " " + app.User.Lastname,
		}

		appsList = append(appsList, appData)
	}

	return appsList, nil
}

func (u *AppUseCaseImpl) GetAllApp() ([]map[string]interface{}, error) {
	apps, err := u.repo.GetAllApp()
	if err != nil {
		return nil, err
	}

	var appsList []map[string]interface{}
	for _, app := range apps {
		appData := map[string]interface{}{
			"id":          app.ID,
			"title":       app.Title,
			"date":        app.Date,
			"start_time":  app.StartTime,
			"building":    app.Building,
			"requirement": app.Requirement,
			"doctor":      app.Doctor,
			"status":      app.Status,
			"user_id":     app.User.ID,
			"name":        app.User.Firstname + " " + app.User.Lastname,
		}

		appsList = append(appsList, appData)
	}

	return appsList, nil
}

func (u *AppUseCaseImpl) UpdateAppByID(id string, app *entities.Appointment) (*entities.Appointment, error) {
	existingApp, err := u.repo.GetAppByID(id)
	if err != nil {
		return nil, err
	}

	existingApp.Title = app.Title
	existingApp.Date = app.Date
	existingApp.StartTime = app.StartTime
	existingApp.Building = app.Building
	existingApp.Requirement = app.Requirement
	existingApp.Doctor = app.Doctor
	existingApp.Status = app.Status
	return u.repo.UpdateAppByID(existingApp)
}

func (u *AppUseCaseImpl) DeleteAppByID(id string) error {
	return u.repo.DeleteAppByID(id)
}
