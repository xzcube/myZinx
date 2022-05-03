package znet

import (
	"myZinx/ziface"
)

// 实现router时，先嵌入之歌BaseRouter基类，然后根据这个基类的方法进行重写
type BaseRouter struct {}

// 这里之所以BaseRouter的方法都为空，是因为有的router不希望有所有的方法业务
// 所以Router全部继承Base的好处就是，不需要实现所有的方法，只需要重写需要的方法
func (b *BaseRouter) PreHandle(request ziface.IRequest) {}

func (b *BaseRouter) Handle(request ziface.IRequest) {}

func (b *BaseRouter) PostHandle(request ziface.IRequest) {}
