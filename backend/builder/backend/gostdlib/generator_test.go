package gostdlib

import (
	"testing"
	"github.com/worldiety/devdrasil/backend/builder"
)

func TestGoGenerator(*testing.T) {
	gen := &GoStdlibGenerator{}
	ctx := &builder.GeneratorContext{}
	ctx.Model = createTestApp()
	gen.Generate(ctx)
	ctx.Emit()
}

func createTestApp() *builder.App {
	app := &builder.App{}
	app.Name = "de.worldiety.test"
	app.Doc = "a test app for unit test"

	mod := &builder.Module{}
	mod.Doc = "package for persistence"
	mod.Name = "persistence"
	app.AddModule(mod)

	pojo1 := &builder.Class{}
	pojo1.Exported = true
	pojo1.Doc = "a sample pojo class"
	pojo1.Name = "User"
	pojo1.AddStereotype(builder.StereotypePersistenceModel)

	pojo1.CreateField("The name of a user", "Firstname", "string")
	pojo1.CreateField("The last name of user", "Lastname", "string")
	pojo1.CreateField("The age", "Age", "int64")

	mod.AddClass(pojo1)

	return app
}
