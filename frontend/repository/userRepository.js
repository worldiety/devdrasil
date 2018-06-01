import {Fetcher, throwFromHTTP} from "/wwt/components.js";

export {UserRepository, User}

class User {
    constructor(id, properties, plugins, isActive) {
        this.id = id;
        this.properties = properties;
        this.plugins = plugins;
        this.isActive = isActive;
    }
}

class UserRepository {
    constructor(fetcher, sessionRepository) {
        this.fetcher = fetcher;
        this.sessionRepository = sessionRepository;
    }


    /**
     * Returns the user either by getting it from cache or by loading it from the endpoint
     * @returns {Promise<User>}
     */
    async getUser() {

        return new User();
    }

    /**
     * Returns all available users
     *
     * @returns {Promise<[]User>}
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
            'login': user,
            'password': password,
            'client': client,
        },
        cache: 'no-store'
    };
    return fetcher.fetchRaw('/session/auth', cfg)
}

function _requestUsers(fetcher, id) {
    let cfg = {
        method: 'GET',
        headers: {
            'sid': id,
        },
        cache: 'no-store'
    };
    return fetcher.fetchRaw('/users', cfg)
}