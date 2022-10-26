package api

import (
	"context"
	projectapp "learn_zinx/Cobra.mayfly/internal/project/application"
	"learn_zinx/Cobra.mayfly/internal/redis/api/form"
	"learn_zinx/Cobra.mayfly/internal/redis/api/vo"
	"learn_zinx/Cobra.mayfly/internal/redis/application"
	"learn_zinx/Cobra.mayfly/internal/redis/domain/entity"
	"learn_zinx/Cobra.mayfly/pkg/biz"
	"learn_zinx/Cobra.mayfly/pkg/ctx"
	"learn_zinx/Cobra.mayfly/pkg/ginx"
	"learn_zinx/Cobra.mayfly/pkg/utils"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
)

type Redis struct {
	RedisApp   application.Redis
	ProjectApp projectapp.Project
}

func (r *Redis) RedisList(rc *ctx.ReqCtx) {
	g := rc.GinCtx
	m := &entity.Redis{EnvId: uint64(ginx.QueryInt(g, "envId", 0)),
		ProjectId: uint64(ginx.QueryInt(g, "projectId", 0)),
	}
	m.CreatorId = rc.LoginAccount.Id
	rc.ResData = r.RedisApp.GetPageList(m, ginx.GetPageParam(rc.GinCtx), new([]vo.Redis))
}

func (r *Redis) Save(rc *ctx.ReqCtx) {
	form := &form.Redis{}
	ginx.BindJsonAndValid(rc.GinCtx, form)

	redis := new(entity.Redis)
	utils.Copy(redis, form)

	// 密码解密，并使用解密后的赋值
	originPwd, err := utils.DefaultRsaDecrypt(redis.Password, true)
	biz.ErrIsNilAppendErr(err, "解密密码错误: %s")
	redis.Password = originPwd

	// 密码脱敏记录日志
	form.Password = "****"
	rc.ReqParam = form

	redis.SetBaseInfo(rc.LoginAccount)
	r.RedisApp.Save(redis)
}

// 获取redis实例密码，由于数据库是加密存储，故提供该接口展示原文密码
func (r *Redis) GetRedisPwd(rc *ctx.ReqCtx) {
	rid := uint64(ginx.PathParamInt(rc.GinCtx, "id"))
	re := r.RedisApp.GetById(rid, "Password")
	re.PwdDecrypt()
	rc.ResData = re.Password
}

func (r *Redis) DeleteRedis(rc *ctx.ReqCtx) {
	r.RedisApp.Delete(uint64(ginx.PathParamInt(rc.GinCtx, "id")))
}

func (r *Redis) RedisInfo(rc *ctx.ReqCtx) {
	ri := r.RedisApp.GetRedisInstance(uint64(ginx.PathParamInt(rc.GinCtx, "id")))

	var res string
	var err error

	ctx := context.Background()
	if ri.Mode == "" || ri.Mode == entity.RedisModeStandalone || ri.Mode == entity.RedisModeSentinel {
		res, err = ri.Cli.Info(ctx).Result()
	} else if ri.Mode == entity.RedisModeCluster {
		host := rc.GinCtx.Query("host")
		biz.NotEmpty(host, "集群模式host信息不能为空")
		clusterClient := ri.ClusterCli
		var redisClient *redis.Client
		// 遍历集群的master节点找到该redis client
		clusterClient.ForEachMaster(ctx, func(ctx context.Context, client *redis.Client) error {
			if host == client.Options().Addr {
				redisClient = client
			}
			return nil
		})
		if redisClient == nil {
			// 遍历集群的slave节点找到该redis client
			clusterClient.ForEachSlave(ctx, func(ctx context.Context, client *redis.Client) error {
				if host == client.Options().Addr {
					redisClient = client
				}
				return nil
			})
		}
		biz.NotNil(redisClient, "该实例不在该集群中")
		res, err = redisClient.Info(ctx).Result()
	}

	biz.ErrIsNilAppendErr(err, "获取redis info失败: %s")

	datas := strings.Split(res, "\r\n")
	i := 0
	length := len(datas)

	parseMap := make(map[string]map[string]string)
	for {
		if i >= length {
			break
		}
		if strings.Contains(datas[i], "#") {
			key := utils.SubString(datas[i], strings.Index(datas[i], "#")+1, utils.StrLen(datas[i]))
			i++
			key = strings.Trim(key, " ")

			sectionMap := make(map[string]string)
			for {
				if i >= length || !strings.Contains(datas[i], ":") {
					break
				}
				pair := strings.Split(datas[i], ":")
				i++
				if len(pair) != 2 {
					continue
				}
				sectionMap[pair[0]] = pair[1]
			}
			parseMap[key] = sectionMap
		} else {
			i++
		}
	}
	rc.ResData = parseMap
}

