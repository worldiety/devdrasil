import {
    AppBar,
    Body1,
    Box,
    Button,
    Card,
    Dialog,
    Drawer,
    Icon,
    PermissionDeniedException,
    showMessage as showMsgLib,
    UserInterfaceState,
} from "/wwt/components.js";

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

        this.getApplication().getUserRepository().getSessionUser().then(user => {
            this.drawer.setCaption(user.firstname + " " + user.lastname);
        });


        this.topBar.setDrawer(this.drawer);
    }

    onCreateSideMenu() {
        this.getApplication().onCreateDefaultSideMenu(this.drawer).then(_ => {


        });

    }

    handleDefaultError(err) {
        if (err.message.indexOf("delete the super user") >= 0) {
            this.showMessage(this.getString("cannot_delete_root"));
            return;
        }


        if (err instanceof PermissionDeniedException) {

            this.showMessage(this.getString("login_failed"));

        } else {
            this.showMessage(err.toString());
        }

    }

    showMessage(text) {
        showMsgLib(text, this.getString("ok"));
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

    getTopBar() {
        return this.topBar;
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

        this.getElement().style.maxWidth = "960px";
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