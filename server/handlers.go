package server

import (
	"encoding/json"
	"log"
	"phoenixcache/distributed"
	"phoenixcache/internal"

	"strconv"
	"time"

	"github.com/valyala/fasthttp"
)

//********************************************************************
// Handlers para el servidor fasthttp
//********************************************************************

// handleSet almacena un valor en la caché (key y ttl por GET, value por BODY)
func HandleSet(peerManager *distributed.PeerManager, cache *internal.Cache, ctx *fasthttp.RequestCtx) {
	key := string(ctx.QueryArgs().Peek("key"))
	ttlStr := string(ctx.QueryArgs().Peek("ttl"))

	if key == "" || ttlStr == "" {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		ctx.SetContentType("application/json")
		ctx.SetBodyString(`{"error": "❌ Parámetros 'key' y 'ttl' son requeridos"}`)
		return
	}

	ttl, err := strconv.Atoi(ttlStr)
	if err != nil {
		log.Print("mierder!")
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		ctx.SetContentType("application/json")
		ctx.SetBodyString(`{"error": "❌ TTL debe ser un número válido"}`)
		return
	}

	value := ctx.PostBody()
	timeTtl := time.Duration(ttl) * time.Second
	cache.Set(key, string(value), timeTtl)
	distributed.PropagateChange(distributed.SyncMessage{Action: "set", Key: key, Value: string(value), TTL: timeTtl}, peerManager)
	ctx.SetStatusCode(fasthttp.StatusOK)
}

// handleGet obtiene un valor de la caché
func HandleGet(cache *internal.Cache, ctx *fasthttp.RequestCtx) {
	key := string(ctx.QueryArgs().Peek("key"))

	if key == "" {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		ctx.SetBodyString(`{"error": "❌ Parámetro 'key' requerido"}`)
		return
	}

	value, found := cache.Get(key)
	if !found {
		ctx.SetStatusCode(fasthttp.StatusNotFound)
		return
	}

	ctx.SetContentType("application/json")
	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.SetBodyString(value.(string))
}

// handleList devuelve un listado de la caché
func HandleList(cache *internal.Cache, ctx *fasthttp.RequestCtx) {
	items := cache.List()
	jsonResponse, _ := json.Marshal(items)
	ctx.SetContentType("application/json")
	ctx.SetBody(jsonResponse)
}

// handleFlushAll borra toda la caché
func HandleFlushAll(peerManager *distributed.PeerManager, cache *internal.Cache, ctx *fasthttp.RequestCtx) {
	cache.FlushAll()
	distributed.PropagateChange(distributed.SyncMessage{Action: "flush", Key: "", Value: nil, TTL: 0}, peerManager)
	ctx.SetStatusCode(fasthttp.StatusOK)
}

// handle que elimina una Key concreta
func HandleRemoveKey(peerManager *distributed.PeerManager, cache *internal.Cache, ctx *fasthttp.RequestCtx) {
	key := string(ctx.QueryArgs().Peek("key"))
	if key == "" {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		ctx.SetContentType("application/json")
		ctx.SetBodyString(`{"error": "❌ Parámetro 'key' es requerido"}`)
		return
	}

	cache.RemoveKey(key)
	distributed.PropagateChange(distributed.SyncMessage{Action: "remove", Key: key, Value: nil, TTL: 0}, peerManager)
	ctx.SetStatusCode(fasthttp.StatusOK)
}

// handle que liminakeys de la cache en funcion de un parametro
func HandleDeleteByPattern(peerManager *distributed.PeerManager, cache *internal.Cache, ctx *fasthttp.RequestCtx) {
	pattern := string(ctx.QueryArgs().Peek("key"))
	if pattern == "" {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		ctx.SetContentType("application/json")
		ctx.SetBodyString(`{"error": "❌ Parámetro 'key' es requerido"}`)
		return
	}

	deletedKeys := cache.RemovePatternKey(pattern)

	distributed.PropagateChange(distributed.SyncMessage{Action: "removePattern", Key: pattern, Value: nil, TTL: 0}, peerManager)

	jsonResponse, _ := json.Marshal(deletedKeys)

	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.SetContentType("application/json")
	ctx.SetBody(jsonResponse)
}

// handle para obtener un item con la expiración de la cache
func HandleTryGetWithExpire(cache *internal.Cache, ctx *fasthttp.RequestCtx) {
	key := string(ctx.QueryArgs().Peek("key"))
	if key == "" {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		ctx.SetContentType("application/json")
		ctx.SetBodyString(`{"error": "❌ Parámetro 'key' es requerido"}`)
		return
	}

	value, expTime, found := cache.GetWithExpiry(key)
	if !found {
		ctx.SetStatusCode(fasthttp.StatusNotFound)
		return
	}

	timeLeft := expTime.(time.Time).Sub(time.Now())
	if timeLeft <= 0 {
		cache.RemoveKey(key)
		ctx.SetStatusCode(fasthttp.StatusNotFound)
		ctx.SetContentType("application/json")
		ctx.SetBodyString(`{"error": "❌ Clave expirada"}`)
		return
	}

	response := map[string]interface{}{
		"value":      value,
		"expires_in": timeLeft.Seconds(),
	}

	jsonResponse, _ := json.Marshal(response)

	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.SetContentType("application/json")
	ctx.SetBody(jsonResponse)

}

// handle que obtiene todas los values de las keys indicadas
func HandleGetKeys(cache *internal.Cache, ctx *fasthttp.RequestCtx) {
	var keys []string

	// Parsear el JSON recibido
	err := json.Unmarshal(ctx.PostBody(), &keys)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		ctx.SetBody([]byte(`{"error": "Invalid request body"}`))
		return
	}

	// Respuesta con valores encontrados
	response := make(map[string]struct {
		Value      interface{} `json:"value"`
		Expiration time.Time   `json:"expiration"`
	})

	internal.CacheMutex.Lock()
	defer internal.CacheMutex.Unlock()

	for _, key := range keys {
		val, found := cache.Get(key)
		val, expTime, found := cache.GetWithExpiry(key)
		timeLeft := expTime.(time.Time)
		if found {
			response[key] = struct {
				Value      interface{} `json:"value"`
				Expiration time.Time   `json:"expiration"`
			}{
				Value:      val,
				Expiration: timeLeft,
			}
		}
	}

	// Serializar la respuesta
	data, err := json.Marshal(response)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		ctx.SetBody([]byte(`{"error": "Error encoding response"}`))
		return
	}

	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.SetContentType("application/json")
	ctx.SetBody(data)
}
