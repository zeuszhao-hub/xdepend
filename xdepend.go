package xdepend

import (
	"context"
	"reflect"
)
import "golang.org/x/sync/errgroup"
import "encoding/json"

type fun func(ctx context.Context, retCollector *values) (interface{}, error)

// base 包转fn的执行状态和返回值，rets是依赖对象的返回值列表
type base struct {
	ret    *ret
	fn     fun
	status chan int
}

// ret service结果存储
type ret struct {
	ret interface{}
	err error
}

// values service结果收集器
type values []*ret

// Values 设置service结果到target
func (r *values) Values(target ...interface{}) error {
	for _, ret := range *r {
		errRet := ret.Error()
		if errRet != nil {
			return errRet
		}

		err := ret.value(target...)
		if err != nil {
			return err
		}
	}
	return nil
}

// value 获取依赖对象的返回值，对`target`赋值
func (r *ret) value(target ...interface{}) error {
	t := reflect.TypeOf(r.ret)
	for i, _ := range target {
		tt := reflect.TypeOf(target[i])
		if "*"+t.String() == tt.String() {
			source, err := json.Marshal(r.ret)
			if err != nil {
				return err
			}
			err = json.Unmarshal(source, target[i])
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (r *ret) Error() error {
	return r.err
}

// NewService 创建一个`base`对象
func NewService() *base {
	return &base{}
}

// Handle 增加执行函数上下文
func (b *base) Handle(fn fun) *base {
	b.fn = fn
	b.ret = new(ret)
	b.status = make(chan int)
	return b
}

//Value 获取当前`base`的返回值
func (b *base) Value(target ...interface{}) error {
	t := reflect.TypeOf(b.ret.ret)
	for i, v := range target {
		tt := reflect.TypeOf(v)
		if "*"+t.String() == tt.String() {
			source, err := json.Marshal(b.ret.ret)
			if err != nil {
				return err
			}
			err = json.Unmarshal(source, target[i])
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (b *base) Error() error {
	return b.ret.err
}

// do 并发执行所有注册的base.fn，根据依赖关系，源对象会一直堵塞直到被依赖对象完成执行（被依赖对象的status channel被关闭），最终同步返回。
func do(ctxp context.Context, depdesc ...describe) error {
	g, ctx := errgroup.WithContext(ctxp)
	for i, _ := range depdesc {
		ii := i
		g.Go(func() error {
			return func(ii int) error {
				if len(depdesc[ii].dep) != 0 {
					func() {
						values := make(values, 0)
						for _, dep := range depdesc[ii].dep {
							if dep == nil {
								continue
							}
							<-dep.status
							values = append(values, dep.ret)
						}
						depdesc[ii].me.ret.ret, depdesc[ii].me.ret.err = depdesc[ii].me.fn(ctx, &values)
					}()
				} else {
					retfnNil := make(values, 0)
					depdesc[ii].me.ret.ret, depdesc[ii].me.ret.err = depdesc[ii].me.fn(ctx, &retfnNil)
				}
				close(depdesc[ii].me.status)
				if depdesc[ii].me.ret.err != nil {
					return depdesc[ii].me.ret.err
				}
				return nil
			}(ii)
		})
	}
	return g.Wait()
}

// describe 包装依赖
type describe struct {
	dep []*base
	me  *base
}

// describe 依赖列表
type depend struct {
	dep []describe
}

// NewDepend 创建一个依赖描述对象
func NewDepend() *depend {
	return &depend{
		dep: make([]describe, 0),
	}
}

// AddDescribe 增加依赖表述，b是源对象，target是被依赖对象的列表
// 注意：当A依赖于B和C时，使用.AddDescribe(A,B,C)而不是.AddDescribe(A,B).AddDescribe(A,C)
func (d *depend) AddDescribe(b *base, target ...*base) *depend {
	desc := describe{
		me:  b,
		dep: target,
	}
	d.dep = append(d.dep, desc)
	return d
}

// Do 参考do
func (d *depend) Do(ctxp context.Context) error {
	return do(ctxp, d.dep...)
}
