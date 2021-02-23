package main

import (
	"fmt"
	"log"
	"strconv"
	"github.com/gorilla/mux"
	"github.com/gorilla/handlers"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"os"
)

type Proceso struct {
	Pid int `json:"pid"`
}

func failOnError(err error, msg string){
	if err != nil {
		log.Fatalf(msg, err)
	}
}

func index(w http.ResponseWriter, r *http.Request){
	fmt.Fprintf(w, "Llego al servidor")
}

func getRam(w http.ResponseWriter, r *http.Request){
	fmt.Println("getRam")

	file_data, err := os.Open("/proc/serverRam")
	if err != nil {
		fmt.Println("Error al leer modulo ram")
	}
	defer file_data.Close()

	byteValue, _ := ioutil.ReadAll(file_data)
	
	var result map[string]interface{}
	json.Unmarshal([]byte(byteValue), &result)
	fmt.Println(result["ram"])
	json.NewEncoder(w).Encode(result)

	//fmt.Fprintf(w, string(file_data))
}

func getProcs(w http.ResponseWriter, r *http.Request){
	fmt.Println("getProcs")
	file_data, err := ioutil.ReadFile("/proc/serverProcs2")
	if err != nil {
		fmt.Println("Error al leer modulo procs")
	}
	fmt.Fprintf(w, string(file_data))
}

func killProc(w http.ResponseWriter, r *http.Request){

	fmt.Println("killProc")

	//recuperar proceso
	var p Proceso

	reqbody, errpid := ioutil.ReadAll(r.Body)
	if errpid != nil {
		fmt.Fprintf(w, "Error en leer proceso")
	}
	json.Unmarshal(reqbody, &p)
	fmt.Println("Eliminando proceso", p.Pid)
	proc, err := os.FindProcess(p.Pid)
	if err != nil {
		fmt.Println("Error al eliminar proceso", p.Pid)
	}
	proc.Kill()
	fmt.Fprintf(w, "Proceso " + strconv.Itoa(p.Pid) + " eliminado.")
}

func main(){
	router := mux.NewRouter().StrictSlash(true)
	headers := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization", "application/json"})
	methods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE"})
	origins := handlers.AllowedOrigins([]string{"*"})
	router.HandleFunc("/", index).Methods("GET")
	router.HandleFunc("/getRam", getRam).Methods("GET")
	router.HandleFunc("/getProcs", getProcs).Methods("GET")
	router.HandleFunc("/killProc", killProc).Methods("DELETE")
	fmt.Println("El servidor go a la escucha en puerto 5000")
	http.ListenAndServe(":5000",handlers.CORS(headers, methods, origins)(router))
}

