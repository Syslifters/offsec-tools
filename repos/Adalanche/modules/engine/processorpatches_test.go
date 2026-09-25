package engine

import (
	"slices"
	"testing"
)

func TestNodePatchSetSetFlexMatchesNodeSetFlex(t *testing.T) {
	direct := NewNode()
	direct.SetFlex(
		IgnoreBlanks,
		Name, NV("alpha"),
		Description, NV("desc"),
		Tag, AttributeValues{NV("bench")},
	)

	patched := NewNode()
	var patch NodePatchSet
	patch.SetFlex(
		patched,
		IgnoreBlanks,
		Name, NV("alpha"),
		Description, NV("desc"),
		Tag, AttributeValues{NV("bench")},
	)
	patch.Apply(NewIndexedGraph())

	for _, attr := range []Attribute{Name, Description, Tag} {
		got := patched.Attr(attr)
		want := direct.Attr(attr)
		if !slices.EqualFunc(got, want, func(a, b AttributeValue) bool {
			return a.Compare(b) == 0
		}) {
			t.Fatalf("attribute %v mismatch:\n got: %#v\nwant: %#v", attr, got, want)
		}
	}
}

func TestNodePatchSetApplyPreservesOperationOrder(t *testing.T) {
	node := NewNode()
	var patch NodePatchSet
	patch.Set(node, Description, NV("first"))
	patch.Clear(node, Description)
	patch.Set(node, Description, NV("second"))
	patch.AddTag(node, "bench")
	patch.AddTag(node, "bench")
	patch.Apply(NewIndexedGraph())

	if got := node.OneAttrString(Description); got != "second" {
		t.Fatalf("unexpected description %q", got)
	}
	tags := node.Attr(Tag)
	if tags.Len() != 1 || tags.First().String() != "bench" {
		t.Fatalf("unexpected tags %#v", tags)
	}
}
