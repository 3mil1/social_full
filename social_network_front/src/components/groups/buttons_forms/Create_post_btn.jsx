import { Button, Input, TextareaAutosize } from "@mui/material";
import AddIcon from '@mui/icons-material/Add';
import CloseIcon from '@mui/icons-material/Close';
import { useState } from "react";
import GroupService from "../../../utilities/group_service";
import * as helper from '../../../helpers/HelperFuncs';
import "./group_buttons.scss"

const Create_post = ({id}) => {
    const group_service = GroupService()
    const [isOpen,setIsOpen] = useState(false)
    const [img,setImg] = useState(null)
    const [subject,setSubject] = useState(null)
    const [content,setContent] = useState(null)
    
    const data  =  {
        group_id    : Number(id), 
        subject     : subject,
        content     : content,
        image : img
    }

    const convertImg = async (image) => { 
        if(image.length !== 0) {
            if(helper.checkImage(image)){
                const resp = await helper.getBase64(image[0]).then(base64 => base64)
                return resp
            }
        }
    }

    
    const handleSubmit = () => { 
        if(data == null) return
        if( helper.handleInputs("subject",data.subject) && helper.handleInputs("content",data.content)) group_service.makeGroupPost(data);
    }
 
    return (
    <>  
        {!isOpen && <Button onClick={() => setIsOpen(!isOpen)}>Create Post <AddIcon/></Button>}

        {isOpen && 
        <form id="postForm" >
             <Button
              variant="contained"
              className="back_btn"
              onClick={() => {
                setIsOpen(false)
                setImg(null)
              }
            }>
              <CloseIcon />
            </Button>
            <div className="input">
                <label htmlFor="subject">Subject* : </label>
                <Input type="text" id="subject" name="subject" onClick={()=> helper.handleAfterErrorClick("subject")} onChange={(e)=>{setSubject(e.target.value)}}></Input>
            </div>
            <div className="input">
                <label htmlFor="content">Content* : </label>
                <TextareaAutosize
                    id="content"
                    type="text"
                    margin="normal"
                    variant="standard"
                    placeholder="Pla pla plapla plal....."
                    minRows={4}
                    onChange={(e)=>{setContent(e.target.value)}}
                    onClick={() => {helper.handleAfterErrorClick("content")}}
                ></TextareaAutosize>
            </div>
            <div className="inputExtra">
            <Button sx={{fontSize:"16px" }} type={"submit"} 
                    id="create_btn"
                    onClick={(e) => { 
                        e.preventDefault();
                        if(data.subject && data.content) setIsOpen(false)
                        handleSubmit()
                    }}> CREATE </Button>
                <div>
                <label className="image_btn" htmlFor="image">{!img ? "PICK IMAGE" : "IMAGE ADDED"} </label>
                <input type="file" id="image" name="image" 
                    onChange={()=>{convertImg(document.getElementById("image").files).then(res => setImg(res))}}
                    />
                </div>
            </div>
        </form>
        }
    </>
    )
}

export default Create_post