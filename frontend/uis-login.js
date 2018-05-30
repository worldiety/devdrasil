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

export {UISLogin}

class UISLogin extends DefaultUserInterfaceState {

    constructor(app) {
        super(app);
    }

    apply() {
        let appName = "dev";
        this.setTitle(this.getString("login_title", appName));
        document.body.style.backgroundColor = "#f1f1f1";

        let card = new Card();
        card.getElement().style.maxWidth = "25rem";
        card.getElement().style.padding = "1rem";
        card.getElement().style.margin = "auto";

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
            this.doLogin(username, password);
        });

        this.setContent(card);
    }

    doLogin(user, password) {
        let headers = new Headers();
        headers.append("login", user);
        headers.append("password", password);
        headers.append("client", "devdrasil");

        let params = {
            headers: headers,
        };

        let request = new Request("/session/auth");


        let promise = fetch(request, params)
        promise.then(res => {
            if (res.status != 200) {

            }
            console.log(res);
        });
        this.attachPromise(promise);
    }
}

