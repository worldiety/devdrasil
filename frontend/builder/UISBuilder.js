import {
    AbsComponent,
    Box,
    CenterBox,
    CircularProgressIndicator,
    LayoutGrid,
    LinearLayout,
    loadScript
} from "/wwt/components.js";
import {DefaultUserInterfaceState, Main} from "/frontend/uis-default.js";
import {AppModelController} from "./AppModelController.js";
import {ViewTree} from "./ViewTree.js";
import {ViewCodeEditor} from "./ViewCodeEditor.js";
import {Class} from "./AppModel.js";
import {ViewClassEditor} from "./ViewClassEditor.js";

export {UISBuilder}

class UISBuilder extends DefaultUserInterfaceState {

    static NAME() {
        return "/builder";
    }

    constructor(app) {
        super(app);
    }

    apply() {
        super.apply();
        this.topBar.setTitle(this.getString("builder_title"));


        this.appModelController = new AppModelController();
        let viewTree = new ViewTree(this);
        viewTree.bind(this.appModelController);


        let src = "function foo(items) {\n" +
            "    var i;\n" +
            "    for (i = 0; i < items.length; i++) {\n" +
            "        alert(\"Ace Rocks \" + items[i]);\n" +
            "    }\n" +
            "}";

        let editor = new ViewCodeEditor();
        editor.setSourceCode(src);
        editor.lockLines([0, 5]);


        this.editorBox = new ViewEditorBox();
        let leftBox = new ViewLeftEditorBox();
        leftBox.add(viewTree);

        let gridView = new LinearLayout();
        gridView.addWithWeight(leftBox, "25%");
        gridView.addWithWeight(this.editorBox, "calc(75% - 1px)");


        viewTree.setOnSelectedListener(entity => {
            this.showEntity(entity);
        });

        this.setContent(gridView);
    }

    /**
     *
     * @param {Class} entity
     */
    showEntity(entity) {
        for (let c of this.editorBox.removeAll()) {
            c.destroy();
        }

        if (entity instanceof Class) {
            let view = new ViewClassEditor(this);
            view.bind(this.appModelController, entity.asType());
            this.editorBox.add(view);
            return;
        }
    }
}

class ViewEditorBox extends Box {
    constructor() {
        super();
        this.getElement().style.height = "calc(100vh - 64px)";
        this.getElement().style.backgroundColor = "white";
        this.getElement().style.borderLeft = "#d0d0d0 solid 1px";
    }
}

class ViewLeftEditorBox extends Box {
    constructor() {
        super();
        this.getElement().style.height = "calc(100vh - 64px)";
        this.getElement().style.backgroundColor = "white";
    }
}

