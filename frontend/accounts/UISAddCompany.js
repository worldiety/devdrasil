import {DefaultUserInterfaceState, Main} from "/frontend/uis-default.js";
import {Company} from "/frontend/repository/CompanyRepository.js";
import {ViewCompanyForm} from "./ViewCompanyForm.js";

export {UISAddCompany}

class UISAddCompany extends DefaultUserInterfaceState {

    static NAME() {
        return "add-company";
    }

    constructor(app) {
        super(app);
    }

    apply() {
        super.apply();

        let company = new Company();

        let card = new Main();
        this.topBar.setTitle(this.getString("new_company"));

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

            this.getApplication().getCompanyRepository().add(company).then(updatedGroup => {
                window.history.back()
            }).catch(err => {
                companyForm.handleFormError(err);

            });

        });


        this.setContent(card);
    }

}

