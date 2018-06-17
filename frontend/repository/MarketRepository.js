import {Fetcher, throwFromHTTP} from "/wwt/components.js";
import {SessionProvider} from "./DefaultRepository.js";

export {MarketRepository, Plugin, MarketIndex, Vendor}


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