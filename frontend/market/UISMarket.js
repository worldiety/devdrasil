import {CenterBox, CircularProgressIndicator,} from "/wwt/components.js";
import {DefaultUserInterfaceState, Main} from "/frontend/uis-default.js";
import {ViewSearchList} from "./ViewSearchList.js";

export {UISMarket}

class UISMarket extends DefaultUserInterfaceState {

    static NAME() {
        return "/market";
    }

    constructor(app) {
        super(app);
    }

    apply() {
        super.apply();
        this.topBar.setTitle(this.getString("market_store"));

        this.getApplication().getMarketRepository().getIndex().then(index => {
            let card = new Main();
            let list = new ViewSearchList(this);
            let query = this.getNavigation().getSearchParam("q");
            if (query == null) {
                query = "";
            }
            list.setModel(index, query);


            card.add(list);
            this.setContent(card);
        }).catch(err => this.handleDefaultError(err));


        let box = new CenterBox();
        let spinner = new CircularProgressIndicator();
        box.add(spinner);
        this.setContent(box);
    }

}


