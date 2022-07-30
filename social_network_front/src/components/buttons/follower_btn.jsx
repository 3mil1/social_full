import { Button } from "@mui/material"
import PersonAddIcon from '@mui/icons-material/PersonAdd';
import PersonRemoveIcon from '@mui/icons-material/PersonRemove';
import FollowerService from "../../utilities/follower_service";
import { useSelector } from "react-redux";

const Follow_btn = ({isPrivate}) => {
  const follower_service = FollowerService()
  const current_id =useSelector(state => state.followers.currentUserId)
  const isFollowing = follower_service.isFollowing(current_id)
  const requestSent = follower_service.isRequested(current_id)

  return (
    <div>
      { !isPrivate  ?
         <div>
           {isFollowing && <Button  className="flex" onClick={() =>follower_service.handleFollowerBtn(false)}>Stop Following <PersonRemoveIcon />  </Button>}
           {!isFollowing && <Button className="flex" onClick={() =>follower_service.handleFollowerBtn(true)} >Follow user <PersonAddIcon /> </Button>}
         </div>
       :
         <div>
            {(!isFollowing && !requestSent)  ? <Button className="flex" onClick={() =>follower_service.handleFollowerBtn(true,isPrivate)} >Follow Private user <PersonAddIcon /> </Button> :
             <div> 
              {!isFollowing ? 
              `Request has been sent` : <Button  className="flex" onClick={() =>follower_service.handleFollowerBtn(false)}>Stop Following <PersonRemoveIcon />  </Button>
            }</div>}
         </div>
         }
    </div>
  )
}

export default Follow_btn