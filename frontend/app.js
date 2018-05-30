import {UISDashboard} from "/frontend/uis-dashboard.js";
import {Application} from "/wwt/components.js";
import {UISLogin} from "/frontend/uis-login.js";

class ExampleApp extends Application {

    onCreate() {
        this.validateSession();
    }

    validateSession() {
        let sessionId = this.getSessionId();
        if (sessionId == null) {
            this.getNavigation().forward(UISLogin.NAME());
        } else {
            this.getNavigation().forward(UISDashboard.NAME());
        }
    }

    setSessionId(id) {
        localStorage.setItem("sessionId", id);
    }

    /*
        @return string|null
     */
    getSessionId() {
        return localStorage.getItem("sessionId");
    }
}


let app = new ExampleApp();
app.getNavigation().registerUserInterfaceState(UISLogin.NAME(), app => new UISLogin(app));
app.getNavigation().registerUserInterfaceState(UISDashboard.NAME(), app => new UISDashboard(app));
app.setLocale("de");
app.addTranslation("de", "/frontend/values-de/strings.xml");
app.addTranslation("en", "/frontend/values-en/strings.xml");
app.create();