package main

import (
	"phoenixcache/configuration"
	"phoenixcache/distributed"
	"phoenixcache/internal"
	"phoenixcache/security"
	"phoenixcache/server"

	"time"
)

func main() {
	// Cargar configuración
	config := configuration.LoadConfig("config.json")

	// Inicializar caché
	cache := internal.NewCache(config.NumCounters, config.MaxCost, config.BufferItems)

	//Iniciamos el modulo de seguridad
	security.InitModule(&config)

	var peerManager *distributed.PeerManager

	if (config.Peers != nil) && (len(config.Peers) > 0) {
		peerManager = distributed.NewPeerManager(config.Peers, time.Duration(config.HeartBeatInterval)*time.Second, config.RetriesToDisabledNode)
		distributed.RecoverCacheFromPeer(peerManager, cache)
	}

	// Iniciar servidor
	server.StartServer(&config, peerManager, cache)
}
