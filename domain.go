package shift

type Domain struct {
	Path   string
	Router *Router
	parent *Domain
}

func newDomain(path string, parent *Domain) *Domain {
	domain := &Domain{
		Path:   path,
		Router: newRouter(),
		parent: parent,
	}

	if parent != nil {
		parent.Router.mount(path, domain.Router)
	}

	return domain
}
