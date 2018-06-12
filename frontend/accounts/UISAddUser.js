import {DefaultUserInterfaceState, Main} from "/frontend/uis-default.js";
import {User} from "/frontend/repository/UserRepository.js";
import {ViewUserForm} from "./ViewUserForm.js";

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

        let userForm = new ViewUserForm(this, user);
        card.add(userForm);


        userForm.btnLogin.setOnClick(e => {
            userForm.clearTips();


            //quick evaluation
            let hasError = userForm.checkDefault();

            if (userForm.pwd1.getText().trim().length === 0) {
                userForm.pwd1.setHelperText(this.getString("field_is_empty"), true);
                hasError = true;
            }

            if (userForm.pwd2.getText().trim().length === 0) {
                userForm.pwd2.setHelperText(this.getString("field_is_empty"), true);
                hasError = true;
            }

            if (userForm.pwd1.getText() !== userForm.pwd2.getText()) {
                userForm.pwd1.setHelperText(this.getString("passwords_unmatch"), true);
                userForm.pwd2.setHelperText(this.getString("passwords_unmatch"), true);
                hasError = true;
            }

            if (hasError) {
                return
            }


            userForm.updateModel();

            this.getApplication().getUserRepository().add(user).then(updatedUser => {
                window.history.back()
            }).catch(err => {
                userForm.handleFormError(err);

            });

        });


        this.setContent(card);
    }

}

