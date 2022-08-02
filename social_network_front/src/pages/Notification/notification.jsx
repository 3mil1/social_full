import { useEffect } from "react"
import NotificationList from "../../components/notifications/NotificationList"
import FollowerService from "../../utilities/follower_service";
import GroupService from "../../utilities/group_service";

const Notification = () => {
 const follower_service = FollowerService();
 const group_service = GroupService();

  useEffect(()=>{
      follower_service.getMyFollowers();
      group_service.getCreatedGroups();
      group_service.getJoinedGroups();
  },[])
  
  return (
    <div><NotificationList /></div>
  )
}

export default Notification