package distributed

import (
	"encoding/json"
	"phoenixcache/internal"
	"phoenixcache/utils"

	"github.com/valyala/fasthttp"
)

// Handler de sincronización, cuando se propagan los cambios, manejamos las diferentes acciones
func SyncHandler(cache *internal.Cache, ctx *fasthttp.RequestCtx) {
	var msg SyncMessage
	if err := json.Unmarshal(ctx.PostBody(), &msg); err != nil {
		ctx.Error("❌ Error en el payload", fasthttp.StatusBadRequest)
		return
	}

	switch msg.Action {
	case "set":
		cache.Set(msg.Key, msg.Value, msg.TTL)
	case "remove":
		cache.RemoveKey(msg.Key)
	case "removePattern":
		cache.RemovePatternKey(msg.Key)
	case "flush":
		cache.FlushAll()
	}
	ctx.Response.Header.Set("Connection", "close")
	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.SetBodyString("Ok")
}

// Handle encargado de exportar la cache para recuperarla en otro servidor
func HandleExportCache(cache *internal.Cache, ctx *fasthttp.RequestCtx) {
	internal.CacheMutex.Lock()
	defer internal.CacheMutex.Unlock()

	items := cache.GetAll(false)
	if items == nil {
		ctx.SetStatusCode(fasthttp.StatusNoContent)
		return
	}

	data, err := json.Marshal(items)

	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		ctx.SetBody([]byte(`{"error": "Error exportando cache"}`))
		return
	}

	compressed := utils.CompressData(data)
	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.SetContentType("application/octet-stream")
	ctx.SetBody(compressed)
}

// Handle encargado de crear un diffing de la cache Key / expiration
func HandleDiff(cache *internal.Cache, ctx *fasthttp.RequestCtx) {
	internal.CacheMutex.Lock()
	defer internal.CacheMutex.Unlock()

	diff := cache.GetDiff()
	data, err := json.Marshal(diff)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		ctx.SetBody([]byte(`{"error": "Error obteniendo el diff"}`))
		return
	}

	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.SetBody(data)
}

// HandlePing es el handler para /ping
func HandlePing(ctx *fasthttp.RequestCtx) {
	ctx.SetStatusCode(fasthttp.StatusOK)
}

func HandleSetBatch(peerManager *PeerManager, cache *internal.Cache) {
	RecoverCacheDiff(peerManager, cache)
}
