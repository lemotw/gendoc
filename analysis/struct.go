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

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		jsonTag := field.Tag.Get("json")
		if jsonTag == "" {
			continue
		}
		// get desc
		desc := ""
		pkgDescMap, ok := descMap[t.PkgPath()]
		if ok {
			ttName := t.Name() + "." + field.Name
			if descStr, ok := pkgDescMap[ttName]; ok {
				desc = descStr
			}
		}

		stf := &model.StructField{Name: jsonTag, Req: true, Desc: desc}
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
