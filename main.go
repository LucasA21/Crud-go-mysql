package main

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

func conexionBD() (conexion *sql.DB) {
	Driver := "mysql"
	Usuario := "usuario"
	Contrasenia := "contraseña"
	Nombre := "Sistema"

	conexion, err := sql.Open(Driver, Usuario+":"+Contrasenia+"@tcp(127.0.0.1)/"+Nombre)
	if err != nil {
		panic(err.Error())
	}
	return conexion
}

var templates = template.Must(template.ParseGlob("templates/*"))

/* Con la funcion template.ParseGlob guarda todos los archivos estaticos que
estan en la carpeta templates/ y los compila para poder trabajar con ellos de forma
dinamica guardandolos en la variable templates*/

func Index(w http.ResponseWriter, r *http.Request) {

	templates.ExecuteTemplate(w, "index", nil)
}

func Create(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "create", nil)
}

func Insert(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {

		name := r.FormValue("name")
		email := r.FormValue("email")

		conexionEstablecida := conexionBD()

		insertarRegistros, err := conexionEstablecida.Prepare("INSERT INTO empleados(name,email) VALUES(?,?)")

		if err != nil {
			panic(err.Error())
		}

		insertarRegistros.Exec(name, email)

		http.Redirect(w, r, "/", http.StatusMovedPermanently)

	}
}

/*Gracias a la funcion de arriba template.ParseGlob podemos hacer referencia
al template index y crear un handle para enrutar esa funcion con la ruta*/

func main() {

	http.HandleFunc("/", Index)
	http.HandleFunc("/create", Create)
	http.HandleFunc("/insert", Insert)

	log.Println("Server Runing...")

	http.ListenAndServe(":8080", nil)

}
