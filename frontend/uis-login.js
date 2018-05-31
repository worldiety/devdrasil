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

import {DefaultUserInterfaceState} from "/frontend/uis-default.js";
import {UISDashboard} from "/frontend/uis-dashboard.js";

export {UISLogin}

class UISLogin extends DefaultUserInterfaceState {

    static NAME() {
        return "login";
    }

    constructor(app) {
        super(app);
    }

    apply() {
        super.apply();
        let appName = "dev";
        this.setTitle(this.getString("login_title", appName));


        let card = new Card();
        card.getElement().style.maxWidth = "25rem";
        card.getElement().style.margin = "auto";
        card.getElement().style.marginTop = "1rem";

        let title = new H6();
        title.setText(this.getString("login_title", appName));
        card.add(title);

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
            this.getApplication().getSessionRepository().getSession(username.getText(), password.getText(), "devdrasil").then(session => {
                    this.getNavigation().forward(UISDashboard.NAME());
                }
            ).catch(err => this.handleDefaultError(err));
        });

        this.setContentWithoutToolbar(card);
    }

}

