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
import {User} from "/frontend/repository/userRepository.js";

export {UISAddUser}

class UISAddUser extends DefaultUserInterfaceState {

    static NAME() {
        return "add-user";
    }

    constructor(app) {
        super(app);
    }

    apply() {
        super.apply();

        let user = new User();

        let card = new Main();
        this.topBar.setTitle(this.getString("new_user"));

        let row = new TextFieldRow();

        let login = new TextField();
        login.setCaption(this.getString("login"));
        login.setText(user.login);
        row.add(login);


        let firstname = new TextField();
        firstname.setAutoComplete("new-password");
        firstname.setCaption(this.getString("firstname"));
        firstname.setText(user.firstname);
        row.add(firstname);

        let lastname = new TextField();
        lastname.setCaption(this.getString("lastname"));
        lastname.setAutoComplete("new-password");
        lastname.setText(user.lastname);
        row.add(lastname);

        card.add(row);


        let row2 = new TextFieldRow();

        let pwd1 = new PasswordField();
        pwd1.setCaption(this.getString("password"));
        pwd1.setText(user.password);
        pwd1.setAutoComplete("new-password");
        row2.add(pwd1);


        let pwd2 = new PasswordField();
        pwd2.setCaption(this.getString("password_repeat"));
        pwd2.setText(user.password);
        pwd2.setAutoComplete("new-password");
        row2.add(pwd2);

        card.add(row2);

        let btnRow = new LRLayout();


        let btnCancel = new FlatButton();
        btnCancel.setText(this.getString("cancel"));
        btnCancel.setOnClick(_ => window.history.back());
        btnRow.addRight(btnCancel);

        btnRow.addRight(new HPadding());

        let btnLogin = new RaisedButton();
        btnLogin.setText(this.getString("add"));
        btnRow.addRight(btnLogin);


        btnLogin.setOnClick(e => {
            //clear the tips
            login.setHelperText();
            firstname.setHelperText();
            lastname.setHelperText();
            pwd1.setHelperText();
            pwd2.setHelperText();

            //quick evaluation
            let hasError = false;
            if (login.getText().trim().length < 3) {
                login.setHelperText(this.getString("login_too_short"), true);
                hasError = true;
            }

            if (firstname.getText().trim().length === 0) {
                firstname.setHelperText(this.getString("field_is_empty"), true);
                hasError = true;
            }

            if (lastname.getText().trim().length === 0) {
                lastname.setHelperText(this.getString("field_is_empty"), true);
                hasError = true;
            }

            if (pwd1.getText().trim().length === 0) {
                pwd1.setHelperText(this.getString("field_is_empty"), true);
                hasError = true;
            }

            if (pwd2.getText().trim().length === 0) {
                pwd2.setHelperText(this.getString("field_is_empty"), true);
                hasError = true;
            }

            if (pwd1.getText() !== pwd2.getText()) {
                pwd1.setHelperText(this.getString("passwords_unmatch"), true);
                pwd2.setHelperText(this.getString("passwords_unmatch"), true);
                hasError = true;
            }

            if (hasError) {
                return
            }


            user.login = login.getText().trim();
            user.firstname = firstname.getText().trim();
            user.lastname = lastname.getText().trim();
            user.password = pwd1.getText().trim();
            this.getApplication().getUserRepository().createUser(user).then(updatedUser => {
                window.history.back()
            }).catch(err => {
                if (err instanceof NotUniqueException) {
                    login.setHelperText(this.getString("name_not_unique"), true);
                } else {

                    if (err.message.indexOf("password to weak") >= 0) {
                        pwd1.setHelperText(this.getString("password_weak"), true);
                        pwd2.setHelperText(this.getString("password_weak"), true);
                        return;
                    } else {
                        this.handleDefaultError(err);
                    }

                }

            });

        });


        card.add(btnRow);

        this.setContent(card);
    }

}

