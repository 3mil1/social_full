import { useEffect, useState } from "react";
import { useSelector } from "react-redux";
import { useParams } from "react-router-dom";
import GroupService from "../../utilities/group_service";
import Request from "./Request";

const Requests = () => {
  const group_service = GroupService();
  const [list, setList] = useState([]);
  const update = useSelector((state) => state.groups.updateStatus);
  let { id } = useParams();

  useEffect(() => {
    group_service.getJoinRequests(id).then((res) => {
      setList(res);
    });
  }, [id, update]);

  return (
    <div>
      {!!list ? (
        <div>
          {list.map((user) => (
            <Request key={user.user_id} data={user} />
          ))}
        </div>
      ) : (
        <div>No new group join requests</div>
      )}
    </div>
  );
};

export default Requests;
