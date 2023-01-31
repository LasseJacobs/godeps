package crawler

type Mark struct{}

type Set struct {
	data map[Package]Mark
}

func (s *Set) Empty() bool {
	return len(s.data) == 0
}

func (s *Set) Add(pkgs ...Package) {
	if s.data == nil {
		s.data = make(map[Package]Mark)
	}
	for _, p := range pkgs {
		s.data[p] = struct{}{}
	}
}

func (s *Set) Pop() Package {
	for k := range s.data {
		delete(s.data, k)
		return k
	}
	return ""
}
