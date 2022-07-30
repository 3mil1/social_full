import * as React from "react";
import { useEffect, useState } from "react";
import Checkbox from "@mui/material/Checkbox";
import { Follower } from "./newPost";
import { List, Paper } from "@mui/material";
import ListItem from "@mui/material/ListItem";
import ListItemIcon from "@mui/material/ListItemIcon";
import ListItemText from "@mui/material/ListItemText";
import Grid from "@mui/material/Grid";

interface Inpt {
  followers: Follower[];
  sendBack: (f: readonly Follower[]) => void;
}

export function UserList(props: Inpt) {
  const [checked, setChecked] = useState<readonly Follower[]>([]);

  const setList = (value: Follower) => () => {
    const currentIndex = checked.map((c) => c.id).indexOf(value.id);
    const newChecked = [...checked];
    if (currentIndex === -1) {
      newChecked.push(value);
    } else {
      newChecked.splice(currentIndex, 1);
    }
    // console.log("Checked", newChecked);
    setChecked(newChecked);
  };

  useEffect(() => {
    props.sendBack(checked);
  }, [checked]);

  const customList = (items: readonly Follower[]) => (
    <Paper sx={{ width: 200, height: 230, overflow: "auto" }}>
      <List dense component="div" role="list">
        {items.map((value: Follower) => {
          const labelId = `user-${value.id}-label`;

          return (
            <ListItem
              key={value.id}
              role="listitem"
              // button
              onClick={setList(value)}
            >
              <ListItemIcon>
                <Checkbox
                  checked={checked.some((c) => c.id === value.id)}
                  disableRipple
                  inputProps={{
                    "aria-labelledby": labelId,
                  }}
                />
              </ListItemIcon>
              <ListItemText
                id={labelId}
                primary={` ${value.firstName} ${value.lastName}`}
              />
            </ListItem>
          );
        })}
        <ListItem />
      </List>
    </Paper>
  );

  return (
    <Grid container spacing={2} justifyContent="left" alignItems="center">
      <Grid item>{customList(props.followers)}</Grid>
    </Grid>
  );
}
