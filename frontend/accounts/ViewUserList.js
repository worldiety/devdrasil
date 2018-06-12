import {
    AbsComponent,
    Body1,
    Button,
    Card,
    CircularProgressIndicator,
    Dialog,
    H3,
    Icon,
    ListView,
    LRLayout,
    Menu,
    Span,
    TwoLineLeadingAndTrailingIcon,
    UserInterfaceState
} from "/wwt/components.js";

import {UISEditUser} from "./UISEditUser.js";
import {UISAddUser} from "./UISAddUser.js";

export {ViewUserList}

class ViewUserList extends Card {
    /**
     *
     * @param {DefaultUserInterfaceState} uis
     */
    constructor(uis) {
        super();
        this.uis = uis;

        this.getElement().style.padding = "0";

        this.refresh();
    }

    refresh() {
        this.removeAll();
        let listView = new ListView();

        let addUserItem = new LRLayout();
        addUserItem.addLeft(new Span(this.uis.getString("your_users")));
        let btnAddUser = new Button();
        btnAddUser.setText(this.uis.getString("add"));
        btnAddUser.setOnClick(e => {
            this.uis.getNavigation().forward(UISAddUser.NAME());
        });
        addUserItem.addRight(btnAddUser);


        let spinner = new CircularProgressIndicator();
        spinner.getElement().style.margin = "auto";
        listView.add(spinner);

        this.add(listView);
        this.uis.getApplication().getUserRepository().list().then(users => {
            listView.removeAll();
            listView.add(addUserItem);
            let first = true;
            for (let user of users) {
                if (first) {
                    first = false;
                } else {
                    listView.addSeparator();
                }
                let moreMenu = new Menu();
                moreMenu.add(this.uis.getString("edit"), _ => {
                    this.uis.getNavigation().forward(UISEditUser.NAME(), {"uid": user.id});
                });
                moreMenu.add(this.uis.getString("delete"), _ => {
                    this.delete(user);

                });

                let more = new Button();
                more.setIcon(new Icon("more_vert"));
                more.setOnClick(e => {
                    moreMenu.popup(more);
                });


                let item = new TwoLineLeadingAndTrailingIcon(new Icon("account_circle"), user.firstname + " " + user.lastname, user.login, more);
                listView.add(item);
            }
        });
    }

    /**
     *
     * @param {User} user
     */
    delete(user) {
        let dlg = new Dialog();
        dlg.add(new Body1(this.uis.getString("delete_user_x", user.firstname + " " + user.lastname)));
        let no = new Button(this.uis.getString("cancel"));
        no.setOnClick(_ => {
            dlg.close();
        });
        let yes = new Button(this.uis.getString("delete"));
        yes.setOnClick(_ => {
            this.uis.getApplication().getUserRepository().delete(user.id).then(_ => this.refresh()).catch(err => this.uis.handleDefaultError(err));
            dlg.close();
        });
        dlg.addFooter(no);
        dlg.addFooter(yes);
        dlg.show();
    }
}