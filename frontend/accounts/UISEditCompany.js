import {DefaultUserInterfaceState, Main} from "/frontend/uis-default.js";
import {ViewCompanyForm} from "./ViewCompanyForm.js";
import {Box, Card, CenterBox, CircularProgressIndicator,} from "/wwt/components.js";

export {UISEditCompany}

class UISEditCompany extends DefaultUserInterfaceState {

    static NAME() {
        return "edit-company";
    }

    constructor(app) {
        super(app);
    }

    apply() {

        super.apply();


        let uid = this.getNavigation().getSearchParam("cid");

        this.getApplication().getCompanyRepository().get(uid).then(company => {
            let card = new Main();
            this.topBar.setTitle(company.name);

            let companyForm = new ViewCompanyForm(this, company);
            card.add(companyForm);

            companyForm.btnLogin.setOnClick(e => {
                companyForm.clearTips();


                //quick evaluation
                let hasError = companyForm.checkDefault();


                if (hasError) {
                    return
                }

                companyForm.updateModel();

                this.getApplication().getCompanyRepository().update(company).then(updatedGroup => {
                    //we succeeded
                    window.history.back()

                }).catch(err => {
                    if (!companyForm.handleFormError(err)) {
                        this.handleDefaultError(err);
                    }

                });

            });


            this.setContent(card);

        }).catch(err => {
            this.handleDefaultError(err);
        });


        let box = new CenterBox();
        let spinner = new CircularProgressIndicator();
        box.add(spinner);
        this.setContent(box);


    }

}

