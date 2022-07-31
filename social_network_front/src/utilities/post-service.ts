import http from "./http-common";
import { NewPostForm } from "../components/posts/newPost";
import { PostInterface } from "../components/posts/PostList";
import { parseDate } from "../helpers/parseDate";

export default {
  async addNewPost(post: NewPostForm): Promise<PostInterface> {
    // console.log("Get post: ", post);
    try {
      const response = await http.post("post/new", {
        subject: post.title,
        content: post.content,
        image: post.imgString,
        privacy: post.privacy,
        parent_id: post.parent_id,
        access: post.userList,
      });
      // console.log("response after add new post: ", response);
      return {
        id: response.data.id,
        user_id: response.data.user_id,
        user_firstname: response.data.user_firstname,
        user_lastname: response.data.user_lastname,
        title: response.data.subject,
        content: response.data.content,
        image: response.data.image,
        created_at: parseDate(response.data.created_at),
        privacy: response.data.privacy,
      };
    } catch (e) {
      console.log(e);
      throw e;
    }
  },

  async getAllMyPosts() {
    try {
      const response = await http.get("post/oneuser");
      return response.data;
    } catch (e) {
      console.error(e);
      throw e;
    }
  },

  async getOverviewPosts() {
    try {
      const response = await http.get("post/all");
      return response.data;
    } catch (e) {
      console.error(e);
      throw e;
    }
  },

  async getAllPosts(userId: string) {
    try {
      const response = await http.get(`post/oneuser?id=${userId}`);
      // console.log("user's posts", response.data);
      return response.data;
    } catch (e) {
      console.error(e);
      throw e;
    }
  },

  async showPost(id: string) {
    try {
      const response = await http.get(`post/${id}`);
      // console.log(response.data);
      return response.data;
    } catch (e) {
      console.error(e);
      throw e;
    }
  },
};
