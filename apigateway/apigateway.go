package apigateway

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
)

// Clase para alamcenar los datos de la ruta a interceptar y el destino
type Target struct {
	Route           string // ruta a interceptar
	DestinationHost string // host de destino a redirecionar
}

// Configuración para cargar los datos de los tarjet a inteceptar
type Configuration struct {
	Targets []Target // array de tarjets a interceptar y a redireccionar
}

// Respuesta del middleware
type ResponseMiddleware struct {
	Err    error
	Status int
}

// middleware función de punto intermedio para añadir función de seguridad entre medias
func middleware(path string, next http.Handler, fnMiddleware func(route string, w http.ResponseWriter, r *http.Request) ResponseMiddleware) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO auth method
		// For test error
		//http.Error(w, "Unautoriced", 401)
		// return
		// Test for add header
		// r.Header.Add("TEST", "This is a test header")
		// TODO add user authentication if needed in this part or in auth method

		if fnMiddleware != nil {
			responseMiddleware := fnMiddleware(path, w, r)

			if responseMiddleware.Err != nil {
				http.Error(w, responseMiddleware.Err.Error(), responseMiddleware.Status)
				return
			}

		}

		next.ServeHTTP(w, r)
	})
}

// LoadConfiguration : Método para cargar la configuración de las rutas a internceptar y host de destino a redirecionar
func LoadConfiguration(pathJsonFile string, fnMiddleware func(route string, w http.ResponseWriter, r *http.Request) ResponseMiddleware) error {
	// Cargamos la configuración de las url a interceptar y el host de destino
	file, _ := os.Open(pathJsonFile)
	defer file.Close()
	decoder := json.NewDecoder(file)
	configuration := Configuration{}

	err := decoder.Decode(&configuration)

	if err == nil {

		// Recorremos las urls a interceptar y a enviar al host de destino
		for _, tarjet := range configuration.Targets {

			log.Println("Cargando target. Ruta: " + tarjet.Route + ", host de destino: " + tarjet.DestinationHost)

			destinationHost, err := url.Parse(tarjet.DestinationHost)

			if err == nil {
				http.Handle(tarjet.Route, middleware(tarjet.Route, httputil.NewSingleHostReverseProxy(destinationHost), fnMiddleware))
			} else {
				log.Fatal(err.Error())
			}
		}
	}

	return err
}

// Método para arrancar el api gateway cargando la configuración y esuchando en el puerto indicado
func Start(pathJsonFile string, port int, fnMiddleware func(route string, w http.ResponseWriter, r *http.Request) ResponseMiddleware) {
	// Cargamos configuración
	LoadConfiguration(pathJsonFile, fnMiddleware)

	// Inicamos el servidor
	log.Fatal(http.ListenAndServe(":"+fmt.Sprint(port), nil))
}
