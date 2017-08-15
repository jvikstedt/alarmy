package api

import (
	"net/http"

	"github.com/jvikstedt/alarmy/internal/model"
)

type TriggerRequest struct {
	model.Trigger
	ProtectedID   interface{} `json:"id,omitempty"`
	OmitCreatedAt interface{} `json:"created_at,omitempty"`
	OmitUpdatedAt interface{} `json:"updated_at,omitempty"`
}

func (p *TriggerRequest) Bind(r *http.Request) error {
	return nil
}

// TriggerAll handler for getting all triggers
func (a *Api) TriggerAll(w http.ResponseWriter, r *http.Request) {
	//triggers, err := a.store.Trigger().All()
	//if stop := a.CheckErr(w, r, err, http.StatusInternalServerError); stop {
	//	return
	//}

	//render.JSON(w, r, triggers)
}

// TriggerCreate handler for creating a trigger
func (a *Api) TriggerCreate(w http.ResponseWriter, r *http.Request) {
	//data := &TriggerRequest{}
	//err := render.Bind(r, data)
	//if stop := a.CheckErr(w, r, err, http.StatusUnprocessableEntity); stop {
	//	return
	//}

	//// Validations
	//errors := data.Trigger.Errors()

	//_, err = a.store.Job().GetOne(data.Trigger.JobID)
	//if err != nil {
	//	errors["job_id"] = append(errors["job_id"], err.Error())
	//}

	//if len(errors) > 0 {
	//	a.Printf(r.Context(), "%v", errors)
	//	http.Error(w, http.StatusText(http.StatusUnprocessableEntity), http.StatusUnprocessableEntity)
	//	render.JSON(w, r, errors)
	//	return
	//}

	//trigger := data.Trigger
	//err = a.store.Trigger().Create(&trigger)
	//if stop := a.CheckErr(w, r, err, http.StatusInternalServerError); stop {
	//	return
	//}

	//render.Status(r, http.StatusCreated)
	//render.JSON(w, r, trigger)
}

// TriggerGetOne handler to get single trigger by id
func (a *Api) TriggerGetOne(w http.ResponseWriter, r *http.Request) {
	//triggerID, err := a.URLParamInt(r, "triggerID")
	//if stop := a.CheckErr(w, r, err, http.StatusInternalServerError); stop {
	//	return
	//}

	//trigger, err := a.store.Trigger().GetOne(triggerID)
	//if stop := a.CheckErr(w, r, err, http.StatusInternalServerError); stop {
	//	return
	//}

	//render.Status(r, http.StatusOK)
	//render.JSON(w, r, trigger)
}

// TriggerDestroy delete a single trigger by id
func (a *Api) TriggerDestroy(w http.ResponseWriter, r *http.Request) {
	//triggerID, err := a.URLParamInt(r, "triggerID")
	//if stop := a.CheckErr(w, r, err, http.StatusInternalServerError); stop {
	//	return
	//}

	//err = a.store.Trigger().Destroy(triggerID)
	//if stop := a.CheckErr(w, r, err, http.StatusInternalServerError); stop {
	//	return
	//}

	//render.Status(r, http.StatusOK)
}

// TriggerUpdate update a trigger by id
func (a *Api) TriggerUpdate(w http.ResponseWriter, r *http.Request) {
	//triggerID, err := a.URLParamInt(r, "triggerID")
	//if stop := a.CheckErr(w, r, err, http.StatusInternalServerError); stop {
	//	return
	//}

	//data := &TriggerRequest{}
	//err = render.Bind(r, data)
	//if stop := a.CheckErr(w, r, err, http.StatusInternalServerError); stop {
	//	return
	//}
	//data.Trigger.ID = triggerID

	//if errors := data.Trigger.Errors(); len(errors) > 0 {
	//	a.HandleError(w, r, errors, http.StatusUnprocessableEntity)
	//	render.JSON(w, r, errors)
	//	return
	//}

	//trigger := data.Trigger
	//err = a.store.Trigger().Update(&trigger)
	//if stop := a.CheckErr(w, r, err, http.StatusInternalServerError); stop {
	//	return
	//}

	//render.Status(r, http.StatusOK)
	//render.JSON(w, r, trigger)
}
