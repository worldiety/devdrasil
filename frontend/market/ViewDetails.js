import {
    Body1,
    Box,
    Card,
    CenterBox,
    CircularProgressIndicator,
    FlatButton,
    H4,
    H5,
    H6,
    Icon,
    Image,
    LayoutGrid,
    ListView,
    P,
    TextField,
    TwoLineLeadingAndTrailingIcon
} from "/wwt/components.js";
import {DefaultUserInterfaceState, Main} from "/frontend/uis-default.js";
import {UISMarket} from "./UISMarket.js";

export {ViewDetails}

class ViewDetails extends Box {
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

    /**
     * @param {Plugin} plugin
     */
    setModel(plugin) {
        this.plugin = plugin;
        this.removeAll();

        let header = new LayoutGrid();
        let icon = new Image();
        icon.getElement().style.width = "100%";
        icon.setSrc(this.plugin.getIcon());


        let box = new Box();
        box.add(new H4(this.plugin.getName()));

        //add cross links for vendor and categories
        box.add(new Body1());
        let crossLinks = [];
        crossLinks.push(this.plugin.getVendor().getName());
        for (let category of this.plugin.getCategories()) {
            crossLinks.push(category);
        }

        for (let query of crossLinks) {
            let vendorLink = new FlatButton(query);
            vendorLink.setOnClick(_ => {
                this.ctx.getNavigation().forward(UISMarket.NAME(), {q: query});
            });

            box.add(vendorLink);
        }


        header.add(icon, 2);
        header.add(box, 10)

        this.add(header);
    }

}