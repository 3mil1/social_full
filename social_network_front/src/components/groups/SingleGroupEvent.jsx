import { Button } from "@mui/material"
import { useEffect, useState } from "react";
import GroupService from "../../utilities/group_service";
import * as helper from "../../helpers/HelperFuncs"
import { useNavigate } from "react-router-dom";
import { useSelector } from "react-redux";

const SingleGroupEvent = ({data}) => {
  const group_service  = GroupService()
  const name = useSelector(state => state.profile.info.first_name + " " + state.profile.info.last_name);
  const redirect = useNavigate();
  const [past, setPast] = useState(false)
  let joined = group_service.isJoining(data.event_id)

  useEffect(()=>{
      if(!helper.timeManager.isFuture(helper.timeManager.todayDate(),data.day)){
         if(helper.timeManager.calcTime(data.time) < helper.timeManager.calcTime(helper.timeManager.todayTime())) {
          setPast(true)
      }
    }
  },[])

  const handleRequest = (nr) => { 
    group_service.sendEventReply({
      event_id : data.event_id,
      option : nr
    })
  }

  return (
    <div>
      <div className={`group_post ${past ? 'past' : ''}`}>
        <div className='header flex'>
          <div className='subject eventsOff'>{data.title} </div>
          <div
            className='author'
            onClick={e => {
              if (e.target.textContent.slice(1,-1) != name) {
                redirect(`/profile/${data.creator_id}`);
              }
            }}
          >
            ({data.creator_firstname} {data.creator_lastname})
          </div>
          <div className='event_btns'>
            {!past && (
              <>
                <Button
                  className={joined ? 'green' : ''}
                  onClick={() => {
                    if (!joined) handleRequest(1);
                  }}
                >
                  Going
                </Button>
                <Button
                  className={!joined ? 'green' : ''}
                  onClick={() => {
                    if (joined) handleRequest(2);
                  }}
                >
                  Not Going
                </Button>
              </>
            )}
          </div>
          <div className='date'>
            {!past ? (
              <div className='time'>
                <div>Date: {data.created_at} </div>
                <div>Time: {data.time} </div>
              </div>
            ) : (
              <div className='time'>Event is over </div>
            )}
          </div>
        </div>
        <div className='content flex'>{data.description}</div>
      </div>
    </div>
  );
}

export default SingleGroupEvent