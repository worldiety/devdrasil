import {Fetcher, throwFromHTTP} from "/wwt/components.js";
import {SessionProvider} from "./DefaultRepository.js";

export {MarketRepository, Plugin, MarketIndex, Vendor, Rating}


class MarketIndex {
    constructor(json) {
        this.json = json;
    }

    /**
     * @param {string} filter
     * @returns {Array<Plugin>}
     */
    getPlugins(filter = "") {
        let tmp = [];
        for (let plugin of this.json["plugins"]) {
            if (filter === "") {
                tmp.push(new Plugin(this, plugin));
            } else {
                let res = new Set();
                nestSearch(plugin, filter, 1, res);
                if (res.size > 0) {
                    tmp.push(new Plugin(this, plugin));
                } else {
                    //if not found for plugin, check it's vendor
                    nestSearch(this.getVendor(plugin["vendor"]), filter, 1, res);
                    if (res.size > 0) {
                        tmp.push(new Plugin(this, plugin));
                    }
                }
            }


        }
        return tmp;
    }

    /**
     *
     * @param id
     * @returns {Plugin|null}
     */
    getPlugin(id) {
        for (let plugin of this.json["plugins"]) {
            if (plugin["id"] === id) {
                return new Plugin(this, plugin);
            }
        }
        return null;
    }

    /**
     *
     * @param id
     * @returns {Vendor|null}
     */
    getVendor(id) {
        for (let vendor of this.json["vendors"]) {
            if (vendor["id"] === id) {
                return new Vendor(this, vendor);
            }
        }
        return null;
    }

    /**
     * finds all words which contains text based on anything in the index
     *
     * @param {string} text
     * @param {int} max
     * @return {Set<string>}
     */
    findWord(text, max) {
        let res = new Set();
        nestSearch(this.json, text, max, res);
        return res
    }
}

function nestSearch(root, text, max, res) {
    text = text.toLowerCase();
    for (let key in root) {
        let value = root[key];
        if (typeof value === 'string') {
            if (value.startsWith("http")) {
                continue;
            }
            if (value === "") {
                continue;
            }
            if (value.toLowerCase().indexOf(text) >= 0) {
                res.add(value);
            }
        } else {
            if (Array.isArray(value)) {
                for (let entry of value) {
                    if (typeof value === 'string') {
                        if (value.startsWith("http")) {
                            continue;
                        }
                        if (value === "") {
                            continue;
                        }
                        if (value.toLowerCase().indexOf(text) >= 0) {
                            res.add(value);
                        }
                    } else {
                        nestSearch(value, text, max, res);
                    }
                    if (res.length >= max) {
                        return res
                    }
                }
            } else {
                nestSearch(value, text, max, res);
            }

        }
        if (res.length >= max) {
            return res
        }
    }
}

class Vendor {
    constructor(parent, json) {
        this.parent = parent;
        this.json = json;
    }

    getName() {
        return this.json["name"];
    }
}

class Plugin {
    constructor(parent, json) {
        this.parent = parent;
        this.json = json;
    }

    getName() {
        return this.json["name"];
    }

    getIcon() {
        return this.json["iconUrl"];
    }

    getId() {
        return this.json["id"];
    }

    /**
     *
     * @returns {Vendor}
     */
    getVendor() {
        return this.parent.getVendor(this.json["vendor"]);
    }

    /**
     *
     * @returns {Array<string>}
     */
    getCategories() {
        return this.json["categories"];
    }

    /**
     *
     * @returns {Rating}
     */
    getRatings() {
        return new Rating(this, this.json["rating"])
    }
}

class Rating {
    /**
     *
     * @param {Plugin} parent
     * @param json
     */
    constructor(parent, json) {
        this.parent = parent;
        this.json = json;
    }


