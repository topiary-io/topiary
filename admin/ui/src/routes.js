import m from "mithril"
import App from "./containers/app"
import SideBarLayout from "./containers/side-bar-layout"
import Loader from "./containers/loader"
import Edit from "./containers/edit"

const Login = {
  view:  () => m("p", "Login... (TODO)")
}

const adminLocation = ADMIN_ROOT.replace(/\/+$/,"")

const admin = {}
admin[adminLocation+"/login"] = Login
admin[adminLocation] = {
  "$container": SideBarLayout,
  "/edit/:filename": Edit,
  "/loader": Loader,
  "/": {view: ()=>m("p", "home")},
  "$default": {view:()=>m("p", "404: page not found")}
}

const routes = {
  "$container": App,
  ...admin
}

export default routes
