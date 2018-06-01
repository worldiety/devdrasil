import {Body1, H3, UserInterfaceState} from "/wwt/components.js";

import {DefaultUserInterfaceState, Main} from "/frontend/uis-default.js";
import {ROLE_LIST_USER} from "/frontend/repository/userRepository.js";

export {UISDashboard}

class UISDashboard extends DefaultUserInterfaceState {

    static NAME() {
        return "dashboard";
    }

    constructor(app) {
        super(app);
    }

    apply() {
        super.apply();

        let card = new Main();

        let hello = new H3();
        hello.setText("dashboard");
        card.add(hello);

        this.topBar.setTitle(this.getString("dashboard"));

        this.getApplication().getUserRepository().getUsers().then(users => {
            for (let user of users) {
                card.add(new Body1(JSON.stringify(user)));
            }
        });

        this.getApplication().getUserRepository().getUser().then(user => {
            if (user.hasProperty(ROLE_LIST_USER)){
                card.add(new Body1("account:" + JSON.stringify(user)));
            }

        });

        this.setContent(card);

    }


}

