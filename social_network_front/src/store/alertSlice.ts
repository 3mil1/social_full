import { createSlice } from "@reduxjs/toolkit";

export const alertSlice = createSlice({
  name: "alert",
  initialState: {
    isOpen: false,
    text: "",
    severity: undefined,
  },
  reducers: {
    setAlert(state, action) {
      state.isOpen = action.payload.text != "";
      state.text = action.payload.text;
      state.severity = action.payload.severity;
    },
  },
});

export default alertSlice.reducer;
export const { setAlert } = alertSlice.actions;
