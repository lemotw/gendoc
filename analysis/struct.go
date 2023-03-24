package analysis

import (
	"reflect"
	"strings"

	"github.com/lemotw/gendoc/model"
	"golang.org/x/net/html"
)

func ReflectStruct(intf interface{}, descMap map[string]map[string]string) (*model.StructDef, []*model.StructDef) {
	structDefs := make(map[reflect.Type]*model.StructDef)

	t := reflect.TypeOf(intf)
	if t.Kind() != reflect.Struct {
		return nil, nil
	}

	sd := &model.StructDef{
		Name:   t.Name(),
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
			if _, ok := structDefs[fType]; !ok {
				subSd, subStructDefs := ReflectStruct(reflect.New(fType).Elem().Interface(), descMap)
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

	theadNode := FindNode(node, "thead")
	for n := theadNode.FirstChild; n != nil; n = n.NextSibling {
		var tdList []*html.Node
		trNode := FindNode(n, "tr")
		for td := trNode.FirstChild; td != nil; td = td.NextSibling {
			tdList = append(tdList, td)
		}

		if len(tdList) == 4 {
			f := &model.StructField{}

			nameTNodes := SearchNodes(tdList[0], html.TextNode)
			for i := 0; i < len(nameTNodes); i++ {
				f.Name += nameTNodes[i].Data
			}

			typeTNodes := SearchNodes(tdList[1], html.TextNode)
			for i := 0; i < len(typeTNodes); i++ {
				f.Type += typeTNodes[i].Data
			}

			reqStr := ""
			reqTNodes := SearchNodes(tdList[2], html.TextNode)
			for i := 0; i < len(reqTNodes); i++ {
				reqStr += reqTNodes[i].Data
			}
			f.Req = reqStr == "Y"

			decTNodes := SearchNodes(tdList[3], html.TextNode)
			for i := 0; i < len(decTNodes); i++ {
				f.Desc += decTNodes[i].Data
			}

			list = append(list, f)
		}
	}

	return list
}