import {Fetcher, throwFromHTTP} from "/wwt/components.js";
import {DefaultRepository} from "./DefaultRepository.js";

export {UserRepository, User}


class User {
    /**
     *
     * @param {string} id
     * @param {string} login
     * @param {string} firstname
     * @param {string} lastname
     * @param {boolean} active
     * @param {string} avatar
     * @param {Array<string>} emails
     * @param {string} company
     * @param {Array<string>} groups
     */
    constructor(id = "", login = "", firstname = "", lastname = "", active = true, avatar = "", emails = [], company = "", groups = []) {
        this.id = id;
        this.login = login;
        this.firstname = firstname;
        this.lastname = lastname;
        this.active = active;
        this.avatar = avatar;
        this.emails = emails;
        this.company = company;
        this.groups = groups;
        this.password = "";
    }


    /**
     *
     * @param {string} gid
     * @returns {boolean}
     */
    hasGroup(gid) {
        for (let id of this.groups) {
            if (id === gid) {
                return true;
            }
        }
        return false;
    }

    /**
     *
     * @param {string} gid
     * @returns {boolean}
     */
    hasCompany(cid) {
        return cid === this.company;
    }
}

class UserPermissions {
    /**
     *
     * @param {boolean} listUsers
     * @param {boolean} createUser
     * @param {boolean} deleteUser
     * @param {boolean} updateUser
     * @param {boolean} getUser
     */
    constructor(listUsers, createUser, deleteUser, updateUser, getUser) {
        this.listUsers = listUsers;
        this.createUser = createUser;
        this.deleteUser = deleteUser;
        this.updateUser = updateUser;
        this.getUser = getUser;

    }
}

class UserRepository extends DefaultRepository {
    constructor(fetcher, sessionProvider) {
        super(fetcher, "users", sessionProvider);
    }


    /**
     * Returns the user either by getting it from cache or by loading it from the endpoint
     * @param {string} id
     * @returns {PromiseLike<UserPermissions>}
     */
    async getUserPermissions(id) {
        let session = await this.sessionProvider.getSession();
        return _requestUserPermissions(this.fetcher, session.sid, id).then(raw => {
            return throwFromHTTP(raw).then(raw => raw.json());
        }).then(json => {
            return new UserPermissions(json["ListUsers"], json["CreateUser"], json["DeleteUser"], json["UpdateUser"], json["GetUser"]);
        });

    }

    /**
     *  Returns the user either by getting it from cache or by loading it from the endpoint
     * @returns {PromiseLike<User>}
     */
    async getSessionUser() {
        let session = await this.sessionProvider.getSession();
        return this.get(session.uid);
    }


    /**
     *
     * @param json
     * @return User
     */
    fromJson(json) {
        return new User(json["Id"], json["Login"], json["Firstname"], json["Lastname"], json["Active"], json["Avatar"], json["EMailAddresses"], json["Company"], json["Groups"] == null ? [] : json["Groups"])
    }

    /**
     * Requests all users from the given id list
     * @param {Array<string>} ids
     * @returns {Promise<Array<User>>}
     */
    async getUsersByIds(ids) {
        let users = [];
        for (let uid of ids) {
            let user = await this.ctx.getApplication().getUserRepository().get(uid);
            users.push(user);
        }
        return users;
    }
}


function _requestUserPermissions(fetcher, sid, uid) {
    let cfg = {
        method: 'GET',
        headers: {
            'sid': sid,
        },
        cache: 'no-store'
    };
    return fetcher.fetchRaw('/users/permissions/' + uid, cfg)
}

