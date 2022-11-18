package command

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"syscall"
	"time"

	"github.com/YunzeGao/fire/framework"
	"github.com/YunzeGao/fire/framework/cobra"
	"github.com/YunzeGao/fire/framework/contract"
	"github.com/YunzeGao/fire/framework/util"

	"github.com/erikdubbelboer/gspt"

	"github.com/sevlyar/go-daemon"
)

// app启动地址
var appAddress = ""
var appDaemon = false

// initAppCommand 初始化app命令和其子命令
func initAppCommand() *cobra.Command {
	appStartCommand.Flags().BoolVarP(&appDaemon, "daemon", "d", false, "start app daemon")
	appStartCommand.Flags().StringVarP(&appAddress, "addr", "a", "", "设置APP启动地址，默认为8080")
	appCommand.AddCommand(appStartCommand)
	appCommand.AddCommand(appStateCommand)
	appCommand.AddCommand(appStopCommand)
	appCommand.AddCommand(appRestartCommand)
	return appCommand
}

// AppCommand 是命令行参数第一级为app的命令，它没有实际功能，只是打印帮助文档
var appCommand = &cobra.Command{
	Use:   "app",
	Short: "业务应用控制命令",
	Long:  "业务应用控制命令，其包含业务启动，关闭，重启，查询等功能",
	RunE: func(c *cobra.Command, args []string) error {
		// 打印帮助文档
		_ = c.Help()
		return nil
	},
}

