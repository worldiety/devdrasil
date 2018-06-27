import {Observable} from "/wwt/components.js";
import {AppModel} from "./AppModel.js";
import {GoBackend} from "./GoBackend.js";
import {ES6Frontend} from "./ES6Frontend.js";


export {AppModelController}

class AppModelController extends Observable {

    constructor() {
        super();
        this.setValue(new AppModel("", "", new GoBackend(), new ES6Frontend()));
    }
}