import {CenterBox, CircularProgressIndicator,} from "/wwt/components.js";
import {DefaultUserInterfaceState, Main} from "/frontend/uis-default.js";
import {AppModelController} from "./AppModelController.js";
import {ViewTree} from "./ViewTree.js";

export {UISBuilder}

class UISBuilder extends DefaultUserInterfaceState {

    static NAME() {
        return "/builder";
    }

    constructor(app) {
        super(app);
    }

    apply() {
        super.apply();
        this.topBar.setTitle(this.getString("builder_title"));

        let appModelController = new AppModelController();
        let viewTree = new ViewTree(this, appModelController);
        this.setContent(viewTree);
    }

}


