package configuration

import (
	"encoding/json"
	"io"
	"log"
	"os"
)

// Config representa la configuración de la aplicación
type Config struct {
	//Configuración del servidor
	Port string `json:"port"`

	//Configuración de la cache
	NumCounters        int64 `json:"num_counters"`
	MaxCost            int64 `json:"max_cost"`
	BufferItems        int64 `json:"buffer_items"`
	ReadTimeout        int   `json:"read_timeout"`
	WriteTimeout       int   `json:"write_timeout"`
	MaxConnsPerIP      int   `json:"max_conns_per_ip"`
	MaxRequestsPerConn int   `json:"max_requests_per_conn"`

	//Configuración de los nodos y sincronización
	Peers                 []string `json:"peers"`
	RetriesToDisabledNode int      `json:"max_retries_to_disabled_node"`
	HeartBeatInterval     int      `json:"heart_beat_interval_in_seconds"`

	//Fichero de configuración de la whitelist de los nodos.
	WhiteListFilePath string `json:"white_list_file_path"`
}

// LoadConfig carga la configuración desde un archivo JSON
func LoadConfig(path string) Config {
	file, err := os.Open(path)
	if err != nil {
		log.Fatalf("❌ Error al abrir el archivo de configuración: %v", err)
	}
	defer file.Close()

	bytes, err := io.ReadAll(file)
	if err != nil {
		log.Fatalf("❌ Error al leer el archivo de configuración: %v", err)
	}

	var config Config
	if err := json.Unmarshal(bytes, &config); err != nil {
		log.Fatalf("❌ Error al parsear JSON de configuración: %v", err)
	}

	//Si no vienen seteadas o vienen a 0
	// le metemos datos por defecto...
	if config.HeartBeatInterval == 0 {
		config.HeartBeatInterval = 5
	}
	if config.RetriesToDisabledNode == 0 {
		config.RetriesToDisabledNode = 3
	}

	return config
}