func startAppServe(server *http.Server, container framework.IContainer) error {
	// 这个goroutine是启动服务的goroutine
	go func() {
		_ = server.ListenAndServe()
	}()

	// 当前的goroutine等待信号量
	quit := make(chan os.Signal)
	// 监控信号：SIGINT, SIGTERM, SIGQUIT
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	// 这里会阻塞当前goroutine等待信号
	<-quit

	// 调用Server.Shutdown graceful结束
	timeoutCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(timeoutCtx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	return nil
}

// appStartCommand 启动一个Web服务
var appStartCommand = &cobra.Command{
	Use:   "start",
	Short: "启动一个Web服务",
	RunE: func(c *cobra.Command, args []string) error {
		// 从Command中获取服务容器
		container := c.GetContainer()
		// 从服务容器中获取kernel的服务实例
		kernelService := container.MustMake(contract.KernelKey).(contract.IKernel)
		// 从kernel服务实例中获取引擎
		core := kernelService.HttpEngine()

		if appAddress == "" {
			envService := container.MustMake(contract.EnvKey).(contract.Env)
			appAddress = envService.Get("ADDR")
			if appAddress == "" {
				configService := container.MustMake(contract.ConfigKey).(contract.IConfig)
				appAddress = configService.GetString("app.address")
			}
			if appAddress == "" {
				appAddress = ":8080"
			}
		}
		// 创建一个Server服务
		server := &http.Server{
			Handler: core,
			Addr:    appAddress,
		}

		appService := container.MustMake(contract.AppKey).(contract.App)
		pidFolder := appService.RuntimeFolder()
		if !util.Exists(pidFolder) {
			if err := os.Mkdir(pidFolder, os.ModePerm); err != nil {
				return err
			}
		}
		logFolder := appService.LogFolder()
		if !util.Exists(logFolder) {
			if err := os.Mkdir(logFolder, os.ModePerm); err != nil {
				return err
			}
		}
		serverPidFile := filepath.Join(pidFolder, "app.pid")
		serverLogFile := filepath.Join(logFolder, "app.log")
		currentFolder := util.GetExecDirectory()
		if appDaemon {
			daemonCtx := &daemon.Context{
				PidFileName: serverPidFile,
				PidFilePerm: 0664,
				LogFileName: serverLogFile,
				LogFilePerm: 0640,
				WorkDir:     currentFolder,
				// 设置所有设置文件的mask，默认为750
				Umask: 027,
				// 子进程的参数，按照这个参数设置，子进程的命令为 ./fire app start --daemon=true
				Args: []string{"", "app", "start", "--daemon=true"},
			}
			d, err := daemonCtx.Reborn()
			if err != nil {
				return err
			}
			if d != nil {
				// 父进程直接打印启动成功信息，不做任何操作
				fmt.Println("app启动成功，pid:", d.Pid)
				fmt.Println("日志文件:", serverLogFile)
				return nil
			}
			defer func(daemonCtx *daemon.Context) {
				_ = daemonCtx.Release()
			}(daemonCtx)
			fmt.Println(time.Now(), "daemon start app")
			gspt.SetProcTitle("fire app")
			if err := startAppServe(server, container); err != nil {
				fmt.Println(err)
			}
			return nil
		}
		// 非daemon模式，直接执行
		content := strconv.Itoa(os.Getpid())
		fmt.Println("[PID]", content)
		err := os.WriteFile(serverPidFile, []byte(content), 0644)
		if err != nil {
			return err
		}
		gspt.SetProcTitle("fire app")

		fmt.Println("app serve url:", appAddress)
		if err := startAppServe(server, container); err != nil {
			fmt.Println(err)
		}
		return nil
	},
}

// 获取启动的app的pid
var appStateCommand = &cobra.Command{
	Use:   "state",
	Short: "获取启动的app的pid",
	RunE: func(cmd *cobra.Command, args []string) error {
		container := cmd.GetContainer()
		appService := container.MustMake(contract.AppKey).(contract.App)

		// 获取pid
		serverPidFile := filepath.Join(appService.RuntimeFolder(), "app.pid")

		content, err := os.ReadFile(serverPidFile)
		if err != nil {
			return err
		}

		if content != nil && len(content) > 0 {
			pid, err := strconv.Atoi(string(content))
			if err != nil {
				return err
			}
			if util.CheckProcessExist(pid) {
				fmt.Println("app服务已经启动, pid:", pid)
				return nil
			}
		}
		fmt.Println("没有app服务存在")
		return nil
	},
}

// 停止一个已经启动的app服务
var appStopCommand = &cobra.Command{
	Use:   "stop",
	Short: "停止一个已经启动的app服务",
	RunE: func(cmd *cobra.Command, args []string) error {
		container := cmd.GetContainer()
		appService := container.MustMake(contract.AppKey).(contract.App)

		// GetPid
		serverPidFile := filepath.Join(appService.RuntimeFolder(), "app.pid")

		content, err := os.ReadFile(serverPidFile)
		if err != nil {
			return err
		}

		if content != nil && len(content) != 0 {
			pid, err := strconv.Atoi(string(content))
			if err != nil {
				return err
			}
			// 发送SIGTERM命令
			if err := syscall.Kill(pid, syscall.SIGTERM); err != nil {
				return err
			}
			if err := os.WriteFile(serverPidFile, []byte{}, 0644); err != nil {
				return err
			}
			fmt.Println(time.Now(), "停止进程:", pid)
		}
		return nil
	},
}

// 重新启动一个app服务
var appRestartCommand = &cobra.Command{
	Use:   "restart",
	Short: "重新启动一个app服务",
	RunE: func(cmd *cobra.Command, args []string) error {
		container := cmd.GetContainer()
		appService := container.MustMake(contract.AppKey).(contract.App)

		// GetPid
		serverPidFile := filepath.Join(appService.RuntimeFolder(), "app.pid")

		content, err := os.ReadFile(serverPidFile)
		if err != nil {
			return err
		}

		if content != nil && len(content) != 0 {
			pid, err := strconv.Atoi(string(content))
			if err != nil {
				return err
			}
			if util.CheckProcessExist(pid) {
				// 杀死进程
				if err := syscall.Kill(pid, syscall.SIGTERM); err != nil {
					return err
				}

				// 获取closeWait
				closeWait := 5
				configService := container.MustMake(contract.ConfigKey).(contract.IConfig)
				if configService.IsExist("app.close_wait") {
					closeWait = configService.GetInt("app.close_wait")
				}

				// 确认进程已经关闭,每秒检测一次， 最多检测closeWait * 2秒
				for i := 0; i < closeWait*2; i++ {
					if util.CheckProcessExist(pid) == false {
						break
					}
					time.Sleep(1 * time.Second)
				}

				// 如果进程等待了2*closeWait之后还没结束，返回错误，不进行后续的操作
				if util.CheckProcessExist(pid) == true {
					fmt.Println(time.Now(), "结束进程失败:"+strconv.Itoa(pid), "请查看原因")
					return errors.New("结束进程失败")
				}
				if err := os.WriteFile(serverPidFile, []byte{}, 0644); err != nil {
					return err
				}

				fmt.Println(time.Now(), "结束进程成功:"+strconv.Itoa(pid))
			}
		}

		appDaemon = true
		// 直接daemon方式启动apps
		return appStartCommand.RunE(cmd, args)
	},
}
