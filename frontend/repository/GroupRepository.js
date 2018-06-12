import {Fetcher, throwFromHTTP} from "/wwt/components.js";
import {DefaultRepository, SessionProvider} from "./DefaultRepository.js";

export {GroupRepository, Group}


class Group {
    constructor(id = "", name = "", users = []) {
        this.id = id;
        this.name = name;
        this.users = users;
    }
}

/**
 * A group repository represents the REST endpoint for groups.
 */
class GroupRepository extends DefaultRepository {
    /**
     * Creates a new repository
     * @param {Fetcher} fetcher
     * @param {SessionProvider} sessionProvider
     */
    constructor(fetcher, sessionProvider) {
        super(fetcher, "groups", sessionProvider)
    }


    fromJson(json) {
        return new Group(json["Id"], json["Name"], json["Users"]);
    }
}

