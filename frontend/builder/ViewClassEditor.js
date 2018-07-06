import {
    Body1,
    Body2,
    Box,
    Card,
    CenterBox,
    CircularProgressIndicator,
    FlatButton,
    FloatingActionButton,
    H4,
    H5,
    H6,
    HPadding,
    HR,
    Icon,
    Image,
    LayoutGrid,
    ListView,
    LRLayout,
    Menu,
    NotFoundException,
    P,
    PullRightBox,
    RaisedButton,
    RoundedIcon,
    Select,
    SelectModelEntry,
    showConfirmationDialog,
    showTextInputDialog,
    Span,
    StarBox,
    TextField,
    TwoLineLeadingAndTrailingIcon
} from "/wwt/components.js";
import {Class, Field, Method, STEREOTYPE, TYPE, Variable,SourceCode} from "./AppModel.js";
import {ViewCodeEditor} from "./ViewCodeEditor.js";


export {ViewClassEditor}

class ViewClassEditor extends Box {
    /**
     *
     * @param {UserInterfaceState} ctx
     */
    constructor(ctx) {
        super();
        this.ctx = ctx;
    }

    /**
     * @param {Observable<App>} observable
     * @param {App} appModel
     * @param {Class|null} clazz
     */
    setModel(observable, appModel, clazz) {
        this.removeAll();
        if (clazz == null) {
            return;
        }

        let grid = new LayoutGrid();


        let nameBox = new CNameField(this.ctx, observable, clazz);
        grid.add(nameBox, 12);

        let stereotype = new StereotypeSelector(this.ctx, observable, clazz);
        stereotype.setLabel(this.ctx.getString("builder_stereotype"));
        grid.add(stereotype, 12);

        grid.add(new HR(), 12);

        for (let field of clazz.fields) {
            let fieldView = new VariableEditor(this.ctx, observable, field);
            fieldView.name.setCaption(this.ctx.getString("builder_field"));
            grid.add(fieldView, 12);
        }

        grid.add(new HR(), 12);

        for (let method of clazz.methods) {
            let methodEditor = new MethodEditor(this.ctx, observable, method);
            grid.add(methodEditor, 12);
        }

        this.add(grid);
    }

    /**
     *
     * @param {Observable<App>}observable
     * @param {Type} type
     */
    bind(observable, type) {
        observable.addObserver(this.getLifecycle(), appModel => {
            let entity = appModel.resolveType(type);
            this.setModel(observable, appModel, entity);
        });
    }
}

class CNameField extends Box {
    constructor(ctx, observable, clazz) {
        super();

        let etName = new TextField();
        etName.getElement().style.display = "inline-flex";
        etName.setEnabled(true);
        etName.setCaption(ctx.getString("builder_class"));
        etName.setText(clazz.name);
        etName.onFocusLostAndChanged(textField => {
            clazz.name = textField.getText();
            observable.notifyValueChanged();

        });


        let btnRemoveClass = new FlatButton();
        btnRemoveClass.getElement().style.display = "inline-flex";
        btnRemoveClass.setIcon(new Icon("remove_circle_outline"));
        btnRemoveClass.setOnClick(_ => {
            showConfirmationDialog(this.uis, ctx.getString("delete_x", clazz.name), ctx.getString("cancel"), ctx.getString("delete"), () => {
                clazz.remove();
                observable.notifyValueChanged();
            });
        });

        let menu = new Menu();
        menu.add(ctx.getString("builder_new_field"), _ => {
            clazz.addField(new Field(clazz.generateFieldName(), TYPE.String));
            observable.notifyValueChanged();
        });


        menu.add(ctx.getString("builder_new_method"), _ => {
            let method = new Method(clazz.generateMethodName());
            clazz.addMethod(method);
            observable.notifyValueChanged();
        });


        let btnAdd = new FlatButton();
        btnAdd.getElement().style.display = "inline-flex";
        btnAdd.setIcon(new Icon("add_circle_outline"));
        btnAdd.setOnClick(_ => {
            menu.popup(btnAdd);
        });


        this.add(etName);
        this.add(btnRemoveClass);
        this.add(btnAdd);
    }
}


class VariableEditor extends Box {
    /**
     * @param {UserInterfaceState} ctx
     * @param {Variable} variable
     * @param {Observable<App>} observable
     */
    constructor(ctx, observable, variable) {
        super();

        this.name = new TextField();
        this.name.setText(variable.name);
        this.name.onFocusLostAndChanged(textField => {
            variable.name = textField.getText();
            observable.notifyValueChanged();
        });

        this.name.getElement().style.display = "inline-flex";

        this.add(this.name);

        this.add(new HPadding());
        this.type = new GenericPicker(ctx, observable.getValue().getTypes(false, true, true), variable.type);
        this.type.setOnTypeSelectedListener((select, v) => {
            variable.type = v;
            observable.notifyValueChanged();
        });
        this.type.setLabel(ctx.getString("builder_type"));
        this.type.getElement().style.display = "inline-flex";
        this.add(this.type);

        this.btnRemoveField = new FlatButton();
        this.btnRemoveField.setIcon(new Icon("remove_circle_outline"));
        this.btnRemoveField.setOnClick(_ => {
            showConfirmationDialog(this.uis, ctx.getString("delete_x", variable.name), ctx.getString("cancel"), ctx.getString("delete"), () => {
                variable.remove();
                observable.notifyValueChanged();
            });
        });
        this.add(this.btnRemoveField);
    }
}


