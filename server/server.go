package server

import (
	"log"
	"time"

	"phoenixcache/configuration"
	"phoenixcache/distributed"
	"phoenixcache/internal"

	"github.com/valyala/fasthttp"
)

// StartServer inicia el servidor
func StartServer(config *configuration.Config, peerManager *distributed.PeerManager, cache *internal.Cache) {

	server := &fasthttp.Server{
		Handler:            SetupRouter(config, peerManager, cache),
		Name:               "UltraFastServer",
		ReadTimeout:        time.Duration(config.ReadTimeout) * time.Second,
		WriteTimeout:       time.Duration(config.WriteTimeout) * time.Second,
		MaxConnsPerIP:      config.MaxConnsPerIP,
		MaxRequestsPerConn: config.MaxRequestsPerConn,
	}

	log.Printf("🐦‍🔥 Servidor corriendo en %s", config.Port)
	log.Fatal(server.ListenAndServe(config.Port))
}
