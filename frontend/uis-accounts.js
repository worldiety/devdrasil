import {
    AbsComponent,
    Body1,
    Button,
    Card,
    H3,
    Icon,
    ListView,
    LRLayout,
    Span,
    TwoLineLeadingAndTrailingIcon,
    UserInterfaceState
} from "/wwt/components.js";

import {DefaultUserInterfaceState, Main} from "/frontend/uis-default.js";
import {ROLE_LIST_USER} from "/frontend/repository/userRepository.js";

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
     * @param {UserInterfaceState} uis
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
            alert("add user");
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
                let more = new Button("test");
                more.setOnClick(e => {
                    alert("mehr");
                });
                let item = new TwoLineLeadingAndTrailingIcon(new Icon("account_circle"), user.id, tmp, more);
                listView.add(item);
            }
        });
    }
}