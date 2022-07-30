import { createSlice } from '@reduxjs/toolkit';

export const groupSlice = createSlice({
  name: 'groups',
  initialState: {
    updateStatus: false,
    // currentUserId: null,
    currentGroupInfo: {},
    createdGroups: [],
    sentRequests: [],
    joinedGroups: [],
    joinedEvents: [],
    allGroups: [],
  },
  reducers: {
    updateCurrentGroup: (state, action) => {
      state.currentGroupInfo = action.payload;
    },
    updateCreatedGroups: (state, action) => {
      state.createdGroups = action.payload;
    },
    updateSentRequests: (state, action) => {
      let arr = state.sentRequests;
      if (!arr.includes(action.payload)) arr.push(action.payload);
      // state.sentRequests = action.payload;
      state.sentRequests = arr;
    },
    updateJoinedGroups: (state, action) => {
      state.joinedGroups = action.payload;
    },
    updateJoinedEvents: (state, action) => {
      state.joinedEvents = action.payload;
    },
    updateStatus: (state, action) => {
      state.updateStatus = action.payload;
    },
    addAllGroups: (state, action) => {
      state.allGroups = action.payload
    },
  },
});

export const {
  updateCurrentGroup,
  updateCreatedGroups,
  updateSentRequests,
  updateJoinedGroups,
  updateJoinedEvents,
  updateStatus,
  addAllGroups,
} = groupSlice.actions;
export default groupSlice.reducer;
