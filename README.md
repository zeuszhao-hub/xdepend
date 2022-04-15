# xdepend
解决多任务并发过程中，任务间相互数据依赖问题带来的代码复杂度问题。

## 如何使用
[例子](https://github.com/zeuszhao-hub/xdepend/blob/main/xdepend_test.go)

[timeline]:data:image/png;base64,data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAjEAAAEPCAYAAACz95W0AAAAAXNSR0IArs4c6QAABu50RVh0bXhmaWxlACUzQ214ZmlsZSUyMGhvc3QlM0QlMjJDaHJvbWUlMjIlMjBtb2RpZmllZCUzRCUyMjIwMjItMDQtMTVUMDMlM0E1NiUzQTMwLjY1MFolMjIlMjBhZ2VudCUzRCUyMjUuMCUyMChNYWNpbnRvc2glM0IlMjBJbnRlbCUyME1hYyUyME9TJTIwWCUyMDEwXzE1XzcpJTIwQXBwbGVXZWJLaXQlMkY1MzcuMzYlMjAoS0hUTUwlMkMlMjBsaWtlJTIwR2Vja28pJTIwQ2hyb21lJTJGOTkuMC40ODQ0LjgzJTIwU2FmYXJpJTJGNTM3LjM2JTIyJTIwdmVyc2lvbiUzRCUyMjE2LjIuMiUyMiUyMGV0YWclM0QlMjJJOWUwQkRDQTViVTBsZWJ2SVhkMyUyMiUyMHR5cGUlM0QlMjJkZXZpY2UlMjIlM0UlM0NkaWFncmFtJTIwaWQlM0QlMjIwdkxJN2RUVmpoTlRJNzlHT2g2YyUyMiUyMG5hbWUlM0QlMjIlRTclQUMlQUMlMjAxJTIwJUU5JUExJUI1JTIyJTNFN1ZsZGI1c3dGUDAxZnV6RWx3MDhob1N1bGRwdFdpZXQycHNESmtFak9BS25TZnJyWjhDUWdNMlN0b0dYVktwU09MYXY0Wng3cjMwTk1LZXIzZGNNcjVlUE5DUUpNTFJ3Qjh3Wk1BeFhzJTJGaHZBZXdyQU5sR0JTeXlPS3dnJTJGUUE4eGE5RWdKcEFOM0ZJOGxaSFJtbkM0blViREdpYWtvQzFNSnhsZE52dUZ0R2tQZXNhTDRnRVBBVTRrZEhmY2NpV0Zlb1k5Z0clMkZJJTJGRmlXYyUyQnNJN2RxV2VHNnMzaVRmSWxEdWoyQ1RCJTJCWTA0eFNWbDJ0ZGxPU0ZOelZ2RlRqYm50YW13ZkxTTXJPR1hBVDM2R2Y3bks2JTJCWFA3JTJGQjFseHQ0SzBVMXQ1Z1VuRyUyRkhHNG1uWnZxWWdvNXMwSklVVkRaamVkaGt6OHJUR1FkRzY1WnB6Yk1sV0NiJTJGVCUyQldWRVV6YWxDYzNLc2FhdHUzRHFjanhuR2YxTDZwYVVwbnk0SiUyQlltR1NPNzNyZlNHNjY0anhHNklpemI4eTVpUU9Nb3dyOTBDS3Y3N1VFdGlFU2Y1WkZTaGkxQUxEeGswZGclMkJrTWd2Qkk5djRQVGlsTVpKY2tScDVBUWtDQ1JLZWN2Y2dSYlVMc1FyYXZOcWFKckVxJTJCVW9lTldkb1hnMWh1V1Y2Q0VrdG9wWEY5a21SdVB4YXFqOEZRMUZxem13dTBhUm9YYlhFTTBSSElwV0pOTnFXbVBTQ29lbE5jVEVpWlMwb3NBaDgyZ1lXazFEUWFzcUN3eEdxelcwdHhMVTQ2MjJPOWVHU3E2dUlnbllZOUtLVHROSzBuQlNiS2o0WFpEZ1BJJTJCRE5wTnQyaVglMkJtcDFBWlppRTByN3JKSHZIUzdxQ25CckxTSUpaJTJGTkkycjJKTXpQQ0R4bnppZ3pobXglMkJlN0MxcE9OMWxBeEtqakRkY3BRMmJIRU1QWmdqREpVQ2xnODlydjE5VCUyQjFMUlBpbWJWJTJGYWltelFPT3BLa2phJTJCcGJ3UE9CT3dPJTJCQXp3SVhBUjhDQndJUEJQNGR2SHJlWkx1UEhXeHR0TEslMkZmbHhjaFFRVHVKRldyZ0xGNVp3M0NzU1ljenJwWWxvV01WaFdFeWpUTHh0ZiUyQnFyRlM2UllRM1kxc20ycEF5ckslMkZ5dEc2S1hLd3EwejJqc1clMkY2czkwYWpiWjB3TkhBMDZxcFNEJTJCRlY0ZXJwUEYlMkJYV21obFBNNkFDOHRROVlxJTJGQXBtQWlRNThCRndlcER4bTNmTENLbU4yQXJ4YjJWRDN2dXM5WE84SFBDZEoyMm5PRDlpTTVQRXJucGYyQ2dkYkY3U1ZSRUlQd05sN0FsYWN5Z2lqb0RrTE9YYkYlMkY0UkxiM2hyWHpRYmZzd1g2eTQwaW5JeWpIZklCZXZrJTJCakl4Y25xVzNxTmNvNnAzaDh2RWNzRjdoUXVrZGNZNXhMaXl5SlhkOVBwa2daMFZUWFdPTWE0czhqbkc3UHBrc2V5T0xJcUNmVnhaNUlMZHZ6NVpZR2R0VVIxUGpTcUxJZSUyRnlIeWZQSFBoMSUyRjhqMTBSN3V2eFglMkZ5aGFOYndaNDRXWjQ1OFRUSlhkeGVNTm9YbjFvZTR0UERCbGVSbHRIM1paMTFKVWZjZDR1Skw4OWZIT3J0bnFIRDVlbSUyRnc4JTNEJTND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

![描述][timeline]

## 伪代码
```text
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
