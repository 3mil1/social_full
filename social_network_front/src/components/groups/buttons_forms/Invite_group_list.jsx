import { useEffect, useState } from "react";
import "./group_buttons.scss";

const Invite_group_list = ({list,setList}) => {
  const [checkedState , setCheckedState] = useState(new Array(list.length).fill(false))

  const handleOnChange  = (position) => { 
    const updatedCheckedState = checkedState.map((item,index)=> 
      index === position ? !item : item
    )
    setCheckedState(updatedCheckedState)
  }

  useEffect(()=>{
    let list = document.getElementById("invite_list");
    let arr = []
    list.querySelectorAll('input').forEach(e => {
      if(e.checked) arr.push(e.value)
    })
    setList(arr)
  },[checkedState])

  return (
    <div className="invite_box">
        <div className="header">
          <div>User</div>
          <div>Invite</div>
        </div>
      <div id="invite_list">
        {list && 
        (list.map((user,index)=> (
          <div key={user.user_id} className="listed_user">
            <label >
            {user.first_name}
            </label>
            <input type="checkbox" value={user.user_id} checked={checkedState[index]} 
            onChange={()=> handleOnChange(index)}
            />
          </div>
          )))
        }
      </div>
    </div>
  )
}

export default Invite_group_list