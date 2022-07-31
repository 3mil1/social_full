import http from "./http-common";
import { useDispatch } from "react-redux";
import { update, addAllUsers, updateAuth } from "../store/profileSlice";
import * as helper from "../helpers/HelperFuncs";
import { useNavigate } from "react-router-dom";
import FollowerService from "./follower_service";

const ProfileService = () => {
  const dispatch = useDispatch();
  const redirect = useNavigate();
  const follower_service = FollowerService();

  const checkAuth = () => {
    if (!localStorage.getItem("accessToken")) return redirect("/");
    dispatch(updateAuth(true));
    follower_service.getMyFollowers();
    getMyInfo();
  };

  const getMyInfo = async () => {
    try {
      // console.log("%cGETTING MY INFO", "color:orange");
      const response = await http.get("/user/me");
      dispatch(update({ ...response.data, id: helper.getTokenId() }));
    } catch (err) {
      helper.checkError(err);
    }
  };

  const updateProfileInfo = async (data) => {
    try {
      // console.log("%cUpdateing Profile Info", "color:orange");
      await http.put("/user/me", {
        nickname: data.nickname,
        about_me: data.about_me,
        user_img: data.user_img,
        is_private: data.is_private,
      });
      getMyInfo();
    } catch (err) {
      helper.checkError(err);
    }
  };

  const getAllUsers = async () => {
    try {
      // console.log("%cFetching All users list", "color:orange");
      const response = await http.get("user/all");
      dispatch(addAllUsers(response.data));
    } catch (err) {
      helper.checkError(err);
    }
  };


  // I need to send id(My info should be empty id) to get one user info = http://localhost:8080/user/oneuser?id=380c54e8-7560-4055-aea4-f6d7b2282d4d
  const getUserInfo = async (id) => {
    try {
      // console.log("Fetching user profile data -->", id);
      const response = await http.get(`user/oneuser?id=${id}`);
      return response.data;
    } catch (err) {
      helper.checkError(err);
    }
  };

  return {
    getMyInfo,
    updateProfileInfo,
    getAllUsers,
    getUserInfo,
    checkAuth,
  };
};

export default ProfileService;
