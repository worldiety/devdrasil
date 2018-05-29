import {UISDemo} from "/frontend/uis-demo.js";
import {Application} from "/wwt/components.js";

class ExampleApp extends Application{

    onCreate(){
        new UISDemo(this).apply()
    }

    validateSession(){

    }
}


let app = new ExampleApp();
app.setLocale("de");
app.addTranslation("de","/frontend/values-de/strings.xml");
app.addTranslation("en","/frontend/values-en/strings.xml");
app.create();