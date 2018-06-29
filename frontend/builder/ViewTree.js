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
import {Entity, View, ViewModel} from "./AppModel.js";


export {ViewTree}

class ViewTree extends Box {
    /**
     * @param {UserInterfaceState} ctx
     */
    constructor(ctx) {
        super();
        this.ctx = ctx;
        this.onSelectedListener = null;
    }

    /**
     *
     * @param {AppModel} model
     * @param {Observable<AppModel>}observable
     */
    setModel(model, observable) {
        this.removeAll();

        let toolbar = new Toolbar(this.ctx, observable);
        this.add(toolbar);

        this.listView = new ListView();
        this.listView.setInteractive(true);
        for (let view of model.entities) {
            let entry = this.listView.add(view.name, true);
            entry.onclick = evt => {
                let idx = this.listView.indexOf(evt.target);
                this.listView.setSelected(idx);

                if (this.onSelectedListener != null) {
                    this.onSelectedListener(view);
                }
            }
        }

        this.add(this.listView);
    }

    /**
     *
     * @param {Observable<AppModel>}observable
     */
    bind(observable) {
        observable.addObserver(this.getLifecycle(), model => this.setModel(model, observable));
    }


    /**
     *
     * @param {function(Class)} onSelected
     */
    setOnSelectedListener(onSelected) {
        this.onSelectedListener = onSelected;
    }
}


class Toolbar extends Box {
    /**
     * @param {UserInterfaceState} ctx
     * @param {Observable<AppModel>} observable
     */
    constructor(ctx, observable) {
        super();


        this.getElement().style.backgroundColor = "#f0f0f0";
        this.getElement().style.borderBottom = "#d0d0d0 solid 1px";

        let menu = new Menu();
        menu.add(ctx.getString("builder_new_entity"), _ => {
            new EntityModelCreator(ctx, observable).showNewDialog();
        });

        menu.add(ctx.getString("builder_new_view"), _ => {
            new ViewModelCreator(ctx, observable).showNewDialog();
        });

        let btnAdd = new FlatButton();
        btnAdd.setIcon(new Icon("add"));
        btnAdd.setOnClick(_ => {
            menu.popup(btnAdd, true);
        });
        this.add(new PullRightBox(btnAdd));
    }

}


class NamedElementCreator {
    /**
     * @param {UserInterfaceState} ctx
     * @param {Observable<AppModel>} observable
     */
    constructor(ctx, observable) {
        this.ctx = ctx;
        this.observable = observable;
    }

    showNewDialog() {
        let prefix = this.getPrefix();
        showTextInputDialog(this.ctx, this.ctx.getString("builder_create_x", this.getEntityName()), "", text => {
            this.addToModel(this.observable.getValue(), text);
            this.observable.notifyValueChanged();
        }, textField => {
            textField.setHelperText("");

            if (textField.getText().length === 0) {
                textField.setHelperText(this.ctx.getString("builder_invalid_identifier"), true);
                return false;
            }
            if (prefix.length > 0 && !textField.getText().startsWith(prefix)) {
                textField.setHelperText(this.ctx.getString("builder_identifier_prefix_x", prefix), true);
                return false;
            }

            if (textField.getText() === prefix) {
                textField.setHelperText(this.ctx.getString("builder_invalid_identifier"), true);
                return false;
            }

            if (this.observable.getValue().nameExists(textField.getText())) {
                textField.setHelperText(this.ctx.getString("builder_identifier_notunique"), true);
                return false;
            }

            let regex = new RegExp("^[A-Z][A-Za-zd_]*$");

            if (!regex.test(textField.getText())) {
                textField.setHelperText(this.ctx.getString("builder_invalid_identifier"), true);
                return false;
            }

            return true;
        });
    }

    /**
     * @return {string}
     */
    getPrefix() {
        throw new Error("abstract method");
    }

    /**
     * @return {string}
     */
    getEntityName() {
        return this.getPrefix();
    }

    /**
     *
     * @param {AppModel} appModel
     * @param {string} name
     */
    addToModel(appModel, name) {
        throw new Error("abstract method");
    }
}


class ViewModelCreator extends NamedElementCreator {

    getPrefix() {
        return "View";
    }

    addToModel(appModel, name) {
        appModel.entities.push(new View(name, new ViewModel()));
    }
}

class EntityModelCreator extends NamedElementCreator {

    getPrefix() {
        return "";
    }

    getEntityName() {
        return "Entity";
    }

    addToModel(appModel, name) {
        appModel.entities.push(new Entity(name, new ViewModel()));
    }
}