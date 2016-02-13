import m from "mithril"
import { connect } from "malatium"

class NavLink {
  controller({ href }) {
    return {
      active: href === m.route(),
      config: m.route
    }
  }

  view({ active, config }, { href, key }, children) {
    return active ?
      m("li", { key },
        m("span", children)
      ) :
      m("li", { key }, 
        m("a", { href, config }, children)
      )
  }
}
const navLink = new NavLink

function navItems (items) {
  return items.map((link, key) => {
    const href = link.uri
    return m.component(navLink, { href, key }, link.text)
  })
}

class SideBar {
  view (ctrl, { sideBar }, children) {
    if (!sideBar || !sideBar.length) return m("nav", m(".loader", "Loading"))
    return m("nav",
      // break menu out into own component, sidebar retains being the
      // "connect" component / optionally: connect 
      // sidebarlayout instead 
      sideBar.map((group, key) =>
        [ m("h2", { key }, group.name),
          m("ul", navItems(group.options))]))  
  }
}

export default connect(
  ({ sideBar }) => ({ sideBar })
)(SideBar)
