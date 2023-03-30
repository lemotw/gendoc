package analysis

import (
	"reflect"
	"strings"
	"time"

	"github.com/lemotw/gendoc/model"
	"golang.org/x/net/html"
)

func ReflectAny(intf interface{}, descMap map[string]map[string][]string) []*model.StructDef {
	t := reflect.TypeOf(intf)
	namePrefix := ""

unwrap:
	switch t.Kind() {
	case reflect.Slice:
		namePrefix += "[]"
		t = t.Elem()
		goto unwrap
	case reflect.Ptr:
		namePrefix += "*"
		t = t.Elem()
		goto unwrap
	}

	structDef := ReflectStructDFS(t, descMap)
	if len(structDef) > 0 {
		structDef[0].Name = namePrefix + structDef[0].Name
	}

	return structDef
}

type WrapStructField struct {
	// for desc
	PName string
	Pkg   string
	Field reflect.StructField
}

func (fild *WrapStructField) GetDesc(descMap map[string]map[string][]string) []string {
	if pkgDescMap, ok := descMap[fild.Pkg]; ok {
		ttName := fild.PName + "." + fild.Field.Name
		if descStr, ok := pkgDescMap[ttName]; ok {
			return descStr
		}
	}
	return []string{}
}

func getFileds(start reflect.Type) []*WrapStructField {
	stack := []*WrapStructField{}
	ret := []*WrapStructField{}

	for i := 0; i < start.NumField(); i++ {
		field := start.Field(i)
		stack = append(stack, &WrapStructField{
			PName: start.Name(),
			Pkg:   start.PkgPath(),
			Field: field,
		})
	}

	for len(stack) > 0 {
		if stack[len(stack)-1] == nil {
			continue
		}
		wrapfield := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		if wrapfield.Field.Tag.Get("gendoc") != "" {
			ret = append(ret, wrapfield)
		}

		if wrapfield.Field.Tag.Get("inheritance") == "true" {
			for j := 0; j < wrapfield.Field.Type.NumField(); j++ {
				rawFieldType := wrapfield.Field.Type
				for rawFieldType.Kind() == reflect.Ptr || rawFieldType.Kind() == reflect.Slice {
					rawFieldType = rawFieldType.Elem()
				}

				fieldj := wrapfield.Field.Type.Field(j)
				stack = append(stack, &WrapStructField{
					PName: rawFieldType.Name(),
					Pkg:   rawFieldType.PkgPath(),
					Field: fieldj,
				})
			}
		}
	}

	return ret
}

func ReflectStructDFS(t reflect.Type, descMap map[string]map[string][]string) []*model.StructDef {
	ret := []*model.StructDef{}
	typeStack := []reflect.Type{t}
	analysisedStruct := map[reflect.Type]struct{}{}

	for len(typeStack) > 0 {
		tp := typeStack[len(typeStack)-1]
		typeStack = typeStack[:len(typeStack)-1]

		fieldList := getFileds(tp)
		structDefind := &model.StructDef{
			Name:   tp.Name(),
			Fields: make([]*model.StructField, 0),
		}

		for i := 0; i < len(fieldList); i++ {
			fieldType := fieldList[i].Field.Type

			sfield := &model.StructField{}
			jsonArr := strings.Split(fieldList[i].Field.Tag.Get("gendoc"), ",")
			for ind := 0; ind < len(jsonArr); ind++ {
				switch ind {
				case 0:
					sfield.Name = jsonArr[ind]
				case 1:
					sfield.Source = jsonArr[ind]
				case 2:
					sfield.Req = (jsonArr[ind] == "Y" || jsonArr[ind] == "y")
				}
			}

		writeTypeName:
			switch fieldType.Kind() {
			case reflect.Ptr:
				sfield.Type += "*"
				fieldType = fieldType.Elem()
				goto writeTypeName
			case reflect.Slice:
				sfield.Type += "[]"
				fieldType = fieldType.Elem()
				goto writeTypeName
			case reflect.Struct:
				if _, ok := analysisedStruct[fieldType]; !ok {
					if fieldType == reflect.TypeOf(time.Time{}) {
						break
					}
					analysisedStruct[fieldType] = struct{}{}
					typeStack = append(typeStack, fieldType)
				}
				sfield.Type += fieldType.Name()
				sfield.Req = !strings.HasPrefix(sfield.Type, "*")
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				sfield.Type += "int"
				sfield.Req = !strings.HasPrefix(sfield.Type, "*")
			case reflect.Float32, reflect.Float64:
				sfield.Type += "float"
				sfield.Req = !strings.HasPrefix(sfield.Type, "*")
			default:
				sfield.Type += fieldType.Name()
				sfield.Req = !strings.HasPrefix(sfield.Type, "*")
			}

			// get desc
			sfield.Desc = fieldList[i].GetDesc(descMap)

			structDefind.Fields = append(structDefind.Fields, sfield)
		}

		ret = append(ret, structDefind)
	}

	return ret
}

func NodeToFieldList(node *html.Node) model.StructTable {
	var list []*model.StructField

	trlist := SearchNodes(node, html.ElementNode, "tr")
	if len(trlist) == 0 {
		return nil
	}

	for i := 0; i < len(trlist); i++ {
		tdList := SearchNodes(trlist[i], html.ElementNode, "td")
		if len(tdList) == 4 {
			f := &model.StructField{}

			nameTNodes := SearchNodes(tdList[3], html.TextNode, "")
			for i := 0; i < len(nameTNodes); i++ {
				f.Name += nameTNodes[i].Data
			}

			typeTNodes := SearchNodes(tdList[2], html.TextNode, "")
			for i := 0; i < len(typeTNodes); i++ {
				f.Type += typeTNodes[i].Data
			}

			reqStr := ""
			reqTNodes := SearchNodes(tdList[1], html.TextNode, "")
			for i := 0; i < len(reqTNodes); i++ {
				reqStr += reqTNodes[i].Data
			}
			f.Req = reqStr == "Y"

			desc := ""
			decTNodes := SearchNodes(tdList[0], html.TextNode, "")
			for i := 0; i < len(decTNodes); i++ {
				desc += decTNodes[i].Data
			}
			f.Desc[0] = desc

			list = append(list, f)
		}
	}

	return list
}
