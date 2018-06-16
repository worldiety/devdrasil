import {
    Body1,
    Box,
    Card,
    CenterBox,
    CircularProgressIndicator,
    Icon,
    ListView,
    TextField,
    TwoLineLeadingAndTrailingIcon
} from "/wwt/components.js";
import {DefaultUserInterfaceState, Main} from "/frontend/uis-default.js";

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
            let list = new SearchListView(this, index);
            list.refresh();

            card.add(list);
            this.setContent(card);
        }).catch(err => this.handleDefaultError(err));


        let box = new CenterBox();
        let spinner = new CircularProgressIndicator();
        box.add(spinner);
        this.setContent(box);
    }

}


class SearchListView extends Box {
    /**
     *
     * @param {DefaultUserInterfaceState} ctx
     * @param {MarketIndex} index
     */
    constructor(ctx, index) {
        super();
        this.ctx = ctx;
        this.index = index;
    }

    refresh() {
        this.removeAll();

        this.listView = new ListView();
        this.listView.setInteractive(true);


        this.searchView = new TextField();
        this.searchView.setCaption(this.ctx.getString("search"));
        this.searchView.widthMatchParent();
        this.searchView.setIcon(new Icon("search"));
        this.searchView.setAutoCompleteHandler(text => {
            this.rebuildListView(text);
            let tmp = [];
            for (let word of this.index.findWord(text, 10)) {
                let entry = new Body1(word);
                entry.getElement().onclick = evt => {
                    this.searchView.setText(word);
                    this.rebuildListView(text);
                };
                tmp.push(entry);
            }
            return tmp;
        });
        this.add(this.searchView);
        this.rebuildListView();

        this.add(this.listView);
    }

    rebuildListView(text = "") {
        this.listView.removeAll();
        for (let plugin of this.index.getPlugins(text)) {
            let item = new TwoLineLeadingAndTrailingIcon(new Icon("account_circle"), plugin.getName(), plugin.getName(), null);
            this.listView.add(item, true);
        }
    }
}