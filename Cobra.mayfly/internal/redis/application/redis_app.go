package application

import (
	"context"
	"fmt"
	"learn_zinx/Cobra.mayfly/internal/constant"
	machineapp "learn_zinx/Cobra.mayfly/internal/machine/application"
	"learn_zinx/Cobra.mayfly/internal/machine/infrastructure/machine"
	"learn_zinx/Cobra.mayfly/internal/redis/domain/entity"
	"learn_zinx/Cobra.mayfly/internal/redis/domain/repository"
	"learn_zinx/Cobra.mayfly/pkg/biz"
	"learn_zinx/Cobra.mayfly/pkg/cache"
	"learn_zinx/Cobra.mayfly/pkg/global"
	"learn_zinx/Cobra.mayfly/pkg/model"
	"learn_zinx/Cobra.mayfly/pkg/utils"
	"net"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
)

type Redis interface {
	// 分页获取机器脚本信息列表
	GetPageList(condition *entity.Redis, pageParam *model.PageParam, toEntity interface{}, orderBy ...string) *model.PageResult

	Count(condition *entity.Redis) int64

	// 根据id获取
	GetById(id uint64, cols ...string) *entity.Redis

	// 根据条件获取
	GetRedisBy(condition *entity.Redis, cols ...string) error

	Save(entity *entity.Redis)

	// 删除数据库信息
	Delete(id uint64)

	// 获取数据库连接实例
	GetRedisInstance(id uint64) *RedisInstance
}

func newRedisApp(redisRepo repository.Redis) Redis {
	return &redisAppImpl{
		redisRepo: redisRepo,
	}
}

type redisAppImpl struct {
	redisRepo repository.Redis
}

// 分页获取机器脚本信息列表
func (r *redisAppImpl) GetPageList(condition *entity.Redis, pageParam *model.PageParam, toEntity interface{}, orderBy ...string) *model.PageResult {
	return r.redisRepo.GetRedisList(condition, pageParam, toEntity, orderBy...)
}

func (r *redisAppImpl) Count(condition *entity.Redis) int64 {
	return r.redisRepo.Count(condition)
}

// 根据id获取
func (r *redisAppImpl) GetById(id uint64, cols ...string) *entity.Redis {
	return r.redisRepo.GetById(id, cols...)
}

// 根据条件获取
func (r *redisAppImpl) GetRedisBy(condition *entity.Redis, cols ...string) error {
	return r.redisRepo.GetRedis(condition, cols...)
}

func (r *redisAppImpl) Save(re *entity.Redis) {
	// ’修改信息且密码不为空‘ or ‘新增’需要测试是否可连接
	if (re.Id != 0 && re.Password != "") || re.Id == 0 {
		TestRedisConnection(re)
	}

	// 查找是否存在该库
	oldRedis := &entity.Redis{Host: re.Host, Db: re.Db}
	err := r.GetRedisBy(oldRedis)

	if re.Id == 0 {
		biz.IsTrue(err != nil, "该库已存在")
		re.PwdEncrypt()
		r.redisRepo.Insert(re)
	} else {
		// 如果存在该库，则校验修改的库是否为该库
		if err == nil {
			biz.IsTrue(oldRedis.Id == re.Id, "该库已存在")
		}
		// 先关闭数据库连接
		CloseRedis(re.Id)
		re.PwdEncrypt()
		r.redisRepo.Update(re)
	}
}

// 删除Redis信息
func (r *redisAppImpl) Delete(id uint64) {
	CloseRedis(id)
	r.redisRepo.Delete(id)
}

