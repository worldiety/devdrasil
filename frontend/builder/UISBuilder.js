import {AbsComponent, CenterBox, CircularProgressIndicator, LayoutGrid, loadScript} from "/wwt/components.js";
import {DefaultUserInterfaceState, Main} from "/frontend/uis-default.js";
import {AppModelController} from "./AppModelController.js";
import {ViewTree} from "./ViewTree.js";

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


        let appModelController = new AppModelController();
        let viewTree = new ViewTree(this, appModelController);
        this.setContent(viewTree);


        let src = "function foo(items) {\n" +
            "    var i;\n" +
            "    for (i = 0; i < items.length; i++) {\n" +
            "        alert(\"Ace Rocks \" + items[i]);\n" +
            "    }\n" +
            "}";

        let editor = new CodeEditor();
        editor.setSourceCode(src);
        editor.lockLines([0, 5]);


        let gridView = new LayoutGrid();
        gridView.getElement().style.padding = "0px";
        gridView.add(viewTree, 3);
        gridView.add(editor, 9);


        this.setContent(gridView);
    }

}


class CodeEditor extends AbsComponent {
    constructor() {
        super("pre");


        this.loadPromise = new Promise(resolve => {
            this._asyncInitAceEditor(resolve);
        });
    }

    _asyncInitAceEditor(resolve) {
        loadScript("/frontend/builder/ace/ace.js", () => {

            this.editor = ace.edit(this.getElement());
            this.editor.setTheme("ace/theme/xcode");
            this.editor.session.setMode("ace/mode/javascript");

            this.getElement().style.height = "calc(100vh - 64px)";
            this.getElement().style.margin = "0px";

            this._asyncInitAceBeautify(resolve);
        });
    }

    _asyncInitAceBeautify(resolve) {
        loadScript("/frontend/builder/beautify/beautify.js", () => {
            resolve(this.editor);
        });
    }

    setSourceCode(str) {
        this.getEditor().then(editor => {
            let out = js_beautify(str, {indent_size: 2, space_in_empty_paren: true});
            editor.setValue(out, -1);
        });
    }

    /**
     *
     * @param {Array<int>} lockLines
     */
    lockLines(lockLines) {
        this.getEditor().then(editor => {
            var Range = ace.require('ace/range').Range // get reference to ace/range


            editor.commands.on("exec", function (e) {
                let rowCol = editor.selection.getCursor();
                for (let lockedLine of lockLines) {
                    if ((rowCol.row === lockedLine)) {
                        e.preventDefault();
                        e.stopPropagation();

                        editor.session.addMarker(
                            new Range(lockedLine, 0, lockedLine, 0), "ace_active-line", "fullLine"
                        );
                    }
                }

            });
        });
    }

    async getSourceCode() {
        let editor = await this.getEditor();
        return editor.getValue();
    }

    /**
     * @return {Promise<Editor>}
     */
    getEditor() {
        return this.loadPromise;
    }


}

