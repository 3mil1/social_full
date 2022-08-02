import { styled } from "@mui/material/styles";
import SearchIcon from "@mui/icons-material/Search";
import InputBase from "@mui/material/InputBase";
import { Avatar, Button, Container } from "@mui/material";
import ProfileService from "../utilities/profile_service";
import { useSelector } from "react-redux";
import { useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";
import GroupService from "../utilities/group_service";
import "./styles/searchbar.scss";

const Search = styled("div")(({ theme }) => ({
  position: "relative",
  marginRight: theme.spacing(2),
  marginLeft: 0,
  width: "100%",
}));

const SearchIconWrapper = styled("div")(({ theme }) => ({
  padding: theme.spacing(0, 2),
  height: "100%",
  position: "absolute",
  pointerEvents: "none",
  display: "flex",
  alignItems: "center",
  justifyContent: "center",
}));

const StyledInputBase = styled(InputBase)(({ theme }) => ({
  color: "inherit",
  "& .MuiInputBase-input": {
    padding: theme.spacing(1, 1, 1, 0),
    paddingLeft: `calc(1em + ${theme.spacing(4)})`,
    transition: theme.transitions.create("width"),
    width: "100%",
    [theme.breakpoints.up("md")]: {
      width: "20ch",
    },
  },
}));

const Searchbar = () => {
  const profile_service = ProfileService();
  const group_service = GroupService();
  const allUsers = useSelector((state) => state.profile.allUsers);
  const allGroups = useSelector((state) => state.groups.allGroups);
  const [groups, setGroups] = useState([]);
  const [users, setUsers] = useState([]);
  const [fetched, setFetched] = useState(false);
  const [input, setInput] = useState("");
  const [option, setOption] = useState(0)
  const options = ["all", "users", "groups"]
  let redirect = useNavigate();
  
  const HandleChange = (e) => {
    setInput(e.target.value);
    if (fetched === false) {
      profile_service.getAllUsers()
      group_service.getAllGroups()
      setFetched(true);
    }
  };

  const mapArray = arr => {
    return arr.map((btn, index) => (
      <Button
        key={index}
        className={option == index ? 'btn' : ''}
        onClick={() => {
          setOption(index);
        }}
      >
        {btn}
      </Button>
    ));
  };

  useEffect(() => {
    if (input.trim() == "") {
      document.querySelector('.searched_objects').classList.add('hide');
      setFetched(false)
    } else {
      document.querySelector('.searched_objects').classList.remove('hide');
      switch(option){
        case 0 :
          setUsers(allUsers.filter(user => user.first_name.toLowerCase().includes(input.toLowerCase())))
          setGroups(allGroups.filter(group => group.title.toLowerCase().includes(input.toLowerCase())))
          break;
          case 1 : 
          setGroups([])
          setUsers(allUsers.filter(user => user.first_name.toLowerCase().includes(input.toLowerCase())))
          break;
          case 2 : 
          setUsers([])
          setGroups(allGroups.filter(group => group.title.toLowerCase().includes(input.toLowerCase())))
          break;
        }
      }
  }, [input,option]);
  
  return (
    <div>
      <div className='filter'>{mapArray(options)}</div>

      <Search className='search_wrapper'>
        <SearchIconWrapper>
          <SearchIcon />
        </SearchIconWrapper>
        <StyledInputBase
          id='search_input'
          placeholder='Search…'
          autoComplete='off'
          inputProps={{ 'aria-label': 'search' }}
          onChange={e => HandleChange(e)}
          value={input}
        />
      </Search>

      <Container className='searched_objects hide'>
        {users &&
          users.map((user, index) => (
            <div
              className='user flex'
              key={index}
              id={user.ID}
              onClick={() => {
                setInput('');
                redirect(`/profile/${user.ID}`);
              }}
            >
              {
                <Avatar
                  sx={{ width: 30, height: 30 }}
                  alt=''
                  src={user.user_img}
                />
              }
              <p>
                {user.first_name} {user.last_name}
              </p>
            </div>
          ))}
        {groups &&
          groups.map((group, index) => (
            <div
              className='user flex'
              key={index}
              id={group.id}
              onClick={() => {
                setInput('');
                redirect(`/group/${group.id}`);
              }}
            >
              <p>
                {group.title} <span className='mini'>(group)</span>
              </p>
            </div>
          ))}
      </Container>
    </div>
  );
};

export default Searchbar;
