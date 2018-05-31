import {throwFromHTTP} from "/wwt/components.js";

export {Session, SessionRepository}

const KEY_ID = "sessionId";
const KEY_LOGIN = "sessionLogin";

/**
 * The session pojo
 */
class Session {
    constructor(id) {
        this.id = id;
    }
}

/**
 * The session repository which caches and requests sessions
 */
class SessionRepository {

    constructor(fetcher) {
        this.fetcher = fetcher;
        this.memcache = null;
        this.login = localStorage.getItem(KEY_LOGIN);
    }

    /**
     * Removes the session
     */
    async deleteSession() {
        localStorage.removeItem(KEY_LOGIN);
        localStorage.removeItem(KEY_ID);
        this.memcache = null;
        this.login = null;
    }


    /**
     * Returns or requests the session. If login is not empty validates if the ondisk cache needs to get purged (e.g. login a different user)
     *
     * @param login the login name {string|null}
     * @param password the password {string|null}
     * @param client the client token {string|null}
     * @throws {IOException|PermissionDeniedException}
     * @returns {Promise<Session>}
     */
    async getSession(login = null, password = "", client = "") {
        if ((login == null || login === this.login) && this.login != null) {
            if (this.memcache == null) {
                this.memcache = localStorage.getItem(KEY_ID);
                if (this.memcache != null) {
                    return new Session(this.memcache);
                }
            } else {
                return new Session(this.memcache);
            }
        }
        //requires update
        let sessionPromise = _requestSession(this.fetcher, login, password, client);
        return sessionPromise.then(raw => {
            throwFromHTTP(raw);
            return raw.json();
        }).then(json => {
            this.memcache = json.SessionId;
            this.login = login;
            localStorage.setItem(KEY_ID, this.memcache);
            localStorage.setItem(KEY_LOGIN, this.login);
            return new Session(this.memcache);
        });

    }


}

function _requestSession(fetcher, user, password, client) {
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