func (r *Redis) ClusterInfo(rc *ctx.ReqCtx) {
	g := rc.GinCtx
	ri := r.RedisApp.GetRedisInstance(uint64(ginx.PathParamInt(g, "id")))
	biz.IsEquals(ri.Mode, entity.RedisModeCluster, "非集群模式")
	info, _ := ri.ClusterCli.ClusterInfo(context.Background()).Result()
	nodesStr, _ := ri.ClusterCli.ClusterNodes(context.Background()).Result()

	nodesRes := make([]map[string]string, 0)
	nodes := strings.Split(nodesStr, "\n")
	for _, node := range nodes {
		if node == "" {
			continue
		}
		nodeInfos := strings.Split(node, " ")
		node := make(map[string]string)
		node["nodeId"] = nodeInfos[0]
		// ip:port1@port2：port1指redis服务器与客户端通信的端口，port2则是集群内部节点间通信的端口
		node["ip"] = nodeInfos[1]
		node["flags"] = nodeInfos[2]
		// 如果节点是slave，并且已知master节点，则为master节点ID；否则为符号"-"
		node["masterSlaveRelation"] = nodeInfos[3]
		// 最近一次发送ping的时间，这个时间是一个unix毫秒时间戳，0代表没有发送过
		node["pingSent"] = nodeInfos[4]
		// 最近一次收到pong的时间，使用unix时间戳表示
		node["pongRecv"] = nodeInfos[5]
		// 节点的epoch值（如果该节点是从节点，则为其主节点的epoch值）。每当节点发生失败切换时，都会创建一个新的，独特的，递增的epoch。
		// 如果多个节点竞争同一个哈希槽时，epoch值更高的节点会抢夺到
		node["configEpoch"] = nodeInfos[6]
		// node-to-node集群总线使用的链接的状态，我们使用这个链接与集群中其他节点进行通信.值可以是 connected 和 disconnected
		node["linkState"] = nodeInfos[7]
		// slave节点没有插槽信息
		if len(nodeInfos) > 8 {
			// slot：master节点第9位为哈希槽值或者一个哈希槽范围，代表当前节点可以提供服务的所有哈希槽值。如果只是一个值,那就是只有一个槽会被使用。
			// 如果是一个范围，这个值表示为起始槽-结束槽，节点将处理包括起始槽和结束槽在内的所有哈希槽。
			node["slot"] = nodeInfos[8]
		}
		nodesRes = append(nodesRes, node)
	}
	rc.ResData = map[string]interface{}{
		"clusterInfo":  info,
		"clusterNodes": nodesRes,
	}
}

// scan获取redis的key列表信息
func (r *Redis) Scan(rc *ctx.ReqCtx) {
	g := rc.GinCtx
	ri := r.RedisApp.GetRedisInstance(uint64(ginx.PathParamInt(g, "id")))
	biz.ErrIsNilAppendErr(r.ProjectApp.CanAccess(rc.LoginAccount.Id, ri.ProjectId), "%s")

	form := &form.RedisScanForm{}
	ginx.BindJsonAndValid(rc.GinCtx, form)

	cmd := ri.GetCmdable()
	ctx := context.Background()

	kis := make([]*vo.KeyInfo, 0)
	var cursorRes map[string]uint64 = make(map[string]uint64)

	if ri.Mode == "" || ri.Mode == entity.RedisModeStandalone || ri.Mode == entity.RedisModeSentinel {
		redisAddr := ri.Cli.Options().Addr
		keys, cursor := ri.Scan(form.Cursor[redisAddr], form.Match, form.Count)
		cursorRes[redisAddr] = cursor

		var keyInfoSplit []string
		if len(keys) > 0 {
			keyInfosLua := `local result = {}
							-- KEYS[1]为第1个参数,lua数组下标从1开始
							for i = 1, #KEYS do
								local ttl = redis.call('ttl', KEYS[i]);
								local keyType = redis.call('type', KEYS[i]);
								table.insert(result, string.format("%d,%s", ttl, keyType['ok']));
							end;
							return table.concat(result, ".");`
			// 通过lua获取 ttl,type.ttl2,type2格式，以便下面切割获取ttl和type。避免多次调用ttl和type函数
			keyInfos, err := cmd.Eval(ctx, keyInfosLua, keys).Result()
			biz.ErrIsNilAppendErr(err, "执行lua脚本获取key信息失败: %s")
			keyInfoSplit = strings.Split(keyInfos.(string), ".")
		}

		for i, k := range keys {
			ttlType := strings.Split(keyInfoSplit[i], ",")
			ttl, _ := strconv.Atoi(ttlType[0])
			ki := &vo.KeyInfo{Key: k, Type: ttlType[1], Ttl: int64(ttl)}
			kis = append(kis, ki)
		}
	} else if ri.Mode == entity.RedisModeCluster {
		var keys []string

		mu := &sync.Mutex{}
		// 遍历所有master节点，并执行scan命令，合并keys
		ri.ClusterCli.ForEachMaster(ctx, func(ctx context.Context, client *redis.Client) error {
			redisAddr := client.Options().Addr
			ks, cursor, _ := client.Scan(ctx, form.Cursor[redisAddr], form.Match, form.Count).Result()
			// 遍历节点的内部回调函数使用异步调用，如不加锁会导致集合并发错误
			mu.Lock()
			cursorRes[redisAddr] = cursor
			keys = append(keys, ks...)
			mu.Unlock()
			return nil
		})

		// 因为redis集群模式执行lua脚本key必须位于同一slot中，故单机获取的方式不适合
		// 使用lua获取key的ttl以及类型，减少网络调用
		keyInfoLua := `local ttl = redis.call('ttl', KEYS[1]);
					   local keyType = redis.call('type', KEYS[1]);
					   return string.format("%d,%s", ttl, keyType['ok'])`
		for _, k := range keys {
			keyInfo, err := cmd.Eval(ctx, keyInfoLua, []string{k}).Result()
			biz.ErrIsNilAppendErr(err, "执行lua脚本获取key信息失败: %s")
			ttlType := strings.Split(keyInfo.(string), ",")
			ttl, _ := strconv.Atoi(ttlType[0])
			ki := &vo.KeyInfo{Key: k, Type: ttlType[1], Ttl: int64(ttl)}
			kis = append(kis, ki)
		}
	}

	size, _ := cmd.DBSize(context.TODO()).Result()
	rc.ResData = &vo.Keys{Cursor: cursorRes, Keys: kis, DbSize: size}
}

