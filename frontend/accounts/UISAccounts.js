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
import {User} from "/frontend/repository/UserRepository.js";
import {ViewUserList} from "./ViewUserList.js"
import {ViewGroupList} from "./ViewGroupList.js";
import {ViewCompanyList} from "./ViewCompanyList.js";

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

        box.add(new Body1(this.getString("manage_companies_hint")));
        let companyList = new ViewCompanyList(this);
        box.add(companyList);

        box.add(new Body1(this.getString("manage_groups_hint")));
        let groupList = new ViewGroupList(this);
        box.add(groupList);

        box.add(new Body1(this.getString("manage_users_hint")));

        let userList = new ViewUserList(this);
        box.add(userList);


        this.setContent(box);

    }


}


