package server

import (
	"phoenixcache/configuration"
	"phoenixcache/distributed"
	"phoenixcache/internal"
	"phoenixcache/security"

	"github.com/valyala/fasthttp"
)

// SetupRouter configura las rutas del servidor
func SetupRouter(config *configuration.Config, peerManager *distributed.PeerManager, cache *internal.Cache) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {

		if !isAllowedNode(config, ctx) {
			return
		}

		switch string(ctx.Path()) {
		case "/set":
			HandleSet(peerManager, cache, ctx)
		case "/get":
			HandleGet(cache, ctx)
		case "/trygetwithexpire":
			HandleTryGetWithExpire(cache, ctx)
		case "/getKeys":
			HandleGetKeys(cache, ctx)
		case "/list":
			HandleList(cache, ctx)
		case "/flush":
			HandleFlushAll(peerManager, cache, ctx)
		case "/removeallkeys":
			HandleDeleteByPattern(peerManager, cache, ctx)
		case "/remove":
			HandleRemoveKey(peerManager, cache, ctx)
		case "/sync":
			distributed.SyncHandler(cache, ctx)
		case "/ping":
			distributed.HandlePing(ctx)
		case "/export":
			distributed.HandleExportCache(cache, ctx)
		case "/diff":
			distributed.HandleDiff(cache, ctx)
		case "/set_batch":
			distributed.HandleSetBatch(peerManager, cache)
		default:
			ctx.SetStatusCode(fasthttp.StatusNotFound)
		}
	}
}

func isAllowedNode(config *configuration.Config, ctx *fasthttp.RequestCtx) bool {

	ip := ctx.RemoteIP().String()
	host := string(ctx.Request.Host())

	// Comprobar si la IP está en la lista blanca
	if !security.IsAllowedNode(ip) && !security.IsAllowedNode(host) {
		ctx.SetStatusCode(fasthttp.StatusForbidden)
		ctx.SetBody([]byte("⛔ Acceso denegado"))
		return false
	}
	// Continuar con la lógica normal...
	return true
}
