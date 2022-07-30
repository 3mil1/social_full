import SingleNotification from "./SingleNotification"
import { useSelector } from "react-redux"
import { Box } from "@mui/system"
import { Typography } from "@mui/material"

const NotificationList = () => {
    const notifications = useSelector(state=> state.notifications.notifications)
    if(notifications == null) notifications = [];

    let responseRequired  = [];
    let allNotifications  = [];

    for ( let item in notifications){
        let type = notifications[item].action_type
        if(type == "friend request" || type == "new group member request" || type == "group invitation"){
            responseRequired.push(notifications[item])
        }else{
            allNotifications.push(notifications[item])
        }
    }

    const mapArray = (arr) => { 
        return (arr.map((notification) =>( 
                <SingleNotification key={notification.data.notif_id} data={notification}/>
            )))
    }

    return (
      <div className='notificationList_wrapper'>
        <Box>
          <Typography className='header' variant='h6'>
            {' '}
            Notifications:{' '}
          </Typography>
          {mapArray(allNotifications)}
        </Box>

        <Box>
          <Typography className='header' variant='h6'>
            Waiting Responses : 
          </Typography>
          {mapArray(responseRequired)}
        </Box>
      </div>
    );
}

export default NotificationList