import m from "mithril"
import SideBar from "./side-bar"
import Footer from "./header-footer/footer"

class SideBarLayout {
  view(ctrl, props, children) {
    return m("div.push-right",
      SideBar,
      m("main", children),
      Footer
    )
  }
}

export default new SideBarLayout
