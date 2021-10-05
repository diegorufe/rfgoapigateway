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

// middleware función de punto intermedio para añadir función de seguridad entre medias
func middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO auth method
		// For test error
		//http.Error(w, "Unautoriced", 401)
		// return
		next.ServeHTTP(w, r)
	})
}

// LoadConfiguration : Método para cargar la configuración de las rutas a internceptar y host de destino a redirecionar
func LoadConfiguration(pathJsonFile string) error {
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
				http.Handle(tarjet.Route, middleware(httputil.NewSingleHostReverseProxy(destinationHost)))
				// http.Handle(tarjet.Route, httputil.NewSingleHostReverseProxy(destinationHost))
			} else {
				break
			}
		}
	}

	return err
}

// Método para arrancar el api gateway cargando la configuración y esuchando en el puerto indicado
func Start(pathJsonFile string, port int) {
	// Cargamos configuración
	LoadConfiguration(pathJsonFile)

	// Inicamos el servidor
	log.Fatal(http.ListenAndServe(":"+fmt.Sprint(port), nil))
}
