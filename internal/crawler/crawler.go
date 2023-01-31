package crawler

import (
	"bytes"
	"fmt"
	"github.com/bitfield/script"
	"godeps/internal/slices"
	"math"
	"strings"
)

var ErrSyntaxWrong = fmt.Errorf("packages list invalid")
var ErrGoListFailure = fmt.Errorf("go list failed")

type Config struct {
	Depth int
}

type Package string

type Crawler struct {
}

func (c *Crawler) Crawl(module string, config Config) map[Package][]Package {
	if config.Depth == 0 {
		config.Depth = math.MaxInt
	}
	var importsMap = make(map[Package][]Package)

	var scanner ImportScanner
	var workingSet, backlogSet Set
	backlogSet.Add(Package(module))

	var currentDepth = 0
	for !backlogSet.Empty() {
		currentDepth++
		workingSet = backlogSet
		for !workingSet.Empty() {
			var pkg = workingSet.Pop()
			if _, ok := importsMap[pkg]; ok {
				continue
			}
			imports, err := scanner.ListImports(pkg)
			if err != nil {
				continue
			}

			imports = slices.Filter(imports, func(p Package) bool { return strings.HasPrefix(string(p), module) })
			importsMap[pkg] = imports
			if currentDepth < config.Depth {
				backlogSet.Add(imports...)
			}
		}
	}
	return importsMap
}

type ImportScanner struct {
}

func (s *ImportScanner) ListImports(pkg Package) ([]Package, error) {
	p := script.Exec(fmt.Sprintf("go list -f '{{ .Imports }}' %s", pkg))
	b, err := p.Bytes()
	if err != nil {
		return nil, ErrGoListFailure
	}
	if p.ExitStatus() != 0 {
		return nil, ErrGoListFailure
	}
	return listOutputToArray(string(b))
}

func listOutputToArray(out string) ([]Package, error) {
	out = strings.TrimRight(out, "\n")

	var packages []Package
	var buffer bytes.Buffer
	for i := range out {
		switch b := out[i]; b {
		case '[':
			if i != 0 {
				return nil, ErrSyntaxWrong
			}
		case ']':
			if i != len(out)-1 {
				return nil, ErrSyntaxWrong
			}
			return append(packages, Package(buffer.String())), nil
		case ' ':
			packages = append(packages, Package(buffer.String()))
			buffer.Reset()
		default:
			buffer.WriteByte(b)
		}
	}
	return nil, ErrSyntaxWrong
}
