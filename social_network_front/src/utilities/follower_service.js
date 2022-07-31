import http from "./http-common";
import { useDispatch, useSelector } from "react-redux";
import {
  updateFollowers,
  updateStalkers,
  updateSentRequests,
  updateCurrentUserId,
} from "../store/followerSlice";
import * as helper from "../helpers/HelperFuncs";

const FollowerService = () => {
  const dispatch = useDispatch();
  const storeInfo = useSelector((state) => state);

  const getMyFollowers = async () => {
    // console.log("%cFetching my followers", "color:orange");
    try {
      const gotFollowers = await http.get("/follower/");
      const gotStalkers = await http.get("/follower/back");
      if (!gotFollowers.data) gotFollowers.data = [];
      if (!gotStalkers.data) gotStalkers.data = [];
      dispatch(updateFollowers(helper.filterDoubles(gotFollowers.data)));
      dispatch(updateStalkers(helper.filterDoubles(gotStalkers.data)));
    } catch (err) {
      helper.checkError(err);
    }
  };

  const getUserFollowers = async (id) => {
    // console.log('Fetching users followers',"color:orange");
    try {
      const gotFollowers = await http.get(`/follower/?id=${id}`);
      // console.log('%c Back responsed with =>','color:orange',gotFollowers.data,);
      return gotFollowers.data;
    } catch (err) {
      helper.checkError(err);
    }
  };

  const getUserStalkers = async (id) => {
    // console.log('%cFetching users stalkers', "color:orange");
    try {
      const gotStalkers = await http.get(`/follower/back?id=${id}`);
      return gotStalkers.data;
    } catch (err) {
      helper.checkError(err);
    }
  };

  const sendFollowerRequest = async (id) => {
    // console.log("%c Sending follower request =>", "color:orange");
    try {
      await http.post("/follower/", { target_id: `${id}` });
      getMyFollowers();
    } catch (err) {
      helper.checkError(err);
    }
  };

  const sendUnFollowRequest = async (id) => {
    try {
      await http.delete(`/follower/?id=${id}`);
      let followers = storeInfo.followers.followers.filter(
        (obj) => obj.user_id != storeInfo.followers.currentUserId
      );
      dispatch(updateFollowers(followers));
    } catch (err) {
      helper.checkError(err);
    }
  };

  const changeFollowerStatus = async (id) => {
    try {
      const response = await http.put("/follower/", {
        target_id: `${id}`,
        status: 2,
      });
      // console.log("%c changeFollowerStatus =>", "color:orange", response.data);
    } catch (err) {
      helper.checkError(err);
    }
  };

  const changeFollowerStatusInNotification = async (data) => {
    try {
      // console.log('%c Sending follower reply in notification --> ', 'color:orange', data);
      await http.put("/follower/", data);
      // console.log("%c changeFollowerStatus =>", "color:orange", response.data);
    } catch (err) {
      helper.checkError(err);
    }
  };

  const isFollowing = () => {
    //  if there is object returned (!!), same as Boolean(storeInfo.followers.followers)
    return !!storeInfo.followers.followers.find(
      (user) => user.user_id == storeInfo.followers.currentUserId
    );
  };

  const isRequested = (id) => {
    return storeInfo.followers.sentRequests.includes(id);
  };

  const setCurrentUserId = async (id) => {
    dispatch(updateCurrentUserId(id));
  };

  const handleFollowerBtn = (isFollowing, profileStatus) => {
    if (isFollowing) {
      sendFollowerRequest(storeInfo.followers.currentUserId);
      if (profileStatus)
        dispatch(updateSentRequests(storeInfo.followers.currentUserId));
    } else {
      sendUnFollowRequest(storeInfo.followers.currentUserId);
    }
  };

  return {
    getMyFollowers,
    getUserFollowers,
    getUserStalkers,
    sendFollowerRequest,
    changeFollowerStatusInNotification,
    changeFollowerStatus,
    sendUnFollowRequest,
    setCurrentUserId,
    handleFollowerBtn,
    isFollowing,
    isRequested,
  };
};

export default FollowerService;
