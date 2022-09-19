import React from "react";
import dynamic from "next/dynamic";

const NoSSR = (props: any) => <React.Fragment>{props.children}</React.Fragment>;
export default dynamic(() => Promise.resolve(NoSSR), {
  ssr: false,
});
