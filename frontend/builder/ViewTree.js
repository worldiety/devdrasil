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
import {View, ViewModel} from "./AppModel.js";


export {ViewTree}

class ViewTree extends Box {
    /**
     * @param {UserInterfaceState} ctx
     * @param {AppModelController} appViewModel
     */
    constructor(ctx, appViewModel) {
        super();
        this.ctx = ctx;
        this.appViewModel = appViewModel;
        appViewModel.addObserver(ctx.getLifecycle(), model => this.rebuildUI(model));
    }

    /**
     *
     * @param {AppModel} model
     */
    rebuildUI(model) {
        this.removeAll();

        let toolbar = new Toolbar(this.ctx, this.appViewModel);
        this.add(toolbar);

        let listView = new ListView();
        listView.setInteractive(true);
        for (let view of model.entities) {
            listView.add(view.name, true);
        }

        this.add(listView);
    }

}


class Toolbar extends Box {
    /**
     * @param {AppModelController} modelController
     */
    constructor(ctx, modelController) {
        super();

        let menu = new Menu();
        menu.add("New View", _ => {
            new ViewModelCreator(ctx, modelController).showNewDialog();

        });

        let btnAdd = new FlatButton();
        btnAdd.setIcon(new Icon("+"));
        btnAdd.setOnClick(_ => {
            menu.popup(btnAdd);
        });
        this.add(btnAdd);
    }

}


class NamedElementCreator {
    /**
     * @param {UserInterfaceState} ctx
     * @param {AppModelController} appViewModel
     */
    constructor(ctx, appViewModel) {
        this.ctx = ctx;
        this.appViewModel = appViewModel;
    }

    showNewDialog() {
        let prefix = this.getPrefix();
        showTextInputDialog(this.ctx, this.ctx.getString("builder_create_x", prefix), "", text => {
            this.addToModel(this.appViewModel.getValue(), text);
            this.appViewModel.notifyValueChanged();
        }, textField => {
            textField.setHelperText("");

            if (textField.getText().length === 0) {
                textField.setHelperText(this.ctx.getString("builder_invalid_identifier"), true);
                return false;
            }
            if (!textField.getText().startsWith(prefix)) {
                textField.setHelperText(this.ctx.getString("builder_identifier_prefix_x", prefix), true);
                return false;
            }

            if (textField.getText() === prefix) {
                textField.setHelperText(this.ctx.getString("builder_invalid_identifier"), true);
                return false;
            }

            if (this.appViewModel.getValue().nameExists(textField.getText())) {
                textField.setHelperText(this.ctx.getString("builder_identifier_notunique"), true);
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