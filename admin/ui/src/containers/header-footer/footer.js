import m from "mithril"

class Footer {
  view () {
    return m("footer", {}, m.trust("&copy; 2016"))
  }
}

export default new Footer
