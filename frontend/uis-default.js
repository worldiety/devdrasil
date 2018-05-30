import {UserInterfaceState} from "/wwt/components.js";

export {DefaultUserInterfaceState}

class DefaultUserInterfaceState extends UserInterfaceState {
    constructor(app) {
        super(app);
    }

    apply(){
        super.apply();
        document.body.style.backgroundColor = "#f1f1f1";
    }

    attachPromise(promise){
        promise.then(res =>{
            if (res.status == 403){
                alert(this.getString("login_failed"));
            }
        });

    }
}