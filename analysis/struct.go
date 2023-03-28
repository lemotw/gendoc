package analysis

import (
	"reflect"
	"strings"
	"time"

	"github.com/lemotw/gendoc/model"
	"golang.org/x/net/html"
)

func ReflectAny(intf interface{}, descMap map[string]map[string]string, prefix string) (*model.StructDef, []*model.StructDef) {
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

	structDef, relateStruct := ReflectStruct(t, descMap, prefix)
	if structDef != nil {
		structDef.Name = namePrefix + structDef.Name
	}

	return structDef, relateStruct
}

type WrapStructField struct {
	// for desc
	PName string
	Field  reflect.StructField
}

func getFileds(start reflect.Type) []*WrapStructField {
	stack := []*WrapStructField{}
	ret := []*WrapStructField{}

	for i:=0; i<start.NumField(); i++ {
		field := start.Field(i)
		stack = append(stack, &WrapStructField{
			PName: start.Name(),
			Field: field,
		})
	}

	for len(stack) > 0 {
		if stack[len(stack)-1] == nil {
			continue
		}
		wrapfield := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		if wrapfield.Field.Tag.Get("json") != "" {
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
					Field: fieldj,
				})
			}
		}
	}

	return ret
}

func ReflectStruct(t reflect.Type, descMap map[string]map[string]string, prefix string) (*model.StructDef, []*model.StructDef) {
	structDefs := make(map[reflect.Type]*model.StructDef)

	if t.Kind() != reflect.Struct || t.Kind() == reflect.Ptr {
		return nil, nil
	}

	sd := &model.StructDef{
		Name:   t.Name(),
		Prefix: prefix,
		Fields: make([]*model.StructField, 0),
	}

	fieldList := getFileds(t)
	for i := 0; i < len(fieldList); i++ {
		field := fieldList[i].Field

		// get desc
		stf := &model.StructField{ Req: true}

		jsonArr := strings.Split(field.Tag.Get("json"), ",")
		if len(jsonArr) > 0 {
			stf.Name = jsonArr[0]
		}
		fType := field.Type

	writeTypeName:
		switch fType.Kind() {
		case reflect.Ptr:
			stf.Type += "*"
			fType = fType.Elem()
			goto writeTypeName
		case reflect.Slice:
			stf.Type += "[]"
			fType = fType.Elem()
			goto writeTypeName
		case reflect.Struct:
			if _, ok := structDefs[fType]; !ok && fType != reflect.TypeOf(time.Time{}) {
				subSd, subStructDefs := ReflectAny(reflect.New(fType).Elem().Interface(), descMap, "")
				if subSd == nil {
					return nil, nil
				}
				structDefs[fType] = subSd
				for _, subStructDef := range subStructDefs {
					if _, ok := structDefs[reflect.TypeOf(subStructDef)]; !ok {
						structDefs[reflect.TypeOf(subStructDef)] = subStructDef
					}
				}
			}
			stf.Type += fType.Name()
			stf.Req = !strings.HasPrefix(stf.Type, "*")
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			stf.Type += "int"
			stf.Req = !strings.HasPrefix(stf.Type, "*")
		case reflect.Float32, reflect.Float64:
			stf.Type += "float"
			stf.Req = !strings.HasPrefix(stf.Type, "*")
		default:
			stf.Type += fType.Name()
			stf.Req = !strings.HasPrefix(stf.Type, "*")
		}

		if pkgDescMap, ok := descMap[t.PkgPath()]; ok {
			ttName := fieldList[i].PName + "." + field.Name
			if descStr, ok := pkgDescMap[ttName]; ok {
				stf.Desc = descStr
			}
		}

		sd.Fields = append(sd.Fields, stf)
	}

	structList := make([]*model.StructDef, 0, len(structDefs))
	for _, subSd := range structDefs {
		structList = append(structList, subSd)
	}

	return sd, structList
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

			decTNodes := SearchNodes(tdList[0], html.TextNode, "")
			for i := 0; i < len(decTNodes); i++ {
				f.Desc += decTNodes[i].Data
			}

			list = append(list, f)
		}
	}

	return list
}
