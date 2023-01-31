package presentation

import (
	"godeps/internal/crawler"
	"io"
)

func AsDotNotation(importsMap map[crawler.Package][]crawler.Package, writer io.Writer) error {
	_, err := writer.Write([]byte("digraph {\n"))
	if err != nil {
		return err
	}
	for src, imports := range importsMap {
		for _, imprt := range imports {
			_, err = writer.Write([]byte("\"" + src + "\"" + " -> " + "\"" + imprt + "\"" + "\n"))
			if err != nil {
				return err
			}
		}
	}
	_, err = writer.Write([]byte("}\n"))
	return err
}
