package security

import (
	"encoding/json"
	"log"
	"os"
	"phoenixcache/configuration"
	"sync"

	"github.com/fsnotify/fsnotify"
)

var (
	allowedNodes = make(map[string]bool)
	mu           sync.RWMutex
)

// LoadWhitelist lee el archivo y actualiza la lista de nodos permitidos
func LoadWhitelist(config *configuration.Config) {
	if config.WhiteListFilePath == "" {
		log.Printf("nein")
		return
	}

	mu.Lock()
	defer mu.Unlock()

	file, err := os.ReadFile(config.WhiteListFilePath)
	if err != nil {
		log.Printf("⚠️ Error leyendo la lista blanca: %v", err)
		return
	}

	var data struct {
		AllowedNodes []string `json:"allowed_nodes"`
	}

	if err := json.Unmarshal(file, &data); err != nil {
		log.Printf("⚠️ Error parseando JSON: %v", err)
		return
	}

	// Limpiar y actualizar la lista
	allowedNodes = make(map[string]bool)
	for _, peers := range data.AllowedNodes {
		allowedNodes[peers] = true
	}

	log.Println("✅ Lista blanca actualizada:", allowedNodes)
}

// IsAllowedNode verifica si un nodo está permitido
func IsAllowedNode(ip string) bool {
	mu.RLock()
	defer mu.RUnlock()
	return allowedNodes[ip]
}

// WatchWhitelistFile observa cambios en el archivo y lo recarga automáticamente
func WatchWhitelistFile(config *configuration.Config) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatalf("❌ Error creando watcher: %v", err)
	}
	defer watcher.Close()

	err = watcher.Add(config.WhiteListFilePath)
	if err != nil {
		log.Fatalf("❌ Error viendo cambios en %s: %v", config.WhiteListFilePath, err)
	}

	log.Printf("👀 Observando cambios en: %s", config.WhiteListFilePath)

	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}
			// Si el archivo se modificó o recreó, recargar la lista
			if event.Op&(fsnotify.Write|fsnotify.Create) != 0 {
				log.Println("🔄 Archivo de lista blanca modificado, recargando...")
				LoadWhitelist(config)
			}

		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			log.Printf("❌ Error en watcher: %v", err)
		}
	}
}
