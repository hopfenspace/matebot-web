import {React} from "../m_react.js";
import SDK from "./m_sdk.js";

let sdk = new SDK();
let ctx = React.createContext({"static": "/static/", "sdk": sdk});
export default ctx;