import http from "./http-common";
import * as helper from "../helpers/HelperFuncs";
import FollowerService from "./follower_service";
import GroupService from "./group_service";
import { useDispatch, useSelector } from "react-redux";
import {
  updateNotifications,
  updateRespondedNotifications,
} from "../store/notificationSlice";

const NotificationService = () => {
  const group_service = GroupService();
  const follower_service = FollowerService();
  const store_notifications = useSelector(
    (state) => state.notifications.notifications
  );
  const dispatch = useDispatch();

  const handleGroupJoinRequest = (data, resp) => {
    group_service.sendGroupJoinReply({
      group_id: data.group_id,
      target_id: data.actor_id,
      status: resp,
    });
  };

  const handleGroupInvite = (data, resp) => {
    group_service.sendGroupInvitationReply({
      actor_id: data.actor_id,
      group_id: data.group_id,
      status: resp,
    });
  };

  const handleFollowerRequest = (data, resp) => {
    follower_service.changeFollowerStatusInNotification({
      target_id: data.actor_id,
      status: resp,
    });
  };

  const handleNotificationSeen = (id, nr) => {
    try {
      // console.log(`%c notifying server of seeing notifications --> ${id}`, 'color:orange');
      http.post(`/user/notification/reply?id=${id}&status=${nr}`);
    } catch (err) {
      helper.checkError(err);
    }
  };

  const updateClicked = (id) => {
    // console.log( "Handeling click and updating store , Id-> ", id);
    let replacmentList = helper.convertListToMutable(store_notifications);
    replacmentList.forEach((obj) => {
      if (obj.data.notif_id == id) obj.data.seen = 2;
    });
    dispatch(updateNotifications(replacmentList));
  };
  
  const handleRequestResponse = (notif_id, response) => {
    dispatch(updateRespondedNotifications([notif_id, response]));
  };

  return {
    handleGroupJoinRequest,
    handleGroupInvite,
    handleFollowerRequest,
    handleNotificationSeen,
    handleRequestResponse,
    updateClicked,
  };
};

export default NotificationService;
