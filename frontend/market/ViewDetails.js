import {
    Body1,
    Body2,
    Box,
    Card,
    CenterBox,
    CircularProgressIndicator,
    FlatButton,
    H4,
    H5,
    H6,
    HR,
    Icon,
    Image,
    LayoutGrid,
    ListView,
    LRLayout,
    NotFoundException,
    P,
    PullRightBox,
    RaisedButton,
    RoundedIcon,
    Span,
    StarBox,
    TextField,
    TwoLineLeadingAndTrailingIcon
} from "/wwt/components.js";
import {DefaultUserInterfaceState, Main} from "/frontend/uis-default.js";
import {UISMarket} from "./UISMarket.js";

export {ViewDetails}

//some random colors from https://material.io/design/color/the-color-system.html#tools-for-picking-colors
const colorTable = ["#B71C1C", "#AD1457", "#6A1B9A", "#9575CD", "#4527A0", "#311B92", "#283593", "#1976D2", "#0288D1", "#00695C", "#388E3C", "#558B2F", "#E65100", "#A1887F", "#757575", "#DD2C00", "#455A64"];

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

        let header = new Box();
        //without overflow, the height with child-float-elements will be zero
        header.getElement().style.overflow = "auto";
        header.getElement().style.width = "100%";
        let icon = new Image();
        icon.getElement().style.width = "72px";
        icon.getElement().style.cssFloat = "left";
        icon.getElement().style.marginRight = "16px";
        icon.setSrc(this.plugin.getIcon());


        //the section at the right of the icon has a fixed with of 70%, if it does not fit it will break, looks fine for smartphones
        let box = new Box();
        box.getElement().style.width = "70%";
        box.getElement().style.cssFloat = "left";
        box.add(new H6(this.plugin.getName()));

        //add cross links for vendor and categories
        let crossLinks = [];
        crossLinks.push(this.plugin.getVendor().getName());
        for (let category of this.plugin.getCategories()) {
            crossLinks.push(category);
        }

        for (let query of crossLinks) {
            let vendorLink = new FlatButton(query);
            vendorLink.getElement().style.fontSize = "13px";
            vendorLink.setOnClick(_ => {
                this.ctx.getNavigation().forward(UISMarket.NAME(), {q: query});
            });

            box.add(vendorLink);
        }


        header.add(icon);
        header.add(box);


        let rating = this.plugin.getRatings();


        let starBox = new StarBox();
        starBox.setModel(rating.asMap());

        starBox.getElement().style.cssFloat = "right";


        this.add(header);


        this.add(new PullRightBox(starBox));
        this.add(new InstallRemoveIdleArea(this.ctx, this.plugin));

        this.add(new HR());

        this.add(new Body1(this.ctx.getString("ratings")));

        for (let comment of rating.getComments()) {
            let viewComment = new ViewComment(this.ctx);
            viewComment.setModel(comment);
            this.add(viewComment);
        }

    }


}

class InstallRemoveIdleArea extends PullRightBox {
    /**
     *
     * @param {DefaultUserInterfaceState} ctx
     * @param {Plugin} plugin
     */
    constructor(ctx, plugin) {
        super();
        this.ctx = ctx;
        this.plugin = plugin;
        this.setStatusUnknown();

        this.ctx.getApplication().getMarketRepository().getInstallInfo(plugin.getId()).then(info => {
            if (info.installed) {
                if (info.repositoryVersionCurrent !== info.repositoryVersionRemote) {
                    this.setStatusUpdateable();
                } else {
                    this.setStatusUninstallable();
                }
            } else {
                this.setStatusInstallable();
            }
        }).catch(err => {
            if (err instanceof NotFoundException) {
                this.setStatusInstallable();
            } else {
                this.ctx.handleDefaultError(err);
            }
        });
    }

    setStatusUnknown() {
        this.removeAll();
        this.add(new CircularProgressIndicator())
    }

    setStatusUpdateable() {
        let btnUpdate = new RaisedButton();
        btnUpdate.setText(this.ctx.getString("update"));
        btnUpdate.setOnClick(evt => {
            btnUpdate.setText(this.ctx.getString("is_updating"));
            btnUpdate.setEnabled(false);
            this.ctx.getApplication().getMarketRepository().update(this.plugin.getId()).then(_ => {
                this.setStatusUninstallable();
            }).catch(err => this.ctx.handleDefaultError(err));
        });

        this.removeAll();
        this.add(btnUpdate);
    }

    setStatusInstallable() {
        let btnInstall = new RaisedButton();
        btnInstall.setText(this.ctx.getString("install"));
        btnInstall.setOnClick(evt => {
            btnInstall.setText(this.ctx.getString("is_installing"));
            btnInstall.setEnabled(false);
            this.ctx.getApplication().getMarketRepository().install(this.plugin.getId()).then(_ => {
                this.setStatusUninstallable();
            }).catch(err => this.ctx.handleDefaultError(err));
        });

        this.removeAll();
        this.add(btnInstall);
    }

    setStatusUninstallable() {
        let btnUninstall = new RaisedButton();
        btnUninstall.setText(this.ctx.getString("uninstall"));

        btnUninstall.setOnClick(_ => {
            btnUninstall.setText(this.ctx.getString("is_uninstalling"));
            btnUninstall.setEnabled(false);
            this.ctx.getApplication().getMarketRepository().remove(this.plugin.getId()).then(_ => {
                this.setStatusInstallable();
            }).catch(err => this.ctx.handleDefaultError(err));
        });

        this.removeAll();
        this.add(btnUninstall);
    }
}


class ViewComment extends Box {
    /**
     *
     * @param {DefaultUserInterfaceState} ctx
     */
    constructor(ctx) {
        super();
        this.ctx = ctx;
        this.getElement().style.overflow = "auto";
    }

    /**
     *
     * @param {Comment} comment
     */
    setModel(comment) {
        this.removeAll();

        let rightSide = new Box();
        rightSide.add(new H6(comment.getFrom()));
        rightSide.getElement().style.cssFloat = "left";

        let ratingBox = new Box();
        for (let i = 0; i < 5; i++) {
            let ico = null;
            if (i + 1 <= comment.getStars()) {
                ico = new Icon("star");
            } else {
                ico = new Icon("star_border");
            }
            ico.getElement().style.verticalAlign = "middle";
            ico.getElement().style.fontSize = "var(--small-text-size)";
            ico.getElement().style.color = "var(--light-text-color)";
            ratingBox.add(ico);
        }
        let when = new Date(comment.getDate() * 1000).toLocaleDateString(this.ctx.getString("locale"), {
            weekday: 'long',
            year: 'numeric',
            month: 'long',
            day: 'numeric'
        });
        let time = new Body2(when);
        time.getElement().style.display = "inline-block";
        time.getElement().style.verticalAlign = "middle";
        time.getElement().style.margin = "3px";
        time.getElement().style.color = "var(--light-text-color)";
        time.getElement().style.fontSize = "var(--small-text-size)";
        ratingBox.add(time);
        rightSide.add(ratingBox);

        let text = new Body2(comment.getText());
        text.getElement().style.marginTop = "0";
        rightSide.add(text);


        let avatar = new RoundedIcon();
        let firstChar = comment.getFrom().substring(0, 1);
        avatar.setIcon(new Span(firstChar));
        avatar.getElement().style.cssFloat = "left";
        avatar.getElement().style.marginRight = "8px";
        avatar.getElement().style.backgroundColor = colorTable[firstChar.charCodeAt(0) % (colorTable.length - 1)];
        this.add(avatar);

        this.add(rightSide);
    }
}


