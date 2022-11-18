package demo

import (
	"fmt"

	"github.com/YunzeGao/fire/framework/cobra"
)

// InitFoo 初始化Foo命令
func InitFoo() *cobra.Command {
	ContainerCommand.AddCommand(ProviderCommand)
	return ContainerCommand
}

var ContainerCommand = &cobra.Command{
	Use:     "container",
	Short:   "当前容器",
	Long:    "",
	Aliases: []string{"nc"},
	Example: "",
	RunE: func(c *cobra.Command, args []string) error {
		container := c.GetContainer()
		fmt.Println("当前容器 = ", container)
		return nil
	},
}

var ProviderCommand = &cobra.Command{
	Use:     "provider",
	Short:   "当前容器所有的providers",
	Long:    "",
	Aliases: []string{"p"},
	Example: "",
	RunE: func(c *cobra.Command, args []string) error {
		container := c.GetContainer()
		fmt.Println("当前容器 = ", container)
		fmt.Println("当期容器的所有服务 = ", container.PrintProviders())
		return nil
	},
}