// 获取数据库连接实例
func (r *redisAppImpl) GetRedisInstance(id uint64) *RedisInstance {
	// Id不为0，则为需要缓存
	needCache := id != 0
	if needCache {
		load, ok := redisCache.Get(id)
		if ok {
			return load.(*RedisInstance)
		}
	}
	// 缓存不存在，则回调获取redis信息
	re := r.GetById(id)
	re.PwdDecrypt()
	biz.NotNil(re, "redis信息不存在")

	redisMode := re.Mode
	var ri *RedisInstance
	if redisMode == "" || redisMode == entity.RedisModeStandalone {
		ri = getRedisCient(re)
		// 测试连接
		_, e := ri.Cli.Ping(context.Background()).Result()
		if e != nil {
			ri.Close()
			panic(biz.NewBizErr(fmt.Sprintf("redis连接失败: %s", e.Error())))
		}
	} else if redisMode == entity.RedisModeCluster {
		ri = getRedisClusterClient(re)
		// 测试连接
		_, e := ri.ClusterCli.Ping(context.Background()).Result()
		if e != nil {
			ri.Close()
			panic(biz.NewBizErr(fmt.Sprintf("redis集群连接失败: %s", e.Error())))
		}
	} else if redisMode == entity.RedisModeSentinel {
		ri = getRedisSentinelCient(re)
		// 测试连接
		_, e := ri.Cli.Ping(context.Background()).Result()
		if e != nil {
			ri.Close()
			panic(biz.NewBizErr(fmt.Sprintf("redis sentinel连接失败: %s", e.Error())))
		}
	}

	global.Log.Infof("连接redis: %s", re.Host)
	if needCache {
		redisCache.Put(re.Id, ri)
	}
	return ri
}

func getRedisCient(re *entity.Redis) *RedisInstance {
	ri := &RedisInstance{Id: re.Id, ProjectId: re.ProjectId, Mode: re.Mode}

	redisOptions := &redis.Options{
		Addr:         re.Host,
		Password:     re.Password, // no password set
		DB:           re.Db,       // use default DB
		DialTimeout:  8 * time.Second,
		ReadTimeout:  -1, // Disable timeouts, because SSH does not support deadlines.
		WriteTimeout: -1,
	}
	if re.EnableSshTunnel == 1 {
		ri.sshTunnelMachineId = re.SshTunnelMachineId
		redisOptions.Dialer = getRedisDialer(re.SshTunnelMachineId)
	}
	ri.Cli = redis.NewClient(redisOptions)
	return ri
}

func getRedisClusterClient(re *entity.Redis) *RedisInstance {
	ri := &RedisInstance{Id: re.Id, ProjectId: re.ProjectId, Mode: re.Mode}

	redisClusterOptions := &redis.ClusterOptions{
		Addrs:       strings.Split(re.Host, ","),
		Password:    re.Password,
		DialTimeout: 8 * time.Second,
	}
	if re.EnableSshTunnel == 1 {
		ri.sshTunnelMachineId = re.SshTunnelMachineId
		redisClusterOptions.Dialer = getRedisDialer(re.SshTunnelMachineId)
	}
	ri.ClusterCli = redis.NewClusterClient(redisClusterOptions)
	return ri
}

func getRedisSentinelCient(re *entity.Redis) *RedisInstance {
	ri := &RedisInstance{Id: re.Id, ProjectId: re.ProjectId, Mode: re.Mode}
	// sentinel模式host为 masterName=host:port,host:port
	masterNameAndHosts := strings.Split(re.Host, "=")
	sentinelOptions := &redis.FailoverOptions{
		MasterName:    masterNameAndHosts[0],
		SentinelAddrs: strings.Split(masterNameAndHosts[1], ","),
		Password:      re.Password, // no password set
		DB:            re.Db,       // use default DB
		DialTimeout:   8 * time.Second,
		ReadTimeout:   -1, // Disable timeouts, because SSH does not support deadlines.
		WriteTimeout:  -1,
	}
	if re.EnableSshTunnel == 1 {
		ri.sshTunnelMachineId = re.SshTunnelMachineId
		sentinelOptions.Dialer = getRedisDialer(re.SshTunnelMachineId)
	}
	ri.Cli = redis.NewFailoverClient(sentinelOptions)
	return ri
}

