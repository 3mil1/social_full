import http from "./http-common";
import { useDispatch, useSelector } from "react-redux";
import * as helper from "../helpers/HelperFuncs";
import {
  updateCreatedGroups,
  updateSentRequests,
  updateJoinedGroups,
  updateStatus,
  updateJoinedEvents,
  addAllGroups,
} from "../store/groupSlice";

//  make new group                                          // /group/new
//  make new group event                                    // /group/event/new
//  make new group post                                     // /group/post

// get all group and show in search bar (group sign)        // /group/all
// get all groups I created                                 // /group/mycreated
// get all groups im in (group sign)                        // /group/joined
// get group information                                    // /group/[SomeNumberHere]
// get all group POSTS                                      // /group/post/all?groupId=${id}`
// get all group EVENTS                                     // /group/event/all?groupId=${id}`
// get specific group posts                                 // /group/post/all?groupId=[some number here]
// get specific group post and comments                     // /group/post?groupId=[number]&postId=[number]
// get group of friends who i haven't send invitation yet   // /group/invite/available?groupId=[someGroupNumberHere]
// get people who wants to join to group                    // /group/join/reply?groupId=[someGroupNumberHere]

// send group invitation to user                            // /group/invite
// send group join request by user                          // /group/join
// send reply to group jon request                          // /group/join/reply
// send reply to group invitation                           // /group/invite/reply
// send reply to group event                                // /group/event/reply

