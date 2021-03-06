// This work is subject to the CC0 1.0 Universal (CC0 1.0) Public Domain Dedication
// license. Its contents can be found at:
// http://creativecommons.org/publicdomain/zero/1.0/

package xmlx

import (
	"os"
	"xml"
	"bytes"
	"fmt"
	"strconv"
)

const (
	NT_ROOT = iota
	NT_DIRECTIVE
	NT_PROCINST
	NT_COMMENT
	NT_ELEMENT
)

type Attr struct {
	Name  xml.Name // Attribute namespace and name.
	Value string   // Attribute value. 
}

type Node struct {
	Type       byte     // Node type.
	Name       xml.Name // Node namespace and name.
	Children   []*Node  // Child nodes.
	Attributes []*Attr  // Node attributes.
	Parent     *Node    // Parent node.
	Value      string   // Node value.
	Target     string   // procinst field.
}

func NewNode(tid byte) *Node {
	n := new(Node)
	n.Type = tid
	n.Children = make([]*Node, 0, 10)
	n.Attributes = make([]*Attr, 0, 10)
	return n
}

// This wraps the standard xml.Unmarshal function and supplies this particular
// node as the content to be unmarshalled.
func (this *Node) Unmarshal(obj interface{}) os.Error {
	return xml.Unmarshal(bytes.NewBuffer(this.Bytes()), obj)
}

// Get node value as string
func (this *Node) S(namespace, name string) string {
	if node := rec_SelectNode(this, namespace, name); node != nil {
		return node.Value
	}
	return ""
}

// Get node value as int
func (this *Node) I(namespace, name string) int {
	if node := rec_SelectNode(this, namespace, name); node != nil && node.Value != "" {
		n, _ := strconv.Atoi(node.Value)
		return n
	}
	return 0
}

// Get node value as int64
func (this *Node) I64(namespace, name string) int64 {
	if node := rec_SelectNode(this, namespace, name); node != nil && node.Value != "" {
		n, _ := strconv.Atoi64(node.Value)
		return n
	}
	return 0
}

// Get node value as uint
func (this *Node) U(namespace, name string) uint {
	if node := rec_SelectNode(this, namespace, name); node != nil && node.Value != "" {
		n, _ := strconv.Atoui(node.Value)
		return n
	}
	return 0
}

// Get node value as uint64
func (this *Node) U64(namespace, name string) uint64 {
	if node := rec_SelectNode(this, namespace, name); node != nil && node.Value != "" {
		n, _ := strconv.Atoui64(node.Value)
		return n
	}
	return 0
}

// Get node value as float32
func (this *Node) F32(namespace, name string) float32 {
	if node := rec_SelectNode(this, namespace, name); node != nil && node.Value != "" {
		n, _ := strconv.Atof32(node.Value)
		return n
	}
	return 0
}

// Get node value as float64
func (this *Node) F64(namespace, name string) float64 {
	if node := rec_SelectNode(this, namespace, name); node != nil && node.Value != "" {
		n, _ := strconv.Atof64(node.Value)
		return n
	}
	return 0
}

// Get node value as bool
func (this *Node) B(namespace, name string) bool {
	if node := rec_SelectNode(this, namespace, name); node != nil && node.Value != "" {
		n, _ := strconv.Atob(node.Value)
		return n
	}
	return false
}

// Get attribute value as string
func (this *Node) As(namespace, name string) string {
	for _, v := range this.Attributes {
		if (namespace == "*" || namespace == v.Name.Space) && name == v.Name.Local {
			return v.Value
		}
	}
	return ""
}

// Get attribute value as int
func (this *Node) Ai(namespace, name string) int {
	if s := this.As(namespace, name); s != "" {
		n, _ := strconv.Atoi(s)
		return n
	}
	return 0
}

// Get attribute value as uint
func (this *Node) Au(namespace, name string) uint {
	if s := this.As(namespace, name); s != "" {
		n, _ := strconv.Atoui(s)
		return n
	}
	return 0
}

// Get attribute value as uint64
func (this *Node) Au64(namespace, name string) uint64 {
	if s := this.As(namespace, name); s != "" {
		n, _ := strconv.Atoui64(s)
		return n
	}
	return 0
}

// Get attribute value as int64
func (this *Node) Ai64(namespace, name string) int64 {
	if s := this.As(namespace, name); s != "" {
		n, _ := strconv.Atoi64(s)
		return n
	}
	return 0
}

// Get attribute value as float32
func (this *Node) Af32(namespace, name string) float32 {
	if s := this.As(namespace, name); s != "" {
		n, _ := strconv.Atof32(s)
		return n
	}
	return 0
}

// Get attribute value as float64
func (this *Node) Af64(namespace, name string) float64 {
	if s := this.As(namespace, name); s != "" {
		n, _ := strconv.Atof64(s)
		return n
	}
	return 0
}

