import {useEffect, useState } from "react"
import { useSelector } from "react-redux"
import GroupService from "../../utilities/group_service"
import SingleGroupPost from "./SingleGroupPost"


const GroupPosts = ({id}) => {
  const [posts,setPosts ] = useState([])
  const group_service = GroupService()

  // const storeInfo = useSelector(state => state)
  const update  = useSelector(state =>  state.groups.updateStatus)
  useEffect(()=>{
    group_service.getGroupPosts(id).then(res => {
      setPosts(res)
    })
  },[id,update])

  return (
    <div>
    {posts && posts.map((post) => (
        <SingleGroupPost key={post.post_id} data={post}  />
      ))}
    </div>
  )
}

export default GroupPosts