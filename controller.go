package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type convocatoria struct {
	Idconvocatoria   int    `json:"idconvocatoria"`
	Nomconvocatoria  string `json:"nom_convocatoria"`
	Infoconvocatoria string `json:"info_convocatoria"`
	Nvacantes        int    `json:"n_vacantes"`
	Fechainicio      string `json:"fecha_inicio"`
	Fechafin         string `json:"fecha_fin"`
	Estado           string `json:"estado"`
}

var convocatoriaDB = db()

func getAllConvocatorias(w http.ResponseWriter, r *http.Request) {
	err := verifyToken(w, r)
	if err != nil {
		log.Println(err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	convocatorias := []convocatoria{}
	response, err := convocatoriaDB.Query("SELECT idconvocatoria, nom_convocatoria, info_convocatoria, n_vacantes, fecha_inicio, fecha_fin, estado FROM convocatorias")
	if err != nil {
		log.Fatal(err)
	}
	defer response.Close()

	for response.Next() {
		var c convocatoria
		err = response.Scan(
			&c.Idconvocatoria,
			&c.Nomconvocatoria,
			&c.Infoconvocatoria,
			&c.Nvacantes,
			&c.Fechainicio,
			&c.Fechafin,
			&c.Estado,
		)
		if err != nil {
			log.Fatal(err)
		}
		convocatorias = append(convocatorias, c)
	}
	json.NewEncoder(w).Encode(convocatorias)
}
