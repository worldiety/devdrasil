import {Fetcher, throwFromHTTP} from "/wwt/components.js";

export {UserRepository, User}


class User {
    constructor(id, isActive) {
        this.id = id;
        this.isActive = isActive;
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

class UserRepository {
    constructor(fetcher, sessionRepository) {
        this.fetcher = fetcher;
        this.sessionRepository = sessionRepository;
    }


    /**
     * Returns the user either by getting it from cache or by loading it from the endpoint
     * @param {string} id
     * @returns {PromiseLike<User>}
     */
    async getUser(id) {
        let session = await this.sessionRepository.getSession();
        return _requestUser(this.fetcher, session.sid, id).then(raw => {
            throwFromHTTP(raw);
            return raw.json();
        }).then(json => {
            console.log(json);
            return new User(json.Id, json.Active);
        });

    }

    /**
     * Returns the user either by getting it from cache or by loading it from the endpoint
     * @param {string} id
     * @returns {PromiseLike<UserPermissions>}
     */
    async getUserPermissions(id) {
        let session = await this.sessionRepository.getSession();
        return _requestUserPermissions(this.fetcher, session.sid, id).then(raw => {
            throwFromHTTP(raw);
            return raw.json();
        }).then(json => {
            console.log(json);
            return new UserPermissions(json.ListUsers, json.CreateUser, json.DeleteUser, json.UpdateUser, json.GetUser);
        });

    }

    /**
     *  Returns the user either by getting it from cache or by loading it from the endpoint
     * @returns {Promise<User>}
     */
    async getSessionUser() {
        let session = await this.sessionRepository.getSession();
        return _requestUser(this.fetcher, session.sid, session.uid).then(raw => {
            throwFromHTTP(raw);
            return raw.json();
        }).then(json => {
            return new User(json.Id, json.Active);
        });
    }

    /**
     * Returns all available users
     *
     * @returns {!PromiseLike<[]User>}
     */
    async getUsers() {
        let session = await this.sessionRepository.getSession();
        return _requestUsers(this.fetcher, session.sid).then(raw => {
            throwFromHTTP(raw);
            return raw.json();
        }).then(json => {
            let users = [];
            for (let user of json.Users) {
                users.push(new User(user.Id, user.Properties, user.Plugins, user.Active));
            }
            return users;
        });

    }
}


function _requestUser(fetcher, sid, uid) {
    let cfg = {
        method: 'GET',
        headers: {
            'sid': sid,
        },
        cache: 'no-store'
    };
    return fetcher.fetchRaw('/users/' + uid, cfg)
}

function _requestUsers(fetcher, sid) {
    let cfg = {
        method: 'GET',
        headers: {
            'sid': sid,
        },
        cache: 'no-store'
    };
    return fetcher.fetchRaw('/users', cfg)
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