package distributed

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/valyala/fasthttp"
)

// PeerManager gestiona los peers y su estado
type PeerManager struct {
	peers         map[string]int // Mapa de peers con fallos consecutivos
	mu            sync.Mutex
	maxFailures   int           // Número máximo de fallos antes de marcar un nodo como inactivo
	checkInterval time.Duration // Intervalo entre checks
}

// NewPeerManager crea un nuevo gestor de peers
func NewPeerManager(peers []string, checkInterval time.Duration, maxFailures int) *PeerManager {
	pm := &PeerManager{
		peers:         make(map[string]int),
		maxFailures:   maxFailures,
		checkInterval: checkInterval,
	}

	// Inicializamos la lista de peers
	for _, peer := range peers {
		pm.peers[peer] = 0 // 0 fallos al inicio
	}

	// Iniciamos el heartbeat
	go pm.startHeartbeat()

	return pm
}

// startHeartbeat inicia la verificación periódica de los peers
func (pm *PeerManager) startHeartbeat() {
	for {
		pm.checkPeers()
		time.Sleep(pm.checkInterval)
	}
}

// checkPeers verifica el estado de cada peer
func (pm *PeerManager) checkPeers() {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	for peer := range pm.peers {
		if pm.pingPeer(peer) {

			//No estaba activo y ahora si lo está...
			if !IsActive(peer, pm) {
				// Recuperar datos faltantes
				callSetBatch(peer)
			}
			//Poniendo el contador a cero se marca el peer activo :-)
			pm.peers[peer] = 0
		} else {

			if pm.peers[peer] < pm.maxFailures {
				pm.peers[peer]++
			}
		}
	}
}

func callSetBatch(peer string) {
	fasthttp.Get(nil, fmt.Sprintf("%s/set_batch", peer))
}

func IsActive(peer string, pm *PeerManager) bool {
	return pm.peers[peer] < pm.maxFailures
}

// pingPeer hace una solicitud al endpoint /ping de un peer
func (pm *PeerManager) pingPeer(peer string) bool {
	statusCode, _, err := fasthttp.Get(nil, fmt.Sprintf("%s/ping", peer))

	if err != nil || statusCode != http.StatusOK {
		return false
	}
	return true
}

// GetActivePeers devuelve una lista de nodos activos
func (pm *PeerManager) GetActivePeers() []string {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	var activePeers []string
	for peer, failures := range pm.peers {
		if failures < pm.maxFailures {
			activePeers = append(activePeers, peer)
		}
	}
	return activePeers
}
