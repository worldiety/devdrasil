import {Observable} from "/wwt/components.js";
import {App} from "./AppModel.js";


export {AppModelController}

/**
 *  @type {Observable<App>}
 */
class AppModelController extends Observable {

    constructor() {
        super();
        this.setValue(new App());
    }
}