class StereotypeSelector extends Select {
    /**
     *
     * @param ctx
     * @param {Class} clazz
     * @param {Observable<App>} observable
     */
    constructor(ctx, observable, clazz) {
        super();
        this.add(new SelectModelEntry(ctx.getString("builder_undefined"), " "));
        for (let type of STEREOTYPE.VALUES()) {
            this.add(new SelectModelEntry(type, type, false, clazz.hasStereotype(type)));
        }
        this.addOnSelectionChangedListener((select, i, s) => {
            clazz.setStereotype(this.getSelectedType());
            observable.notifyValueChanged();
        });
    }

    /**
     *
     * @return {String|null}
     */
    getSelectedType() {
        if (this.getSelectedIndex() < 0) {
            return null;
        }
        return STEREOTYPE.VALUES()[this.getSelectedIndex() - 1];
    }
}


class MethodEditor extends Box {
    /**
     *
     * @param {UserInterfaceState} ctx
     * @param {Observable<App>} observable
     * @param {Method} method
     */
    constructor(ctx, observable, method) {
        super();

        this.name = new TextField();
        this.name.setCaption(ctx.getString("builder_method"));
        this.name.setText(method.name);
        this.name.getElement().style.display = "inline-flex";
        this.name.onFocusLostAndChanged(textField => {
            method.name = textField.getText();
            observable.notifyValueChanged();
        });


        this.btnRemoveField = new FlatButton();
        this.btnRemoveField.setIcon(new Icon("remove_circle_outline"));
        this.btnRemoveField.setOnClick(_ => {
            showConfirmationDialog(this.uis, ctx.getString("delete_x", method.name), ctx.getString("cancel"), ctx.getString("delete"), () => {
                method.remove();
                observable.notifyValueChanged();
            });
        });

        let menu = new Menu();
        menu.add(ctx.getString("builder_new_method_parameter"), _ => {
            method.addParameter(new Variable(method.generateParameterName(), TYPE.String));
            observable.notifyValueChanged();
        });


        menu.add(ctx.getString("builder_new_method_result"), _ => {
            method.addReturn(new Variable(method.generateReturnName(), TYPE.String));
            observable.notifyValueChanged();
        });


        let btnAdd = new FlatButton();
        btnAdd.getElement().style.display = "inline-flex";
        btnAdd.setIcon(new Icon("add_circle_outline"));
        btnAdd.setOnClick(_ => {
            menu.popup(btnAdd);
        });

        this.add(this.name);
        this.add(this.btnRemoveField);
        this.add(btnAdd);


        for (let p of method.parameters) {
            let varBox = new VariableBox(ctx, observable, p);
            varBox.ico.setName("arrow_forward");
            this.add(varBox);
        }

        for (let r of method.returns) {
            let varBox = new VariableBox(ctx, observable, r);
            varBox.ico.setName("arrow_back");
            this.add(varBox);
        }

        let clazz = method.parent;
        let editor = new ViewCodeEditor();
        if (method.implementations.length > 0) {
            editor.setSourceCode(method.implementations[0].code);
        }
        editor.mainElement.addEventListener("focusout", _ => {

            if (clazz.stereotypes.length > 0) {
                let platform = observable.getValue().getPlatformForStereotype(clazz.stereotypes[0]);
                editor.getSourceCode().then(text => {
                    method.putImplementation(new SourceCode(platform, text));
                    observable.notifyValueChanged();
                });

            }

        });
        this.add(editor);
    }
}

class VariableBox extends Box {
    /**
     * @param {UserInterfaceState} ctx
     * @param {Variable} variable
     * @param {Observable<App>} observable
     */
    constructor(ctx, observable, variable) {
        super();
        this.ico = new Icon("");
        this.ico.getElement().style.marginRight = "8px";
        this.ico.getElement().style.marginBottom = "12px";
        this.ico.getElement().style.verticalAlign = "bottom";
        this.add(this.ico);
        let variableEditor = new VariableEditor(ctx, observable, variable);
        variableEditor.getElement().style.display = "inline-block";
        this.add(variableEditor);

    }
}


/**
 * @template {T} the actual picker type
 */
class GenericPicker extends Select {
    /**
     *
     * @param ctx
     * @param {Array<T>} values
     * @param {T|null} defaultValue
     */
    constructor(ctx, values, defaultValue) {
        super();
        this.values = values;
        this.callback = null;
        this.add(new SelectModelEntry(ctx.getString("builder_undefined"), " "));
        for (let v of this.values) {
            this.add(new SelectModelEntry(v.toString(), v.toString(), false, defaultValue != null && defaultValue.toString() === v.toString()));
        }
        this.addOnSelectionChangedListener((select, i, s) => {
            if (this.callback != null) {
                this.callback(this, this.getSelectedValue());
            }
        });
    }

    /**
     *
     * @return {string|null}
     */
    getSelectedValue() {
        if (this.getSelectedIndex() < 0) {
            return null;
        }
        return this.values[this.getSelectedIndex() - 1];
    }

    /**
     *
     * @param {function(Select,T)} callback
     */
    setOnTypeSelectedListener(callback) {
        this.callback = callback;
    }
}