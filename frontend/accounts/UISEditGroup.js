import {DefaultUserInterfaceState, Main} from "/frontend/uis-default.js";
import {Group} from "/frontend/repository/GroupRepository.js";
import {ViewGroupForm} from "./ViewGroupForm.js";
import {Box, Card, CenterBox, CircularProgressIndicator,} from "/wwt/components.js";

export {UISEditGroup}

class UISEditGroup extends DefaultUserInterfaceState {

    static NAME() {
        return "edit-group";
    }

    constructor(app) {
        super(app);
    }

    apply() {

        super.apply();


        let uid = this.getNavigation().getSearchParam("gid");

        this.getApplication().getGroupRepository().get(uid).then(group => {
            let card = new Main();
            this.topBar.setTitle(group.name);

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

                this.getApplication().getGroupRepository().update(group).then(updatedGroup => {
                    //we succeeded
                    window.history.back()

                }).catch(err => {
                    if (!groupForm.handleFormError(err)) {
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

