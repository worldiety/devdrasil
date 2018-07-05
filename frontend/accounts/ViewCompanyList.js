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
    showConfirmationDialog,
    Span,
    TwoLineLeadingAndTrailingIcon,
    UserInterfaceState
} from "/wwt/components.js";
import {UISAddCompany} from "./UISAddCompany.js";
import {UISEditCompany} from "./UISEditCompany.js";

export {ViewCompanyList}

class ViewCompanyList extends Card {
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
        addUserItem.addLeft(new Span(this.uis.getString("your_companies")));
        let btnAddUser = new Button();
        btnAddUser.setText(this.uis.getString("add"));
        btnAddUser.setOnClick(e => {
            this.uis.getNavigation().forward(UISAddCompany.NAME());
        });
        addUserItem.addRight(btnAddUser);


        let spinner = new CircularProgressIndicator();
        spinner.getElement().style.margin = "auto";
        listView.add(spinner);

        this.add(listView);
        this.uis.getApplication().getCompanyRepository().list().then(companies => {
            listView.removeAll();
            listView.add(addUserItem);
            let first = true;
            for (let company of companies) {
                if (first) {
                    first = false;
                } else {
                    listView.addSeparator();
                }
                let moreMenu = new Menu();
                moreMenu.add(this.uis.getString("edit"), _ => {
                    this.uis.getNavigation().forward(UISEditCompany.NAME(), {"cid": company.id});
                });
                moreMenu.add(this.uis.getString("delete"), _ => {
                    this.delete(company);

                });

                let more = new Button();
                more.setIcon(new Icon("more_vert"));
                more.setOnClick(e => {
                    moreMenu.popup(more);
                });


                let item = new TwoLineLeadingAndTrailingIcon(new Icon("business"), company.name, this.uis.getString("x_employees", company.users.length + ""), more);
                listView.add(item);
            }
        });
    }

    /**
     *
     * @param {Group} group
     */
    delete(group) {
        showConfirmationDialog(this.uis, this.uis.getString("delete_x", group.name), this.uis.getString("cancel"), this.uis.getString("delete"), () => {
            this.uis.getApplication().getCompanyRepository().delete(group.id).then(_ => this.refresh()).catch(err => this.uis.handleDefaultError(err));
        });
    }
}