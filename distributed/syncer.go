package distributed

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"phoenixcache/internal"
	"phoenixcache/utils"

	"github.com/valyala/fasthttp"
)

// Estructura de sincronización
type SyncMessage struct {
	Action string        `json:"action"`
	Key    string        `json:"key"`
	Value  interface{}   `json:"value,omitempty"`
	TTL    time.Duration `json:"ttl,omitempty"`
}

// Propaga los cambios a los diferentes servidores asignados
func PropagateChange(msg SyncMessage, peerManager *PeerManager) {
	if peerManager == nil {
		return
	}

	data, _ := json.Marshal(msg)
	for _, peer := range peerManager.GetActivePeers() {
		go func(peer string) {
			resp := fasthttp.AcquireResponse()
			req := fasthttp.AcquireRequest()
			defer fasthttp.ReleaseResponse(resp)
			defer fasthttp.ReleaseRequest(req)

			req.SetRequestURI(peer + "/sync")
			req.Header.SetMethod(fasthttp.MethodPost)
			req.Header.SetContentType("application/json")
			req.Header.Set("Connection", "close")

			req.SetBody(data) // <-- Usa data directamente en SetBody

			if err := fasthttp.Do(req, resp); err != nil {
				log.Printf("⚠️ Error sincronizando con %s: %v", peer, err)
			}

		}(peer)
	}
}

// Recuperamos la cache del primer servidor activo de caches
func RecoverCacheFromPeer(peerManager *PeerManager, cache *internal.Cache) {

	peer := peerManager.GetActivePeers()[0]
	if peer == "" {
		return
	}

	url := fmt.Sprintf("%s/export", peer)

	statusCode, compressedData, err := fasthttp.Get(nil, url)

	if err != nil {
		log.Printf("⚠️ Error al obtener la caché de %s: %s", url, err.Error())
		return
	}

	if statusCode == 0 {
		log.Printf("⚠️ Respuesta vacía o conexión rechazada de %s", peer)
		return
	}

	if statusCode == fasthttp.StatusNoContent {
		log.Printf("⚠️ No hay caché de %s", peer)
		return
	}
	if statusCode != fasthttp.StatusOK {
		log.Printf("⚠️ No se pudo recuperar la caché de %s", peer)
		return
	}

	decompressedData, err := utils.DecompressData(compressedData)
	if err != nil {
		log.Printf("⚠️ Error al descomprimir datos de %s: %v", peer, err)
		return
	}

	var recoveredCache []internal.CacheEntry

	err = json.Unmarshal(decompressedData, &recoveredCache)
	if err != nil {
		log.Printf("⚠️ Error al parsear JSON de %s: %v", peer, err)
		return
	}

	// Bloqueamos el acceso a la caché antes de modificarla
	internal.CacheMutex.Lock()
	defer internal.CacheMutex.Unlock()

	// Limpiar la caché antes de poblarla con la nueva data
	cache.FlushAll()

	for _, entry := range recoveredCache {
		duration, err := time.ParseDuration(entry.ExpiresIn)
		if err != nil {
			log.Printf("⚠️ Error al parsear duración para clave %s: %v", entry.Key, err)
			continue
		}
		cache.Set(entry.Key, entry.Value, duration)
	}

	log.Println("✅ Caché recuperada con éxito desde", peer)
}

func RecoverCacheDiff(peerManager *PeerManager, cache *internal.Cache) {
	peer := peerManager.GetActivePeers()[0]
	if peer == "" {
		return
	}

	url := fmt.Sprintf("%s/diff", peer)
	statusCode, diffData, err := fasthttp.Get(nil, url)

	if err != nil || statusCode != fasthttp.StatusOK {
		log.Printf("⚠️ No se pudo recuperar el diff de %s: %v", peer, err)
		return
	}

	var remoteDiff map[string]int64
	err = json.Unmarshal(diffData, &remoteDiff)
	if err != nil {
		log.Printf("⚠️ Error al parsear JSON del diff de %s: %v", peer, err)
		return
	}

	// Obtener diff local
	localDiff := cache.GetDiff()

	// Detectar claves desactualizadas o faltantes
	var missingKeys []string
	for key, remoteExp := range remoteDiff {
		localExp, exists := localDiff[key]

		if !exists || remoteExp > localExp {
			missingKeys = append(missingKeys, key)
		}
	}

	// Detectar claves expiradas localmente que ya no existen en el peer
	for key := range localDiff {
		if _, exists := remoteDiff[key]; !exists {
			cache.RemoveKey(key) // Eliminar claves expiradas
		}
	}

	// Si hay claves desactualizadas, pedir sus valores
	if len(missingKeys) > 0 {
		FetchAndUpdateKeys(peer, cache, missingKeys)
	}
}

func FetchAndUpdateKeys(peer string, cache *internal.Cache, keys []string) {
	url := fmt.Sprintf("%s/getKeys", peer)
	body, _ := json.Marshal(keys)

	req := fasthttp.AcquireRequest()
	req.SetRequestURI(url)
	req.Header.SetMethod("POST")
	req.SetBody(body)

	resp := fasthttp.AcquireResponse()
	err := fasthttp.Do(req, resp)

	fasthttp.ReleaseRequest(req)
	defer fasthttp.ReleaseResponse(resp)

	if err != nil || resp.StatusCode() != fasthttp.StatusOK {
		log.Printf("⚠️ No se pudo recuperar claves de %s: %v", peer, err)
		return
	}

	var recoveredData map[string]struct {
		Value      interface{} `json:"value"`
		Expiration int64       `json:"expiration"`
	}

	err = json.Unmarshal(resp.Body(), &recoveredData)
	if err != nil {
		log.Printf("⚠️ Error al parsear claves de %s: %v", peer, err)
		return
	}

	internal.CacheMutex.Lock()
	defer internal.CacheMutex.Unlock()

	for key, entry := range recoveredData {
		timeTtl := time.Duration(entry.Expiration) * time.Second
		cache.Set(key, entry.Value, timeTtl)
	}

	log.Println("✅ Claves sincronizadas desde", peer)
}
