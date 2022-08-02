import http from "./http-common";
import { RegisterForm } from "../pages/Register/register";
import {
  removeAccessToken,
  removeRefreshToken,
  setAccessToken,
  setRefreshToken,
} from "./token";

export default {
  async login(email: string, pwd: string) {
    try {
      const response = await http.post("user/signin", {
        email: email,
        password: pwd,
      });

      setAccessToken(response.data.access_token);
      setRefreshToken(response.data.refresh_token);
    } catch (err) {
      throw err;
    }
  },

  async register(user: RegisterForm) {
    try {
      // console.log(
      //   "%c Sending user registration data to server",
      //   "color:green",
      //   user
      // );
      await http.post("user/signup", {
        email: user.email,
        password: user.password,
        nickname: user.nickname,
        first_name: user.first_name,
        last_name: user.last_name,
        birth_day: user.dob,
        about_me: user.desc,
        user_img: user.image_path,
      });
    } catch (err) {
      console.log("Error caught");
      throw err;
    }
  },

  async logout() {
    http
      .delete("user/signout")
      .then(() => {
        removeAccessToken();
        removeRefreshToken();
      })
      .catch((err) => {
        console.log(err);
      });
  },
};
