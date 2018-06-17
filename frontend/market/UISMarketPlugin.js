import {CenterBox, CircularProgressIndicator,} from "/wwt/components.js";
import {DefaultUserInterfaceState, Main} from "/frontend/uis-default.js";
import {ViewDetails} from "./ViewDetails.js";

export {UISMarketPlugin}

class UISMarketPlugin extends DefaultUserInterfaceState {

    static NAME() {
        return "/market/plugin";
    }

    constructor(app) {
        super(app);
    }

    apply() {
        super.apply();
        let id = this.getNavigation().getSearchParam("id");
        this.topBar.setTitle(this.getString("market_store"));

        this.getApplication().getMarketRepository().getIndex().then(index => {
            let plugin = index.getPlugin(id);
            this.topBar.setTitle(plugin.getName());
            let card = new Main();


            let details = new ViewDetails(this, index);
            details.setModel(plugin);
            card.add(details);
            this.setContent(card);
        }).catch(err => this.handleDefaultError(err));


        let box = new CenterBox();
        let spinner = new CircularProgressIndicator();
        box.add(spinner);
        this.setContent(box);
    }

}


