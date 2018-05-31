import {UISDashboard} from "/frontend/uis-dashboard.js";
import {Application} from "/wwt/components.js";
import {UISLogin} from "/frontend/uis-login.js";
import {SessionRepository} from "/frontend/repository/sessionRepository.js";
import {UserRepository} from "/frontend/repository/userRepository.js";

class ExampleApp extends Application {

    onCreate() {
        this.validateSession();
    }

    getSessionRepository() {
        if (this.sessionRepository == null) {
            this.sessionRepository = new SessionRepository(this.getFetcher());
        }
        return this.sessionRepository;
    }

    getUserRepository() {
        if (this.userRepository == null) {
            this.userRepository = new UserRepository(this.getFetcher(), this.getSessionRepository());
        }
        return this.userRepository;
    }

    validateSession() {
        let session = this.getSessionRepository().getSession().then(session => {
            this.getNavigation().forward(UISDashboard.NAME());
        }).catch(err => {
            this.getNavigation().forward(UISLogin.NAME());
        });

    }


}


let app = new ExampleApp();
app.getNavigation().registerUserInterfaceState(UISLogin.NAME(), app => new UISLogin(app));
app.getNavigation().registerUserInterfaceState(UISDashboard.NAME(), app => new UISDashboard(app));
app.setLocale("de");
app.addTranslation("de", "/frontend/values-de/strings.xml");
app.addTranslation("en", "/frontend/values-en/strings.xml");
app.create();