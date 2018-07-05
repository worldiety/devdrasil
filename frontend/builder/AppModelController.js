import {Observable} from "/wwt/components.js";
import {App} from "./AppModel.js";
import {Asset} from "./ProjectRepository.js";


export {AppModelController}

/**
 *  @type {Observable<App>}
 */
class AppModelController extends Observable {

    /**
     *
     * @param {ProjectRepository} projectRepository
     * @param {string} id
     */
    constructor(projectRepository, id) {
        super();

        let path = "/app.json";
        this.addObserver(null, val => {
            if (val == null) {
                return;
            }
            projectRepository.putAsset(new Asset(id, path, val)).then(_ => {
                console.log(`saved ${id}${path}`);
            }).catch(err => {
                console.log(`failed to save ${id}${path}:${err}`);
            });
        });


        projectRepository.getAsset(id, path).then(asset => {
            console.log(`loaded successfully app model from ${id}${path}`);
            let app = new App();
            app.fromObject(asset.payload);
            this.setValue(app);
        }).catch(err => {
            console.log(`failed to load app model from ${id}${path}:${err}`);
            this.setValue(new App());
        });
    }
}