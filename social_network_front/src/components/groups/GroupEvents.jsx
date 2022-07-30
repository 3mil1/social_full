import {useEffect, useState } from "react"
import { useSelector } from "react-redux"
import GroupService from "../../utilities/group_service"
import SingleGroupEvent from "./SingleGroupEvent"
import * as helper from '../../helpers/HelperFuncs';


const GroupEvents = ({id}) => {
  const [events,setEvents] = useState([])
  const group_service = GroupService()
  const update  = useSelector(state =>  state.groups.updateStatus)
  let s 
  useEffect(()=>{
    group_service.getGroupEvents(id).then(res => {
      setEvents(res.reverse());})
    group_service.getJoinedEvents()
  },[id,update])

  return (
    <div>
    {events && events.map((event) => (
        <SingleGroupEvent key={event.event_id} data={event} />
      ))}
    </div>
  )
}

export default GroupEvents