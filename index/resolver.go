package index

import "errors"

type Resolver struct {
	indices   []Index
	providers map[string]string
	nameToPkg map[string]Package
}

func (r *Resolver) AddIndex(i Index) {
	r.indices = append(r.indices, i)

	for _, pkg := range i.Packages {
		if _, ok := r.nameToPkg[pkg.Name]; !ok {
			r.nameToPkg[pkg.Name] = pkg
		}

		for _, provides := range pkg.Provides.PkgConfig {
			if _, ok := r.providers[provides]; !ok {
				r.providers[provides] = pkg.Name
			}
		}
		for _, provides := range pkg.Provides.PkgConfig32 {
			if _, ok := r.providers[provides]; !ok {
				r.providers[provides] = pkg.Name
			}
		}
	}
}

func (r *Resolver) Query(pkgs []string) (res []string, err error) {
	visited := make(map[string]bool)

	var dfs func(name string) error
	dfs = func(name string) error {
		if _, ok := r.providers[name]; ok {
			name = r.providers[name]
		}

		if visited[name] {
			return nil
		}

		if _, ok := r.nameToPkg[name]; !ok {
			return errors.New("Unable to found provider or package " + name)
		}
		visited[name] = true

		res = append(res, name)
		for _, dep := range r.nameToPkg[name].RuntimeDependencies {
			err = dfs(dep.Name)
			if err != nil {
				return err
			}
		}

		return nil
	}

	for _, pkg := range pkgs {
		err = dfs(pkg)
		if err != nil {
			return
		}
	}

	return
}
