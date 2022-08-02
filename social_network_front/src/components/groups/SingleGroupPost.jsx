import { useSelector } from "react-redux";
import { useNavigate, useParams } from "react-router-dom";

export const SingleGroupPost = ({data}) => {
  const name = useSelector(state => state.profile.info.last_name + " " + state.profile.info.first_name);
  const redirect = useNavigate();
  let onGroupPage = window.location.href.split("/").indexOf("post") < 0;
  let {id} = useParams();
  
  return (
    <>
    {onGroupPage &&
    <div className="group_post" >
        <div className="header flex" >
            <div className="subject" onClick={() => { redirect(`/group/${id}/post/${data.post_id}`); }}>{data.subject}  </div>
            <div className="author" onClick={(e) => {
                if (e.target.textContent.slice(2, -1) != name){
                  redirect(`/profile/${data.user_id}`);}
                }
            }> ({data.user_firstname} {data.User_lastname})</div>
            <div className="time">{data.created_at == "" ? "???" : data.created_at} </div>
        </div>
        <div className="content flex">
            {data.image && <img className="image" src={`${data.image}`} alt="picture" />}
             {data.content}
        </div>
    </div>
    }
    </>
  )
}

export default SingleGroupPost;