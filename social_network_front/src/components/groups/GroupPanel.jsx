import {Typography } from "@mui/material"
import {useEffect, useState } from "react"
import {useParams } from "react-router-dom"
import GroupService from "../../utilities/group_service"
import StarIcon from '@mui/icons-material/Star';
import Join_group_btn from "./buttons_forms/Join_group_btn";
import Requests from "./RequestList";
import Invite_group_btn from "./buttons_forms/Invite_group_btn";
import "./group.scss";

const GroupPanel = ({isAdmin,isMember}) => {
    const group_service = GroupService()
    const [info,setInfo ] = useState({})
    let [count,setCount ] = useState(0)
    let {id} = useParams()

    useEffect(()=>{
        group_service.getGroupInfo(id).then(res => {
            setInfo(res)
            if(isAdmin){
                setCount(res.Members.length)
            }else{
                setCount(res.members)
            }
        })
    },[id])

    return (
        <>
        {isAdmin &&
        <div className="admin_panel flex">
            <div className="flex">
                <StarIcon fontSize="large" sx={{color:"yellow",margin:"0.3em",padding:"0.2em",background: 'black',     borderRadius: "50%"}}/>
                <Typography variant="h6">
                    ADMIN PANEL 
                </Typography>
            </div>
            <Requests />
        </div>}

        <h1 className="flex">Group Info</h1>
        <div className="group_panel">
            <div className="header ">
                <div className="left">
                    <Typography variant="h4">{info.title}</Typography>
                    <Typography variant="h6" className="creator">({info.creator_first_name}  {info.creator_last_name})</Typography>
                </div>
                <div className="right">
                    <Typography variant="h6">Members: {count}</Typography>
                {(!isMember && !isAdmin) ? <Join_group_btn/> : <Invite_group_btn />}
                </div>
            </div>
            <Typography variant="h6" className="introduction">Introduction: </Typography>
            <Typography variant="p"  className="content"> " {info.description}  "</Typography>
        </div>
        </>
    )
}

export default GroupPanel