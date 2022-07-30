import {
  Button,
  Grid,
  Input,
  TextareaAutosize,
} from "@mui/material";
import AddIcon from "@mui/icons-material/Add";
import CloseIcon from "@mui/icons-material/Close";
import { useState } from "react";
import * as helper from "../../../helpers/HelperFuncs";
import Invite_group_list from "./Invite_group_list";
import "./group_buttons.scss";
import GroupService from "../../../utilities/group_service";
import { useSelector } from "react-redux";

export default function Make_group_btn() {
  const group_service = GroupService();
  const storeInfo = useSelector((state) => state);
  const [open, setOpen] = useState(false);
  const [list, setList] = useState([]);
  const [title, setTitle] = useState("");
  const [description, setDescription] = useState("");
  const [error, setError] = useState(false)

  const data = {
    title,
    description,
  };


  const handleSubmit = async (data) => {
    if (data == null) return;
    if (
      helper.handleInputs("title", data.title) &&
      helper.handleInputs("description", data.description)
    ) {
      const resp = await group_service.makeNewGroupRequest(data);
      if (resp.status == 200) {
        if (list.length != 0) list.forEach((userId) =>
        group_service.sendGroupInvitation(Number(resp.data), userId)
        );
        setOpen(false)
      }
      if(Object.is(resp,"Error")){
        setError(true)
        document.getElementById("title").value = ""
        setTitle("")
      }
    }
    group_service.getCreatedGroups();
  };

  return (
    <div className="make_group">
      <Button
        className="make_group_btn"
        onClick={() => {
          setTitle("");
          setDescription("");
          setOpen(true);
          setError(false)
        }}
      >
        Create New Group <AddIcon />
      </Button>

      {open && (
        <Grid className="make_group_form" container spacing={2}>
          <Grid item xs={12}>
            <label htmlFor="title">*Group Title : </label>
              <Input
              id="title"
              required
              fullWidth
              onInput={(e) => {setTitle(e.target.value); }}
              onClick={() => {helper.handleAfterErrorClick("title")}}
            />
          </Grid>
          <Grid item xs={12}>
            <label htmlFor="description">*Group Description : </label>
            <TextareaAutosize
              id="description"
              type="text"
              margin="normal"
              variant="standard"
              placeholder="Pla pla plapla pla plplaaaplall plapal....."
              style={{ width: 400 }}
              minRows={8}
              onInput={(e) => setDescription(e.target.value)}
              onClick={() => {helper.handleAfterErrorClick("description")}}
            ></TextareaAutosize>
          </Grid>
          <Grid item xs={8}>
            <Invite_group_list
              list={storeInfo.followers.followers}
              setList={setList}
            />
          </Grid>
          <Grid item xs={8}>
            {error && <div className="error">Name Unavailable</div>}
            <Button
              id="create_btn"
              onClick={() => {handleSubmit(data);}}
            >
              CREATE
            </Button>

          </Grid>
          <Button
            variant="contained"
            className="back_btn"
            onClick={() => setOpen(false)}
          >
            <CloseIcon />
          </Button>
        </Grid>
      )}
    </div>
  );
}
