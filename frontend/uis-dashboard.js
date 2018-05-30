import {H1, UserInterfaceState} from "/wwt/components.js";

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
        let hello = new H1();
        hello.setText("dashboard");

        this.setContent(hello);

    }


}

