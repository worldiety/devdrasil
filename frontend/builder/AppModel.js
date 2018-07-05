export {
    Method,
    SourceCode,
    Variable,
    Type,
    TYPE,
    Field,
    Module,
    PLATFORM,
    App,
    Class,
    STEREOTYPE,
    TYPE_FACTORY
}


/**
 * A type refers to the dot-separated unique name of a class or build-in type, e.g. modulea.moduleb.MyClass
 */
class Type {
    /**
     *
     * @param {string} fullQualifiedName
     */
    constructor(fullQualifiedName = "") {
        /**
         * The id of the type which is the FullQualifiedName, like "my.module.MyType" or a base type like "void", "string", "float64" etc.
         *
         * @type {string}
         */
        this.fullQualifiedName = fullQualifiedName;
        /**
         * A type may have other generics, like a list or a map.
         *
         * @type {Array<Type>}
         */
        this.generics = [];
    }

    /**
     * @return {Object}
     */
    toObject() {
        let obj = {};
        obj["id"] = this.fullQualifiedName;
        obj["generics"] = [];
        for (let g of this.generics) {
            obj["generics"].push(g.toObject());
        }
        return obj;
    }

    /**
     *
     * @param {Object} obj
     */
    fromObject(obj) {
        this.fullQualifiedName = obj["id"];
        this.generics = [];
        for (let o of obj["generics"]) {
            let g = new Type();
            g.fromObject(o);
            this.generics.push(g);
        }
    }

    toString() {
        if (this.generics.length === 0) {
            return this.fullQualifiedName;
        } else {
            return this.fullQualifiedName + "<" + this.generics.join(",") + ">";
        }
    }
}

/**
 * A module has private and exported classes. Child modules are always exported.
 */
class Module {
    constructor() {
        /**
         * A per parent unique name, something like "mymodule".  May not contain dots (.)
         * @type {string}
         */
        this.name = "";
        /**
         * Something like "This module contains the logic for handling..."
         * @type {string}
         */
        this.doc = "";

        /**
         * A module can contain other (public) child modules.
         *
         *
         * @type {Array<Module>}
         */
        this.modules = [];

        /**
         * A module can contain public and private classes.
         * @type {Array<Class>}
         */
        this.classes = [];

        /**
         * The parent of this module is always a module or null, if it is in the root module.
         * @type {Module|null}
         */
        this.parent = null;
    }

    /**
     * Returns a child module by name. Only finds direct children.
     * @param name
     * @return {Module|null}
     */
    getModule(name) {
        for (let mod of this.modules) {
            if (mod.name === name) {
                return mod;
            }
        }
        return null;
    }

    /**
     * Returns a child class by name. Only finds direct children.
     * @param name
     * @return {Class|null}
     */
    getClass(name) {
        for (let cl of this.classes) {
            if (cl.name === name) {
                return cl;
            }
        }
        return null;
    }

    /**
     * We want full control over serialization and deserialization process.
     * @return {Object}
     */
    toObject() {
        let json = {};
        json["name"] = this.name;
        json["doc"] = this.doc;
        json["modules"] = [];
        for (let mod of this.modules) {
            json["modules"].push(mod.toObject());
        }
        json["classes"] = [];
        for (let cl of this.classes) {
            json["classes"].push(cl.toObject());
        }
        return json;
    }

    /**
     * We want custom types and restore cycles. So we need a custom deserialize method.
     * @param {Object} obj
     */
    fromObject(obj) {
        this.name = obj["name"];
        this.doc = obj["doc"];
        this.modules = [];
        for (let o of obj["modules"]) {
            let mod = new Module();
            mod.fromObject(o);
            this.modules.push(mod);
        }

        this.classes = [];
        for (let o of obj["classes"]) {
            let cl = new Class();
            cl.fromObject(o);
            cl.parent = this;
            this.classes.push(cl);
        }
    }


    /**
     * loops over all classes defined by this module and it's children modules.
     * @param {function(Class)} closure
     */
    forEachClass(closure) {
        for (let cl of this.classes) {
            closure(cl);
        }

        for (let mod of this.modules) {
            mod.forEachClass(closure);
        }
    }

