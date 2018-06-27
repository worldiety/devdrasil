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
            showTextInputDialog(ctx, "hallo", "", text => {
                modelController.getValue().entities.push(new View(text, new ViewModel()));
                modelController.notifyValueChanged();
            }, textField => {
                textField.setHelperText("");

                if (textField.getText().length === 0) {
                    textField.setHelperText("Darf nicht leer sein", true);
                    return false;
                }
                if (!textField.getText().startsWith("View")) {
                    textField.setHelperText("Muss mit View anfangen", true);
                    return false;
                }

                if (textField.getText() === "View") {
                    textField.setHelperText("Darf nicht View heißen", true);
                    return false;
                }

                if (modelController.getValue().nameExists(textField.getText())) {
                    textField.setHelperText("Name muss eindeutig sein", true);
                    return false;
                }

                return true;
            });

        });

        let btnAdd = new FlatButton();
        btnAdd.setIcon(new Icon("+"));
        btnAdd.setOnClick(_ => {
            menu.popup(btnAdd);
        });
        this.add(btnAdd);
    }

}

