import {
  SET_SIDE_BAR,
  EDIT_FILE
} from "../actions"

// please disregard the ugly state of the state right now :blush:

export const application = (state = {}, action) => {
  console.log(action, state)
  return state
}

export const sideBar = (state = [], action) => {
  switch (action.type) {
    case SET_SIDE_BAR: return action.data 
  }
  return state
}

export const openFile = (state = {}, action) => {
  switch (action.type) {
    case EDIT_FILE: 
      const { filename, content } = action
      return { filename, content }
  }
  return state
} 
