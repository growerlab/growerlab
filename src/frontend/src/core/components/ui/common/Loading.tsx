import React from "react";
import { EuiLoadingContent } from "@elastic/eui";
import { LineRange } from "@elastic/eui/src/components/loading/loading_content";
import { useTitle } from "../../../global/state";

interface Props {
  lines?: LineRange;
}

export default function Loading(props: Props) {
  useTitle("Loading...");

  return (
    <div>
      <EuiLoadingContent lines={props.lines} />
    </div>
  );
}
