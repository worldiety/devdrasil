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
    showTextInputDialog,
    Span,
    StarBox,
    TextField,
    TwoLineLeadingAndTrailingIcon
} from "/wwt/components.js";


export {ViewEntity}

class ViewEntity extends Box {
    /**
     *
     * @param {UserInterfaceState} ctx
     */
    constructor(ctx) {
        super();
        this.ctx = ctx;
    }

    /**
     *
     * @param {Entity|null} entity
     */
    setModel(entity) {
        this.removeAll();
        if (entity == null) {
            return;
        }

        let grid = new LayoutGrid();

        let etName = new TextField();
        etName.setEnabled(false);
        etName.setText(entity.name);
        grid.add(etName, 12);


        this.add(grid);
    }

    /**
     *
     * @param {Observable<AppModel>}observable
     * @param {string} entityName
     */
    bind(observable, entityName) {
        observable.addObserver(this.getLifecycle(), appModel => {
            let entity = appModel.findEntity(entityName);
            this.setModel(entity);
        });
    }
}