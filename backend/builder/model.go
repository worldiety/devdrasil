package builder

import "encoding/json"

//unmarshals the given bytes as json into an app. Re-creates the parent structure
func Unmarshal(b []byte) (*App, error) {
	r := &App{}
	err := json.Unmarshal(b, r)
	rebuildParents(r)
	return r, err
}

func rebuildParents(root ModuleLike) {
	for _, class := range root.GetClasses() {
		class.Parent = root
		for _, method := range class.Methods {
			method.Parent = class
			for _, p := range method.Parameters {
				p.Parent = method
			}
			for _, p := range method.Returns {
				p.Parent = method
			}
		}
		for _, field := range class.Fields {
			field.Parent = class
		}
	}

	for _, mod := range root.GetModules() {
		mod.SetParent(root)
		rebuildParents(mod)
	}
}

/*
A stereotype defines the role a class takes in an application. It directly influences the source code generator and the platform to use, e.g. for backend and frontend.
It also influences lifecycle and injection capabilities.
 */
type Stereotype string

const (
	//Frontend: Represents a view, e.g. like a form. It is the smallest and most reusable unit in the frontend.
	//It's lifecycle is usually at most to the lifecycle of a {@link STEREOTYPE.USER_INTERFACE_STATE}
	StereotypeView Stereotype = "VIEW"

	//Backend: Represents a controller, which is a singleton. It incorporates the actual server side application logic.
	StereotypeController Stereotype = "CONTROLLER"

	/*
	Frontend: In the MVVM world, this is the ViewModel and provides a (potential) two way binding to a {@link STEREOTYPE.VIEW}.
	On the other hand, it just may be a simple controller which is used directly (MVC world). In both cases it is a singleton
    at application level.
	 */
	StereotypeViewController Stereotype = "VIEW_CONTROLLER"
	/*
	Frontend: Represents a logical state which allows a forward/backward navigation. It uses typically {@link STEREOTYPE.VIEW}s and applies them.
    It provides a basic create/destroy lifecycle and provides the infrastructure for MVVM or other callback bindings to resolve leaks reliable.
	 */
	StereotypeUserInterfaceState Stereotype = "USER_INTERFACE_STATE"
	/*
	Frontend & Backend: Represents a model which is persistent and supports CRUD (create read update delete). Each model
    is made available by a repository. The backend holds the actual truth, whereas the frontend repository just
    represents the client's end, perhaps providing additional caching purposes. A {@link STEREOTYPE.VIEW_CONTROLLER} typically imports
    a repository to make data available to the view.
	 */
	StereotypePersistenceModel Stereotype = "PERSISTENCE_MODEL"
	/*
	Frontend & Backend: An abstract component, which encapsulates a model with optional methods. It can be used to either
    represent common intermediate models (without persistence) or separation of concerns (e.g. splitted controller logic).
    It is shared across the frontend and backend. See also {@link STEREOTYPE.FRONTEND_COMPONENT} and {@link STEREOTYPE.BACKEND_COMPONENT}.
	 */
	StereotypeComponent Stereotype = "COMPONENT"
	//Frontend: a component which is only applicable to the frontend. See also {@link STEREOTYPE.COMPONENT}.
	StereotypeFrontendComponent Stereotype = "FRONTEND_COMPONENT"
	//Backend: a component which is only applicable to the backend. See also {@link STEREOTYPE.COMPONENT}.
	StereotypeBackendComponent Stereotype = "BACKEND_COMPONENT"
)

type ModuleLike interface {
	GetDoc() string
	GetName() string
	GetModules() []ModuleLike
	GetClasses() []*Class
	GetParent() ModuleLike
	SetParent(like ModuleLike)
}

//A type refers to the dot-separated unique name of a class or build-in type, e.g. modulea.moduleb.MyClass
type Type struct {
	//The id of the type which is the FullQualifiedName, like "my.module.MyType" or a base type like "void", "string", "float64" etc.
	FullQualifiedName string `json:"id"`

	//A type may have other generics, like a list or a map.
	Generics []*Type `json:"Generics"`
}

//A module has private and exported classes. Child modules are always exported.
type Module struct {
	//A per parent unique name, something like "mymodule".  May not contain dots (.)
	Name string `json:"name"`
	// Something like "This module contains the logic for handling..."
	Doc string `json:"doc"`

	// A module can contain other (public) child modules.
	Modules []*Module `json:"modules"`;

	// A module can contain public and private classes.
	Classes []*Class `json:"classes"`

	Parent *Module `json:"-"`
}

