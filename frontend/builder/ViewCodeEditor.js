import {AbsComponent, loadScript} from "/wwt/components.js";

export {ViewCodeEditor};

class ViewCodeEditor extends AbsComponent {
    constructor() {
        super("pre");
        this.lockedMarker = [];

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
            let Range = ace.require('ace/range').Range;

            for (let markerId of this.lockedMarker) {
                editor.session.removeMarker(markerId);
            }

            for (let lockedLine of lockLines) {
                let markerId = editor.session.addMarker(new Range(lockedLine, 0, lockedLine, 1000), "ace_active-line", "fullLine");
                this.lockedMarker.push(markerId);
            }


            editor.commands.on("exec", e => {

                let rowCol = editor.selection.getCursor();
                for (let lockedLine of lockLines) {
                    if ((rowCol.row === lockedLine)) {
                        e.preventDefault();
                        e.stopPropagation();
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