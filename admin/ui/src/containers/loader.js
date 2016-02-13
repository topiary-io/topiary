import m from "mithril"

class Loader {
  view () {
    return m("div.loader", "Loading...")
  }
}

export default new Loader