func getRedisDialer(machineId uint64) func(ctx context.Context, network, addr string) (net.Conn, error) {
	sshTunnel := machineapp.GetMachineApp().GetSshTunnelMachine(machineId)
	return func(_ context.Context, network, addr string) (net.Conn, error) {
		if sshConn, err := sshTunnel.GetDialConn(network, addr); err == nil {
			// 将ssh conn包装，否则redis内部设置超时会报错,ssh conn不支持设置超时会返回错误: ssh: tcpChan: deadline not supported
			return &utils.WrapSshConn{Conn: sshConn}, nil
		} else {
			return nil, err
		}
	}
}

//------------------------------------------------------------------------------

// redis客户端连接缓存，指定时间内没有访问则会被关闭
var redisCache = cache.NewTimedCache(constant.RedisConnExpireTime, 5*time.Second).
	WithUpdateAccessTime(true).
	OnEvicted(func(key interface{}, value interface{}) {
		global.Log.Info(fmt.Sprintf("删除redis连接缓存 id = %d", key))
		value.(*RedisInstance).Close()
	})

// 移除redis连接缓存并关闭redis连接
func CloseRedis(id uint64) {
	redisCache.Delete(id)
}

func init() {
	machine.AddCheckSshTunnelMachineUseFunc(func(machineId uint64) bool {
		// 遍历所有redis连接实例，若存在redis实例使用该ssh隧道机器，则返回true，表示还在使用中...
		items := redisCache.Items()
		for _, v := range items {
			if v.Value.(*RedisInstance).sshTunnelMachineId == machineId {
				return true
			}
		}
		return false
	})
}

func TestRedisConnection(re *entity.Redis) {
	var cmd redis.Cmdable
	if re.Mode == "" || re.Mode == entity.RedisModeStandalone {
		rcli := getRedisCient(re)
		defer rcli.Close()
		cmd = rcli.Cli
	} else if re.Mode == entity.RedisModeCluster {
		ccli := getRedisClusterClient(re)
		defer ccli.Close()
		cmd = ccli.ClusterCli
	} else if re.Mode == entity.RedisModeSentinel {
		rcli := getRedisSentinelCient(re)
		defer rcli.Close()
		cmd = rcli.Cli
	}

	// 测试连接
	_, e := cmd.Ping(context.Background()).Result()
	biz.ErrIsNilAppendErr(e, "Redis连接失败: %s")
}

// redis实例
type RedisInstance struct {
	Id                 uint64
	ProjectId          uint64
	Mode               string
	Cli                *redis.Client
	ClusterCli         *redis.ClusterClient
	sshTunnelMachineId uint64
}

// 获取命令执行接口的具体实现
func (r *RedisInstance) GetCmdable() redis.Cmdable {
	redisMode := r.Mode
	if redisMode == "" || redisMode == entity.RedisModeStandalone || r.Mode == entity.RedisModeSentinel {
		return r.Cli
	}
	if redisMode == entity.RedisModeCluster {
		return r.ClusterCli
	}
	return nil
}

func (r *RedisInstance) Scan(cursor uint64, match string, count int64) ([]string, uint64) {
	keys, newcursor, err := r.GetCmdable().Scan(context.Background(), cursor, match, count).Result()
	biz.ErrIsNilAppendErr(err, "scan失败: %s")
	return keys, newcursor
}

func (r *RedisInstance) Close() {
	if r.Mode == entity.RedisModeStandalone || r.Mode == entity.RedisModeSentinel {
		if err := r.Cli.Close(); err != nil {
			global.Log.Errorf("关闭redis单机实例[%d]连接失败: %s", r.Id, err.Error())
		}
		r.Cli = nil
	}
	if r.Mode == entity.RedisModeCluster {
		if err := r.ClusterCli.Close(); err != nil {
			global.Log.Errorf("关闭redis集群实例[%d]连接失败: %s", r.Id, err.Error())
		}
		r.ClusterCli = nil
	}
}
