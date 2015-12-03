package result

import (
	"encoding/xml"

	"github.com/ChrisTrenkamp/goxpath/parser/result/element"
	"github.com/ChrisTrenkamp/goxpath/parser/result/namespace"
	"github.com/ChrisTrenkamp/goxpath/parser/result/pathexpr"
	"github.com/ChrisTrenkamp/goxpath/parser/result/pathres"
	"github.com/ChrisTrenkamp/goxpath/xconst"
)

type findFunc func(pathres.PathRes, *pathexpr.PathExpr, *[]pathres.PathRes)

var findMap = map[string]findFunc{
	xconst.AxisAncestor:         findAncestor,
	xconst.AxisAncestorOrSelf:   findAncestorOrSelf,
	xconst.AxisAttribute:        findAttribute,
	xconst.AxisChild:            findChild,
	xconst.AxisDescendent:       findDescendent,
	xconst.AxisDescendentOrSelf: findDescendentOrSelf,
	xconst.AxisFollowing:        findFollowing,
	xconst.AxisFollowingSibling: findFollowingSibling,
	xconst.AxisNamespace:        findNamespace,
	xconst.AxisParent:           findParent,
	xconst.AxisPreceding:        findPreceding,
	xconst.AxisPrecedingSibling: findPrecedingSibling,
	xconst.AxisSelf:             findSelf,
}

//Find finds nodes based on the pathexpr.PathExpr
func Find(x pathres.PathRes, p pathexpr.PathExpr) []pathres.PathRes {
	ret := []pathres.PathRes{}

	if p.Axis == "" {
		findChild(x, &p, &ret)
		return ret
	}

	f := findMap[p.Axis]
	f(x, &p, &ret)

	return ret
}

func findAncestor(x pathres.PathRes, p *pathexpr.PathExpr, ret *[]pathres.PathRes) {
	if x.GetParent().EvalPath(p) {
		*ret = append(*ret, x.GetParent())
	}
	if x.GetParent() != x {
		findAncestor(x.GetParent(), p, ret)
	}
}

func findAncestorOrSelf(x pathres.PathRes, p *pathexpr.PathExpr, ret *[]pathres.PathRes) {
	findSelf(x, p, ret)
	findAncestor(x, p, ret)
}

func findAttribute(x pathres.PathRes, p *pathexpr.PathExpr, ret *[]pathres.PathRes) {
	if ele, ok := x.(*element.PathResElement); ok {
		for i := range ele.Attrs {
			if ele.Attrs[i].EvalPath(p) {
				*ret = append(*ret, ele.Attrs[i])
			}
		}
	}
}

func findChild(x pathres.PathRes, p *pathexpr.PathExpr, ret *[]pathres.PathRes) {
	ch := x.GetChildren()
	for i := range ch {
		if ch[i].EvalPath(p) {
			*ret = append(*ret, ch[i])
		}
	}
}

func findDescendent(x pathres.PathRes, p *pathexpr.PathExpr, ret *[]pathres.PathRes) {
	ch := x.GetChildren()
	for i := range ch {
		if ch[i].EvalPath(p) {
			*ret = append(*ret, ch[i])
		}
		findDescendent(ch[i], p, ret)
	}
}

func findDescendentOrSelf(x pathres.PathRes, p *pathexpr.PathExpr, ret *[]pathres.PathRes) {
	findSelf(x, p, ret)
	findDescendent(x, p, ret)
}

func findFollowing(x pathres.PathRes, p *pathexpr.PathExpr, ret *[]pathres.PathRes) {
	if x == x.GetParent() {
		return
	}
	par := x.GetParent()
	ch := par.GetChildren()
	i := 0
	for x != ch[i] {
		i++
	}
	i++
	for i < len(ch) {
		findDescendentOrSelf(ch[i], p, ret)
		i++
	}
	findFollowing(par, p, ret)
}

func findFollowingSibling(x pathres.PathRes, p *pathexpr.PathExpr, ret *[]pathres.PathRes) {
	if x == x.GetParent() {
		return
	}
	par := x.GetParent()
	ch := par.GetChildren()
	i := 0
	for x != ch[i] {
		i++
	}
	i++
	for i < len(ch) {
		findSelf(ch[i], p, ret)
		i++
	}
}

func findNamespace(x pathres.PathRes, p *pathexpr.PathExpr, ret *[]pathres.PathRes) {
	if ele, ok := x.(*element.PathResElement); ok {
		for k, v := range ele.NS {
			ns := &namespace.PathResNamespace{
				Value:  xml.Attr{Name: k, Value: v},
				Parent: ele,
			}
			if ns.EvalPath(p) {
				*ret = append(*ret, ns)
			}
		}
	}
}

func findParent(x pathres.PathRes, p *pathexpr.PathExpr, ret *[]pathres.PathRes) {
	if x.GetParent() != x && x.GetParent().EvalPath(p) {
		*ret = append(*ret, x.GetParent())
	}
}

func findPreceding(x pathres.PathRes, p *pathexpr.PathExpr, ret *[]pathres.PathRes) {
	if x == x.GetParent() {
		return
	}
	par := x.GetParent()
	ch := par.GetChildren()
	i := len(ch) - 1
	for x != ch[i] {
		i--
	}
	i--
	for i >= 0 {
		findDescendentOrSelf(ch[i], p, ret)
		i--
	}
	findPreceding(par, p, ret)
}

func findPrecedingSibling(x pathres.PathRes, p *pathexpr.PathExpr, ret *[]pathres.PathRes) {
	if x == x.GetParent() {
		return
	}
	par := x.GetParent()
	ch := par.GetChildren()
	i := len(ch) - 1
	for x != ch[i] {
		i--
	}
	i--
	for i >= 0 {
		findSelf(ch[i], p, ret)
		i--
	}
}

func findSelf(x pathres.PathRes, p *pathexpr.PathExpr, ret *[]pathres.PathRes) {
	if x.EvalPath(p) {
		*ret = append(*ret, x)
	}
}