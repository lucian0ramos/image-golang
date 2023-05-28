package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/lucian0ramos/image-golang/src/models"

	"github.com/rs/zerolog/hlog"
)

func ManageErrors(w http.ResponseWriter, r *http.Request) {
	mR := models.MyResponse{Code: 1}

	i := models.InputManageErrors{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&i)
	defer r.Body.Close()

	if err != nil {
		mR.Message = "Unexpected error, probably the data type you have added in json data is wrong. Check documentation and try again."
		hlog.FromRequest(r).Error().Caller().Err(err).Msg(mR.Message)
		models.GenerateResponse(w, mR, http.StatusInternalServerError)
		return
	}

	//Check de parámetros
	err = models.CheckManageErrorsParams(i)
	if err != nil {
		mR.Message = "Params error in JSON Data. Check documentation and try again."
		hlog.FromRequest(r).Error().Caller().Err(err).Msg(mR.Message)
		models.GenerateResponse(w, mR, http.StatusInternalServerError)
		return
	}

	//Call another service (esto dará error, si se quiere comprobar si ginputrda bien en bd comentar esta función)
	err = models.CallServiceManageErrors(i, "kubernetes.localhost.com")
	if err != nil {
		mR.Message = "Error calling endpoint"
		hlog.FromRequest(r).Error().Caller().Err(err).Msg(mR.Message)
		models.GenerateResponse(w, mR, http.StatusInternalServerError)
		return
	}

	//Save data in bd
	err = models.SaveManageErrorsInDb(i)
	if err != nil {
		mR.Message = "Error saving data in db "
		hlog.FromRequest(r).Error().Caller().Err(err).Msg(mR.Message)
		models.GenerateResponse(w, mR, http.StatusInternalServerError)
		return
	}

	mR.Code = 0
	mR.Message = "Ok"
	models.GenerateResponse(w, mR, http.StatusOK)
}
