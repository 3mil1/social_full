import { createSlice } from "@reduxjs/toolkit";

let initial = {
  id: "",
  first_name: "",
  last_name: "",
  email: "",
  birth_day: "",
  nickname: "",
  about_me: "",
  user_img: "",
  is_private: false,
  logout: false,
};

export const profileSlice = createSlice({
  name: "profile",
  initialState: {
    auth: false,
    info: initial,
    allUsers: [],
  },
  reducers: {
    update: (state, action) => {
      state.info = action.payload;
    },
    addAllUsers: (state, action) => {
      state.allUsers = action.payload.filter(
        (user) => state.info.id != user.ID
      );
    },
    updateAuth: (state, action) => {
      state.auth = action.payload;
    },
  },
});

export const { update, addAllUsers, updateAuth } = profileSlice.actions;
export default profileSlice.reducer;