func (r *Redis) DeleteKey(rc *ctx.ReqCtx) {
	g := rc.GinCtx
	key := g.Query("key")
	biz.NotEmpty(key, "key不能为空")

	ri := r.RedisApp.GetRedisInstance(uint64(ginx.PathParamInt(g, "id")))
	biz.ErrIsNilAppendErr(r.ProjectApp.CanAccess(rc.LoginAccount.Id, ri.ProjectId), "%s")

	rc.ReqParam = key
	ri.GetCmdable().Del(context.Background(), key)
}

func (r *Redis) checkKey(rc *ctx.ReqCtx) (*application.RedisInstance, string) {
	g := rc.GinCtx
	key := g.Query("key")
	biz.NotEmpty(key, "key不能为空")

	ri := r.RedisApp.GetRedisInstance(uint64(ginx.PathParamInt(g, "id")))
	biz.ErrIsNilAppendErr(r.ProjectApp.CanAccess(rc.LoginAccount.Id, ri.ProjectId), "%s")

	return ri, key
}

func (r *Redis) GetStringValue(rc *ctx.ReqCtx) {
	ri, key := r.checkKey(rc)
	str, err := ri.GetCmdable().Get(context.TODO(), key).Result()
	biz.ErrIsNilAppendErr(err, "获取字符串值失败: %s")
	rc.ResData = str
}

func (r *Redis) SetStringValue(rc *ctx.ReqCtx) {
	g := rc.GinCtx
	keyValue := new(form.StringValue)
	ginx.BindJsonAndValid(g, keyValue)

	ri := r.RedisApp.GetRedisInstance(uint64(ginx.PathParamInt(g, "id")))
	biz.ErrIsNilAppendErr(r.ProjectApp.CanAccess(rc.LoginAccount.Id, ri.ProjectId), "%s")

	str, err := ri.GetCmdable().Set(context.TODO(), keyValue.Key, keyValue.Value, time.Second*time.Duration(keyValue.Timed)).Result()
	biz.ErrIsNilAppendErr(err, "保存字符串值失败: %s")
	rc.ResData = str
}

func (r *Redis) Hscan(rc *ctx.ReqCtx) {
	ri, key := r.checkKey(rc)
	g := rc.GinCtx
	count := ginx.QueryInt(g, "count", 10)
	match := g.Query("match")
	cursor := ginx.QueryInt(g, "cursor", 0)
	contextTodo := context.TODO()

	cmdable := ri.GetCmdable()
	keys, nextCursor, err := cmdable.HScan(contextTodo, key, uint64(cursor), match, int64(count)).Result()
	biz.ErrIsNilAppendErr(err, "hcan err: %s")
	keySize, err := cmdable.HLen(contextTodo, key).Result()
	biz.ErrIsNilAppendErr(err, "hlen err: %s")

	rc.ResData = map[string]interface{}{
		"keys":    keys,
		"cursor":  nextCursor,
		"keySize": keySize,
	}
}

func (r *Redis) Hdel(rc *ctx.ReqCtx) {
	ri, key := r.checkKey(rc)
	field := rc.GinCtx.Query("field")

	delRes, err := ri.GetCmdable().HDel(context.TODO(), key, field).Result()
	biz.ErrIsNilAppendErr(err, "hdel err: %s")
	rc.ResData = delRes
}

