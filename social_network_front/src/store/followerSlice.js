import { createSlice } from '@reduxjs/toolkit';

export const followerSlice = createSlice({
  name: 'follower',
  initialState: {
    updateStatus: false,
    currentUserId: null,
    followers: [],
    stalkers: [],
    sentRequests: [],
  },
  reducers: {
    updateFollowers: (state, action) => {
      state.followers = action.payload.filter(obj => {return obj.status == 1});
    },
    updateStalkers: (state, action) => {
      state.stalkers = action.payload.filter(obj => {return obj.status == 1;});
    },
    updateSentRequests: (state, action) => {
      let arr = state.sentRequests;
      if (!arr.includes(action.payload)) arr.push(action.payload);
      state.sentRequests = arr;
    },
    updateCurrentUserId: (state, action) => {
      if (action.payload == 'id') {
        state.currentUserId = '';
      } else {
        state.currentUserId = action.payload;
      }
    },
    updateStatus: (state, action) => {
      state.updateStatus = action.payload;
    },
  },
});

export const {
  updateFollowers,
  updateStalkers,
  updateSentRequests,
  updateCurrentUserId,
  updateStatus,
} = followerSlice.actions;
export default followerSlice.reducer;