    /**
     * @return {Map<int,int>}
     */
    asMap() {
        let res = new Map();
        res.set(1, this.json["s1"]);
        res.set(2, this.json["s2"]);
        res.set(3, this.json["s3"]);
        res.set(4, this.json["s4"]);
        res.set(5, this.json["s5"]);
        return res;
    }

    /**
     *
     * @return {Array<Comment>}
     */
    getComments() {
        let tmp = [];
        for (let comment of this.json["comments"]) {
            tmp.push(new Comment(this, comment));
        }
        return tmp;
    }
}

class Comment {
    /**
     *
     * @param {Rating} parent
     * @param json
     */
    constructor(parent, json) {
        this.parent = parent;
        this.json = json;
    }

    getStars() {
        return this.json["stars"];
    }

    getDate() {
        return this.json["date"];
    }

    getText() {
        return this.json["text"];
    }

    getFrom() {
        return this.json["from"];
    }
}


class PluginInstallInfo {
    /**
     *
     * @param {string} id
     * @param {string} version
     */
    constructor(id, version) {
        this.id = id;
        this.version = version;
    }
}

/**
 * A group repository represents the REST endpoint for groups.
 */
class MarketRepository {
    /**
     * Creates a new repository
     * @param {Fetcher} fetcher
     * @param {SessionProvider} sessionProvider
     */
    constructor(fetcher, sessionProvider) {
        this.fetcher = fetcher;
        this.sessionProvider = sessionProvider;
    }

    /**
     * Requests the market index.
     *
     * @return MarketIndex
     */
    async getIndex() {
        let session = await this.sessionProvider.getSession();
        return restIndex(this.fetcher, session.sid).then(raw => {
            return throwFromHTTP(raw).then(raw => raw.json());
        }).then(json => {
            return this.fromJson(json);
        });
    }


    /**
     * Tries to install the given plugin id
     * @param {string} pluginId
     * @return {Promise<void>}
     */
    async install(pluginId) {
        let session = await this.sessionProvider.getSession();
        return restInstall(this.fetcher, session.sid, pluginId).then(raw => {
            return throwFromHTTP(raw);
        });
    }

    /**
     * Tries to remove the given plugin id
     * @param {string} pluginId
     * @return {Promise<void>}
     */
    async remove(pluginId) {
        let session = await this.sessionProvider.getSession();
        return restDelete(this.fetcher, session.sid, pluginId).then(raw => {
            return throwFromHTTP(raw);
        });
    }

    /**
     * Tries to get the current install info for the plugin
     * @param pluginId
     * @return {PromiseLike<PluginInstallInfo>}
     */
    async getInstallInfo(pluginId) {
        let session = await this.sessionProvider.getSession();
        return restInstallInfo(this.fetcher, session.sid, pluginId).then(raw => {
            return throwFromHTTP(raw).then(raw => raw.json());
        }).then(json => {
            return new PluginInstallInfo(json["Id"], json["Version"])
        });
    }

    fromJson(json) {
        return new MarketIndex(json);
    }
}

function restIndex(fetcher, sid) {
    let cfg = {
        method: 'GET',
        headers: {
            'sid': sid,
        },
        cache: 'no-store'
    };
    return fetcher.fetchRaw("/market/index", cfg)
}

function restInstall(fetcher, sid, pid) {
    let cfg = {
        method: 'POST',
        headers: {
            'sid': sid,
        },
        cache: 'no-store'
    };
    return fetcher.fetchRaw("/plugins/" + pid, cfg)
}

function restDelete(fetcher, sid, pid) {
    let cfg = {
        method: 'DELETE',
        headers: {
            'sid': sid,
        },
        cache: 'no-store'
    };
    return fetcher.fetchRaw("/plugins/" + pid, cfg)
}

function restInstallInfo(fetcher, sid, pid) {
    let cfg = {
        method: 'GET',
        headers: {
            'sid': sid,
        },
        cache: 'no-store'
    };
    return fetcher.fetchRaw("/plugins/" + pid, cfg)
}