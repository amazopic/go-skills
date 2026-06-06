// Package prototype is an example of the Prototype Pattern.
//
// The Prototype pattern creates new objects by cloning an existing
// instance (the prototype) via a Clone method, rather than constructing
// them from scratch. This is useful when object creation is expensive or
// when the concrete type to instantiate is decided at runtime.
//
// Here Prototyper declares Clone, and ConcreteProduct.Clone returns an
// independent copy of the receiver.
package prototype

// Prototyper provides a cloning interface.
type Prototyper interface {
	Clone() Prototyper
	GetName() string
}

// ConcreteProduct implements product "A"
type ConcreteProduct struct {
	name string // Имя продукта
}

// NewConcreteProduct is the Prototyper constructor.
func NewConcreteProduct(name string) Prototyper {
	return &ConcreteProduct{
		name: name,
	}
}

// GetName returns product name
func (p *ConcreteProduct) GetName() string {
	return p.name
}

// Clone returns a cloned object.
func (p *ConcreteProduct) Clone() Prototyper {
	return &ConcreteProduct{p.name}
}
