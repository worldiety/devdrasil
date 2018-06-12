import {DefaultUserInterfaceState, Main} from "/frontend/uis-default.js";
import {Group} from "/frontend/repository/GroupRepository.js";
import {ViewGroupForm} from "./ViewGroupForm.js";

export {UISAddGroup}

class UISAddGroup extends DefaultUserInterfaceState {

    static NAME() {
        return "add-group";
    }

    constructor(app) {
        super(app);
    }

    apply() {
        super.apply();

        let group = new Group();

        let card = new Main();
        this.topBar.setTitle(this.getString("new_group"));

        let groupForm = new ViewGroupForm(this, group);
        card.add(groupForm);


        groupForm.btnLogin.setOnClick(e => {
            groupForm.clearTips();


            //quick evaluation
            let hasError = groupForm.checkDefault();


            if (hasError) {
                return
            }


            groupForm.updateModel();

            this.getApplication().getGroupRepository().add(group).then(updatedGroup => {
                window.history.back()
            }).catch(err => {
                groupForm.handleFormError(err);

            });

        });


        this.setContent(card);
    }

}

