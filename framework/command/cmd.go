package command

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/YunzeGao/fire/framework/cobra"
	"github.com/YunzeGao/fire/framework/contract"
	"github.com/YunzeGao/fire/framework/util"

	"github.com/jianfengye/collection"

	"github.com/AlecAivazis/survey/v2"
	"github.com/pkg/errors"
)

// 初始化command相关命令
func initCmdCommand() *cobra.Command {
	cmdCommand.AddCommand(cmdListCommand)
	cmdCommand.AddCommand(cmdCreateCommand)
	return cmdCommand
}

var cmdCommand = &cobra.Command{
	Use:     "command",
	Short:   "控制台命令相关",
	Aliases: []string{"cmd"},
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			_ = cmd.Help()
		}
		return nil
	},
}

// cmdListCommand 列出容器内的所有服务
var cmdListCommand = &cobra.Command{
	Use:     "list",
	Aliases: []string{"l"},
	Short:   "展示容器内的所有服务",
	RunE: func(cmd *cobra.Command, args []string) error {
		cmds := cmd.Root().Commands()
		var ps [][]string
		for _, cmd := range cmds {
			line := []string{cmd.Name(), cmd.Short}
			ps = append(ps, line)
		}
		util.PrettyPrint(ps)
		return nil
	},
}

// cmdCreateCommand 创建一个新的服务，包括服务提供者，服务接口协议，服务实例
var cmdCreateCommand = &cobra.Command{
	Use:     "new",
	Aliases: []string{"n", "c", "create", "init"},
	Short:   "创建一个服务",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("开始创建控制台命令...")
		var name string
		{
			prompt := &survey.Input{
				Message: "请输入控制台命令名称:",
			}
			err := survey.AskOne(prompt, &name)
			if err != nil {
				return err
			}
		}
		var folder string
		{
			prompt := &survey.Input{
				Message: "请输入文件夹名称(默认: 同控制台命令):",
			}
			err := survey.AskOne(prompt, &folder)
			if err != nil {
				return err
			}
		}

		if folder == "" {
			folder = name
		}

		container := cmd.GetContainer()

		// 判断文件不存在
		app := container.MustMake(contract.AppKey).(contract.App)
		preFolder := app.CommandFolder()
		subFolders, err := util.SubDir(preFolder)
		if err != nil {
			return err
		}
		subColl := collection.NewStrCollection(subFolders)
		if !subColl.Contains(folder) {
			// 开始创建文件
			if err := os.Mkdir(filepath.Join(preFolder, folder), 0700); err != nil {
				return err
			}
			return nil
		}

		type Data struct {
			Name   string
			Folder string
		}

		// 创建title这个模版方法
		funcs := template.FuncMap{"title": strings.Title}
		{
			//  创建name.go
			file := filepath.Join(preFolder, folder, name+".go")
			f, err := os.Create(file)
			if err != nil {
				return errors.Cause(err)
			}

			// 使用contractTmp模版来初始化template，并且让这个模版支持title方法，即支持{{.|title}}
			t := template.Must(template.New("cmd").Funcs(funcs).Parse(cmdTemplate))
			// 将name传递进入到template中渲染，并且输出到contract.go 中
			if err := t.Execute(f, map[string]string{
				"Name":   name,
				"Folder": folder,
			}); err != nil {
				return errors.Cause(err)
			}
		}

		fmt.Println("创建新命令行工具成功，路径:", filepath.Join(preFolder, folder))
		fmt.Println("请记得开发完成后将命令行工具挂载到 console/kernel.go")
		return nil
	},
}

var cmdTemplate string = `package {{.Folder}}
import (
	"github.com/YunzeGao/fire/framework/cobra"
	"fmt"
)
var {{.Name|title}}Command = &cobra.Command{
	Use:   "{{.Name}}",
	Short: "{{.Name}}",
	RunE: func(cmd *cobra.Command, args []string) error {
        container := cmd.GetContainer()
		fmt.Println(container)
		return nil
	},
}
`
