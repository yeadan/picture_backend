package data

var clienteCache CacheProvider = nil

func GetCacheClient() CacheProvider {
	if clienteCache == nil {
		clienteCache = NewGoCacheClient()
	}
	return clienteCache
}
