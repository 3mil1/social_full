import { createSlice } from "@reduxjs/toolkit";

export const notificationSlice = createSlice({
  name: "notifications",
  initialState: {
    notifications: [],
    respondedNotifications: [],
    updateStatus: false,
    messages: [],
  },
  reducers: {
    updateNotifications: (state, action) => {
      let notificationList = action.payload;
      state.notifications = notificationList.reverse();
    },
    updateRespondedNotifications: (state, action) => {
      let obj = {
        id: action.payload[0],
        response: action.payload[1],
      };
      state.respondedNotifications.push(obj);
    },
    addNotification: (state, action) => {
      state.messages.push(`${action.payload}`);
    },
    removeNotification: (state, action) => {
      state.messages = state.messages.filter((s) => {
        return s !== `${action.payload}`;
      });
    },
  },
});

export const {
  updateNotifications,
  updateRespondedNotifications,
  addNotification,
  removeNotification,
} = notificationSlice.actions;
export default notificationSlice.reducer;
