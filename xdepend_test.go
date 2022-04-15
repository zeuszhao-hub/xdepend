package xdepend

import (
	"context"
	"errors"
	"testing"
	"time"
)

type A struct {
	Name string
	B    B
	C    C
	D    D
	E    E
}

type B struct {
	Name string
}

type C struct {
	Name string
	D    D
}

type D struct {
	Name string
}

type E struct {
	Name string
	D    D
}

func getAName(ctx context.Context, b bool) (string, error) {
	if !b {
		select {
		case <-ctx.Done():
			return "", errors.New("a被取消")
		case <-time.After(time.Millisecond * 500):
		}
	}
	return "a", nil
}

func getBName(ctx context.Context, b bool) (string, error) {
	if !b {
		select {
		case <-ctx.Done():
			return "", errors.New("b被取消")
		case <-time.After(time.Millisecond * 100):
		}
	}
	return "b", nil
}

func getCName(ctx context.Context, b bool) (string, error) {
	if !b {
		select {
		case <-ctx.Done():
			return "", errors.New("c被取消")
		case <-time.After(time.Millisecond * 300):
		}
	}
	return "c", nil
}

func getDName(ctx context.Context, b bool) (string, error) {
	if !b {
		select {
		case <-ctx.Done():
			return "", errors.New("d被取消")
		case <-time.After(time.Millisecond * 200):
		}
	}
	return "d", nil
}

func getEName(ctx context.Context, b bool) (string, error) {
	if !b {
		select {
		case <-ctx.Done():
			return "", errors.New("e被取消")
		case <-time.After(time.Millisecond * 400):
		}
	}
	return "e", nil
}

func TestNewDepend(t *testing.T) {
	a := NewService().Handle(func(ctx context.Context, retCollector *RetCollector) (interface{}, error) {
		ab := B{}
		ac := C{}
		ad := D{}
		ae := E{}
		err := retCollector.Values(&ab, &ac, &ad, &ae)
		if err != nil {
			return nil, err
		}

		// 获取a的名字
		aname, err := getAName(ctx, false)
		if err != nil {
			return nil, err
		}
		return A{
			Name: aname,
			B:    ab,
			C:    ac,
			D:    ad,
			E:    ae,
		}, nil
	})

	b := NewService().Handle(func(ctx context.Context, retCollector *RetCollector) (interface{}, error) {
		// 获取b的名字
		bname, err := getBName(ctx, false)
		if err != nil {
			return nil, err
		}
		return B{Name: bname}, nil
	})

	c := NewService().Handle(func(ctx context.Context, retCollector *RetCollector) (interface{}, error) {
		cd := D{}
		err := retCollector.Values(&cd)
		if err != nil {
			return nil, err
		}

		// 获取c的名字
		cname, err := getCName(ctx, false)
		if err != nil {
			return nil, err
		}
		return C{
			Name: cname,
			D:    cd,
		}, nil
	})

	d := NewService().Handle(func(ctx context.Context, retCollector *RetCollector) (interface{}, error) {
		// 获取b的名字
		dname, err := getDName(ctx, false)
		if err != nil {
			return nil, err
		}
		return D{Name: dname}, nil
	})

	e := NewService().Handle(func(ctx context.Context, retCollector *RetCollector) (interface{}, error) {
		ed := D{}
		err := retCollector.Values(&ed)
		if err != nil {
			return nil, err
		}

		// 获取c的名字
		ename, err := getEName(ctx, false)
		if err != nil {
			return nil, err
		}
		return E{
			Name: ename,
			D:    ed,
		}, nil
	})

	ctx, _ := context.WithTimeout(context.TODO(), 5*time.Second)
	err := NewDepend().AddDescribe(a, b, c, d, e).
		AddDescribe(b).
		AddDescribe(c, d).
		AddDescribe(d).
		AddDescribe(e, d).
		Do(ctx)
	if err != nil {
		t.Fatal(err)
	}
	aa := A{}
	a.Value(&aa)
	if aa.E.D.Name != "d" {
		t.Fatal()
	}
}

func BenchmarkNewDepend(bb *testing.B) {
	a := NewService().Handle(func(ctx context.Context, retCollector *RetCollector) (interface{}, error) {
		ab := B{}
		ac := C{}
		ad := D{}
		ae := E{}
		err := retCollector.Values(&ab, &ac, &ad, &ae)
		if err != nil {
			return nil, err
		}

		// 获取a的名字
		aname, err := getAName(ctx, true)
		if err != nil {
			return nil, err
		}
		return A{
			Name: aname,
			B:    ab,
			C:    ac,
			D:    ad,
			E:    ae,
		}, nil
	})

	b := NewService().Handle(func(ctx context.Context, retCollector *RetCollector) (interface{}, error) {
		// 获取b的名字
		bname, err := getBName(ctx, true)
		if err != nil {
			return nil, err
		}
		return B{Name: bname}, nil
	})

	c := NewService().Handle(func(ctx context.Context, retCollector *RetCollector) (interface{}, error) {
		cd := D{}
		err := retCollector.Values(&cd)
		if err != nil {
			return nil, err
		}

		// 获取c的名字
		cname, err := getCName(ctx, true)
		if err != nil {
			return nil, err
		}
		return C{
			Name: cname,
			D:    cd,
		}, nil
	})

	d := NewService().Handle(func(ctx context.Context, retCollector *RetCollector) (interface{}, error) {
		// 获取b的名字
		dname, err := getDName(ctx, true)
		if err != nil {
			return nil, err
		}
		return D{Name: dname}, nil
	})

	e := NewService().Handle(func(ctx context.Context, retCollector *RetCollector) (interface{}, error) {
		ed := D{}
		err := retCollector.Values(&ed)
		if err != nil {
			return nil, err
		}

		// 获取e的名字
		ename, err := getEName(ctx, true)
		if err != nil {
			return nil, err
		}
		return E{
			Name: ename,
			D:    ed,
		}, nil
	})

	err := NewDepend().AddDescribe(a, b, c, d, e).
		AddDescribe(b).
		AddDescribe(c, d).
		AddDescribe(d).
		AddDescribe(e, d).
		Do(context.TODO())
	if err != nil {
		bb.Fatal(err)
	}
	aa := A{}
	a.Value(&aa)
	if aa.E.D.Name != "d" {
		bb.Fatal()
	}
}