const GroupService = () => {
  const dispatch = useDispatch();
  const storeInfo = useSelector((state) => state);

  const makeNewGroupRequest = async (data) => {
    try {
      // console.log("%c Posting new group --> ","color:orange", data );
      const response = await http.post("/group/new", data);
      dispatch(updateStatus(!storeInfo.groups.updateStatus));
      return response;
    } catch (err) {
      return err.name;
    }
  };

  const makeGroupPost = (data) => {
    // console.log('%c Posting new post to group --> ', 'color:orange', data);
    const response = http.post("/group/post", data);
    dispatch(updateStatus(!storeInfo.groups.updateStatus));
  };
  
  const makeCommentToPost = (data) => {
    // console.log("%c Comment post --> ", "color:orange", data);
    const response = http.post("/group/post", data);
    // dispatch(updateStatus(!storeInfo.groups.updateStatus));
  };

  const makeEvent = (data) => {
    // console.log("%c Posting new event to group --> ", "color:orange", data);
    const response = http.post("/group/event/new", data);
    dispatch(updateStatus(!storeInfo.groups.updateStatus));
  };

  const getAllGroups = async () => {
    try {
      // console.log("%c Fetching all groups --> ", "color:orange");
      const response = await http.get("/group/all");
      if (!response.data) response.data = [];
      dispatch(addAllGroups(response.data));
      return response.data;
    } catch (err) {
      helper.checkError(err);
    }
  };

  const getCreatedGroups = async () => {
    try {
      // console.log("%c Fetching my created groups --> ", "color:orange");
      const response = await http.get("/group/mycreated");
      if (response.data) dispatch(updateCreatedGroups(response.data));
    } catch (err) {
      helper.checkError(err);
    }
  };

  const getJoinedGroups = async () => {
    try {
      // console.log("%c Fetching my joined groups --> ", "color:orange");
      const response = await http.get("/group/joined");
      if (response.data) dispatch(updateJoinedGroups(response.data));
    } catch (err) {
      helper.checkError(err);
    }
  };

  const getJoinedEvents = async () => {
    try {
      // console.log("%c Fetching my joined events --> ", "color:orange");
      const response = await http.get("/group/event/joined");
      if (response.data) dispatch(updateJoinedEvents(response.data));
    } catch (err) {
      helper.checkError(err);
    }
  };

  const getGroupInfo = async (id) => {
    try {
      // console.log("%c Fetching specific group info --> ", "color:orange");
      const response = await http.get(`/group/${id}`);
      // console.log('%c getting group info--> ','color:coral',response);
      // dispatch(updateCurrentGroup(response.data));
      return response.data;
    } catch (err) {
      helper.checkError(err);
    }
  };

  const getGroupPosts = async (id) => {
    try {
      console.log("%c Fetching specific group posts --> ", "color:orange");
      const response = await http.get(`/group/post/all?groupId=${id}`);
      return response.data;
    } catch (err) {
      helper.checkError(err);
    }
  };

  const getSpecificGroupPost = async (groupId, postId) => {
    try {
      // console.log("%c Fetching specific post in group--> ", "color:orange");
      // /group/post?groupId=[number]&postId=[number]
      const response = await http.get(
        `/group/post?groupId=${groupId}&postId=${postId}`
      );
      return response.data;
    } catch (err) {
      helper.checkError(err);
    }
  };

  const getGroupEvents = async (id) => {
    try {
      console.log("%c Fetching specific group events --> ", "color:orange");
      const response = await http.get(`/group/event/all?groupId=${id}`);
      return response.data;
    } catch (err) {
      helper.checkError(err);
    }
  };

  const getAvailableFriends = async (id) => {
    try {
      console.log(
        "%c Fetching available friends to send invites--> ",
        "color:orange"
      );
      const response = await http.get(`/group/invite/available?groupId=${id}`);
      return response.data;
    } catch (err) {
      helper.checkError(err);
    }
  };

  const getJoinRequests = async (id) => {
    try {
      // console.log("%c Fetching group join requests --> ", "color:orange");
      const response = await http.get(`/group/join/reply?groupId=${id}`);
      // console.log('%c Group requests response --> ',response.data ,  'color:orange');
      return response.data;
    } catch (err) {
      helper.checkError(err);
    }
  };

  const sendGroupInvitation = async (groupId, userId) => {
    try {
      // console.log("%c Sending group invitation to user --> ", "color:orange");
      await http.post(`/group/invite`, {
        group_id: groupId,
        target_id: userId,
      });
      dispatch(updateStatus(!storeInfo.groups.updateStatus));
    } catch (err) {
      helper.checkError(err);
    }
  };

  const sendGroupInvitationReply = async (data) => {
    try {
      // console.log(
      //   "%c Sending group invitation reply--> ",
      //   "color:orange",
      //   data
      // );
      await http.put(`/group/invite/reply`, data);
      // console.log("%c group join response--> ", "color:coral", response);
      // dispatch(updateStatus(!storeInfo.groups.updateStatus));
    } catch (err) {
      helper.checkError(err);
    }
  };

  const sendGroupJoinRequest = async (id) => {
    try {
      // console.log("%c Sending join request to group --> ", "color:orange");
      await http.post(`/group/join`, {
        group_id: id,
      });
      dispatch(updateSentRequests(id));
      // console.log("%c sending group join request--> ", "color:coral", response);
    } catch (err) {
      helper.checkError(err);
    }
  };

  const sendGroupJoinReply = async (data) => {
    try {
      // console.log("%c Sending group join reply--> ", "color:orange", data);
      await http.put(`/group/join/reply`, data);
      // console.log("%c group join response--> ", "color:coral", response);
      dispatch(updateStatus(!storeInfo.groups.updateStatus));
    } catch (err) {
      helper.checkError(err);
    }
  };

  const sendEventReply = async (data) => {
    // console.log(data);
    try {
      // console.log("%c Sending event reply-> ", "color:orange", data);
      await http.post(`/group/event/reply`, data);
      // console.log("%c event reply response--> ", "color:coral", response);
      dispatch(updateStatus(!storeInfo.groups.updateStatus));
      if (data.option == 2)
        dispatch(
          updateJoinedEvents(
            storeInfo.groups.joinedEvents.filter(
              (obj) => obj.event_id != data.event_id
            )
          )
        );
    } catch (err) {
      helper.checkError(err);
    }
  };

  const isAdmin = (id) => {
    return !!storeInfo.groups.createdGroups.find((group) => group.id == id);
  };

  const isMember = (id) => {
    return !!storeInfo.groups.joinedGroups.find((group) => group.id == id);
  };

  const isRequested = (id) => {
    // return !!storeInfo.groups.sentRequests.find(group => group.id == id)
    return !!storeInfo.groups.sentRequests.includes(parseInt(id));
  };

  const isJoining = (id) => {
    return !!storeInfo.groups.joinedEvents.find(
      (event) => event.event_id == id
    );
  };

  return {
    makeNewGroupRequest,
    makeGroupPost,
    makeCommentToPost,
    makeEvent,
    getAllGroups,
    getCreatedGroups,
    getJoinedGroups,
    getJoinedEvents,
    getGroupInfo,
    getGroupPosts,
    getSpecificGroupPost,
    getGroupEvents,
    getAvailableFriends,
    getJoinRequests,
    sendGroupInvitation,
    sendGroupInvitationReply,
    sendGroupJoinRequest,
    sendGroupJoinReply,
    sendEventReply,
    isAdmin,
    isMember,
    isRequested,
    isJoining,
  };
};

export default GroupService;
