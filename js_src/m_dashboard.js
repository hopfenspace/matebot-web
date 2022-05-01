import {React, toast} from "./m_react.js";
import ctx from "./lib/m_ctx.js";

export default class DashboardView extends React.Component {
    static contextType = ctx;

    constructor() {
        super();

        this.state = {

        };
    }

    render() {
        return <div>
            <div className="">
                <div className=""
                     onClick={(v) => {
                         this.context.sdk.logout().then((v) => {
                            if (v) {
                                this.context.setPath({"logged_in": false, "path": ["login"]});
                            }
                         });
                     }}>
                    Logout
                </div>
            </div>
        </div>;
    }
}