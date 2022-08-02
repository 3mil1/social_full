import { useEffect, useState } from "react";
import { useSelector } from "react-redux";
import GroupService from "../../utilities/group_service";
import SingleGroupEvent from "./SingleGroupEvent";

const GroupEvents = ({ id }) => {
  const [events, setEvents] = useState([]);
  const group_service = GroupService();
  const update = useSelector((state) => state.groups.updateStatus);

  useEffect(() => {
    group_service.getGroupEvents(id).then((res) => {
      if (res !== null) {
        setEvents(res.reverse());
      }
    });
    group_service.getJoinedEvents();
  }, [id, update]);

  return (
    <div>
      {events &&
        events.map((event) => (
          <SingleGroupEvent key={event.event_id} data={event} />
        ))}
    </div>
  );
};

export default GroupEvents;
