import GroupCard from "./GroupCard"
import "./group.scss"
const GroupList = ({group,myInfo}) => {
    return (
        <div className='group_list'>
            <div className='group_container'>
                { group.map(group => (
                    <GroupCard 
                    key={group.id} 
                    data={group}
                    myInfo={myInfo}
                    />
                    ))}
            </div>
        </div>
  )
}

export default GroupList