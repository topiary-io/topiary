
export const SET_SIDE_BAR = "SET_SIDE_BAR"
export const EDIT_FILE = "EDIT_FILE"

export function setSideBar (data) {
  return {
    type: SET_SIDE_BAR,
    data
  }
}

export function openFileToEdit ({filename, content}) {
  return {
    type: EDIT_FILE,
    filename,
    content
  }
}

