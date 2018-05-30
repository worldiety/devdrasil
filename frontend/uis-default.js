import {UserInterfaceState} from "/wwt/components.js";

export {DefaultUserInterfaceState}

class DefaultUserInterfaceState extends UserInterfaceState {
    constructor(app) {
        super(app);
    }

    attachPromise(promise){
        promise.then(res =>{
            if (res.status == 403){
                alert(this.getString("login_failed"));
            }
        });

    }
}