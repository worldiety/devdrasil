import {Fetcher, throwFromHTTP} from "/wwt/components.js";

export {UserRepository, User, ROLE_ADD_USER, ROLE_LIST_USER}

const ROLE_ADD_USER = "ROLE_ADD_USER";
const ROLE_LIST_USER = "ROLE_LIST_USER";

class User {
    constructor(id, properties, plugins, isActive) {
        this.id = id;
        this.properties = properties;
        this.plugins = plugins;
        this.isActive = isActive;
    }

    hasProperty(propertyName) {
        if (this.properties == null) {
            return false;
        }
        return propertyName in this.properties;
    }
}

class UserRepository {
    constructor(fetcher, sessionRepository) {
        this.fetcher = fetcher;
        this.sessionRepository = sessionRepository;
    }


    /**
     * Returns the user either by getting it from cache or by loading it from the endpoint
     * @returns {!PromiseLike<User>}
     */
    async getUser() {
        let session = await this.sessionRepository.getSession();
        return _requestUser(this.fetcher, session.id).then(raw => {
            throwFromHTTP(raw);
            return raw.json();
        }).then(json => {
            return new User(json.Id, json.Properties, json.Plugins, json.Active);
        });

    }

    /**
     * Returns all available users
     *
     * @returns {!PromiseLike<[]User>}
     */
    async getUsers() {
        let session = await this.sessionRepository.getSession();
        return _requestUsers(this.fetcher, session.id).then(raw => {
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

function _requestUser(fetcher, id) {
    let cfg = {
        method: 'GET',
        headers: {
            'sid': id,
        },
        cache: 'no-store'
    };
    return fetcher.fetchRaw('/users/account', cfg)
}

function _requestUsers(fetcher, id) {
    let cfg = {
        method: 'GET',
        headers: {
            'sid': id,
        },
        cache: 'no-store'
    };
    return fetcher.fetchRaw('/users/all', cfg)
}