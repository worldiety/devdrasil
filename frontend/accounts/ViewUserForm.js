import {
    Box,
    Card,
    CenterBox,
    CircularProgressIndicator,
    fadeIn,
    fadeOut,
    FlatButton,
    H6,
    HPadding,
    LayoutGrid,
    LRLayout,
    NotUniqueException,
    PasswordField,
    RaisedButton,
    TextField,
    TextFieldRow,
    UserInterfaceState
} from "/wwt/components.js";

export {ViewUserForm}

class ViewUserForm extends Box {
    /**
     *
     * @param {UserInterfaceState} ctx
     * @param {User} user
     */
    constructor(ctx, user) {
        super();
        this.uis = ctx;
        this.user = user;
        let row = new LayoutGrid();
        row.widthMatchParent();


        this.firstname = new TextField();
        this.firstname.setCaption(ctx.getString("firstname"));
        this.firstname.setText(user.firstname);
        this.firstname.widthMatchParent();
        row.add(this.firstname, 6);

        this.lastname = new TextField();
        this.lastname.setCaption(ctx.getString("lastname"));
        this.lastname.setText(user.lastname);
        this.lastname.widthMatchParent();
        row.add(this.lastname, 6);

        this.add(row);

        let btnRow = new LRLayout();


        this.btnCancel = new FlatButton();
        this.btnCancel.setText(ctx.getString("cancel"));
        this.btnCancel.setOnClick(_ => window.history.back());
        btnRow.addRight(this.btnCancel);

        btnRow.addRight(new HPadding());

        this.btnLogin = new RaisedButton();
        this.btnLogin.setText(ctx.getString("save"));
        btnRow.addRight(this.btnLogin);


        let row2 = new LayoutGrid();
        row2.widthMatchParent();

        this.login = new TextField();
        this.login.setCaption(ctx.getString("login"));
        this.login.setText(user.login);
        this.login.widthMatchParent();
        row2.add(this.login, 4);

        this.pwd1 = new PasswordField();
        this.pwd1.setCaption(ctx.getString("password"));
        this.pwd1.setText(user.password);
        this.pwd1.setAutoComplete("new-password");
        this.pwd1.widthMatchParent();
        row2.add(this.pwd1, 4);


        this.pwd2 = new PasswordField();
        this.pwd2.setCaption(ctx.getString("password_repeat"));
        this.pwd2.setText(user.password);
        this.pwd2.setAutoComplete("new-password");
        this.pwd2.widthMatchParent();
        row2.add(this.pwd2, 4);

        this.add(row2);

        this.add(btnRow);
    }

    clearTips() {
        this.login.setHelperText();
        this.firstname.setHelperText();
        this.lastname.setHelperText();
        this.pwd1.setHelperText();
        this.pwd2.setHelperText();
    }

    /**
     *
     * @returns {boolean}
     */
    checkDefault() {
        let hasError = false;
        if (this.login.getText().trim().length < 3) {
            this.login.setHelperText(this.uis.getString("login_too_short"), true);
            hasError = true;
        }

        if (this.firstname.getText().trim().length === 0) {
            this.firstname.setHelperText(this.uis.getString("field_is_empty"), true);
            hasError = true;
        }

        if (this.lastname.getText().trim().length === 0) {
            this.lastname.setHelperText(this.uis.getString("field_is_empty"), true);
            hasError = true;
        }
        return hasError;
    }

    updateModel() {
        this.user.login = this.login.getText().trim();
        this.user.firstname = this.firstname.getText().trim();
        this.user.lastname = this.lastname.getText().trim();
        this.user.password = this.pwd1.getText().trim();

        return this.user;
    }

    /**
     *
     * @param {Error} err
     * @returns {boolean} true if handled, otherwise false
     */
    handleFormError(err) {
        if (err instanceof NotUniqueException) {
            this.login.setHelperText(this.uis.getString("name_not_unique"), true);
            return true;
        } else {

            if (err.message.indexOf("password to weak") >= 0) {
                this.pwd1.setHelperText(this.uis.getString("password_weak"), true);
                this.pwd2.setHelperText(this.uis.getString("password_weak"), true);
                return true;
            }

        }
        return false;
    }
}