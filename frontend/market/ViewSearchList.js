import {
    Body1,
    Box,
    Card,
    CenterBox,
    CircularProgressIndicator,
    Icon,
    Image,
    ListView,
    P,
    TextField,
    TwoLineLeadingAndTrailingIcon
} from "/wwt/components.js";
import {DefaultUserInterfaceState, Main} from "/frontend/uis-default.js";
import {UISMarketPlugin} from "./UISMarketPlugin.js";

export {ViewSearchList}

class ViewSearchList extends Box {
    /**
     *
     * @param {DefaultUserInterfaceState} ctx
     *
     */
    constructor(ctx) {
        super();
        this.ctx = ctx;

    }

    /**
     * @param {MarketIndex} index
     * @param {string} query
     */
    setModel(index, query = "") {
        this.index = index;
        this.removeAll();

        this.listView = new ListView();
        this.listView.setInteractive(true);


        this.searchView = new TextField();
        if (query !== ""){
            this.searchView.setText(query);
        }
        this.searchView.setCaption(this.ctx.getString("search"));
        this.searchView.widthMatchParent();
        this.searchView.setIcon(new Icon("search"));
        this.searchView.setAutoCompleteHandler(text => {
            this.rebuildListView(text);
            let tmp = [];
            for (let word of this.index.findWord(text, 10)) {
                let entry = new P(word);
                entry.getElement().className = "mdc-list-item autocomplete-entry";
                entry.getElement().onclick = evt => {
                    this.searchView.setText(word);
                    this.rebuildListView(text);
                };
                tmp.push(entry);
            }
            return tmp;
        });
        this.add(this.searchView);
        this.rebuildListView(query);


        this.add(this.listView);
    }

    rebuildListView(text = "") {
        this.listView.removeAll();
        for (let plugin of this.index.getPlugins(text)) {
            let logo = new Image();
            logo.setSrc(plugin.getIcon());
            let item = new TwoLineLeadingAndTrailingIcon(logo, plugin.getName(), plugin.getVendor().getName(), null);


            //rewrite logo style
            logo.getElement().className = "";
            logo.getElement().style.width = "30pt";
            logo.getElement().style.height = "30pt";
            logo.getElement().style.marginRight = "16px";
            let listEntry = this.listView.add(item, true);
            listEntry.onclick = evt => {
                this.ctx.getNavigation().forward(UISMarketPlugin.NAME(), {id: plugin.getId()});
            };
        }
    }
}