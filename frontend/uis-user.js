import {
    Box,
    Card,
    H6,
    LayoutGrid,
    PasswordField,
    RaisedButton,
    TextField,
    UserInterfaceState
} from "/wwt/components.js";

import {DefaultUserInterfaceState, Main} from "/frontend/uis-default.js";
import {UISDashboard} from "/frontend/uis-dashboard.js";

export {UISUser}

class UISUser extends DefaultUserInterfaceState {

    static NAME() {
        return "user";
    }

    constructor(app) {
        super(app);
    }

    apply() {
        super.apply();


        let card = new Main();
        let uid = this.getNavigation().getSearchParam("uid");

        this.getApplication().getUserRepository().getUser(uid).then(user => {

            card.add(new H6("Sie betrachten jetzt "+user.firstname));

            let username = new TextField();
            username.setCaption(this.getString("username"));
            card.add(username);

            let password = new PasswordField();
            password.setCaption(this.getString("password"));
            card.add(password);

            let btnLogin = new RaisedButton();
            btnLogin.getElement().style.maxWidth = "5rem";
            btnLogin.getElement().style.marginLeft = "auto";
            btnLogin.getElement().style.marginRight = "auto";
            btnLogin.setText(this.getString("login"));
            card.add(btnLogin);

            btnLogin.setOnClick(e => {
                this.getApplication().getSessionRepository().deleteSession();
                this.getApplication().getSessionRepository().getSession(username.getText(), password.getText(), "web-client-1.0").then(session => {
                        this.getNavigation().forward(UISDashboard.NAME());
                    }
                ).catch(err => this.handleDefaultError(err));
            });
        }).catch(err=>{
            this.handleDefaultError(err)
        });


        this.setContent(card);
    }

}

