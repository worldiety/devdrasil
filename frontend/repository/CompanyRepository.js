import {Fetcher, throwFromHTTP} from "/wwt/components.js";
import {DefaultRepository, SessionProvider} from "./DefaultRepository.js";

export {CompanyRepository, Company}


class Company {
    constructor(id = "", name = "", users = [], themePrimaryColor = "") {
        this.id = id;
        this.name = name;
        this.users = users;
        this.themePrimaryColor = themePrimaryColor;
    }
}

/**
 * A group repository represents the REST endpoint for groups.
 */
class CompanyRepository extends DefaultRepository {
    /**
     * Creates a new repository
     * @param {Fetcher} fetcher
     * @param {SessionProvider} sessionProvider
     */
    constructor(fetcher, sessionProvider) {
        super(fetcher, "companies", sessionProvider)
    }


    fromJson(json) {
        return new Company(json["Id"], json["Name"], json["Users"], json["ThemePrimaryColor"]);
    }
}

