import{A as e}from"./Api.1663901876151.js";const s={redisList:e.create("/redis","get"),getRedisPwd:e.create("/redis/{id}/pwd","get"),redisInfo:e.create("/redis/{id}/info","get"),clusterInfo:e.create("/redis/{id}/cluster-info","get"),saveRedis:e.create("/redis","post"),delRedis:e.create("/redis/{id}","delete"),scan:e.create("/redis/{id}/scan","post"),getStringValue:e.create("/redis/{id}/string-value","get"),saveStringValue:e.create("/redis/{id}/string-value","post"),getHashValue:e.create("/redis/{id}/hash-value","get"),hscan:e.create("/redis/{id}/hscan","get"),hget:e.create("/redis/{id}/hget","get"),hdel:e.create("/redis/{id}/hdel","delete"),saveHashValue:e.create("/redis/{id}/hash-value","post"),getSetValue:e.create("/redis/{id}/set-value","get"),saveSetValue:e.create("/redis/{id}/set-value","post"),del:e.create("/redis/{id}/scan/{cursor}/{count}","delete"),delKey:e.create("/redis/{id}/key","delete"),getListValue:e.create("/redis/{id}/list-value","get"),saveListValue:e.create("/redis/{id}/list-value","post"),setListValue:e.create("/redis/{id}/list-value/lset","post")};export{s as r};
