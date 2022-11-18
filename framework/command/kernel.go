package command

import "github.com/YunzeGao/fire/framework/cobra"

func AddKernelCommands(root *cobra.Command) {
	root.AddCommand(DemoCommand)
	// app
	root.AddCommand(initAppCommand())
	// 环境
	root.AddCommand(initEnvCommand())
	// provider
	root.AddCommand(initProviderCommand())
	// cmd
	root.AddCommand(initCmdCommand())
	// middleware
	root.AddCommand(initMiddlewareCommand())
}
