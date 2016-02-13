import m from "mithril"
import { createStore, combineReducers } from "redux"
import Malatium from "malatium"

import routes from "./routes"
import * as reducers from "./reducers"
import { setSideBar } from "./actions"

const store = createStore(combineReducers(reducers))

// get initial app state, and then render ui
m.request({
  method: "GET",
  url: "http://localhost:3000"+ADMIN_ROOT+"api/side-bar"
}).then(function (sideBar) {
  store.dispatch(setSideBar(sideBar))
  Malatium
    .init(m, store)
    .route(document.body, ADMIN_ROOT, routes, "pathname")
})

