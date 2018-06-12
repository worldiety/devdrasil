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
import {UISAddGroup} from "./UISAddGroup.js";
import {UISEditGroup} from "./UISEditGroup.js";

export {ViewGroupList}

class ViewGroupList extends Card {
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
        addUserItem.addLeft(new Span(this.uis.getString("your_groups")));
        let btnAddUser = new Button();
        btnAddUser.setText(this.uis.getString("add"));
        btnAddUser.setOnClick(e => {
            this.uis.getNavigation().forward(UISAddGroup.NAME());
        });
        addUserItem.addRight(btnAddUser);


        let spinner = new CircularProgressIndicator();
        spinner.getElement().style.margin = "auto";
        listView.add(spinner);

        this.add(listView);
        this.uis.getApplication().getGroupRepository().list().then(groups => {
            listView.removeAll();
            listView.add(addUserItem);
            let first = true;
            for (let group of groups) {
                if (first) {
                    first = false;
                } else {
                    listView.addSeparator();
                }
                let moreMenu = new Menu();
                moreMenu.add(this.uis.getString("edit"), _ => {
                    this.uis.getNavigation().forward(UISEditGroup.NAME(), {"gid": group.id});
                });
                moreMenu.add(this.uis.getString("delete"), _ => {
                    this.delete(group);

                });

                let more = new Button();
                more.setIcon(new Icon("more_vert"));
                more.setOnClick(e => {
                    moreMenu.popup(more);
                });


                let item = new TwoLineLeadingAndTrailingIcon(new Icon("supervised_user_circle"), group.name, this.uis.getString("x_members", group.users.length + ""), more);
                listView.add(item);
            }
        });
    }

    /**
     *
     * @param {Group} group
     */
    delete(group) {
        let dlg = new Dialog();
        dlg.add(new Body1(this.uis.getString("delete_group_x", group.name)));
        let no = new Button(this.uis.getString("cancel"));
        no.setOnClick(_ => {
            dlg.close();
        });
        let yes = new Button(this.uis.getString("delete"));
        yes.setOnClick(_ => {
            this.uis.getApplication().getGroupRepository().delete(group.id).then(_ => this.refresh()).catch(err => this.uis.handleDefaultError(err));
            dlg.close();
        });
        dlg.addFooter(no);
        dlg.addFooter(yes);
        dlg.show();
    }
}