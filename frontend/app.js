import {UISDashboard} from "/frontend/uis-dashboard.js";
import {Application} from "/wwt/components.js";
import {UISLogin} from "./uis-login.js";

class ExampleApp extends Application{

    onCreate(){
        this.validateSession();
    }

    validateSession(){
        let sessionId = this.getSessionId();
        if (sessionId == null){
            new UISLogin(this).apply();
        }else{
            new UISDashboard(this).apply();
        }
    }

    setSessionId(id){
        localStorage.setItem("sessionId",id);
    }

    /*
        @return string|null
     */
    getSessionId(){
        return localStorage.getItem("sessionId");
    }
}


let app = new ExampleApp();
app.setLocale("de");
app.addTranslation("de","/frontend/values-de/strings.xml");
app.addTranslation("en","/frontend/values-en/strings.xml");
app.create();