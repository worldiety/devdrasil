import {UISDashboard} from "/frontend/uis-dashboard.js";
import {UISAccounts} from "/frontend/accounts/UISAccounts.js";
import {UISEditUser} from "/frontend/accounts/UISEditUser.js";
import {UISAddUser} from "/frontend/accounts/UISAddUser.js";
import {Application, Icon} from "/wwt/components.js";
import {UISLogin} from "/frontend/uis-login.js";
import {SessionRepository} from "/frontend/repository/sessionRepository.js";
import {UserRepository} from "./repository/UserRepository.js";
import {GroupRepository} from "./repository/GroupRepository.js";
import {CompanyRepository} from "./repository/CompanyRepository.js";
import {UISAddGroup} from "./accounts/UISAddGroup.js";
import {UISEditGroup} from "./accounts/UISEditGroup.js";
import {UISEditCompany} from "./accounts/UISEditCompany.js";
import {UISAddCompany} from "./accounts/UISAddCompany.js";
import {MarketRepository} from "./repository/MarketRepository.js";
import {UISMarket} from "/frontend/market/UISMarket.js";
import {UISMarketPlugin} from "/frontend/market/UISMarketPlugin.js";
import {UISBuilder} from "./builder/UISBuilder.js";
import {ProjectRepository} from "./builder/ProjectRepository.js";

class ExampleApp extends Application {

    constructor() {
        super();
        this.menuEnrichers = [];

        this.getNavigation().registerUserInterfaceState("", app => new UISDashboard(app));
        this.getNavigation().registerUserInterfaceState(UISLogin.NAME(), app => new UISLogin(app));
        this.getNavigation().registerUserInterfaceState(UISDashboard.NAME(), app => new UISDashboard(app));
        this.getNavigation().registerUserInterfaceState(UISAccounts.NAME(), app => new UISAccounts(app));
        this.getNavigation().registerUserInterfaceState(UISEditUser.NAME(), app => new UISEditUser(app));
        this.getNavigation().registerUserInterfaceState(UISAddUser.NAME(), app => new UISAddUser(app));
        this.getNavigation().registerUserInterfaceState(UISAddGroup.NAME(), app => new UISAddGroup(app));
        this.getNavigation().registerUserInterfaceState(UISEditGroup.NAME(), app => new UISEditGroup(app));
        this.getNavigation().registerUserInterfaceState(UISEditCompany.NAME(), app => new UISEditCompany(app));
        this.getNavigation().registerUserInterfaceState(UISAddCompany.NAME(), app => new UISAddCompany(app));
        this.getNavigation().registerUserInterfaceState(UISMarket.NAME(), app => new UISMarket(app));
        this.getNavigation().registerUserInterfaceState(UISMarketPlugin.NAME(), app => new UISMarketPlugin(app));
        this.getNavigation().registerUserInterfaceState(UISBuilder.NAME(), app => new UISBuilder(app));
        this.setLocale("de");
        this.addTranslation("de", "/frontend/values-de/strings.xml");
        this.addTranslation("en", "/frontend/values-en/strings.xml");

        //add the default bootstrapping entries
        this.getMenuEnrichers().push(drawer => {
            return this.getUserRepository().getSessionUser().then(user => {
                //present logoff, if we are logged in
                let logoutItem = drawer.addMenuEntry("#" + UISLogin.NAME(), new Icon("exit_to_app"), this.getString("logout"), false);
                logoutItem.onclick = e => {
                    this.getSessionRepository().deleteSession().then(_ => {
                        this.getNavigation().forward(UISLogin.NAME());
                        //perform a force reload
                        location.reload(true);
                    });

                };

                this.getUserRepository().getUserPermissions(user.id).then(permissions => {
                    if (permissions.listUsers) {
                        let selected = this.getNavigation().getPendingName() === UISAccounts.NAME();
                        drawer.addMenuEntry("#" + UISAccounts.NAME(), new Icon("supervisor_account"), this.getString("accounts"), selected);
                    }

                    if (permissions.listMarket) {
                        let selected = this.getNavigation().getPendingName() === UISMarket.NAME();
                        drawer.addMenuEntry("#" + UISMarket.NAME(), new Icon("extension"), this.getString("market_store"), selected);
                    }
                });


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

    getGroupRepository() {
        if (this.groupRepository == null) {
            this.groupRepository = new GroupRepository(this.getFetcher(), this.getSessionRepository());
        }
        return this.groupRepository;
    }

    getCompanyRepository() {
        if (this.companyRepository == null) {
            this.companyRepository = new CompanyRepository(this.getFetcher(), this.getSessionRepository());
        }
        return this.companyRepository;
    }

    /**
     *
     * @return {MarketRepository}
     */
    getMarketRepository() {
        if (this.marketRespository == null) {
            this.marketRespository = new MarketRepository(this.getFetcher(), this.getSessionRepository());
        }
        return this.marketRespository;
    }

    /**
     *
     * @return {ProjectRepository}
     */
    getProjectRepository() {
        if (this.projectRepository == null) {
            this.projectRepository = new ProjectRepository();
        }
        return this.projectRepository;
    }

    validateSession() {
        let session = this.getUserRepository().getSessionUser().then(user => {
            //try to navigate directly to the input link, if registered
            let targetName = window.location.hash.substring(1);
            if (this.getNavigation().hasName(targetName)) {
                this.getNavigation().reload();
            } else {
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