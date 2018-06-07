import {
    Box,
    Card,
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

export {UISUser}

class UISUser extends DefaultUserInterfaceState {

    static NAME() {
        return "user";
    }

    constructor(app) {
        super(app);
    }

    apply() {
        super.apply();


        let card = new Main();
        let uid = this.getNavigation().getSearchParam("uid");

        this.getApplication().getUserRepository().getUser(uid).then(user => {

            this.topBar.setTitle(user.firstname + " " + user.lastname);

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
            btnLogin.setText(this.getString("save"));
            btnRow.addRight(btnLogin);

            card.add(btnRow);

            btnLogin.setOnClick(e => {
                //clear the tips
                login.setHelperText();
                firstname.setHelperText();
                lastname.setHelperText();

                //quick evaluation
                let hasError = false;
                if (login.getText().length < 3) {
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


                user.login = login.getText();
                user.firstname = firstname.getText();
                user.lastname = lastname.getText();
                this.getApplication().getUserRepository().updateUser(user).then(updatedUser => {
                    //login case may have been rewritten
                    user.login = updatedUser.login;

                    //we succeeded
                    window.history.back()

                }).catch(err => {
                    if (err instanceof NotUniqueException) {
                        login.setHelperText(this.getString("name_not_unique"), true);
                    } else {
                        this.handleDefaultError(err);
                    }

                });

            });
        }).catch(err => {

            this.handleDefaultError(err)
        });


        this.setContent(card);
    }

}

