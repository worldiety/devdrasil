import {H3, UserInterfaceState} from "/wwt/components.js";

import {DefaultUserInterfaceState} from "/frontend/uis-default.js";

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
        let hello = new H3();
        hello.setText("dashboard");

        this.topBar.setTitle(this.getString("dashboard"));

        this.setContent(hello);

    }


}

