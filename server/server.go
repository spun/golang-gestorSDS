package server

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/bertus193/gestorSDS/config"
	"github.com/bertus193/gestorSDS/server/database"
	"github.com/bertus193/gestorSDS/utils"
)

// Launch lanza el servidor
func Launch() {

	// suscripción SIGINT
	stopChan := make(chan os.Signal)
	signal.Notify(stopChan, os.Interrupt)

	// Rutas disponibles
	mux := http.NewServeMux()
	mux.Handle("/usuario/login", http.HandlerFunc(loginUsuario))
	mux.Handle("/usuario/registro", http.HandlerFunc(registroUsuario))
	mux.Handle("/usuario/eliminar", http.HandlerFunc(eliminarUsuario))
	mux.Handle("/usuario/detalles", http.HandlerFunc(detallesUsuario))
	mux.Handle("/a2f/activar", http.HandlerFunc(activarA2F))
	mux.Handle("/a2f/desactivar", http.HandlerFunc(desactivarA2F))
	mux.Handle("/a2f/desbloquear", http.HandlerFunc(desbloquearA2F))
	mux.Handle("/vault", http.HandlerFunc(listarEntradas))
	mux.Handle("/vault/nueva", http.HandlerFunc(crearEntrada))
	mux.Handle("/vault/detalles", http.HandlerFunc(detallesEntrada))
	mux.Handle("/vault/eliminar", http.HandlerFunc(eliminarEntrada))

	srv := &http.Server{Addr: config.SecureServerPort, Handler: mux}

	go func() {
		if err := srv.ListenAndServeTLS("cert.pem", "key.pem"); err != nil {
			log.Printf("listen: %s\n", err)
		}
	}()

	<-stopChan // espera señal SIGINT
	log.Println("Apagando servidor ...")

	// Guarda la información de la BD en un fichero
	database.After()

	//Guarda logs en fichero
	utils.AfterLogs()

	// Apaga el servidor de forma segura
	ctx, fnc := context.WithTimeout(context.Background(), 5*time.Second)
	fnc()
	srv.Shutdown(ctx)

	log.Println("Servidor detenido correctamente")
}
