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
    LRLayout,
    P,
    RaisedButton,
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
        header.add(box, 10);


        let rateAndInstallBox = new Box();
        let rating = this.plugin.getRatings();

        let fontSize = "13px";
        let starBox = new Box();
        starBox.getElement().style.color = "#616161";

        for (let i = 0; i < 5; i++) {
            let ico = new Icon("star");
            ico.getElement().style.verticalAlign = "middle";
            ico.getElement().style.fontSize = fontSize;
            starBox.add(ico);
        }
        let numOfStars = new Body1(rating.countStars() + "");
        numOfStars.getElement().style.display = "inline-block";
        numOfStars.getElement().style.verticalAlign = "middle";
        numOfStars.getElement().style.fontSize = fontSize;
        numOfStars.getElement().style.marginLeft = "3px";
        numOfStars.getElement().style.marginRight = "3px";
        starBox.add(numOfStars);

        let ico = new Icon("person");
        ico.getElement().style.verticalAlign = "middle";
        ico.getElement().style.fontSize = fontSize;
        starBox.add(ico);
        rateAndInstallBox.add(starBox);

        let btnInstall = new RaisedButton();
        btnInstall.setText("install");
        btnInstall.getElement().style.cssFloat = "right";
        rateAndInstallBox.add(btnInstall);

        let alignRight = new LRLayout();
        alignRight.addRight(rateAndInstallBox);
        box.add(alignRight);

        this.add(header);
    }

}