    /**
     * Just like {@link #forEachClass}, but returns a list of classes.
     *
     * @return {Array<Class>}
     */
    getClasses() {
        let tmp = [];
        this.forEachClass(cl => tmp.push(cl));
        return tmp;
    }


    /**
     * Adds the given class. If another class with that name already exists, the existing one is removed.
     * @param {Class} clazz
     */
    addClass(clazz) {
        this.classes = this.classes.filter(value => value.name !== clazz.name);
        this.classes.push(clazz);
        clazz.parent = this;
    }
}

/**
 * The receipt to create apps or plugins. An app is basically just a {@link Module}.
 */
class App extends Module {
    constructor() {
        super();
    }

    /**
     * Checks if the name does exist in any (global) configuration.
     * @param {string} str
     * @return {boolean}
     */
    nameExists(str) {
        for (let e of this.classes) {
            if (e.name === str) {
                return true;
            }
        }
        return false;
    }


    /**
     * Returns all available types, including the basic ones and all self defined.
     * @param {boolean} includeVoid, if true the void type is included
     * @param {boolean} includeLists, if true, generates List-Types with all
     * @param {boolean} includeMap, if true, includes maps with string as key and all other types as values.
     * @return {Array<Type>}
     */
    getTypes(includeVoid = false, includeLists = false, includeMap = false) {
        let res = [];


        for (let key in TYPE) {
            if (!includeVoid && TYPE[key].fullQualifiedName === TYPE.Void.fullQualifiedName) {
                continue;
            }
            res.push(TYPE[key]);
        }

        this.forEachClass(cl => res.push(cl.asType()));

        let tmp = [...res];
        if (includeLists) {
            for (let type of tmp) {
                if (type.id === TYPE.Void.fullQualifiedName) {
                    continue;
                }
                res.push(TYPE_FACTORY.NewList(type));
            }
        }

        if (includeMap) {
            for (let type of tmp) {
                if (type.id === TYPE.Void.fullQualifiedName) {
                    continue;
                }
                res.push(TYPE_FACTORY.NewMap(TYPE.String, type));
            }
        }


        return res;
    }

    /**
     * Tries to resolve a class in the context of the given app.
     *
     * @param {Type} type
     * @return {Class|null}
     */
    resolveType(type) {
        let parent = this;
        let tokens = type.fullQualifiedName.split(".");
        for (let i = 0; i < tokens.length - 1; i++) {
            parent = parent.getModule(tokens[i]);
            if (parent == null) {
                return null;
            }
        }
        let className = tokens[tokens.length - 1];
        return parent.getClass(className);
    }

    /**
     * Checks if the given type is a build-in type.
     * @param {Type} type
     */
    isBuildInType(type) {
        for (let bin of TYPE) {
            if (type.fullQualifiedName === bin.id) {
                return true;
            }
        }
        return false;
    }


}


/**
 * A class is a blue print to describe specific fields and methods for instances of this type.
 * Inheritance is not possible and not planned. You have to use composition using fields.
 *
 * TODO introduce interface and polymorphism definition. Think about trait/default implementations.
 * TODO think about go-like type composition
 */
class Class {
    constructor(name = "") {
        /**
         * A per parent unique name, something like "MyClass". May not contain dots (.)
         * @type {string}
         */
        this.name = name;

        /**
         * Something like "This class represents the entity ..."
         * @type {string}
         */
        this.doc = "";

        /**
         * The parent of this class is always a module or null, if it is in the root module.
         * @type {Module|null}
         */
        this.parent = null;

        /**
         * The visibility of a class. True is for all, otherwise module wise only. Does not make it available to sub module.
         * @type {boolean}
         */
        this.exported = true;

        /**
         * A class contains a bunch of fields.
         * @type {Array<Field>}
         */
        this.fields = [];

        /**
         * A class contains a bunch of methods
         * @type {Array<Method>}
         */
        this.methods = [];

        /**
         * A class should define stereotypes, so that it's role is always explicit and unambiguous.
         * @type {Array<string>}
         */
        this.stereotypes = [];
    }

