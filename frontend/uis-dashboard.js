import {H3, UserInterfaceState} from "/wwt/components.js";

import {DefaultUserInterfaceState, Main} from "/frontend/uis-default.js";

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

        this.setContent(card);

    }


}

