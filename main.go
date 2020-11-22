package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"

	_ "github.com/go-sql-driver/mysql"
)

type Convocatoria struct {
	Idconvocatoria   int    `json:"idconvocatoria"`
	Nomconvocatoria  string `json:"nom_convocatoria"`
	Infoconvocatoria string `json:"info_convocatoria"`
	Nvacantes        int    `json:"n_vacantes"`
	Fechainicio      string `json:"fecha_inicio"`
	Fechafin         string `json:"fecha_fin"`
	Estado           string `json:"estado"`
}

// type allConvos []convo

// var data = allConvos

func obtenerBaseDeDatos() (db *sql.DB, e error) {
	usuario := "uysv9t18xsovpgmn"
	pass := "lqXoWZythDxqtkIJbcOp"
	host := "tcp(birvxmqrxnbzpdvm0ogz-mysql.services.clever-cloud.com)"
	nombreBaseDeDatos := "birvxmqrxnbzpdvm0ogz"
	// Debe tener la forma usuario:contraseña@host/nombreBaseDeDatos
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@%s/%s", usuario, pass, host, nombreBaseDeDatos))
	if err != nil {
		return nil, err
	}
	return db, nil
}

func conexionSuccess() {
	db, err := obtenerBaseDeDatos()
	if err != nil {
		fmt.Printf("Error obteniendo base de datos: %v", err)
		return
	}
	// Terminar conexión al terminar función
	defer db.Close()

	// Ahora vemos si tenemos conexión
	err = db.Ping()
	if err != nil {
		fmt.Printf("Error conectando: %v", err)
		return
	}
	// Listo, aquí ya podemos usar a db!
	fmt.Printf("Conectado correctamente")
}

func main() {
	conexionSuccess()
	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/api", indexRoute)
	router.HandleFunc("/api/convocatoria", getConvocatoria)
	log.Fatal(http.ListenAndServe(":3000", router))
}

func getConvocatoria(w http.ResponseWriter, r *http.Request) {
	convocatorias, err := obtenerContactos()
	if err != nil {
		fmt.Printf("Error obteniendo contactos: %v", err)
		return
	}
	for _, conv := range convocatorias {
		fmt.Printf("%v\n", conv)
	}
	json.NewEncoder(w).Encode(convocatorias)
}

//ruta principal
func indexRoute(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Bienvenido a Convocatoria2")
}

func obtenerContactos() ([]Convocatoria, error) {
	convocatorias := []Convocatoria{}
	db, err := obtenerBaseDeDatos()
	if err != nil {
		return nil, err
	}
	defer db.Close()
	filas, err := db.Query("SELECT idconvocatoria, nom_convocatoria, info_convocatoria, n_vacantes, fecha_inicio, fecha_fin, estado FROM convocatorias")
	if err != nil {
		return nil, err
	}
	// Si llegamos aquí, significa que no ocurrió ningún error
	defer filas.Close()

	// Aquí vamos a "mapear" lo que traiga la consulta en el while de más abajo
	var c Convocatoria
	// Recorrer todas las filas, en un "while"
	for filas.Next() {
		err = filas.Scan(&c.Idconvocatoria, &c.Nomconvocatoria, &c.Infoconvocatoria, &c.Nvacantes, &c.Fechainicio, &c.Fechafin, &c.Estado,)
		// Al escanear puede haber un error
		if err != nil {
			return nil, err
		}
		// Y si no, entonces agregamos lo leído al arreglo
		convocatorias = append(convocatorias, c)
	}
	// Vacío o no, regresamos el arreglo de contactos
	return convocatorias, nil
}

//Intentos
// func main() {

// fmt.Println("Connecting to database...")
// connString := "uysv9t18xsovpgmn:lqXoWZythDxqtkIJbcOp@tcp(sql.birvxmqrxnbzpdvm0ogz-mysql.services.clever-cloud.com)/birvxmqrxnbzpdvm0ogz"
// db, err := sql.Open("mysql", connString)
// if err != nil {
// 	panic(err.Error())
// }
// fmt.Println("Go mysql tutorial")
// db, err := sql.Open("mysql", "uysv9t18xsovpgmn:lqXoWZythDxqtkIJbcOp@/birvxmqrxnbzpdvm0ogz")

// if err != nil {
// 	panic(err.Error())
// }
// defer db.Close()

// fmt.Println("Success Connected to Mysql Database")

// results, err := db.Query("SELECT idconvocatoria, nom_convocatoria FROM convocatorias")
// if err != nil {
// 	panic(err.Error()) // proper error handling instead of panic in your app
// }

// for results.Next() {
// 	var tag Tag
// 	// for each row, scan the result into our tag composite object
// 	err = results.Scan(&tag.Idconvocatoria, &tag.Nomconvocatoria)
// 	if err != nil {
// 		panic(err.Error()) // proper error handling instead of panic in your app
// 	}
// 	// and then print out the tag's Name attribute
// 	log.Printf(tag.Nomconvocatoria)
// }
// }
