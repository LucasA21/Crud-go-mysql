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

type Empleado struct {
	Id    int
	Name  string
	Email string
}

func Index(w http.ResponseWriter, r *http.Request) {

	conexionEstablecida := conexionBD()

	registros, err := conexionEstablecida.Query("SELECT * FROM empleados")

	if err != nil {
		panic(err.Error())
	}

	empleado := Empleado{}
	arrayEmpleado := []Empleado{}

	for registros.Next() {
		var id int
		var name, email string
		err = registros.Scan(&id, &name, &email)

		if err != nil {
			panic(err.Error())
		}
		empleado.Id = id
		empleado.Name = name
		empleado.Email = email

		arrayEmpleado = append(arrayEmpleado, empleado)
	}

	templates.ExecuteTemplate(w, "index", arrayEmpleado)
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

func Update(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {

		id := r.FormValue("id")
		name := r.FormValue("name")
		email := r.FormValue("email")

		conexionEstablecida := conexionBD()

		actualizarRegistros, err := conexionEstablecida.Prepare("UPDATE empleados SET name=? , email=? WHERE id=?")

		if err != nil {
			panic(err.Error())
		}

		actualizarRegistros.Exec(name, email, id)

		http.Redirect(w, r, "/", http.StatusMovedPermanently)

	}
}

func Delete(w http.ResponseWriter, r *http.Request) {
	idEmpleado := r.URL.Query().Get("id")

	conexionEstablecida := conexionBD()

	borrarRegistros, err := conexionEstablecida.Prepare("DELETE FROM empleados WHERE id=?")

	if err != nil {
		panic(err.Error())
	}

	borrarRegistros.Exec(idEmpleado)

	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}

func Edit(w http.ResponseWriter, r *http.Request) {
	idEmpleado := r.URL.Query().Get("id")

	conexionEstablecida := conexionBD()

	registro, err := conexionEstablecida.Query("SELECT * FROM empleados WHERE id=?", idEmpleado)

	if err != nil {
		panic(err.Error())
	}

	empleado := Empleado{}

	for registro.Next() {
		var id int
		var name, email string
		err = registro.Scan(&id, &name, &email)

		if err != nil {
			panic(err.Error())
		}

		empleado.Id = id
		empleado.Name = name
		empleado.Email = email
	}

	templates.ExecuteTemplate(w, "edit", empleado)
}

func main() {

	http.HandleFunc("/", Index)
	http.HandleFunc("/create", Create)
	http.HandleFunc("/insert", Insert)
	http.HandleFunc("/delete", Delete)
	http.HandleFunc("/edit", Edit)
	http.HandleFunc("/update", Update)

	log.Println("Server Runing in http://localhost:8080")

	http.ListenAndServe(":8080", nil)

}