func (m *Module) AddModule(mod *Module) {
	m.Modules = append(m.Modules, mod)
	mod.Parent = m
}

func (m *Module) AddClass(c *Class) {
	m.Classes = append(m.Classes, c)
	c.Parent = m
}

func (m *Module) GetDoc() string {
	return m.Doc
}

func (m *Module) GetName() string {
	return m.Name
}

func (m *Module) GetModules() []ModuleLike {
	tmp := make([]ModuleLike, len(m.Modules))
	for i := 0; i < len(tmp); i++ {
		tmp[i] = m.Modules[i]
	}
	return tmp
}

func (m *Module) GetClasses() []*Class {
	return m.Classes
}

func (m *Module) GetParent() ModuleLike {
	//otherwise go will create an interface of type ModuleLike with a nil value, which is not what we want
	if m.Parent == nil {
		return nil
	}
	return m.Parent
}

func (m *Module) SetParent(like ModuleLike) {
	switch v := like.(type) {
	case *Module:
		m.Parent = v
	default:
		panic(v)
	}
}

func (a *Module) EachClass(closure func(class *Class)) {
	for _, class := range a.Classes {
		closure(class)
	}
	for _, mod := range a.Modules {
		mod.EachClass(closure);
	}
}

type App struct {
	Module
}

/*
  A class is a blue print to describe specific fields and methods for instances of this type.
  Inheritance is not possible and not planned. You have to use composition using fields.

  TODO introduce interface and polymorphism definition. Think about trait/default implementations.
  TODO think about go-like type composition
 */
type Class struct {
	// A per parent unique name, something like "MyClass". May not contain dots (.)
	Name string `json:"name"`

	// Something like "This class represents the entity ..."
	Doc string `json:"doc"`

	// The visibility of a class. True is for all, otherwise module wise only. Does not make it available to sub module.
	Exported bool `json:"exported"`

	// A class contains a bunch of fields.
	Fields []*Field `json:"fields"`

	// A class contains a bunch of methods
	Methods []*Method `json:"methods"`

	// A class should define stereotypes, so that it's role is always explicit and unambiguous.
	Stereotypes []Stereotype `json:"stereotypes"`

	Parent ModuleLike
}

func (c *Class) AddField(f *Field) {
	c.Fields = append(c.Fields, f)
	f.Parent = c
}

func (c *Class) CreateField(doc string, name string, fullQualifiedTypeName string) {
	f := &Field{}
	f.Name = name
	f.Doc = doc
	f.Type.FullQualifiedName = fullQualifiedTypeName
	f.Exported = true
	c.AddField(f)
}

func (c *Class) AddStereotype(s Stereotype) {
	c.Stereotypes = append(c.Stereotypes, s)
}

func (c *Class) HasStereotype(stereotype Stereotype) bool {
	for _, s := range c.Stereotypes {
		if s == stereotype {
			return true
		}
	}
	return false
}

/*
 A variable is a named storage location containing a data, usually a reference or a pointer to data but for
 some primitives it may represent the value itself. This is the base type for Parameter and fields.
 */
type Variable struct {
	// The unique name of the variable
	Name string `json:"name"`
	// The type of this variable. Not all types make sense, like the void type, and others may be illegal without generics (like list and maps).
	Type Type `json:"type"`

	// The documentation for this field, like "A 'fieldname' is used to represent ..."
	Doc string `json:"doc"`
}

type Parameter struct {
	Variable
	Parent *Method
}

//A field is just a member of an Entity.
type Field struct {
	Variable
	Exported bool `json:"exported"`
	Parent   *Class
}

// Each class can have a method. It may also provide various source code implementations.
type Method struct {
	// Part of the signature. The name of the method.
	Name string `json:"name"`

	// Part of the signature. A method may have multiple parameters, which are just variables.
	Parameters []*Parameter `json:"parameters"`

	// Part of the signature. A method may return multiple result types. The actual capabilities depend on the concrete source implementation.
	// There languages which support such declarations like go and swift, but others do not even support tuples, like java.
	// In those cases the generator inserts local variables and a generated tuple holder which has to be returned.
	Returns []*Parameter `json:"returns"`

	// The concrete implementation, potentially in different languages.
	Implementations []*SourceCode `json:"implementations"`;
	Parent          *Class
}

//One can attach a concrete hand written source code at multiple places, e.g. on {@link Method}s
type SourceCode struct {
	// A constant from {@link PLATFORM}
	Platform string `json:"platform"`

	// The actual source code snipped. Lines are separated by \n
	Code string `json:"code"`

	Parent *Method
}
