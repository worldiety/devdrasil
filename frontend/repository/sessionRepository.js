import {throwFromHTTP} from "/wwt/components.js";

export {Session, SessionRepository}

const KEY_SID = "_sid";
const KEY_LOGIN = "_login";
const KEY_UID = "_uid";

/**
 * The session pojo
 */
class Session {
    /**
     *
     * @param {string} sid
     * @param {string} uid
     */
    constructor(sid, uid) {
        this.sid = sid;
        this.uid = uid;
    }
}

/**
 * The session repository which caches and requests sessions
 */
class SessionRepository {

    constructor(fetcher) {
        this.fetcher = fetcher;
        this.memCacheSID = localStorage.getItem(KEY_SID);
        this.memCacheUID = localStorage.getItem(KEY_UID);
        this.memCacheLogin = localStorage.getItem(KEY_LOGIN);

        //consistency check
        if (this.memCacheLogin == null || this.memCacheSID == null || this.memCacheUID == null) {
            this.memCacheSID = "";
            this.memCacheLogin = "";
            this.memCacheUID = "";
        }
    }


    /**
     * Removes the session
     */
    deleteSession() {
        localStorage.removeItem(KEY_LOGIN);
        localStorage.removeItem(KEY_UID);
        localStorage.removeItem(KEY_SID);
        this.memCacheSID = "";
        this.memCacheLogin = "";
        this.memCacheUID = "";
    }


    /**
     * Returns or requests the session. If login is not empty validates if the ondisk cache needs to get purged (e.g. login a different user)
     *
     * @param {string} login the login name
     * @param {string} password the password
     * @param {string} client the client token
     * @throws {IOException|PermissionDeniedException}
     * @returns {Promise<Session>}
     */
    async getSession(login = "", password = "", client = "") {
        let requiresLogin = this.memCacheLogin === "" || (login !== this.memCacheLogin && login !== "");

        if (requiresLogin) {
            //requires update
            let sessionPromise = _requestSession(this.fetcher, login, password, client);
            return sessionPromise.then(raw => {
                throwFromHTTP(raw);
                return raw.json();
            }).then(json => {
                this.memCacheSID = json.Id;
                this.memCacheLogin = login;
                this.memCacheUID = json.User;
                localStorage.setItem(KEY_SID, this.memCacheSID);
                localStorage.setItem(KEY_UID, this.memCacheUID);
                localStorage.setItem(KEY_LOGIN, this.memCacheLogin);
                return new Session(this.memCacheSID, this.memCacheUID);
            });
        } else {
            return new Session(this.memCacheSID, this.memCacheUID);
        }


    }


}

function _requestSession(fetcher, user, password, client) {
    let cfg = {
        method: 'POST',
        headers: {
            'login': user,
            'password': password,
            'client': client,
            'User-Agent': navigator.userAgent,
        },
        cache: 'no-store'
    };
    return fetcher.fetchRaw('/sessions', cfg)
}