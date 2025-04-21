package dc

type providerAbstract interface {
	reset()
}

type provider[T any] struct {
	createDependency func() T
	instance         T
	mock             T
	isPending        bool
	hasInstance      bool
	hasMock          bool
}

var providers = make([]providerAbstract, 0)

func Provider[T any](createDependency func() T) *provider[T] {
	p := &provider[T]{
		createDependency: createDependency,
		mock:             *new(T),
		hasMock:          false,
		isPending:        false,
		hasInstance:      false,
	}
	providers = append(providers, p)
	return p
}

func (p *provider[T]) Use() T {
	if p.hasMock {
		return p.mock
	}

	if p.isPending {
		panic("detected circular dependency, please check your code")
	}

	if !p.hasInstance {
		p.isPending = true
		p.instance = p.createDependency()
		p.isPending = false
		p.hasInstance = true
	}

	return p.instance
}

func (p *provider[T]) reset() {
	p.mock = *new(T)
	p.hasMock = false
	p.instance = *new(T)
	p.hasInstance = false
	p.isPending = false
}

func (p *provider[T]) Mock(mockObject T) {
	if p.hasInstance {
		panic("instance already set, please reset the provider first")
	}

	p.mock = mockObject
	p.hasMock = true
}

func Reset() {
	for _, p := range providers {
		p.reset()
	}
}
