package gostdlib

import (
	"github.com/worldiety/devdrasil/backend/builder"
	"strings"
	. "github.com/dave/jennifer/jen"
)

//generates the entire repository with all it's helper methods
func genRepository(class *builder.Class, f *builder.GoFile) {
	entityName := getStructName(class)
	tableName := "Table" + entityName
	repoName := getStructName(class) + "Repository"

	//define the table name
	f.Const().Defs(Id(tableName).Op("=").Lit(strings.ToLower(entityName)))

	//define the repository
	f.Comment("A " + repoName + " creates, reads, updates, lists or deletes entities of " + class.Name)
	f.Type().Id(repoName).Struct(Id("db").Op("*").Add(DB_Database))

	//define the constructor
	f.Commentf("Creates a new instance of %s", repoName)
	f.Func().Id("New" + repoName).Params(Id("d").Add(Op("*"), DB_Database)).Op("*").Qual("", repoName).
		Block(Return(Op("&").Id(repoName).Values(Id("d"))))

	//define FindAll() ([]*Type,error)
	f.Commentf("FindAll collects all instances of %s and returns it.", entityName)
	f.Func().Params(Id("r").Op("*").Id(repoName)).Id("FindAll").Params().Params(Index().Op("*").Id(entityName), Error()).Block(
		Id("tx").Op(":=").Id("r").Dot("db").Dot("Partition").Params(Id(tableName)).Dot("Begin").Params(False()),
		Defer().Id("tx").Dot("Commit").Params(),
		Id("res").Op(":=").Make(Index().Op("*").Id(entityName), Lit(0)),
		Id("cursor").Op(":=").Id("tx").Dot("GetAll").Params(),
		Defer().Id("cursor").Dot("Close").Params(),
		Id("tmp").Op(":=").Op("&").Qual("bytes", "Buffer"),
		For(Id("cursor").Dot("Next").Params()).Block(
			Id("tmp").Dot("reset").Params(),
			Id("entity").Op(":=").Op("&").Id(entityName).Values(),
			List(Id("_"), Id("e")).Op(":=").Id("cursor").Dot("Get").Params(Id("tmp")),
			If(Id("e").Op("!=").Nil()).Block(
				Qual("log", "Println").Params(Id("e")),
				Continue(),
			),
			Id("e").Op("=").Qual("json", "Unmarshal").Params(Id("tmp").Dot("Bytes").Params(), Id("entity")),
			If(Id("e").Op("!=").Nil()).Block(
				Qual("log", "Println").Params(Id("e")),
				Continue(),
			),
			Id("res").Op("=").Append(Id("res"), Id("entity")),
		),
		Return(Id("res"), Nil()),
	)

	//define GetByID(id db.PK)(*Type,error)
	f.Commentf("GetByID performs a primary key lookup and returns the entity. Returns ErrNotFound if no such entity is available.")
	f.Func().Params(Id("r").Op("*").Id(repoName)).Id("GetByID").Params(Id("id").Add(DB_PK)).Params(Op("*").Id(entityName), Error()).Block(
		Id("tx").Op(":=").Id("r").Dot("db").Dot("Partition").Params(Id(tableName)).Dot("Begin").Params(False()),
		Defer().Id("tx").Dot("Commit").Params(),
		Id("tmp").Op(":=").Op("&").Qual("bytes", "Buffer"),
		List(Id("_"), Id("e")).Op(":=").Id("tx").Dot("Get").Params(Id("id"), Id("tmp")),
		If(Id("e").Op("!=").Nil()).Block(
			Return(Nil(), Id("e")),
		),
		Id("entity").Op(":=").Op("*").Id(entityName).Values(),
		Id("e").Op("=").Qual("json", "Unmarshal").Params(Id("tmp").Dot("Bytes").Params(), Id("entity")),
		If(Id("e").Op("!=").Nil()).Block(
			Return(Nil(), Id("e")),
		),
		Return(Id("entity"), Nil()),
	)

	//define DeleteByID(id db.PK) error
	f.Commentf("DeleteByID performs a primary key lookup and ensures that the entity is gone.")
	f.Func().Params(Id("r").Op("*").Id(repoName)).Id("DeleteByID").Params(Id("id").Add(DB_PK)).Params(Op("*").Id(entityName), Error()).Block()

}
