import {
    Body1,
    Body2,
    Box,
    Card,
    CenterBox,
    CircularProgressIndicator,
    FlatButton,
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
import {STEREOTYPE} from "./AppModel.js";


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


        for (let field of clazz.fields) {
            let fieldView = new ViewField(this.ctx, appModel.getTypes(false, true, true), observable, field);
            grid.add(fieldView, 12);
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
        btnRemoveClass.setIcon(new Icon("remove"));
        btnRemoveClass.setOnClick(_=>{
            showConfirmationDialog(this.uis, ctx.getString("delete_x", clazz.name), ctx.getString("cancel"), ctx.getString("delete"), () => {
                clazz.remove();
                observable.notifyValueChanged();
            });
        });


        this.add(etName);
        this.add(btnRemoveClass);
    }
}

class ViewField extends Box {
    /**
     * @param {UserInterfaceState} ctx
     * @param {Array<Type>} types
     * @param {Field} field
     * @param {Observable<App>} observable
     */
    constructor(ctx, types, observable, field) {
        super();
        /**
         *
         * @type {Field}
         */
        this.field = field;

        this.name = new TextField();
        this.name.setCaption(ctx.getString("builder_field"));
        this.name.setText(field.name);
        this.name.getElement().style.display = "inline-flex";
        this.name.onFocusLostAndChanged(textField => {
            field.name = textField.getText();
            observable.notifyValueChanged();
        });
        this.add(this.name);

        this.add(new HPadding());

        this.type = new TypeSelector(ctx, types, observable, field);
        this.type.setLabel(ctx.getString("builder_type"));
        this.type.getElement().style.display = "inline-flex";
        this.add(this.type);

        this.btnRemoveField = new FlatButton();
        this.btnRemoveField.setIcon(new Icon("remove"));
        this.btnRemoveField.setOnClick(_ => {
            showConfirmationDialog(this.uis, ctx.getString("delete_x", field.name), ctx.getString("cancel"), ctx.getString("delete"), () => {
                field.remove();
                observable.notifyValueChanged();
            });
        });
        this.add(this.btnRemoveField);
    }
}

class TypeSelector extends Select {
    /**
     *
     * @param ctx
     * @param {Array<Type>} types
     * @param {Field} field
     * @param {Observable<App>} observable
     */
    constructor(ctx, types, observable, field) {
        super();
        this.types = types;
        this.add(new SelectModelEntry(ctx.getString("builder_undefined"), ""));
        for (let type of types) {
            this.add(new SelectModelEntry(type.toString(), type.toString(), false, type.toString() === field.type.toString()));
        }
        this.addOnSelectionChangedListener((select, i, s) => {
            field.type = this.getSelectedType();
            observable.notifyValueChanged();
        });
    }

    /**
     *
     * @return {Type|null}
     */
    getSelectedType() {
        if (this.getSelectedIndex() < 0) {
            return null;
        }
        return this.types[this.getSelectedIndex() - 1];
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
        this.add(new SelectModelEntry(ctx.getString("builder_undefined"), ""));
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


