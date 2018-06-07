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

import {DefaultUserInterfaceState, Main} from "/frontend/uis-default.js";
import {UISUser} from "/frontend/uis-user.js";
import {User} from "/frontend/repository/userRepository.js";
import {UISAddUser} from "/frontend/uis-user-add.js";

export {UISAccounts}

class UISAccounts extends DefaultUserInterfaceState {

    static NAME() {
        return "accounts";
    }

    constructor(app) {
        super(app);
    }

    apply() {
        super.apply();
        this.getTopBar().setTitle(this.getString("accounts"));

        let box = new Main(false);

        box.add(new Body1(this.getString("manage_groups_hint")));
        let card = new Card();
        box.add(card);

        box.add(new Body1(this.getString("manage_users_hint")));

        let userList = new UserList(this);
        box.add(userList);


        this.setContent(box);

    }


}


class UserList extends Card {
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
        this.uis.getApplication().getUserRepository().getUsers().then(users => {
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
                    this.uis.getNavigation().forward(UISUser.NAME(), {"uid": user.id});
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
            this.uis.getApplication().getUserRepository().deleteUser(user.id).then(_ => this.refresh()).catch(err => this.uis.handleDefaultError(err));
            dlg.close();
        });
        dlg.addFooter(no);
        dlg.addFooter(yes);
        dlg.show();
    }
}