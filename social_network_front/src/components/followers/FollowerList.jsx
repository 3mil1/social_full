import "./follower.scss"
import FollowerCard from "./FollowerCard"

const FollowerList = ({list, label}) => {

    return (
        <div className='FollowerList'>
            <h2>{label}</h2>
            <div className='followers_container'>
            { list.map(follower => (
                <FollowerCard 
                key={follower.user_id} 
                data={follower}
                />
            ))}
            </div>
      </div>
  )
}

export default FollowerList