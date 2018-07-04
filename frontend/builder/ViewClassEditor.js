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
    showTextInputDialog,
    Span,
    StarBox,
    TextField,
    TwoLineLeadingAndTrailingIcon
} from "/wwt/components.js";


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
     * @param {App} appModel
     * @param {Class|null} entity
     */
    setModel(appModel, entity) {
        this.removeAll();
        if (entity == null) {
            return;
        }

        let grid = new LayoutGrid();

        let etName = new TextField();
        etName.setEnabled(false);
        etName.setText(entity.name);
        grid.add(etName, 12);


        for (let field of entity.fields) {
            let fieldView = new ViewField(this.ctx, appModel.getTypes(false, true, true), field);
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
            this.setModel(appModel, entity);
        });
    }
}

class ViewField extends Box {
    /**
     * @param {UserInterfaceState} ctx
     * @param {Array<AbsType>} types
     * @param {Field} field
     */
    constructor(ctx, types, field) {
        super();
        /**
         *
         * @type {Field}
         */
        this.field = field;

        this.name = new TextField();
        this.name.setCaption(ctx.getString("builder_identifier"));
        this.name.setText(field.name);
        this.name.getElement().style.display = "inline-flex";
        this.add(this.name);

        this.add(new HPadding());

        this.type = new TypeSelector(ctx, types, field);
        this.type.setLabel(ctx.getString("builder_type"));
        this.type.getElement().style.display = "inline-flex";
        this.add(this.type);
    }
}

class TypeSelector extends Select {
    /**
     *
     * @param ctx
     * @param {Array<AbsType>} types
     * @param {Field} field
     */
    constructor(ctx, types, field) {
        super();
        this.types = types;
        for (let type of types) {
            this.add(new SelectModelEntry(type.toString(), type.toString(), false, type.toString() === field.type.toString()));
        }
        this.addOnSelectionChangedListener((select, i, s) => field.absType = this.getSelectedType());
    }

    /**
     *
     * @return {AbsType|null}
     */
    getSelectedType() {
        if (this.getSelectedIndex() < 0) {
            return null;
        }
        return this.types[this.getSelectedIndex()];
    }
}