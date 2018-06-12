import {Box, Card, CenterBox, CircularProgressIndicator,} from "/wwt/components.js";
import {DefaultUserInterfaceState, Main} from "/frontend/uis-default.js";
import {ViewUserForm} from "./ViewUserForm.js";

export {UISEditUser}

class UISEditUser extends DefaultUserInterfaceState {

    static NAME() {
        return "user";
    }

    constructor(app) {
        super(app);
    }

    apply() {
        super.apply();


        let uid = this.getNavigation().getSearchParam("uid");

        this.getApplication().getUserRepository().get(uid).then(user => {
            let card = new Main();
            this.topBar.setTitle(user.firstname + " " + user.lastname);

            let userForm = new ViewUserForm(this, user);
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

                this.getApplication().getUserRepository().update(user).then(updatedUser => {
                    //login case may have been rewritten
                    user.login = updatedUser.login;

                    //we succeeded
                    window.history.back()

                }).catch(err => {
                    if (!userForm.handleFormError(err)) {
                        this.handleDefaultError(err);
                    }

                });

            });


            this.setContent(card);

        }).catch(err => {
            this.handleDefaultError(err);
        });


        let box = new CenterBox();
        let spinner = new CircularProgressIndicator();
        box.add(spinner);
        this.setContent(box);
    }

}



