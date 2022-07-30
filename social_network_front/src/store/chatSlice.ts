import { createSlice } from "@reduxjs/toolkit";

export const chatSlice = createSlice({
  name: "chat",
  initialState: {
    msgHistory: [],
    followers: [],
  },
  reducers: {
    loadMsgs: (state, action) => {
      state.msgHistory = action.payload || [];
    },
    addMsg: (state, action) => {
      // @ts-ignore
      state.msgHistory.push(action.payload);
    },
    addToBegining: (state, action) => {
      const prev = state.msgHistory;
      // @ts-ignore
      state.msgHistory = [...(action.payload || []), ...prev];
    },
    setFollowerList: (state, action) => {
      state.followers = action.payload;
    },
  },
});

export const { addMsg, loadMsgs, addToBegining, setFollowerList } =
  chatSlice.actions;
export default chatSlice.reducer;