// Get attribute value as bool
func (this *Node) Ab(namespace, name string) bool {
	if s := this.As(namespace, name); s != "" {
		n, _ := strconv.Atob(s)
		return n
	}
	return false
}

// Returns true if this node has the specified attribute. False otherwise.
func (this *Node) HasAttr(namespace, name string) bool {
	for _, v := range this.Attributes {
		if (namespace == "*" || namespace == v.Name.Space) && name == v.Name.Local {
			return true
		}
	}
	return false
}

// Select single node by name
func (this *Node) SelectNode(namespace, name string) *Node {
	return rec_SelectNode(this, namespace, name)
}

func rec_SelectNode(cn *Node, namespace, name string) *Node {
	// Allow wildcard for namespace names. Meaning we will match any namespace
	// name with a matching local name.
	if (namespace == "*" || cn.Name.Space == namespace) && cn.Name.Local == name {
		return cn
	}

	var tn *Node
	for _, v := range cn.Children {
		if tn = rec_SelectNode(v, namespace, name); tn != nil {
			return tn
		}
	}
	return nil
}

// Select multiple nodes by name
func (this *Node) SelectNodes(namespace, name string) []*Node {
	list := make([]*Node, 0, 16)
	rec_SelectNodes(this, namespace, name, &list)
	return list
}

func rec_SelectNodes(cn *Node, namespace, name string, list *[]*Node) {
	// Allow wildcard for namespace names. Meaning we will match any namespace
	// name with a matching local name.
	if (namespace == "*" || cn.Name.Space == namespace) && cn.Name.Local == name {
		*list = append(*list, cn)
		return
	}

	for _, v := range cn.Children {
		rec_SelectNodes(v, namespace, name, list)
	}
}

// Convert node to appropriate []byte representation based on it's @Type.
// Note that NT_ROOT is a special-case empty node used as the root for a
// Document. This one has no representation by itself. It merely forwards the
// String() call to it's child nodes.
func (this *Node) Bytes() (b []byte) {
	switch this.Type {
	case NT_PROCINST:
		b = this.printProcInst()
	case NT_COMMENT:
		b = this.printComment()
	case NT_DIRECTIVE:
		b = this.printDirective()
	case NT_ELEMENT:
		b = this.printElement()
	case NT_ROOT:
		b = this.printRoot()
	}
	return
}

// Convert node to appropriate string representation based on it's @Type.
// Note that NT_ROOT is a special-case empty node used as the root for a
// Document. This one has no representation by itself. It merely forwards the
// String() call to it's child nodes.
func (this *Node) String() (s string) {
	return string(this.Bytes())
}

func (this *Node) printRoot() []byte {
	var b bytes.Buffer
	for _, v := range this.Children {
		b.WriteString(v.String())
	}
	return b.Bytes()
}

func (this *Node) printProcInst() []byte {
	return []byte("<?" + this.Target + " " + this.Value + "?>")
}

func (this *Node) printComment() []byte {
	return []byte("<!-- " + this.Value + " -->")
}

func (this *Node) printDirective() []byte {
	return []byte("<!" + this.Value + "!>")
}

func (this *Node) printElement() []byte {
	var b bytes.Buffer

	if len(this.Name.Space) > 0 {
		b.WriteRune('<')
		b.WriteString(this.Name.Space)
		b.WriteRune(':')
		b.WriteString(this.Name.Local)
	} else {
		b.WriteRune('<')
		b.WriteString(this.Name.Local)
	}

	for _, v := range this.Attributes {
		if len(v.Name.Space) > 0 {
			b.WriteString(fmt.Sprintf(` %s:%s="%s"`, v.Name.Space, v.Name.Local, v.Value))
		} else {
			b.WriteString(fmt.Sprintf(` %s="%s"`, v.Name.Local, v.Value))
		}
	}

	if len(this.Children) == 0 && len(this.Value) == 0 {
		b.WriteString(" />")
		return b.Bytes()
	}

	b.WriteRune('>')

	for _, v := range this.Children {
		b.WriteString(v.String())
	}

	b.WriteString(this.Value)
	if len(this.Name.Space) > 0 {
		b.WriteString("</")
		b.WriteString(this.Name.Space)
		b.WriteRune(':')
		b.WriteString(this.Name.Local)
		b.WriteRune('>')
	} else {
		b.WriteString("</")
		b.WriteString(this.Name.Local)
		b.WriteRune('>')
	}

	return b.Bytes()
}

// Add a child node
func (this *Node) AddChild(t *Node) {
	if t.Parent != nil {
		t.Parent.RemoveChild(t)
	}
	t.Parent = this
	this.Children = append(this.Children, t)
}

// Remove a child node
func (this *Node) RemoveChild(t *Node) {
	p := -1
	for i, v := range this.Children {
		if v == t {
			p = i
			break
		}
	}

	if p == -1 {
		return
	}

	copy(this.Children[p:], this.Children[p+1:])
	this.Children = this.Children[0 : len(this.Children)-1]

	t.Parent = nil
}
