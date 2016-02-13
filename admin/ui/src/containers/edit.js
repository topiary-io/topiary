import m from "mithril"
import { connect } from "malatium"
import { openFileToEdit } from "../actions"

// TODO : move to global redux state ?

const filename = m.prop("")
const content = m.prop("")

function config(el, isInit) {
  if (!isInit) {
    return m.request({
        method: "GET",
        url: "http://localhost:3000"+ADMIN_ROOT+"api/read/" + filename()
      })
      .then(function (res) {
        content(res.content)
      })
  }
}

function save(e) {
  e.preventDefault()
  m.request({
      method: "POST",
      url: "http://localhost:3000"+ADMIN_ROOT+"api/write/",
      data: {
        filename: filename(),
        content: content()
      }
    })
    .then(function (response) {
      // should gray out save until doc is changed again
      // flag "commit" button to active
      // display message that save was successful
      console.log(response.message)
    })
}

class Edit {
  controller({ openFile, openFileToEdit }) {
    filename(m.route.param("filename"))
    return {
      filename: filename,
      content: content 
    }
  }

  view({ filename, content }) {
    return m("form#edit", 
      m("h2", { config }, "Editing: " + filename()),
      m("textarea", { oninput: m.withAttr("value", content), value: content() }),
      m("button", { onclick: save }, "Save")
    )
  }
}

export default connect(
  ({ openFile }) => ({ openFile }),
  {},
  { openFileToEdit }
)(Edit) 