    /**
     * checks if a stereotype has been assigned
     * @param {string} stereotype
     * @return {boolean}
     */
    hasStereotype(stereotype) {
        return this.stereotypes.indexOf(stereotype) >= 0;
    }

    /**
     * Adds the stereotype and removes all others
     * @param {string} stereotype
     */
    setStereotype(stereotype) {
        this.stereotypes = [];
        this.stereotypes.push(stereotype);
    }

    /**
     * @return {Object}
     */
    toObject() {
        let obj = {};
        obj["name"] = this.name;
        obj["doc"] = this.doc;
        obj["exported"] = this.exported;
        obj["fields"] = [];
        for (let f of this.fields) {
            obj["fields"].push(f.toObject());
        }

        obj["methods"] = [];
        for (let m of this.methods) {
            obj["methods"].push(m.toObject());
        }
        obj["stereotypes"] = this.stereotypes;
        return obj;
    }

    /**
     *
     * @param {Object} obj
     */
    fromObject(obj) {
        this.name = obj["name"];
        this.doc = obj["doc"];
        this.exported = obj["exported"];
        this.fields = [];
        for (let o of obj["fields"]) {
            let field = new Field();
            field.fromObject(o);
            field.parent = this;
            this.fields.push(field);
        }

        this.methods = [];
        for (let m of obj["methods"]) {
            let meth = new Method();
            meth.fromObject(m);
            meth.parent = this;
            this.methods.push(meth);
        }
        this.stereotypes = obj["stereotypes"]
    }


    /**
     * Returns this class as type
     * @return {Type}
     */
    asType() {
        let fullQualifiedName = this.name;
        let root = this.parent;
        while (root != null) {
            if (root.name.length > 0) {
                fullQualifiedName = root.name + "." + fullQualifiedName;
            }
            root = root.parent;
        }
        return new Type(fullQualifiedName);
    }

    /**
     * Removes this class from it's parent module.
     * @return {boolean}, true if successful
     */
    remove() {
        if (this.parent == null) {
            return false;
        }
        this.parent.classes = this.parent.classes.filter(value => value !== this);
        this.parent = null;
        return true;
    }
}

/**
 * The base types, which are always available. Map and List should only occur as concrete instances. Do not modify them.
 * @type {{Void: Type, String: Type, Int64: Type, Float64: Type, Bool: Type, List: Type, Map: Type}}
 */
const TYPE = {
    Void: new Type("void"),
    String: new Type("string"),
    Int64: new Type("int64"),
    Float64: new Type("float64"),
    Bool: new Type("bool"),
    List: new Type("List"),
    Map: new Type("Map"),
};

const TYPE_FACTORY = {
    /**
     * Creates a new list type
     * @param type
     * @return {Type}
     * @constructor
     */
    NewList: function (type) {
        let r = new Type("List");
        r.generics.push(type);
        return r;
    },
    /**
     * Creates a new Map
     * @param keyType
     * @param valueType
     * @return {Type}
     * @constructor
     */
    NewMap: function (keyType, valueType) {
        let r = new Type("Map");
        r.generics.push(keyType);
        r.generics.push(valueType);
        return r;
    },
};


/**
 * Walks up the module hierarchy to get the top most module, which is usually the {@link App}.
 * @param {Module} module
 * @return {Module}
 */
function getRootModule(module) {
    let root = module;
    while (root.parent != null) {
        root = root.parent;
    }
    return root;
}

class Variable {
    /**
     *
     * @param {string} name. The (entity) unique name of this field. Must start with lowercase letter and camel case
     * @param {Type} type. A specific type which describes this field.
     */
    constructor(name = "", type = null) {
        /**
         * The unique name of the variable
         * @type {string}
         */
        this.name = name;
        /**
         * The type of this variable. Not all types make sense, like the void type, and others may be illegal without generics (like list and maps).
         * @type {Type}
         */
        this.type = type;

        /**
         * The documentation for this field, like "A 'fieldname' is used to represent ..."
         *
         * @type {string}
         */
        this.doc = "";

        /**
         * A variable may have different kind of parents, depending on the concrete kind of variable. E.g. a field always has a {@link Class} as a parent.
         * @type {Method|Class}
         */
        this.parent = null;
    }

