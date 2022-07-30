import { createSlice } from '@reduxjs/toolkit';

export const postSlice = createSlice({
  name: 'post',
  initialState: {
    isOpen: false,
    posts: [],
    comments: [],
  },
  reducers: {
    openModal: state => {
      state.isOpen = !state.isOpen;
    },
    loadPosts: (state, action) => {
      state.posts = action.payload;
    },
    updatePosts: (state, action) => {
      state.posts.push(action.payload);
    },
    loadComments: (state, action) => {
      state.comments = action.payload;
    },
    updateComments: (state, action) => {
      state.comments.push(action.payload);
    },
  },
});

export const {
  openModal,
  loadPosts,
  loadComments,
  updatePosts,
  updateComments,
} = postSlice.actions;
export default postSlice.reducer;
