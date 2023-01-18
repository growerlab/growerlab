import React from "react";
import { useTitle } from "react-use";
import { EuiLoadingContent } from "@elastic/eui";
import { LineRange } from "@elastic/eui/src/components/loading/loading_content";

import { getTitle } from "../../common/document";

interface Props {
  lines?: LineRange;
}

export default function Loading(props: Props) {
  useTitle(getTitle("Loading..."));

  return (
    <div>
      <EuiLoadingContent lines={props.lines} />
    </div>
  );
}
