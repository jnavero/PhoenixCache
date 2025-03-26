package security

import "phoenixcache/configuration"

func InitModule(config *configuration.Config) {

	StartWhitelistUpdater(config)
}

func StartWhitelistUpdater(config *configuration.Config) {
	//White list init
	LoadWhitelist(config)
	go func() {
		WatchWhitelistFile(config)
	}()
}
