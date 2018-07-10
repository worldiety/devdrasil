package gostdlib

import (
	"github.com/worldiety/devdrasil/backend/builder"
	"strings"
	"hash/fnv"
	"strconv"
	. "github.com/dave/jennifer/jen"
)

var DB_PK = Qual("github.com/worldiety/devdrasil/db", "PK")
var DB_Database = Qual("github.com/worldiety/devdrasil/db", "Database")

type GoPackageType struct {
	// something like mypackage/otherpackage
	Import string
	// something like otherpackage
	NamedImport string
	//something like mypackage.otherpackage.MyType
	FullQualifiedName string
	//something like MyType
	Name string
}

func (t *GoPackageType) GetIdentifier() string {
	return "*" + t.NamedImport + "." + t.Name
}

type GoStdlibGenerator struct {
}

func (g *GoStdlibGenerator) Generate(ctx *builder.GeneratorContext) error {
	ctx.Model.EachClass(func(class *builder.Class) {
		if class.HasStereotype(builder.StereotypePersistenceModel) {
			//e.g. creates a mymodel/person.go

			f := ctx.NewGoFile(class.Parent.GetName(), getFilename(class, ""))
			genStruct(class, f)

			genRepository(class, f)
		}
	})
	return nil
}


func genStruct(class *builder.Class, f *builder.GoFile) {
	fields := make([]Code, 0)
	fields = append(fields, Comment("The ID is unique and identifies an entity as a primary key.").Line().Id("ID").Add(DB_PK))
	for _, field := range class.Fields {
		code := Comment(field.Doc).Line().Id(field.Name)
		code.Add(genTypeStatement(&field.Type))
		fields = append(fields, code)
	}
	f.Comment(class.Doc)
	f.Type().Id(getStructName(class)).Struct(fields...)
}

func genTypeStatement(t *builder.Type) *Statement {
	var stmt *Statement
	switch t.FullQualifiedName {
	case "int64":
		stmt = Int64()
	case "void":
		stmt = nil
	case "string":
		stmt = String()
	case "float64":
		stmt = Float64()
	case "bool":
		stmt = Bool()
	case "List":
		stmt = Index()
		stmt.Add(genTypeStatement(t.Generics[0]))
	case "Map":
		stmt = Map(genTypeStatement(t.Generics[0]))
		stmt.Add(genTypeStatement(t.Generics[1]))
	default:
		//custom type
		stmt = QualifierFromType(t)
	}
	return stmt
}

func QualifierFromClass(class *builder.Class) *Statement {
	return Qual(getFullPackageName(class.Parent), getStructName(class))
}
func QualifierFromType(t *builder.Type) *Statement {
	if !strings.Contains(t.FullQualifiedName, ".") {
		return Qual("", t.FullQualifiedName)
	}
	segments := strings.Split(t.FullQualifiedName, ".")
	importPath := strings.Join(segments[0:len(segments)-2], "/")
	return Qual(importPath, segments[len(segments)-1])
}

//returns the full package name separated with /
func getFullPackageName(module builder.ModuleLike) string {
	tmp := ""
	root := module
	for root != nil {
		tmp = module.GetName() + "/" + tmp
		root = root.GetParent()
	}
	return tmp
}

// suffix is concated before .go extension is appended. All things are lowercase. E.g. returns mypackage/otherpackage/person.go
func getFilename(class *builder.Class, suffix string) string {
	path := strings.ToLower(class.Name) + suffix + ".go"
	root := class.Parent
	for root != nil {
		path = strings.ToLower(root.GetName()) + "/" + path
		root = root.GetParent()
	}
	return path
}

//returns the (in-packackge) struct name of the class, which respects the correct naming conventions, regarding public/private
func getStructName(class *builder.Class) string {
	if class.Exported {
		up := strings.ToUpper(class.Name[0:1])
		return up + class.Name[1:]
	} else {
		down := strings.ToLower(class.Name[0:1])
		return down + class.Name[1:]
	}
}

//either returns "main" or the lower case of the parent module name
func getPackageName(class *builder.Class) string {
	if class.Parent == nil {
		return "main"
	} else {
		return strings.ToLower(class.Parent.GetName())
	}
}

//returns all import required for the given type, including generics e.g. for slices and maps
func getImport(t *builder.Type) []*GoPackageType {
	res := make([]*GoPackageType, 0)
	res = append(res, getImportFQN(t.FullQualifiedName))
	for _, g := range t.Generics {
		res = append(res, getImport(g)...)
	}
	return res
}

//returns an import from a fullqualified dot (.) separated type name
func getImportFQN(fullQualifiedName string) *GoPackageType {
	if !strings.Contains(fullQualifiedName, ".") {
		return &GoPackageType{Import: "", NamedImport: "", FullQualifiedName: fullQualifiedName}
	}
	segments := strings.Split(fullQualifiedName, ".")
	importPath := strings.Join(segments[0:len(segments)-1], "/")
	namedImport := "_" + segments[len(segments)-2] + strconv.FormatInt(int64(hashcode(importPath)), 10)
	name := segments[len(segments)-1]
	return &GoPackageType{Import: importPath, NamedImport: namedImport, FullQualifiedName: fullQualifiedName, Name: name}
}

func hashcode(s string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(s))
	return h.Sum32()
}
