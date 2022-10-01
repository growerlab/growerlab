import React from "react";

export default function HomeLayout(props: any) {
  return (
    <div className="layout">
      <div style={{ padding: "30px 50px" }}>
        <div className="site-layout-content">{props.children}</div>
      </div>
    </div>
  );
}
