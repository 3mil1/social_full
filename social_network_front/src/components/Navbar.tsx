import { Link } from "react-router-dom";
import { Home } from "@mui/icons-material";
import { Avatar } from "@mui/material";
import "./styles/navbar.scss";
import { useSelector } from "react-redux";
import { RootState } from "../store/store";
import Logout from "./buttons/logout";
import Searchbar from "./Searchbar";
import NotificationsIcon from "@mui/icons-material/Notifications";
import ChatIcon from "@mui/icons-material/Chat";
import InsertEmoticonIcon from "@mui/icons-material/InsertEmoticon";
import NotificationService from "../utilities/notification_service";
import { useEffect, useState } from "react";

const Navbar = () => {
  const storeInfo = useSelector((state: RootState) => state);
  const notification_service = NotificationService();
  let notificationList =
    storeInfo.notifications.notifications != null
      ? storeInfo.notifications.notifications
      : [];
  let [notificationCount, setNotificationCount] = useState(
    notificationList.filter((obj) => !obj["data"]["seen"]).length
  );

  const replyServerOfNotifications = () => {
    try {
      notificationList.forEach((obj) => {
        if (obj["data"]["seen"] == 0) {
          // console.log("SENDING INFO TO SERVER ABOUT NOTIFICATION ->" , obj["data"]["notif_id"]);
          notification_service.handleNotificationSeen(
            obj["data"]["notif_id"],
            1
          );
        }
      });
    } catch (err) {
      console.log("SOME ERROR :", err);
    }
    setNotificationCount(0);
  };

  useEffect(() => {
    setNotificationCount(
      notificationList.filter((obj) => !obj["data"]["seen"]).length
    );
  }, [notificationList]);

  // @ts-ignore
  const storeProfileInfo = useSelector(
    (state: RootState) => state.profile.info
  );

  return (
    <div className="navigation">
      <Searchbar />

      <div className="profile_box">
        <Logout />
        <Avatar
          alt={`${storeProfileInfo.first_name} `}
          src={storeProfileInfo.user_img}
        />
        <p>
          {storeProfileInfo.first_name} {storeProfileInfo.last_name}
        </p>
        <button onClick={() => console.log(storeInfo)}>show storeInfo</button>
      </div>

      <Link className="link" to={"/homepage"}>
        Home <Home sx={{ marginLeft: 1 }} />
      </Link>

      <Link className="link" to={"/profile/me"}>
        Profile <InsertEmoticonIcon sx={{ marginLeft: 1 }} />
      </Link>

      <Link
        className="link"
        to={"/notifications"}
        onClick={replyServerOfNotifications}
      >
        Notifications
        <NotificationsIcon
          sx={{
            color: notificationCount !== 0 ? "red" : "",
            marginLeft: 1,
            marginRight: 1,
          }}
        />
        {notificationCount}
      </Link>

      <Link className={"link"} to={"/chat"}>
        Chat
        <ChatIcon
          sx={{
            color: storeInfo.notifications.messages.length > 0 ? "red" : "",
            marginLeft: 1,
          }}
        />
      </Link>
    </div>
  );
};

export default Navbar;
