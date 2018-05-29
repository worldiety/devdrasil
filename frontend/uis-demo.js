import {Box, RaisedButton, TextField, UserInterfaceState} from "/wwt/components.js";


class UISDemo extends UserInterfaceState {

    apply() {
        this.counter = 0;

        this.setTitle("kitchen sink")

        let box = new Box()
        let firstname = new TextField();
        firstname.setCaption("firstname")
        box.add(firstname)

        let lastname = new TextField()
        lastname.setCaption("lastname")
        box.add(lastname)

        let raisedButton = new RaisedButton()
        raisedButton.setText("button")
        raisedButton.setOnClick(ev=>{
            console.log("blub");
            this.counter++;
            lastname.setText("hallo welt: "+this.counter);
        })
        box.add(raisedButton)

        this.setContent(box)
    }
}

export {UISDemo}