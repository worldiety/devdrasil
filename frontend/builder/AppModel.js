export {AppModel, AbsPlatform, Class, Entity, View, ViewModel}

/**
 * The receipt to create apps or plugins.
 */
class AppModel {

    /**
     *
     * @param {AbsPlatform} frontend. Which platform is in the backend?
     * @param {AbsPlatform} backend. Which platform is in the frontend?
     * @param {string} appId. A unique something like e.g. "com.mydomain.myplugin"
     * @param {string} appDoc. A description of the app e.g. "This plugin models a build server which provides the following ..."
     * @param {Array<Entity>} entities, which are custom data types, e.g. like a "car" or a "motor"
     * @param {Array<Repository>} repositories. Each repository refers to exact to one entity and provides CRUD functionalities.
     */
    constructor(appId, appDoc, backend, frontend, entities = [], repositories = []) {
        this.appId = appId;
        this.appDoc = appDoc;
        this.entities = entities;
        this.repositories = repositories;
    }


    /**
     * Checks if the name does exist in any (global) configuration.
     * @param {string} str
     * @return {boolean}
     */
    nameExists(str) {
        for (let e of this.entities) {
            if (e.name === str) {
                return true;
            }
        }
        return false;
    }
}

/**
 * A platform denotes a language and framework combination. Most languages are platforms anyway, however there may be multiple platforms for each language also.
 */
class AbsPlatform {

}


/**
 * The base for every generated class or struct is this class. It is language agnostic with regards to it's fields
 * but may provide methods of a concrete implementation.
 */
class Class {
    /**
     *
     * @param {string} name. The unique name of a field. Must start with uppercase letter and follow the convention of a public Go struct (camel case)
     * @param {Array<Field>} fields. The fields of this entity.
     * @param {Array<AbsMethod>} methods. The list of methods.
     */
    constructor(name = "", fields = [], methods = []) {
        this.name = name;
        this.fields = fields;
        this.methods = methods;
    }
}


/**
 * A custom data type, which consists of fields which are either build-in types or other entities. References to other Entities are only possible if a Repository is also defined.
 */
class Entity extends Class {
    /**
     *
     * @param {string} name. The unique name of a field. Must start with uppercase letter and follow the convention of a public Go struct (camel case)
     * @param {Array<Field>} fields. The fields of this entity.
     */
    constructor(name = "", fields = []) {
        super(name, fields);
    }
}

/**
 * A field is just a member of an Entity.
 */
class Field {
    /**
     *
     * @param {string} name. The (entity) unique name of this field. Must start with lowercase letter and camel case
     * @param {AbsType} absType. An instance of a concrete descendant of an AbsType.
     */
    constructor(name, absType) {
        this.name = name;
        this.absType = absType;
    }
}

/**
 * Just a field, but it will get autowired from the context. Usually only works for members of a class.
 */
class AutowiredField extends Field {
    constructor(name, absType) {
        super(name, absType);
    }
}

/**
 * An abstract source definition
 */
class AbsSource {
    /**
     * @param {string} source. The actual source code in whatever language
     */
    constructor(source = "") {
        this.source = source;
        //autogenerated source means, that the user may not change the implementation on his own
        this.autogenerated = false;
    }
}

/**
 * An abstract method definition with source
 */
class AbsMethod extends AbsSource {
    /**
     * @param {Array<Field>}parameter
     * @param {AbsType} result. The result of the method
     * @param {string} doc. The documentation of this method.
     * @param {string} source. The actual source code in whatever language
     */
    constructor(parameter = [], result = new VoidType(), doc = "", source = "") {
        super(source);
        this.parameter = parameter;
        this.result = result;
        this.doc = doc;
        this.public = true;
    }
}

/**
 * Represents a method with go(lang) source code
 */
class GoMethodSource extends AbsMethod {
    /**
     * @param {Array<Field>}parameter
     * @param {AbsType} result. The result of the method
     * @param {string} doc. The documentation of this method.
     * @param {string} source. The actual source code in Go(lang)
     */
    constructor(parameter = [], result = new VoidType(), doc = "", source = "") {
        super(parameter, result, doc, source);
        this.source = source;
    }
}

/**
 * Represents a method with javascript source code
 */
class JSMethodSource extends AbsMethod {
    /**
     * @param {Array<Field>}parameter
     * @param {AbsType} result. The result of the method
     * @param {string} doc. The documentation of this method.
     * @param {string} source. The actual source code in Go(lang)
     */
    constructor(parameter = [], result = new VoidType(), doc = "", source = "") {
        super(parameter, result, doc, source);
        this.source = source;
    }
}

/**
 * The non-instantiable super class of all field types
 */
class AbsType {

}

/**
 * A list type wraps another type
 */
class ListType extends AbsType {
    /**
     *
     * @param {AbsType} boxType
     */
    constructor(boxType) {
        super();
        this.boxType = boxType;
    }
}

/**
 * The VoidType represents void, which is only applicable for result values of
 */
class VoidType extends AbsType {
}

class StringType extends AbsType {
}

class Int64Type extends AbsType {
}

class Float64Type extends AbsType {
}

class BoolType extends AbsType {
}

class EntityType extends AbsType {
    /**
     *
     * @param {Entity} entity
     */
    constructor(entity) {
        super();
        this.entity = entity;
    }
}

/**
 * A repository just refers to a bunch of persisted entities all of the same type. It is always a singleton in the backend (e.g. Go) and frontend meaning (e.g. Javascript).
 * It is distinguishable from a service only because it does not contain any other things than working with that kind of entity.
 *
 * Stereotype of the persistence layer.
 */
class Repository extends Class {
    /**
     *
     * @param{string} name. Must be unique across all names
     * @param {Entity} entity. Instance of an entity, which should be organized in a repo.
     * @param {Array<AbsMethod>} methods, used to declare e.g. CRUD methods in the server side repository. So this must of the backend source type. Also used to declare custom filter/query methods. Other methods are accessible by using "this" (whatever that means in a certain language).
     */
    constructor(name, entity, methods = []) {
        super(name, []);
        this.name = name;
        this.entity = entity;
    }

}


/**
 * A service is a singleton in the backend and forwarded to the client, just like a repository but entirely a custom implementation.
 *
 * Stereotype of the service layer.
 */
class Service extends Class {
    /**
     * @see Class
     * @param {string} name
     * @param {Array<Field>} fields
     * @param {Array<AbsMethod>} methods
     * @param {Array<Service|Repository>} injections
     */
    constructor(name, fields, methods, injections) {
        super(name, fields, methods);
        this.injections = injections;
    }
}

/**
 * A view model is like a view controller with observable patterns in most methods.
 */
class ViewModel extends Class {
    /**
     * @see Class.constructor
     * @param {string} name
     * @param {Array<Field>} fields
     * @param {Array<AbsMethod>} methods
     * @param {Array<Service|Repository>} injections
     */
    constructor(name, fields, methods, injections) {
        super(name, fields, methods);
        this.injections = injections;
    }
}

/**
 * A view is always a client side concrete implementation and interacts with a specific view model (view controller), just like in MVVM thinking.
 * You implement methods in your view model and call the async observable methods. Registering to observables is only possible with a given lifecycle.
 */
class View extends Class {
    /**
     * @see Class.constructor
     * @param {string} name
     * @param {ViewModel} viewModel
     */
    constructor(name, viewModel) {
        super(name, [], []);
    }
}