    /**
     * @return {Object}
     */
    toObject() {
        let obj = {};
        obj["name"] = this.name;
        if (this.type != null) {
            obj["type"] = this.type.toObject();
        }
        obj["doc"] = this.doc;
        return obj;
    }

    /**
     *
     * @param {Object} obj
     */
    fromObject(obj) {
        this.name = obj["name"];
        this.doc = obj["doc"];
        this.type = null;
        if (obj["type"] != null) {
            this.type = new Type();
            this.type.fromObject(obj["type"]);
        }
    }
}

/**
 * A field is just a member of an Entity.
 */
class Field extends Variable {
    /**
     *
     * @param {string} name. The (entity) unique name of this field. Must start with lowercase letter and camel case
     * @param {Type} type. A specific type which describes this field.
     */
    constructor(name = "", type = null) {
        super(name, type);

        /**
         * Just like a class can be exported from a module, fields can be as well.
         *
         * @type {boolean}
         */
        this.exported = true;

    }

    toObject() {
        let obj = super.toObject();
        obj["exported"] = this.exported;
        return obj;
    }

    fromObject(obj) {
        super.fromObject(obj);
        this.exported = obj["exported"];
    }

    /**
     * Removes this field from it's parent
     * @return {boolean}, true if successful
     */
    remove() {
        if (this.parent == null) {
            return false;
        }
        this.parent.fields = this.parent.fields.filter(value => value !== this);
        this.parent = null;
        return true;
    }
}


/**
 * Each class can have a method. It may also provide various source code implementations.
 */
class Method {
    constructor(name = "") {
        /**
         * Part of the signature. The name of the method.
         * @type {string}
         */
        this.name = name;

        /**
         * Part of the signature. A method may have multiple parameters, which are just variables.
         * @type {Array<Variable>}
         */
        this.parameters = [];

        /**
         * Part of the signature. A method may return multiple result types. The actual capabilities depend on the concrete source implementation.
         * @type {Array<Type>}
         */
        this.returns = [];

        /**
         * The concrete implementation, potentially in different languages.
         *
         * @type {Array<SourceCode>}
         */
        this.implementations = [];

        /**
         * Currently only class can have methods.
         * @type {Class}
         */
        this.parent = null;
    }

    /**
     * @return {Object}
     */
    toObject() {
        let obj = {};
        obj["name"] = this.name;
        obj["parameters"] = [];
        for (let p of this.parameters) {
            obj["parameters"].push(p.toObject());
        }
        obj["returns"] = [];
        for (let r of this.returns) {
            obj["returns"].push(r.toObject());
        }
        obj["implementations"] = [];
        for (let i of this.implementations) {
            obj["implementations"].push(i.toObject());
        }
        return obj;
    }

    /**
     *
     * @param {Object} obj
     */
    fromObject(obj) {
        this.code = obj["code"];
        this.platform = obj["platform"];
        this.parameters = [];
        this.returns = [];
        this.implementations = [];

        for (let o of obj["parameters"]) {
            let p = new Variable();
            p.fromObject(o);
            p.parent = this;
            this.parameters.push(p);
        }

        for (let o of obj["returns"]) {
            let t = new Type();
            t.fromObject(o);
            this.returns.push(t);
        }

        for (let o of obj["implementations"]) {
            let s = new SourceCode();
            s.fromObject(o);
            s.parent = this;
            this.implementations.push(s);
        }
    }
}

/**
 * Source code can be created for various platforms. Why platform? Because there is no useful language without any base
 * sdk. So when you talk about a language you usually mean a language specification, a concrete implementation
 * and a SDK.
 *
 * @type {{GO_1_x: string, ES6: string}}
 */
