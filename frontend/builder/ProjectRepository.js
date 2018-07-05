export {Project, Asset, ProjectRepository}


const VERSION_1 = 1;
const DB = "projects-db";
const STORE_PROJECT = "project";
const STORE_ASSET = "asset";


/**
 * A project index model.
 */
class Project {
    /**
     *
     * @param {string} id
     * @param {string} name
     */
    constructor(id, name) {
        /**
         * The readable project name
         * @type {string}
         */
        this.name = name;

        /**
         * The unique id of a project
         * @type {string}
         */
        this.id = id;
    }
}

/**
 * Internal Asset model
 */
class Asset {
    /**
     *
     * @param {string} id
     * @param {string} path
     * @param {Object} payload
     */
    constructor(id, path, payload) {
        this.id = id;
        this.path = path;
        this.payload = payload;
    }
}

/**
 * Repository to manage projects and the app's specification.
 */
class ProjectRepository {
    constructor() {
        /**
         * In Indexed db everything is about asynchronously.
         * @type {Promise<IDBDatabase>}
         */
        this.dbPromise = new Promise((resolve, reject) => {
            let request = window.indexedDB.open(DB, VERSION_1);
            request.onerror = evt => {
                reject(evt);
            };

            request.onupgradeneeded = evt => {
                let db = request.result;
                if (evt.oldVersion < VERSION_1) {
                    let prjStore = db.createObjectStore(STORE_PROJECT, {keyPath: "id"});
                    let asStore = db.createObjectStore(STORE_ASSET, {keyPath: ["id", "path"]});
                }

                /*
                if (evt.oldVersion < VERSION_2) {
                    db.deleteObjectStore("store1");
                    db.createObjectStore("store2");
                }*/
            };

            request.onsuccess = () => {
                resolve(request.result);
            };
        });
    }

    /**
     *
     * @return {Promise<IDBDatabase>}
     */
    getDB() {
        return this.dbPromise
    }

    /**
     * Finds all projects.
     *
     * @return {Array<Project>}
     */
    async findAllProjects() {
        let res = [];
        await this.query(item => {
            res.push(new Project(item["id"], item["name"]));
        });
        return res;
    }

    /**
     * Inserts or updates the given project into the store
     * @param {Project} project
     * @return {Promise<IDBObjectStore|IDBIndex|IDBCursor>}
     */
    async putProject(project) {
        let db = await this.getDB();
        let tx = db.transaction(STORE_PROJECT, "readwrite");
        let store = tx.objectStore(STORE_PROJECT);
        return requestAsPromise(store.put(project));
    }

    /**
     * Deletes a project
     * @param {string} id
     * @return {Promise<void>}
     */
    async deleteProject(id) {
        let db = await this.getDB();
        let tx = db.transaction(STORE_PROJECT, "readwrite");
        let store = tx.objectStore(STORE_PROJECT);
        return requestAsPromise(store.delete(id));
    }

    /**
     * Inserts or updates an asset using the given path and id
     * @param {Asset} asset
     * @return {Promise<void>}
     */
    async putAsset(asset) {
        let db = await this.getDB();
        let tx = db.transaction(STORE_ASSET, "readwrite");
        let store = tx.objectStore(STORE_ASSET);
        return requestAsPromise(store.put(asset), false);
    }


    /**
     * Retrieves an asset
     * @param id
     * @param path
     * @return {Promise<Asset>}
     */
    async getAsset(id, path) {
        let db = await this.getDB();
        let tx = db.transaction(STORE_ASSET, "readonly");
        let store = tx.objectStore(STORE_ASSET);
        return requestAsPromise(store.get([id, path])).then(res => new Asset(res["id"], res["path"], res["payload"]));
    }

    /**
     * Performs a query using a cursor
     * @param {IDBKeyRange| number | string | Date | IDBArrayKey} range
     * @param {IDBCursorDirection} direction
     * @param {function(Object)} cursorValueCallback
     * @return {Promise<int>}
     */
    async query(cursorValueCallback, range = null, direction = null) {
        let db = await this.getDB();
        let tx = db.transaction(STORE_PROJECT, "readonly");
        let cursorRequest = tx.objectStore(STORE_PROJECT).openCursor(range, direction);
        return new Promise((resolve, reject) => {
            cursorRequest.onerror = evt => {
                reject(evt);
            };
            let visited = 0;
            cursorRequest.onsuccess = evt => {
                let cursor = evt.result;
                if (cursor) {
                    visited++;
                    cursorValueCallback(cursor.value);
                    cursor.continue();
                } else {
                    resolve(visited);
                }
            }
        });
    }
}

/**
 * @param {IDBRequest} request
 * @param {boolean} assertNotDefined
 * @return {Promise<IDBObjectStore | IDBIndex | IDBCursor|Object>}
 */
function requestAsPromise(request, assertNotDefined = true) {
    let callstack = new Error("result is undefined");
    return new Promise((resolve, reject) => {
        request.onerror = evt => {
            reject(evt);
        };

        request.onsuccess = evt => {
            if (assertNotDefined && request.result === undefined) {
                reject(callstack);
                return;
            }

            resolve(request.result);
        };
    });
}