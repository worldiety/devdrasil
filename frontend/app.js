import {UISDashboard} from "/frontend/uis-dashboard.js";
import {UISAccounts} from "/frontend/uis-accounts.js";
import {Application, Icon} from "/wwt/components.js";
import {UISLogin} from "/frontend/uis-login.js";
import {SessionRepository} from "/frontend/repository/sessionRepository.js";
import {ROLE_LIST_USER, UserRepository} from "/frontend/repository/userRepository.js";

class ExampleApp extends Application {

    constructor() {
        super();
        this.menuEnrichers = [];

        this.getNavigation().registerUserInterfaceState(UISLogin.NAME(), app => new UISLogin(app));
        this.getNavigation().registerUserInterfaceState(UISDashboard.NAME(), app => new UISDashboard(app));
        this.getNavigation().registerUserInterfaceState(UISAccounts.NAME(), app => new UISAccounts(app));
        this.setLocale("de");
        this.addTranslation("de", "/frontend/values-de/strings.xml");
        this.addTranslation("en", "/frontend/values-en/strings.xml");

        //add the default bootstrapping entries
        this.getMenuEnrichers().push(drawer => {
            return this.getUserRepository().getUser().then(user => {
                //present logoff, if we are logged in
                let logoutItem = drawer.addMenuEntry("#" + UISLogin.NAME(), new Icon("exit_to_app"), this.getString("logout"), false);
                logoutItem.onclick = e => {
                    this.getSessionRepository().deleteSession().then(_ => {
                        this.getNavigation().forward(UISLogin.NAME());
                    });

                };

                if (user.hasProperty(ROLE_LIST_USER)) {
                    let selected = this.getNavigation().getPendingName() === UISAccounts.NAME();
                    drawer.addMenuEntry("#" + UISAccounts.NAME(), new Icon("supervisor_account"), this.getString("accounts"), selected);
                }

            });
        });

    }

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
            //try to navigate directly to the input link, if registered
            let targetName = window.location.hash.substring(1);
            if (this.getNavigation().hasName(targetName)){
                this.getNavigation().forward(targetName);
            }else{
                this.getNavigation().forward(UISDashboard.NAME());
            }
        }).catch(err => {
            this.getNavigation().forward(UISLogin.NAME());
        });

    }

    async onCreateDefaultSideMenu(drawer) {
        let tmp = [];
        for (let enricher of this.menuEnrichers) {
            tmp.push(enricher(drawer));
        }
        return await Promise.all(tmp);
    }


    getMenuEnrichers() {
        return this.menuEnrichers;
    }
}


class MenuEnricher {
    /**
     *
     * @param drawer {Drawer}
     */
    async apply(drawer) {
        return null;
    }
}

let app = new ExampleApp();

app.create();