package models

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
)

// Creamos la interfaz de lo que esperamos recibir en Manage Errors
type InputManageErrors struct {
	DeployId int             `json:"deploy_id"`
	AppName  string          `json:"app_name"`
	Type     string          `json:"type"`
	Version  string          `json:"version"`
	Values   json.RawMessage `json:"values"`
}

func CheckManageErrorsParams(input InputManageErrors) error {
	// Validar deployId (number)
	if input.DeployId <= 0 {
		return fmt.Errorf("deployId debe ser un número positivo")
	}

	// Validar appName (string - solo letras)
	match, _ := regexp.MatchString("^[a-zA-Z]+$", input.AppName)
	if !match {
		return fmt.Errorf("appName solo puede contener letras")
	}

	// Validar type (solo "ci/release")
	if input.Type != "ci/release" {
		return fmt.Errorf("type solo puede ser \"ci/release\"")
	}

	// Validar versión (M.m.f)
	versionPattern := regexp.MustCompile(`^\d+\.\d+\.\d+$`)
	if !versionPattern.MatchString(input.Version) {
		return fmt.Errorf("versión debe tener el formato M.m.f (ejemplo: 1.2.3)")
	}

	// Validar values (formato JSON válido)
	var values interface{}
	if err := json.Unmarshal(input.Values, &values); err != nil {
		return fmt.Errorf("values no tiene un formato JSON válido")
	}

	return nil
}

func CallServiceManageErrors(input InputManageErrors, url string) error {

	// Convertir la estructura a JSON
	jsonData, err := json.Marshal(input)
	if err != nil {
		return fmt.Errorf("json error")
	}

	// Realizar la solicitud POST al servicio REST
	response, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("error doing rest call")
	}
	defer response.Body.Close()

	// Verificar el código de respuesta
	if response.StatusCode == http.StatusOK {
		return nil // Devuelve nil si el código de respuesta es 200
	} else {
		return fmt.Errorf("service response is: %d", response.StatusCode)
	}
}

func SaveManageErrorsInDb(input InputManageErrors) error {
	stmt, err := db.Prepare("INSERT INTO manage_errors_endpoint(deploy_id, app_name, type, version, values) VALUES(?, ?, ?, ?, ?)")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(&input.DeployId, &input.AppName, &input.Type, &input.Version, &input.Values)
	if err != nil {
		return err
	}
	return nil
}
