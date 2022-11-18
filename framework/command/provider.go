package command

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/YunzeGao/fire/framework"
	"github.com/YunzeGao/fire/framework/cobra"
	"github.com/YunzeGao/fire/framework/contract"
	"github.com/YunzeGao/fire/framework/util"

	"github.com/AlecAivazis/survey/v2"
	"github.com/jianfengye/collection"
	"github.com/pkg/errors"
)

// 初始化provider相关服务
func initProviderCommand() *cobra.Command {
	providerCommand.AddCommand(providerCreateCommand)
	providerCommand.AddCommand(providerListCommand)
	return providerCommand
}

var providerCommand = &cobra.Command{
	Use:   "provider",
	Short: "服务提供相关命令",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			_ = cmd.Help()
		}
		return nil
	},
}

// providerListCommand 列出容器内的所有服务
var providerListCommand = &cobra.Command{
	Use:     "list",
	Aliases: []string{"l"},
	Short:   "展示容器内的所有服务",
	RunE: func(cmd *cobra.Command, args []string) error {
		container := cmd.GetContainer().(*framework.FireContainer)
		// 打印字符串凭证
		for _, line := range container.NameList() {
			println(line)
		}
		return nil
	},
}

// providerCreateCommand 创建一个新的服务，包括服务提供者，服务接口协议，服务实例
var providerCreateCommand = &cobra.Command{
	Use:     "new",
	Aliases: []string{"n", "c", "create", "init"},
	Short:   "创建一个服务",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("新建一个服务")
		// 指定服务名称
		var name string
		{
			prompt := &survey.Input{
				Message: "请输入服务器名称(服务凭证):",
			}
			if err := survey.AskOne(prompt, &name); err != nil {
				return nil
			}
		}
		container := cmd.GetContainer()
		providers := container.(*framework.FireContainer).NameList()
		providerColl := collection.NewStrCollection(providers)
		if providerColl.Contains(name) {
			fmt.Println("服务名称已经存在！")
			return nil
		}

		// 指定文件路径
		var folder string
		{
			prompt := &survey.Input{
				Message: "请输入服务器所在目录名称(默认同服务名称):",
			}
			if err := survey.AskOne(prompt, &folder); err != nil {
				return nil
			}
		}
		if folder == "" {
			folder = name
		}

		app := container.MustMake(contract.AppKey).(contract.App)
		preFolder := app.ProviderFolder()
		subFolders, err := util.SubDir(preFolder)
		if err != nil {
			return err
		}
		subColl := collection.NewStrCollection(subFolders)
		if subColl.Contains(folder) {
			fmt.Println("目录名称已经存在")
			return nil
		}

		// 开始创建文件
		if err := os.Mkdir(filepath.Join(preFolder, folder), 0700); err != nil {
			return err
		}
		// 创建title这个模版方法
		funcs := template.FuncMap{"title": strings.Title}
		{
			//  创建contract.go
			file := filepath.Join(preFolder, folder, "contract.go")
			f, err := os.Create(file)
			if err != nil {
				return errors.Cause(err)
			}
			// 使用contractTmp模版来初始化template，并且让这个模版支持title方法，即支持{{.|title}}
			t := template.Must(template.New("contract").Funcs(funcs).Parse(contractTmp))
			// 将name传递进入到template中渲染，并且输出到contract.go 中
			if err := t.Execute(f, name); err != nil {
				return errors.Cause(err)
			}
		}
		{
			// 创建provider.go
			file := filepath.Join(preFolder, folder, "provider.go")
			f, err := os.Create(file)
			if err != nil {
				return err
			}
			t := template.Must(template.New("provider").Funcs(funcs).Parse(providerTmp))
			if err := t.Execute(f, name); err != nil {
				return err
			}
		}
		{
			//  创建service.go
			file := filepath.Join(preFolder, folder, "service.go")
			f, err := os.Create(file)
			if err != nil {
				return err
			}
			t := template.Must(template.New("service").Funcs(funcs).Parse(serviceTmp))
			if err := t.Execute(f, name); err != nil {
				return err
			}
		}
		fmt.Println("创建服务成功, 文件夹地址:", filepath.Join(preFolder, folder))
		fmt.Println("请不要忘记挂载新创建的服务")
		return nil
	},
}

var contractTmp = `package {{.}}
const {{.|title}}Key = "{{.}}"
type IService interface {
	// Foo 请在这里定义你的方法
	Foo() string
}
`

var providerTmp string = `package {{.}}
import (
	"github.com/YunzeGao/fire/framework"
)
type {{.|title}}Provider struct {
	framework.IServiceProvider
	container framework.IContainer
}
func (sp *{{.|title}}Provider) Name() string {
	return {{.|title}}Key
}
func (sp *{{.|title}}Provider) Register(c framework.IContainer) framework.NewInstance {
	return New{{.|title}}Service
}
func (sp *{{.|title}}Provider) IsDefer() bool {
	return false
}
func (sp *{{.|title}}Provider) Params(c framework.IContainer) []interface{} {
	return []interface{}{c}
}
func (sp *{{.|title}}Provider) Boot(c framework.IContainer) error {
	return nil
}
`

var serviceTmp string = `package {{.}}
import (
	"github.com/YunzeGao/fire/framework"
)
type {{.|title}}Service struct {
	container framework.IContainer
}
func New{{.|title}}Service(params ...interface{}) (interface{}, error) {
	container := params[0].(framework.IContainer)
	return &{{.|title}}Service{container: container}, nil
}
func (s *{{.|title}}Service) Foo() string {
    return ""
}
`
