import {
    AbsComponent,
    Body1,
    Button,
    Card,
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
import {ROLE_LIST_USER, User} from "/frontend/repository/userRepository.js";

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
        let listView = new ListView();

        let addUserItem = new LRLayout();
        addUserItem.addLeft(new Span(uis.getString("your_users")));
        let btnAddUser = new Button();
        btnAddUser.setText(uis.getString("add"));
        btnAddUser.setOnClick(e => {
            this.uis.showMessage("add user");
        });
        addUserItem.addRight(btnAddUser);
        listView.add(addUserItem);


        this.add(listView);
        this.uis.getApplication().getUserRepository().getUsers().then(users => {
            let first = true;
            for (let user of users) {
                if (first) {
                    first = false;
                } else {
                    listView.addSeparator();
                }
                let tmp = "";
                for (let key in user.properties) {
                    if (tmp !== "") {
                        tmp += ", ";
                    }


                    tmp += key;
                }
                let moreMenu = new Menu();
                moreMenu.add(uis.getString("edit"), _ => {
                    this.uis.showMessage("bearbeiten");
                });
                moreMenu.add(uis.getString("delete"), _ => {
                    this.delete(user);

                });

                let more = new Button();
                more.setIcon(new Icon("more_vert"));
                more.setOnClick(e => {
                    moreMenu.popup(more);
                });


                let item = new TwoLineLeadingAndTrailingIcon(new Icon("account_circle"), user.id, tmp, more);
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
        dlg.add(new Body1(this.uis.getString("delete_user_x", user.id)));
        let no = new Button(this.uis.getString("cancel"));
        no.setOnClick(_ => {
            dlg.close();
        });
        let yes = new Button(this.uis.getString("delete"));
        yes.setOnClick(_ => {
            dlg.close();
        })
        dlg.addFooter(no);
        dlg.addFooter(yes);
        dlg.show();
    }
}