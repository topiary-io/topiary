import m from "mithril"
import Header from "./header-footer/header"
import Footer from "./header-footer/footer"

class App {
  view (ctrl, props, children) {
    return m("#app",
        Header,
        children
    )
  }
}

export default new App
