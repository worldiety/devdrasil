import {
    Body1,
    Box,
    Card,
    CenterBox,
    Checkbox,
    CircularProgressIndicator,
    fadeIn,
    fadeOut,
    FlatButton,
    H6,
    HPadding,
    LayoutGrid,
    ListView,
    LRLayout,
    NotUniqueException,
    PasswordField,
    RaisedButton,
    TextField,
    TextFieldRow,
    UserInterfaceState
} from "/wwt/components.js";

export {ViewGroupForm}

class ViewGroupForm extends Box {
    /**
     *
     * @param {DefaultUserInterfaceState} ctx
     * @param {Group} user
     */
    constructor(ctx, group) {
        super();
        this.uis = ctx;
        this.group = group;
        let row = new LayoutGrid();
        row.widthMatchParent();


        this.name = new TextField();
        this.name.setCaption(ctx.getString("name"));
        this.name.setText(group.name);
        this.name.widthMatchParent();
        row.add(this.name, 6);


        this.add(row);

        this.userList = new SelectedUserList(ctx, group);
        this.userList.refresh();
        this.add(this.userList);


        let btnRow = new LRLayout();
        this.btnCancel = new FlatButton();
        this.btnCancel.setText(ctx.getString("cancel"));
        this.btnCancel.setOnClick(_ => window.history.back());
        btnRow.addRight(this.btnCancel);

        btnRow.addRight(new HPadding());

        this.btnLogin = new RaisedButton();
        this.btnLogin.setText(ctx.getString("save"));
        btnRow.addRight(this.btnLogin);


        this.add(btnRow);
    }

    clearTips() {
        this.name.setHelperText();
    }

    /**
     *
     * @returns {boolean}
     */
    checkDefault() {
        let hasError = false;

        if (this.name.getText().trim().length === 0) {
            this.name.setHelperText(this.uis.getString("field_is_empty"), true);
            hasError = true;
        }

        return hasError;
    }

    updateModel() {
        this.group.name = this.name.getText().trim();
        this.userList.updateModel();
        return this.group;
    }

    /**
     *
     * @param {Error} err
     * @returns {boolean} true if handled, otherwise false
     */
    handleFormError(err) {
        if (err instanceof NotUniqueException) {
            this.name.setHelperText(this.uis.getString("name_not_unique"), true);
            return true;
        }
        return false;
    }
}


class SelectedUserList extends Box {
    /**
     * @param {DefaultUserInterfaceState} ctx
     * @param {Group} group
     */
    constructor(ctx, group) {
        super();
        this.group = group;
        this.ctx = ctx;
    }

    refresh() {
        //show a spinner
        this.removeAll();
        let box = new CenterBox();
        let spinner = new CircularProgressIndicator();
        box.add(spinner);
        this.add(box);

        //perform reload
        this.ctx.getApplication().getUserRepository().list().then(users => {
            this.listView = new ListView();
            for (let user of users) {
                let entry = new SelectableUserEntry(user);
                if (user.hasGroup(this.group.id)) {
                    entry.setSelected(true);
                }
                this.listView.add(entry);
            }
            this.removeAll();
            this.add(this.listView);
        }).catch(err => this.ctx.handleDefaultError(err));
    }

    updateModel() {
        this.group.users = [];
        for (let entry of this.listView.getChildren()) {
            if (entry.isSelected()) {
                this.group.users.push(entry.user.id);
            }
        }
    }

}

class SelectableUserEntry extends Box {
    /**
     *
     * @param {User} user
     */
    constructor(user) {
        super();
        this.user = user;
        this.checkbox = new Checkbox();
        this.checkbox.setCaption(user.firstname + " " + user.lastname + " (" + user.login + ")");
        this.add(this.checkbox);
    }

    isSelected() {
        return this.checkbox.isChecked();
    }

    setSelected(selected) {
        this.checkbox.setChecked(selected);
    }
}