package intfns

import (
	"fmt"
	"math"

	"github.com/bserdar/goxpath/tree"
)

func number(a tree.Adapter, c tree.Ctx, args ...tree.Result) (tree.Result, error) {
	if b, ok := args[0].(tree.IsNum); ok {
		return b.Num(), nil
	}

	return nil, fmt.Errorf("Cannot convert object to a number")
}

func sum(a tree.Adapter, c tree.Ctx, args ...tree.Result) (tree.Result, error) {
	n, ok := args[0].(tree.NodeSet)
	if !ok {
		return nil, fmt.Errorf("Cannot convert object to a node-set")
	}

	ret := 0.0
	for _, i := range n.GetNodes() {
		ret += float64(tree.String(a.StringValue(i)).Num())
	}

	return tree.Num(ret), nil
}

func floor(a tree.Adapter, c tree.Ctx, args ...tree.Result) (tree.Result, error) {
	n, ok := args[0].(tree.IsNum)
	if !ok {
		return nil, fmt.Errorf("Cannot convert object to a number")
	}

	return tree.Num(math.Floor(float64(n.Num()))), nil
}

func ceiling(a tree.Adapter, c tree.Ctx, args ...tree.Result) (tree.Result, error) {
	n, ok := args[0].(tree.IsNum)
	if !ok {
		return nil, fmt.Errorf("Cannot convert object to a number")
	}

	return tree.Num(math.Ceil(float64(n.Num()))), nil
}

func round(a tree.Adapter, c tree.Ctx, args ...tree.Result) (tree.Result, error) {
	isn, ok := args[0].(tree.IsNum)
	if !ok {
		return nil, fmt.Errorf("Cannot convert object to a number")
	}

	n := isn.Num()

	if math.IsNaN(float64(n)) || math.IsInf(float64(n), 0) {
		return n, nil
	}

	if n < -0.5 {
		n = tree.Num(int(n - 0.5))
	} else if n > 0.5 {
		n = tree.Num(int(n + 0.5))
	} else {
		n = 0
	}

	return n, nil
}
