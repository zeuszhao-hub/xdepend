package xdepend

import (
	"context"
	"errors"
	"testing"
	"time"
)

type A struct {
	Name string
	AB   B
	AC   C
}

type B struct {
	Name string
}

type C struct {
	Name string
}

func getAName(ctx context.Context, dep string, b bool) (string, error) {
	if !b {
		select {
		case <-ctx.Done():
			return "", errors.New("a被取消")
		case <-time.After(time.Second * 1):
		}
	}
	return "a" + dep, nil
}

func getBName(ctx context.Context, dep string, b bool) (string, error) {
	if !b {
		select {
		case <-ctx.Done():
			return "", errors.New("b被取消")
		case <-time.After(time.Second * 2):
		}
	}
	return "b" + dep, nil
}

func getCName(ctx context.Context, dep string, b bool) (string, error) {
	if !b {
		select {
		case <-ctx.Done():
			return "", errors.New("c被取消")
		case <-time.After(time.Second * 3):
		}
	}
	return "c" + dep, nil
}

func TestNewDepend(t *testing.T) {
	a := NewService().Handle(func(ctx context.Context, retCollector *RetCollector) (interface{}, error) {
		ab := B{}
		ac := C{}
		err := retCollector.Values(&ab, &ac)
		if err != nil {
			return nil, err
		}

		// 获取a的名字
		aname, err := getAName(ctx, ab.Name+ac.Name, false)
		if err != nil {
			return nil, err
		}
		return A{
			Name: aname,
			AB:   ab,
			AC:   ac,
		}, nil
	})

	b := NewService().Handle(func(ctx context.Context, retCollector *RetCollector) (interface{}, error) {
		// 获取b的名字
		bname, err := getBName(ctx, "", false)
		if err != nil {
			return nil, err
		}
		return B{Name: bname}, nil
	})

	c := NewService().Handle(func(ctx context.Context, retCollector *RetCollector) (interface{}, error) {
		// 获取c的名字
		cname, err := getCName(ctx, "", false)
		if err != nil {
			return nil, err
		}
		return C{Name: cname}, nil
	})

	ctx, _ := context.WithTimeout(context.TODO(), 5*time.Second)
	err := NewDepend().AddDescribe(a, b, c).AddDescribe(b).AddDescribe(c).Do(ctx)
	if err != nil {
		t.Fatal(err)
	}
	aa := A{}
	a.Value(&aa)
	if aa.Name != "abc" {
		t.Fatal()
	}
}

func BenchmarkNewDepend(bb *testing.B) {
	a := NewService().Handle(func(ctx context.Context, retCollector *RetCollector) (interface{}, error) {
		ab := B{}
		ac := C{}
		err := retCollector.Values(&ab, &ac)
		if err != nil {
			return nil, err
		}

		// 获取a的名字
		aname, _ := getAName(ctx, ab.Name+ac.Name, true)
		return A{
			Name: aname,
			AB:   ab,
			AC:   ac,
		}, nil
	})

	b := NewService().Handle(func(ctx context.Context, retCollector *RetCollector) (interface{}, error) {
		// 获取b的名字
		bname, _ := getBName(ctx, "", true)
		return B{Name: bname}, nil
	})

	c := NewService().Handle(func(ctx context.Context, retCollector *RetCollector) (interface{}, error) {
		// 获取c的名字
		cname, _ := getCName(ctx, "", true)
		return C{Name: cname}, nil
	})

	_ = NewDepend().AddDescribe(a, b, c).AddDescribe(b).AddDescribe(c).Do(context.TODO())
	aa := A{}
	a.Value(&aa)
	if aa.Name != "abc" {
		bb.Fatal()
	}
}
