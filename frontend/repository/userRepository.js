import {Fetcher} from "/wwt/components.js";

export {UserRepository, User}

class User {

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
    async getUsers(){

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