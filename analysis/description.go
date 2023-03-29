package analysis

import (
	"errors"
	"go/parser"
	"go/token"
	"io/fs"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"golang.org/x/tools/go/packages"
)

func GetDescriptionMap() (map[string]map[string][]string, error) {
	descriptionMap := make(map[string]map[string][]string)

	// get package root
	currentDir, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	for {
		if _, err := os.Stat(filepath.Join(currentDir, "go.mod")); err == nil {
			break
		}

		if currentDir == filepath.Dir(currentDir) {
			return nil, errors.New("go.mod not found")
		}

		currentDir = filepath.Dir(currentDir)
	}

	// pkg name maintain
	pkgNameMap := make(map[string]string)
	updatePkgNameMap := func(dir string) (string, error) {
		pkgs, err := packages.Load(&packages.Config{
			Mode:  packages.NeedName,
			Dir:   dir,
			Tests: false,
		}, dir)
		if err != nil {
			return "", err
		}

		if len(pkgs) > 0 {
			pkgNameMap[dir] = pkgs[0].ID
			return pkgs[0].ID, nil
		}

		return "", errors.New("package not found")
	}

	// walk all file parse model description
	err = filepath.Walk(currentDir, func(path string, info fs.FileInfo, err error) error {
		if filepath.Ext(path) != ".go" {
			// ignore not go file
			return nil
		}

		// fetch package name
		pkgName, ok := pkgNameMap[filepath.Dir(path)]
		if !ok {
			if pkgName, err = updatePkgNameMap(filepath.Dir(path)); err != nil {
				return err
			}
		}

		file, err := parser.ParseFile(token.NewFileSet(), path, nil, parser.ParseComments)
		if err != nil {
			return err
		}

		pkgDescMap, ok := descriptionMap[pkgName]
		if !ok {
			pkgDescMap = make(map[string][]string)
			descriptionMap[pkgName] = pkgDescMap
		}

		for i := 0; i < len(file.Comments); i++ {
			for j := 0; j < len(file.Comments[i].List); j++ {
				// match with spec("// @") to find description recommend
				if strings.HasPrefix(file.Comments[i].List[j].Text, "// @") {
					descStr := file.Comments[i].List[j].Text[4:]
					descArr := strings.Split(descStr, ":")
					if len(descArr) >= 2 {
						val, err := strconv.Unquote(strings.TrimSpace(strings.Join(descArr[1:], ":")))
						if err != nil {
							return nil
						}
						pkgDescMap[descArr[0]] = append(pkgDescMap[descArr[0]], val)
					}
				}
			}
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return descriptionMap, nil
}
