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

import {DefaultUserInterfaceState, Main} from "/frontend/uis-default.js";
import {UISDashboard} from "/frontend/uis-dashboard.js";

export {UISUser, UserFormLayout}

class UISUser extends DefaultUserInterfaceState {

    static NAME() {
        return "user";
    }

    constructor(app) {
        super(app);
    }

    apply() {
        super.apply();


        let uid = this.getNavigation().getSearchParam("uid");

        this.getApplication().getUserRepository().getUser(uid).then(user => {
            let card = new Main();
            this.topBar.setTitle(user.firstname + " " + user.lastname);

            let userForm = new UserFormLayout(this, user);
            card.add(userForm);

            userForm.btnLogin.setOnClick(e => {
                userForm.clearTips();


                //quick evaluation
                let hasError = userForm.checkDefault();


                //only set pwd if not empty, otherwise keep the old one
                if (userForm.pwd1.getText().trim().length !== 0 || userForm.pwd2.getText().trim().length !== 0) {

                    if (userForm.pwd1.getText() !== userForm.pwd2.getText()) {
                        userForm.pwd1.setHelperText(this.getString("passwords_unmatch"), true);
                        userForm.pwd2.setHelperText(this.getString("passwords_unmatch"), true);
                        hasError = true;
                    }
                }


                if (hasError) {
                    return
                }

                userForm.updateModel();

                this.getApplication().getUserRepository().updateUser(user).then(updatedUser => {
                    //login case may have been rewritten
                    user.login = updatedUser.login;

                    //we succeeded
                    window.history.back()

                }).catch(err => {
                    userForm.handleFormError(err);

                });

            });


            this.setContent(card);

        }).catch(err => {

            this.handleDefaultError(err)
        });


        let box = new CenterBox();
        let spinner = new CircularProgressIndicator();
        box.add(spinner);
        this.setContent(box);
    }

}


class UserFormLayout extends Box {
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

    handleFormError(err) {
        if (err instanceof NotUniqueException) {
            login.setHelperText(this.uis.getString("name_not_unique"), true);
        } else {

            if (err.message.indexOf("password to weak") >= 0) {
                this.pwd1.setHelperText(this.uis.getString("password_weak"), true);
                this.pwd2.setHelperText(this.uis.getString("password_weak"), true);
                return;
            } else {
                this.handleDefaultError(err);
            }

        }
    }
}
