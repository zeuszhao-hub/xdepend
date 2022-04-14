# xdepend
缓解并发编程过程中，由并发实体间数据相互依赖所带来的代码复杂度问题。

## 如何使用
参见 [xdepend_test.go](https://github.com/zeuszhao-hub/xdepend/blob/main/xdepend_test.go)


## 说明
```go
// 依赖关系
A->B
A->C

// 实体业务逻辑代码
a := NewService().Handle(func())
b := NewService().Handle(func())
c := NewService().Handle(func())

// 依赖描述，a依赖于b、c，b独立，c独立
err := NewDepend().AddDescribe(a, b, c).AddDescribe(b).AddDescribe(c).Do(context.TODO())


// 执行时间 = a+max(b,c)
```

## 如何安装

```shell
go get github.com/zeuszhao-hub/xdepend
```
