import PostList from "../../components/posts/PostList";
import "./homePage.scss";

export default function Homepage() {
  return (
    <div className="home-page">
      <PostList className="fullHeight" />
    </div>
  );
}
