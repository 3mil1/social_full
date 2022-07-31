import { useEffect, useState } from "react";
import { useNavigate, useParams } from "react-router-dom";
import FollowerList from "../../components/followers/FollowerList";
import ProfileInfo from "../../components/ProfileInfo";
import Make_group from "../../components/groups/buttons_forms/Make_group_btn";
import GroupList from "../../components/groups/GroupList";
// Redux
import { useSelector } from "react-redux";
// data
import FollowerService from "../../utilities/follower_service";
import GroupService from "../../utilities/group_service";
import { Box, Tab, Tabs } from "@mui/material";
import PropTypes from "prop-types";
import PostList from "../../components/posts/PostList";
import "./profile.scss";

function TabPanel(props) {
  const { children, value, index, ...other } = props;

  return (
    <div
      role="tabpanel"
      hidden={value !== index}
      id={`simple-tabpanel-${index}`}
      aria-labelledby={`simple-tab-${index}`}
      {...other}
    >
      {value === index && <Box sx={{ p: 3 }}>{children}</Box>}
    </div>
  );
}

TabPanel.propTypes = {
  children: PropTypes.node,
  index: PropTypes.number.isRequired,
  value: PropTypes.number.isRequired,
};

function a11yProps(index) {
  return {
    id: `simple-tab-${index}`,
    "aria-controls": `simple-tabpanel-${index}`,
  };
}

const Profile = () => {
  const follower_service = FollowerService();
  const group_service = GroupService();
  const storeInfo = useSelector((state) => state);
  let redirect = useNavigate();
  let updateFollowers = useSelector((state) => state.followers.updateStatus);
  let [myInfo, setMyInfo] = useState(false);
  let [followers, setFollowers] = useState(null);
  let [stalkers, setStalkers] = useState(null);
  let { id } = useParams();

  useEffect(() => {
    follower_service.setCurrentUserId(id);
    setTabValue(0);
    if (id === storeInfo.profile.info.id) {
      redirect("/profile/me");
    }
    if (id == "me") {
      setMyInfo(true);
      follower_service.getMyFollowers();
      group_service.getCreatedGroups();
      group_service.getJoinedGroups();
    } else {
      setMyInfo(false);
      follower_service.getUserFollowers(id).then((res) => {
        setFollowers(res);
      });
      follower_service.getUserStalkers(id).then((res) => {
        setStalkers(res);
      });
    }
  }, [id, updateFollowers]);

  const [tabValue, setTabValue] = useState(0);
  const handleChange = (event, newValue) => {
    setTabValue(newValue);
  };

  return (
    <Box sx={{ width: "100%" }}>
      <Box
        sx={{
          borderBottom: 1,
          borderColor: "divider",
        }}
        className={"tabMenu"}
        >
        <Tabs
          value={tabValue}
          indicatorColor="primary"
          textColor="primary"
          variant="fullWidth"
          onChange={handleChange}
        >
          <Tab label="Profile" {...a11yProps(0)} />
          <Tab label="Posts" {...a11yProps(1)} />
          <Tab label="Followers" {...a11yProps(2)} />
          {myInfo && <Tab label="Groups" {...a11yProps(3)} />}
        </Tabs>
      </Box>

      <TabPanel index={0} value={tabValue}>
        <ProfileInfo />
      </TabPanel>
      <TabPanel index={1} value={tabValue}>
        <PostList />
      </TabPanel>
      <TabPanel index={2} value={tabValue}>
        {myInfo ? (
          <>
            {storeInfo.followers.followers && (
              <FollowerList
                list={storeInfo.followers.followers}
                label={"I spy on"}
              />
            )}
            {storeInfo.followers.stalkers && (
              <FollowerList
                list={storeInfo.followers.stalkers}
                label={"My Stalkers"}
              />
            )}
          </>
        ) : (
          <>
            {followers ? (
              <FollowerList list={followers} label={"User spies on"} />
            ) : (
              <div>User doesn't follow anybody</div>
            )}
            {stalkers ? (
              <FollowerList list={stalkers} label={"User stalked by"} />
            ) : (
              <div>User doesn't have stalkers</div>
            )}
          </>
        )}
      </TabPanel>
      {myInfo && (
        <TabPanel index={3} value={tabValue}>
          <div className="groups_container">
            <div className="header">
              <h2>My created groups</h2>
              <Make_group />
            </div>
            {storeInfo.groups.createdGroups ? (
              <GroupList
                group={storeInfo.groups.createdGroups}
                myInfo={myInfo}
              />
            ) : (
              <div> No groups created</div>
            )}
            <div className="header">
              <h2>Groups I'm in</h2>
            </div>
            {storeInfo.groups.joinedGroups.length != 0 ? (
              <GroupList
                group={storeInfo.groups.joinedGroups}
                myInfo={myInfo}
              />
            ) : (
              <div> No joined groups</div>
            )}
          </div>
        </TabPanel>
      )}
    </Box>
  );
};

export default Profile;
