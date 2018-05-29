import {Box, RaisedButton, TextField, UserInterfaceState} from "/wwt/components.js";
export {UISLogin}

class UISLogin extends UserInterfaceState {

    constructor(app){
        super(app);
    }

    apply() {
        this.setTitle(this.getString("uis_login_title"))
    }

}

