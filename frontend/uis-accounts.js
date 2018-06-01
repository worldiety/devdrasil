import {Body1, H3, UserInterfaceState} from "/wwt/components.js";

import {DefaultUserInterfaceState, Main} from "/frontend/uis-default.js";
import {ROLE_LIST_USER} from "/frontend/repository/userRepository.js";

export {UISAccounts}

class UISAccounts extends DefaultUserInterfaceState {

    static NAME() {
        return "accounts";
    }

    constructor(app) {
        super(app);
    }

    apply() {
        super.apply();

        this.getTopBar().setTitle(this.getString("accounts"));

        let card = new Main();


        this.getApplication().getUserRepository().getUsers().then(users => {
            for (let user of users) {
                card.add(new Body1(JSON.stringify(user)));
            }
        });


        this.setContent(card);

    }


}

