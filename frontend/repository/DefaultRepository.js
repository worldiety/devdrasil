import {Fetcher, throwFromHTTP} from "/wwt/components.js";

export {DefaultRepository, SessionProvider}

/**
 * A simple contract for a session provider.
 */
class SessionProvider {
    /**
     *
     * @returns {Promise<Session>}
     */
    async getSession() {
        throw Error("abstract method must be overridden");
    }
}

/**
 * A default REST repository implementation, supporting a single resource endpoint with List and CRUD features.
 * In addition to that it may cache data for offline purposes (some day).
 */
class DefaultRepository {

    /**
     * Creates a new repository
     * @param {Fetcher} fetcher, the instance to fetch the http endpoints
     * @param {string} resourceName the name of the rest resource
     * @param {SessionProvider} sessionProvider the repository which holds the session
     */
    constructor(fetcher, resourceName, sessionProvider) {
        this.fetcher = fetcher;
        this.resourceName = resourceName;
        this.sessionProvider = sessionProvider;
    }


    /**
     * Converts an untyped object into the concrete typed entity class
     *
     * @param {Object} json
     * @return {Object}
     * @abstract
     */
    fromJson(json) {
        throw Error("abstract method must be overridden");
    }

    /**
     * Returns the entity identified by id
     * @param {string} id
     * @returns {PromiseLike<Object>}
     */
    async get(id) {
        let session = await this.sessionProvider.getSession();
        return restGet(this.fetcher, this.resourceName, session.sid, id).then(raw => {
            return throwFromHTTP(raw).then(raw => raw.json());
        }).then(json => {
            return this.fromJson(json);
        });
    }

    /**
     * Deletes the group identified by id
     * @param {string} id
     * @returns {PromiseLike<void>}
     */
    async delete(id) {
        let session = await this.sessionProvider.getSession();
        return restDelete(this.fetcher, this.resourceName, session.sid, id).then(raw => {
            return throwFromHTTP(raw);
        });
    }

    /**
     * Updates the given group
     * @param {Group} entity with an existing id
     * @returns {PromiseLike<Object>}
     */
    async update(entity) {
        let session = await this.sessionProvider.getSession();
        return restUpdate(this.fetcher, this.resourceName, session.sid, entity).then(raw => {
            return throwFromHTTP(raw).then(raw => raw.json());
        }).then(json => {
            return this.fromJson(json);
        });
    }

    /**
     * Inserts a new entity
     * @param {Group} entity
     * @returns {PromiseLike<Object>}
     */
    async add(entity) {
        let session = await this.sessionProvider.getSession();
        return restAdd(this.fetcher, this.resourceName, session.sid, entity).then(raw => {
            return throwFromHTTP(raw).then(raw => raw.json());
        }).then(json => {
            return this.fromJson(json);
        });
    }


    /**
     * Lists all entities
     * @returns {PromiseLike<Array<Object>>}
     */
    async list() {
        let session = await this.sessionProvider.getSession();
        return restList(this.fetcher, this.resourceName, session.sid).then(raw => {
            return throwFromHTTP(raw).then(raw => raw.json());
        }).then(json => {
            let list = [];
            for (let entry of json["List"]) {
                list.push(this.fromJson(entry));
            }
            return list;
        });
    }

}

function restGet(fetcher, name, sid, id) {
    let cfg = {
        method: 'GET',
        headers: {
            'sid': sid,
        },
        cache: 'no-store'
    };
    return fetcher.fetchRaw("/" + name + "/" + id, cfg)
}

function restDelete(fetcher, name, sid, id) {
    let cfg = {
        method: 'DELETE',
        headers: {
            'sid': sid,
        },
        cache: 'no-store'
    };
    return fetcher.fetchRaw("/" + name + "/" + id, cfg)
}

function restUpdate(fetcher, name, sid, entity) {
    let cfg = {
        method: 'PUT',
        headers: {
            'sid': sid,
        },
        body: JSON.stringify(entity),
        cache: 'no-store'
    };
    return fetcher.fetchRaw("/" + name + "/" + entity.id, cfg)
}

function restList(fetcher, name, sid) {
    let cfg = {
        method: 'GET',
        headers: {
            'sid': sid,
        },
        cache: 'no-store'
    };
    return fetcher.fetchRaw("/" + name, cfg)
}

function restAdd(fetcher, name, sid, entity) {
    let cfg = {
        method: 'POST',
        headers: {
            'sid': sid,
        },
        body: JSON.stringify(entity),
        cache: 'no-store'
    };
    return fetcher.fetchRaw("/" + name, cfg)
}