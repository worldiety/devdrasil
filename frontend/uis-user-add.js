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
        user.login = "";
        user.firstname = "";
        user.lastname = "";

        let card = new Main();
        this.topBar.setTitle(this.getString("new_user"));

        let row = new TextFieldRow();

        let login = new TextField();
        login.setCaption(this.getString("login"));
        login.setText(user.login);
        row.add(login);


        let firstname = new TextField();
        firstname.setCaption(this.getString("firstname"));
        firstname.setText(user.firstname);
        row.add(firstname);

        let lastname = new TextField();
        lastname.setCaption(this.getString("lastname"));
        lastname.setText(user.lastname);
        row.add(lastname);

        card.add(row);

        let btnRow = new LRLayout();


        let btnCancel = new FlatButton();
        btnCancel.setText(this.getString("cancel"));
        btnCancel.setOnClick(_ => window.history.back());
        btnRow.addRight(btnCancel);

        btnRow.addRight(new HPadding());

        let btnLogin = new RaisedButton();
        btnLogin.setText(this.getString("add"));
        btnRow.addRight(btnLogin);

        card.add(btnRow);

        btnLogin.setOnClick(e => {
            //clear the tips
            login.setHelperText();
            firstname.setHelperText();
            lastname.setHelperText();

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

            if (hasError) {
                return
            }


            user.login = login.getText().trim();
            user.firstname = firstname.getText().trim();
            user.lastname = lastname.getText().trim();
            this.getApplication().getUserRepository().createUser(user).then(updatedUser => {
                window.history.back()
            }).catch(err => {
                if (err instanceof NotUniqueException) {
                    login.setHelperText(this.getString("name_not_unique"), true);
                } else {
                    this.handleDefaultError(err);
                }

            });

        });

        this.setContent(card);
    }

}