func (r *Redis) Hget(rc *ctx.ReqCtx) {
	ri, key := r.checkKey(rc)
	field := rc.GinCtx.Query("field")

	res, err := ri.GetCmdable().HGet(context.TODO(), key, field).Result()
	biz.ErrIsNilAppendErr(err, "hget err: %s")
	rc.ResData = res
}

func (r *Redis) SetHashValue(rc *ctx.ReqCtx) {
	g := rc.GinCtx
	hashValue := new(form.HashValue)
	ginx.BindJsonAndValid(g, hashValue)

	ri := r.RedisApp.GetRedisInstance(uint64(ginx.PathParamInt(g, "id")))
	biz.ErrIsNilAppendErr(r.ProjectApp.CanAccess(rc.LoginAccount.Id, ri.ProjectId), "%s")

	cmd := ri.GetCmdable()
	key := hashValue.Key
	contextTodo := context.TODO()
	for _, v := range hashValue.Value {
		res := cmd.HSet(contextTodo, key, v["field"].(string), v["value"])
		biz.ErrIsNilAppendErr(res.Err(), "保存hash值失败: %s")
	}
	if hashValue.Timed != 0 && hashValue.Timed != -1 {
		cmd.Expire(context.TODO(), key, time.Second*time.Duration(hashValue.Timed))
	}
}

func (r *Redis) GetSetValue(rc *ctx.ReqCtx) {
	ri, key := r.checkKey(rc)
	res, err := ri.GetCmdable().SMembers(context.TODO(), key).Result()
	biz.ErrIsNilAppendErr(err, "获取set值失败: %s")
	rc.ResData = res
}

func (r *Redis) SetSetValue(rc *ctx.ReqCtx) {
	g := rc.GinCtx
	keyvalue := new(form.SetValue)
	ginx.BindJsonAndValid(g, keyvalue)

	ri := r.RedisApp.GetRedisInstance(uint64(ginx.PathParamInt(g, "id")))
	biz.ErrIsNilAppendErr(r.ProjectApp.CanAccess(rc.LoginAccount.Id, ri.ProjectId), "%s")
	cmd := ri.GetCmdable()

	key := keyvalue.Key
	// 简单处理->先删除，后新增
	cmd.Del(context.TODO(), key)
	cmd.SAdd(context.TODO(), key, keyvalue.Value...)

	if keyvalue.Timed != -1 {
		cmd.Expire(context.TODO(), key, time.Second*time.Duration(keyvalue.Timed))
	}
}

func (r *Redis) GetListValue(rc *ctx.ReqCtx) {
	ri, key := r.checkKey(rc)
	ctx := context.TODO()
	cmdable := ri.GetCmdable()

	len, err := cmdable.LLen(ctx, key).Result()
	biz.ErrIsNilAppendErr(err, "获取list长度失败: %s")

	g := rc.GinCtx
	start := ginx.QueryInt(g, "start", 0)
	stop := ginx.QueryInt(g, "stop", 10)
	res, err := cmdable.LRange(ctx, key, int64(start), int64(stop)).Result()
	biz.ErrIsNilAppendErr(err, "获取list值失败: %s")

	rc.ResData = map[string]interface{}{
		"len":  len,
		"list": res,
	}
}

func (r *Redis) SaveListValue(rc *ctx.ReqCtx) {
	g := rc.GinCtx
	listValue := new(form.ListValue)
	ginx.BindJsonAndValid(g, listValue)

	ri := r.RedisApp.GetRedisInstance(uint64(ginx.PathParamInt(g, "id")))
	biz.ErrIsNilAppendErr(r.ProjectApp.CanAccess(rc.LoginAccount.Id, ri.ProjectId), "%s")
	cmd := ri.GetCmdable()

	key := listValue.Key
	ctx := context.TODO()
	for _, v := range listValue.Value {
		cmd.RPush(ctx, key, v)
	}

	if listValue.Timed != -1 {
		cmd.Expire(context.TODO(), key, time.Second*time.Duration(listValue.Timed))
	}
}

func (r *Redis) SetListValue(rc *ctx.ReqCtx) {
	g := rc.GinCtx
	listSetValue := new(form.ListSetValue)
	ginx.BindJsonAndValid(g, listSetValue)

	ri := r.RedisApp.GetRedisInstance(uint64(ginx.PathParamInt(g, "id")))
	biz.ErrIsNilAppendErr(r.ProjectApp.CanAccess(rc.LoginAccount.Id, ri.ProjectId), "%s")

	_, err := ri.GetCmdable().LSet(context.TODO(), listSetValue.Key, listSetValue.Index, listSetValue.Value).Result()
	biz.ErrIsNilAppendErr(err, "list set失败: %s")
}