const PLATFORM = {
    /**
     * Go 1.x with gc. At least go 1.11 is required.
     */
    GO_1_x: "Go 1.x",

    /**
     * ES6 or compatible is required.
     */
    ES6: "ES6"
};

const STEREOTYPE = {
    /**
     * Frontend: Represents a view, e.g. like a form. It is the smallest and most reusable unit in the frontend.
     * It's lifecycle is usually at most to the lifecycle of a {@link STEREOTYPE.USER_INTERFACE_STATE}
     */
    VIEW: "VIEW",
    /**
     * Backend: Represents a controller, which is a singleton. It incorporates the actual server side application logic.
     */
    CONTROLLER: "CONTROLLER",
    /**
     * Frontend: In the MVVM world, this is the ViewModel and provides a (potential) two way binding to a {@link STEREOTYPE.VIEW}.
     * On the other hand, it just may be a simple controller which is used directly (MVC world). In both cases it is a singleton
     * at application level.
     */
    VIEW_CONTROLLER: "VIEW_CONTROLLER",

    /**
     * Frontend: Represents a logical state which allows a forward/backward navigation. It uses typically {@link STEREOTYPE.VIEW}s and applies them.
     * It provides a basic create/destroy lifecycle and provides the infrastructure for MVVM or other callback bindings to resolve leaks reliable.
     */
    USER_INTERFACE_STATE: "USER_INTERFACE_STATE",

    /**
     * Frontend & Backend: Represents a model which is persistent and supports CRUD (create read update delete). Each model
     * is made available by a repository. The backend holds the actual truth, whereas the frontend repository just
     * represents the client's end, perhaps providing additional caching purposes. A {@link STEREOTYPE.VIEW_CONTROLLER} typically imports
     * a repository to make data available to the view.
     */
    PERSISTENCE_MODEL: "PERSISTENCE_MODEL",

    /**
     *Frontend & Backend: An abstract component, which encapsulates a model with optional methods. It can be used to either
     * represent common intermediate models (without persistence) or separation of concerns (e.g. splitted controller logic).
     * It is shared across the frontend and backend. See also {@link STEREOTYPE.FRONTEND_COMPONENT} and {@link STEREOTYPE.BACKEND_COMPONENT}.
     */
    COMPONENT: "COMPONENT",

    /**
     * Frontend: a component which is only applicable to the frontend. See also {@link STEREOTYPE.COMPONENT}.
     */
    FRONTEND_COMPONENT: "FRONTEND_COMPONENT",

    /**
     * Backend: a component which is only applicable to the backend. See also {@link STEREOTYPE.COMPONENT}.
     */
    BACKEND_COMPONENT: "BACKEND_COMPONENT",

    VALUES: function () {
        return [STEREOTYPE.VIEW, STEREOTYPE.CONTROLLER]
    }
};

/**
 * One can attach a concrete hand written source code at multiple places, e.g. on {@link Method}s
 */
class SourceCode {
    constructor(platform = "", code = "") {
        /**
         * A constant from {@link PLATFORM}
         *
         * @type {string}
         */
        this.platform = platform;

        /**
         * The actual source code snipped. Lines are separated by \n
         * @type {string}
         */
        this.code = code;


        /**
         * Currently only a method may have exactly 1:1 relation an implementation
         * @type {Method}
         */
        this.parent = null;
    }

    /**
     * Returns the code lines as string array
     * @return {string[]}
     */
    getLines() {
        return this.code.split("\n")
    }

    /**
     * Sets the code as lines separated by line break \n
     * @param {string[]} lines
     */
    setLines(lines) {
        this.code = lines.join("\n");
    }

    /**
     * @return {Object}
     */
    toObject() {
        let obj = {};
        obj["platform"] = this.platform;
        obj["code"] = this.code;
        return obj;
    }

    /**
     *
     * @param {Object} obj
     */
    fromObject(obj) {
        this.code = obj["code"];
        this.platform = obj["platform"];
    }
}

