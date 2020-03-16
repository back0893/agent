package services

import (
	"agent/src"
	"agent/src/agent/iface"
	"agent/src/g"
	"errors"
	"github.com/back0893/goTcp/utils"
	"os"
	"os/exec"
	"syscall"
)

type RedisService struct {
	CurrentStatus string
}

func (r RedisService) Status(map[string]string) bool {
	return g.Status(g.ReadPid("./redisPid"))
}

func NewRedisService() *RedisService {
	return &RedisService{}
}

func (r RedisService) Start(map[string]string) error {
	if g.Status(g.ReadPid("./redisPid")) {
		return errors.New("redis已经运行")
	}
	cmd := exec.Command("bash", "-c", "nohup redis-server >/dev/null 2>&1& echo $!>./redisPid")
	return cmd.Run()
}

func (r RedisService) Stop(map[string]string) error {
	pid := g.ReadPid("./redisPid")
	if pid == 0 {
		return errors.New("redis灭有在运行")
	}
	syscall.Kill(pid, syscall.SIGKILL)
	//参数pid
	os.Remove("./redisPid")
	return nil
}

func (r RedisService) Restart(args map[string]string) error {
	if err := r.Stop(args); err != nil {
		return err
	}
	if err := r.Start(args); err != nil {
		return err
	}
	return nil
}

func (r RedisService) Action(action string, args map[string]string) {
	var str = "未知命令"
	switch action {
	case "start":
		err := r.Start(args)
		if err != nil {
			str = "redis启动失败"
		} else {
			str = "redis启动成功"
		}
	case "stop":
		err := r.Stop(args)
		if err != nil {
			str = "redis停止失败"
		} else {
			str = "redis停止成功"
		}
	case "status":
		status := r.Status(args)
		if status {
			str = "redis正在运行"
		} else {
			str = "redis没有运行"
		}
	case "restart":
		err := r.Restart(args)
		if err != nil {
			str = "redis重启失败"
		} else {
			str = "redis重启成功"
		}
	}
	//无论怎么样,都会启动一个定时器,定时查询状态,以监控
	go r.Watcher()

	pkt := src.NewPkt()
	pkt.Id = g.ServiceResponse
	pkt.Data = []byte(str)
	a := utils.GlobalConfig.Get(g.AGENT).(iface.IAgent)
	err := a.GetCon().Write(pkt)
	if err != nil {
		//todo 发送失败..应该有后续操作
	}
}
func (r *RedisService) Watcher() {
	src.AddTimer(20, func() {
		start := r.Status(nil)

	})
}
