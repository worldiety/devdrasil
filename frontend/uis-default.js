import {AppBar, Box, Card, Drawer, Icon, PermissionDeniedException, UserInterfaceState} from "/wwt/components.js";

export {DefaultUserInterfaceState, Main}

/**
 * All UIS inherit from here, to get the common behavior
 */
class DefaultUserInterfaceState extends UserInterfaceState {
    constructor(app) {
        super(app);

        //the app bar should never be missing
        this.topBar = new AppBar();
        //the title sets the document title automatically, you can manually set it differently with  "this.setTitle("Component Demo");"
        this.topBar.setTitle("app title");

        //each app should have a drawer
        this.drawer = new Drawer();
        this.drawer.setCaption("drawer caption");

        this.topBar.setDrawer(this.drawer);
    }

    onCreateSideMenu() {
        this.drawer.addMenuEntry("#", new Icon("star"), "Star Text and selected", true);
        this.drawer.addMenuEntry("#", new Icon("inbox"), "Inbox Text", false);
        this.drawer.addMenuEntry("#", new Icon("add"), "Add Text", false);

        let logoutItem = this.drawer.addMenuEntry("#", new Icon("exit_to_app"), this.getString("logout"), false);
        logoutItem.onclick = e => {
            this.getApplication().getSessionRepository().deleteSession();
            //avoid dependency cycle
            this.getNavigation().forward("login");
        }
    }

    handleDefaultError(err) {
        if (err instanceof PermissionDeniedException) {
            alert(this.getString("login_failed"))
        } else {
            alert(err);
        }

    }


    apply() {
        this.onCreateSideMenu();
        super.apply();
        document.body.style.backgroundColor = "#f1f1f1";
    }

    setContent(component) {
        this.topBar.setContent(component);
        super.setContent(this.topBar);
    }

    setContentWithoutToolbar(component) {
        super.setContent(component);
    }
}


class Main extends Box {
    constructor(showAsCard = true) {
        super();
        this.showAsCard = showAsCard;
        if (showAsCard) {
            this.card = new Card();
            super.add(this.card);
        }

        this.getElement().style.maxWidth = "75rem";
        this.getElement().style.margin = "auto";
        this.getElement().style.padding = "1rem";

    }

    add(component) {
        if (this.showAsCard) {
            this.card.add(component);
        } else {
            super.add(component);
        }

    }
}