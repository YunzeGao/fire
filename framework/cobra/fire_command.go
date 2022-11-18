package cobra

import "github.com/YunzeGao/fire/framework"

func (c *Command) SetContainer(container framework.IContainer) {
	c.container = container
}

func (c *Command) GetContainer() framework.IContainer {
	return c.Root().container
}
