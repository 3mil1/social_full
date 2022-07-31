import { Button, Input, TextareaAutosize } from "@mui/material";
import AddIcon from "@mui/icons-material/Add";
import CloseIcon from "@mui/icons-material/Close";
import { useEffect, useState } from "react";
import GroupService from "../../../utilities/group_service";
import * as helper from "../../../helpers/HelperFuncs";
import "./group_buttons.scss";

const Create_event = ({ id }) => {
  const group_service = GroupService();
  const [isOpen, setIsOpen] = useState(false);
  const [date, setDate] = useState("");
  const [times, setTimes] = useState("");
  const [title, setTitle] = useState("");
  const [description, setDescription] = useState("");

  const todayTime = () => {
    let currentTime = new Date();
    let hour = currentTime.getHours() + 1;
    let minutes = currentTime.getMinutes();
    return `${hour}:${minutes < 10 ? "0" + minutes : minutes}`;
  };

  const todayDate = () => {
    let currentTime = new Date();
    let month = currentTime.getMonth() + 1;
    month = month < 10 ? "0" + month : month;
    let day = currentTime.getDate();
    day = day < 10 ? "0" + day : day;
    let year = currentTime.getFullYear();
    var todaysDate = year + "-" + month + "-" + day;
    return todaysDate;
  };

  const isFuture = (current, target) => {
    let currentArr = current.split("-").map(Number);
    let targetArr = target.split("-").map(Number);
    if (currentArr[0] < targetArr[0]) return true;
    if (currentArr[0] == targetArr[0] && currentArr[1] < targetArr[1])
      return true;
    if (
      currentArr[0] == targetArr[0] &&
      currentArr[1] == targetArr[1] &&
      currentArr[2] < targetArr[2]
    )
      return true;
    if (current == target) setTimes(todayTime());
    return false;
  };
  const calcTime = (time) => {
    let arr = time.split(":");
    return Number(arr[0] * 60 + arr[1]);
  };

  const data = {
    group_id: Number(id),
    title: title,
    description: description,
    day: date,
    time: times,
    going_status: 1,
  };

  const handleSubmit = () => {
    if (data == null) return;
    if (
      helper.handleInputs("title", data.title) &&
      helper.handleInputs("description", data.description)
    ) {
      group_service.makeEvent(data);
      reset();
    }
  };

  const reset = () => {
    setIsOpen(false);
    setDate(todayDate());
    setTimes(todayTime());
    setTitle("");
    setDescription("");
  };

  useEffect(() => {
    setTimes(todayTime());
    setDate(todayDate());
  }, []);

  return (
    <>
      {!isOpen && (
        <Button onClick={() => setIsOpen(!isOpen)}>
          Create Event <AddIcon />
        </Button>
      )}

      {isOpen && (
        <form id="eventForm">
          <Button
            variant="contained"
            className="back_btn"
            onClick={() => setIsOpen(false)}
          >
            <CloseIcon onClick={reset} />
          </Button>
          <div className="input">
            <label htmlFor="title">Title* : </label>
            <Input
              type="text"
              id="title"
              name="title"
              onClick={() => helper.handleAfterErrorClick("title")}
              onChange={(e) => {
                setTitle(e.target.value);
              }}
            ></Input>
          </div>
          <div className="input">
            <label htmlFor="description">Description* : </label>
            <TextareaAutosize
              id="description"
              type="text"
              margin="normal"
              style={{ width: 260 }}
              variant="standard"
              placeholder="Pla pla plapla plal....."
              minRows={4}
              onChange={(e) => {
                setDescription(e.target.value);
              }}
              onClick={() => {
                helper.handleAfterErrorClick("description");
              }}
            ></TextareaAutosize>
          </div>
          <div className="input">
            <label htmlFor="day">Day* : </label>
            <input
              required
              className="selectionBox"
              id="day"
              min={todayDate()}
              type="date"
              value={date}
              onChange={(e) => {
                setDate(e.target.value);
                data.day = e.target.value;
              }}
            />
          </div>
          <div className="input">
            <label htmlFor="times">Time* : </label>
            <input
              required
              id="times"
              className="selectionBox"
              type="time"
              value={times}
              onChange={(e) => {
                if (!isFuture(todayDate(), date)) {
                  if (calcTime(todayTime()) < calcTime(e.target.value)) {
                    setTimes(e.target.value);
                  } else {
                    setTimes(todayTime());
                  }
                } else {
                  setTimes(e.target.value);
                }
                data.time = times;
              }}
            />
          </div>
          <div className="input">
            <label htmlFor="status">Status: </label>
            <select
              name="status"
              className="selectionBox"
              id="status"
              onChange={(e) => {
                data.going_status = Number(e.target.value);
              }}
            >
              <option value="1" defaultValue={"1"}>
                going
              </option>
              <option value="2">not going</option>
              <option value="3">intrested</option>
            </select>
          </div>
          <Button
            sx={{ fontSize: "16px" }}
            type={"submit"}
            id="create_btn"
            onClick={(e) => {
              e.preventDefault();
              if (data.title && data.description) setIsOpen(false);
              handleSubmit();
            }}
          >
            CREATE
          </Button>
        </form>
      )}
    </>
  );
};

export default Create_event;
