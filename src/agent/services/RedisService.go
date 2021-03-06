package services

import (
	"agent/src"
	"agent/src/agent/iface"
	"agent/src/g"
	"agent/src/g/model"
	"errors"
	"fmt"
	"github.com/back0893/goTcp/utils"
	"log"
	"os"
	"os/exec"
	"syscall"
	"time"
)

type RedisService struct {
	CurrentStatus int
	timerId       int64
}

func (r *RedisService) GetCurrentStatus() int {
	return r.CurrentStatus
}

func (r *RedisService) SetCurrentStatus(status int) {
	r.CurrentStatus = status
}

func (r *RedisService) Status(map[string]string) bool {
	return g.Status(g.ReadPid("./redisPid"))
}

func NewRedisService(status int) *RedisService {
	s := &RedisService{
		CurrentStatus: status,
	}
	s.Upload(map[string]string{})
	return s
}

func (r *RedisService) Start(args map[string]string) error {
	if g.Status(g.ReadPid("./redisPid")) {
		return nil
	}
	r.CurrentStatus = 1
	cmd := exec.Command("bash", "-c", "nohup redis-server >/dev/null 2>&1& echo $!>./redisPid")
	return cmd.Run()
}

func (r *RedisService) Stop(map[string]string) error {
	status := r.Status(nil)
	if status == false {
		return errors.New("redis灭有在运行")
	}
	fmt.Printf("%p\n", r)
	r.CurrentStatus = 0
	syscall.Kill(g.ReadPid("./redisPid"), syscall.SIGKILL)
	//参数pid
	os.Remove("./redisPid")
	fmt.Println("redis stop ok", r.CurrentStatus)
	return nil
}

func (r *RedisService) Restart(args map[string]string) error {
	if err := r.Stop(args); err != nil {
		return err
	}
	if err := r.Start(args); err != nil {
		return err
	}
	return nil
}

func (r *RedisService) Action(action string, args map[string]string) {
	fmt.Printf("%p\n", r)
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

	pkt := src.NewPkt()
	pkt.Id = g.ServiceResponse
	pkt.Data = []byte(str)
	a := utils.GlobalConfig.Get(g.AGENT).(iface.IAgent)
	err := a.GetCon().Write(pkt)
	if err != nil {
		//todo 发送失败..应该有后续操作
	}
}

func (r *RedisService) info() {
	info := model.NewServiceResponse(g.REDISSERVICE, r.CurrentStatus)
	if r.Status(nil) {
		info.Info, _ = g.EncodeData("redis启动中")
	} else {
		info.Status = 0
		info.Info, _ = g.EncodeData("停止服务")
	}
	pkt := src.ServiceResponsePkt(info)
	//todo 收集redis的信息
	a := utils.GlobalConfig.Get(g.AGENT).(iface.IAgent)
	if err := a.GetCon().Write(pkt); err != nil {
		log.Println(err)
	}
}
func (r *RedisService) Upload(args map[string]string) {
	if r.timerId != 0 {
		src.CancelTimer(r.timerId)
	}
	interval := g.GetInterval(args, 60*5)
	r.timerId = src.AddTimer(interval*time.Second, r.info)
}
func (r *RedisService) Watcher() {
	run := r.Status(nil)
	if run == true && r.CurrentStatus == 0 {
		r.CurrentStatus = 1
	} else if r.CurrentStatus == 1 && run == false {
		r.Start(map[string]string{})
	}
}
func (r *RedisService) Cancel() {
	src.CancelTimer(r.timerId)
}
