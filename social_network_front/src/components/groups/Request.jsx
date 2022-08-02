import { Avatar, Button } from "@mui/material";
import GroupService from "../../utilities/group_service";

const Request = ({ data }) => {
  const group_service = GroupService();
  const handleGroupJoinRequest = (nr) => {
    group_service.sendGroupJoinReply({
      group_id: data.group_id,
      target_id: data.user_id,
      status: nr,
    });
  };

  return (
    <div className="user flex">
      <Avatar sx={{ width: 30, height: 30 }} alt="" src={data.user_img} />
      <p>
        {data.user_firstname} {data.user_lastname}
      </p>
      <div>
        <Button onClick={() => handleGroupJoinRequest(1)}>YES</Button>
        <Button onClick={() => handleGroupJoinRequest(2)}>NO</Button>
      </div>
    </div>
  );
};

export default Request;
