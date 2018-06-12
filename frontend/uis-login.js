import {
    Box,
    Card,
    H6,
    LayoutGrid,
    PasswordField,
    PermissionDeniedException,
    RaisedButton,
    TextField,
    TextFieldRow,
    UserInterfaceState,
    CircularProgressIndicator
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
        username.widthMatchParent();
        card.add(username);

        let password = new PasswordField();
        password.setCaption(this.getString("password"));
        password.widthMatchParent();
        card.add(password);


        let btnLogin = new RaisedButton();
        btnLogin.getElement().style.maxWidth = "5rem";
        btnLogin.getElement().style.marginLeft = "auto";
        btnLogin.getElement().style.marginRight = "auto";
        btnLogin.setText(this.getString("login"));

        let spinner = new CircularProgressIndicator();
        spinner.getElement().style.marginLeft = "auto";
        spinner.getElement().style.marginRight = "auto";
        spinner.getElement().style.display = "none";
        card.add(spinner);

        card.add(btnLogin);

        btnLogin.setOnClick(e => {
            spinner.getElement().style.display = "block";
            btnLogin.getElement().style.display = "none";
            username.setHelperText("");
            password.setHelperText("");

            this.getApplication().getSessionRepository().deleteSession();
            this.getApplication().getSessionRepository().getSession(username.getText(), password.getText(), "web-client-1.0").then(session => {
                    this.getNavigation().forward(UISDashboard.NAME());
                }
            ).catch(err => {
                btnLogin.getElement().style.display = "";
                spinner.getElement().style.display = "none";
                if (err instanceof PermissionDeniedException) {
                    username.setHelperText(this.getString("login_failed"), true);
                    password.setHelperText(this.getString("login_failed"), true);
                } else {
                    this.handleDefaultError(err)
                }
            });
        });

        this.setContentWithoutToolbar(card);
    }

}

